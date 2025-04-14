package campaign

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type CachedRepository struct {
	repo   Repository
	cache  *redis.Client
	ttl    time.Duration
	prefix string
}

func NewCachedRepository(repo Repository, cache *redis.Client, ttl time.Duration) *CachedRepository {
	return &CachedRepository{
		repo:   repo,
		cache:  cache,
		ttl:    ttl,
		prefix: "campaign:",
	}
}

func (r *CachedRepository) cacheKey(id string) string {
	return r.prefix + id
}

func (r *CachedRepository) GetByID(ctx context.Context, id string) (*Campaign, error) {
	key := r.cacheKey(id)
	val, err := r.cache.Get(ctx, key).Bytes()
	if err == nil {
		var c Campaign
		if err := json.Unmarshal(val, &c); err == nil {
			return &c, nil
		}
	}

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

	return c, nil
}

func (r *CachedRepository) GetActive(ctx context.Context) ([]Campaign, error) {
	key := r.cacheKey("active")
	val, err := r.cache.Get(ctx, key).Bytes()
	if err == nil {
		var campaigns []Campaign
		if err := json.Unmarshal(val, &campaigns); err == nil {
			return campaigns, nil
		}
	}

	campaigns, err := r.repo.GetActive(ctx)
	if err != nil {
		return nil, err
	}

	val, err = json.Marshal(campaigns)
	if err == nil {
		r.cache.Set(ctx, key, val, r.ttl)
	}

	return campaigns, nil
}

func (r *CachedRepository) GetByIDs(ctx context.Context, ids []string) ([]Campaign, error) {
	return r.repo.GetByIDs(ctx, ids)
}
