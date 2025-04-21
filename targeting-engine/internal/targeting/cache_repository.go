package targeting

import (
	"context"
	"encoding/json"
	"time"

	"targeting-engine/internal/monitoring"

	"github.com/go-redis/redis/v8"
)

type CachedRuleRepository struct {
	repo    RuleRepository
	cache   *redis.Client
	ttl     time.Duration
	prefix  string
	metrics *monitoring.Metrics // Added for metrics tracking
}

// NewCachedRuleRepository creates a new CachedRuleRepository with metrics support
func NewCachedRuleRepository(repo RuleRepository, cache *redis.Client, ttl time.Duration, metrics *monitoring.Metrics) *CachedRuleRepository {
	return &CachedRuleRepository{
		repo:    repo,
		cache:   cache,
		ttl:     ttl,
		prefix:  "targeting:",
		metrics: metrics,
	}
}

func (r *CachedRuleRepository) cacheKey(campaignID string) string {
	return r.prefix + campaignID
}

func (r *CachedRuleRepository) GetByCampaignID(ctx context.Context, campaignID string) ([]Rule, error) {
	key := r.cacheKey(campaignID)
	start := time.Now()

	// Try to get from cache
	val, err := r.cache.Get(ctx, key).Bytes()
	if err == nil {
		var rules []Rule
		if err := json.Unmarshal(val, &rules); err == nil {
			r.metrics.IncrementCacheHits("rule")
			r.metrics.ObserveCacheLatency("rule", time.Since(start).Seconds())
			return rules, nil
		}
	}

	// Cache miss
	if err == redis.Nil {
		r.metrics.IncrementCacheMisses("rule")
	}

	// Fetch from repository
	rules, err := r.repo.GetByCampaignID(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	if len(rules) > 0 {
		val, err := json.Marshal(rules)
		if err == nil {
			r.cache.Set(ctx, key, val, r.ttl)
		}
	}

	r.metrics.ObserveCacheLatency("rule", time.Since(start).Seconds())
	return rules, nil
}

func (r *CachedRuleRepository) GetByCampaignIDs(ctx context.Context, campaignIDs []string) ([]Rule, error) {
	return r.repo.GetByCampaignIDs(ctx, campaignIDs)
}

func (r *CachedRuleRepository) Store(ctx context.Context, rule *Rule) error {
	if err := r.repo.Store(ctx, rule); err != nil {
		return err
	}
	return r.cache.Del(ctx, r.cacheKey(rule.CampaignID)).Err()
}

func (r *CachedRuleRepository) DeleteByCampaign(ctx context.Context, campaignID string) error {
	if err := r.repo.DeleteByCampaign(ctx, campaignID); err != nil {
		return err
	}
	return r.cache.Del(ctx, r.cacheKey(campaignID)).Err()
}
