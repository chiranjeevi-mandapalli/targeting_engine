package targeting

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type CachedRuleRepository struct {
	repo   RuleRepository
	cache  *redis.Client
	ttl    time.Duration
	prefix string
}

func NewCachedRuleRepository(repo RuleRepository, cache *redis.Client, ttl time.Duration) *CachedRuleRepository {
	return &CachedRuleRepository{
		repo:   repo,
		cache:  cache,
		ttl:    ttl,
		prefix: "targeting:",
	}
}

func (r *CachedRuleRepository) cacheKey(campaignID string) string {
	return r.prefix + campaignID
}

func (r *CachedRuleRepository) GetByCampaignID(ctx context.Context, campaignID string) ([]Rule, error) {
	key := r.cacheKey(campaignID)
	val, err := r.cache.Get(ctx, key).Bytes()
	if err == nil {
		var rules []Rule
		if err := json.Unmarshal(val, &rules); err == nil {
			return rules, nil
		}
	}

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
