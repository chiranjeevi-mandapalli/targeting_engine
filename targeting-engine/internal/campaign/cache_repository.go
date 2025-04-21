package campaign

import (
	"context"
	"encoding/json"
	"time"

	"targeting-engine/internal/monitoring"

	"github.com/go-redis/redis/v8"
)

type CachedRepository struct {
	repo    Repository
	cache   *redis.Client
	ttl     time.Duration
	prefix  string
	metrics *monitoring.Metrics // Added for metrics tracking
}

// NewCachedRepository creates a new CachedRepository with metrics support
func NewCachedRepository(repo Repository, cache *redis.Client, ttl time.Duration, metrics *monitoring.Metrics) *CachedRepository {
	return &CachedRepository{
		repo:    repo,
		cache:   cache,
		ttl:     ttl,
		prefix:  "campaign:",
		metrics: metrics,
	}
}

func (r *CachedRepository) cacheKey(id string) string {
	return r.prefix + id
}

func (r *CachedRepository) GetByID(ctx context.Context, id string) (*Campaign, error) {
	key := r.cacheKey(id)
	start := time.Now()

	// Try to get from cache
	val, err := r.cache.Get(ctx, key).Bytes()
	if err == nil {
		var c Campaign
		if err := json.Unmarshal(val, &c); err == nil {
			r.metrics.IncrementCacheHits("campaign")
			r.metrics.ObserveCacheLatency("campaign", time.Since(start).Seconds())
			return &c, nil
		}
	}

	// Cache miss
	if err == redis.Nil {
		r.metrics.IncrementCacheMisses("campaign")
	}

	// Fetch from repository
	c, err := r.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if c != nil {
		val, err := json.Marshal(c)
		if err == nil {
			r.cache.Set(ctx, key, val, r.ttl)
		}
	}

	r.metrics.ObserveCacheLatency("campaign", time.Since(start).Seconds())
	return c, nil
}

func (r *CachedRepository) GetActive(ctx context.Context) ([]Campaign, error) {
	key := r.cacheKey("active")
	start := time.Now()

	// Try to get from cache
	val, err := r.cache.Get(ctx, key).Bytes()
	if err == nil {
		var campaigns []Campaign
		if err := json.Unmarshal(val, &campaigns); err == nil {
			r.metrics.IncrementCacheHits("campaign")
			r.metrics.ObserveCacheLatency("campaign", time.Since(start).Seconds())
			return campaigns, nil
		}
	}

	// Cache miss
	if err == redis.Nil {
		r.metrics.IncrementCacheMisses("campaign")
	}

	// Fetch from repository
	campaigns, err := r.repo.GetActive(ctx)
	if err != nil {
		return nil, err
	}

	val, err = json.Marshal(campaigns)
	if err == nil {
		r.cache.Set(ctx, key, val, r.ttl)
	}

	r.metrics.ObserveCacheLatency("campaign", time.Since(start).Seconds())
	return campaigns, nil
}

func (r *CachedRepository) GetByIDs(ctx context.Context, ids []string) ([]Campaign, error) {
	return r.repo.GetByIDs(ctx, ids)
}

func (r *CachedRepository) CountActiveCampaigns(ctx context.Context) (int, error) {
	return r.repo.CountActiveCampaigns(ctx)
}
