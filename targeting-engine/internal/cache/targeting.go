package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"targeting-engine/internal/targeting"
	"time"
)

type TargetingCache struct {
	cache Cache
	ttl   time.Duration
}

func NewTargetingCache(cache Cache, ttl time.Duration) *TargetingCache {
	return &TargetingCache{
		cache: cache,
		ttl:   ttl,
	}
}

func (t *TargetingCache) GetRules(ctx context.Context, campaignID string) ([]targeting.Rule, error) {
	data, err := t.cache.Get(ctx, rulesKey(campaignID))
	if err != nil {
		return nil, fmt.Errorf("cache get error: %w", err)
	}
	if data == nil {
		return nil, nil
	}

	var rules []targeting.Rule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("cache unmarshal error: %w", err)
	}
	return rules, nil
}

func (t *TargetingCache) SetRules(ctx context.Context, campaignID string, rules []targeting.Rule) error {
	data, err := json.Marshal(rules)
	if err != nil {
		return fmt.Errorf("cache marshal error: %w", err)
	}
	return t.cache.Set(ctx, rulesKey(campaignID), data, t.ttl)
}

func (t *TargetingCache) InvalidateRules(ctx context.Context, campaignID string) error {
	return t.cache.Delete(ctx, rulesKey(campaignID))
}

func rulesKey(campaignID string) string {
	return fmt.Sprintf("targeting_rules:%s", campaignID)
}
