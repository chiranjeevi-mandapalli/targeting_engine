package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"targeting-engine/internal/campaign"
	"time"
)

type CampaignCache struct {
	cache Cache
	ttl   time.Duration
}

func NewCampaignCache(cache Cache, ttl time.Duration) *CampaignCache {
	return &CampaignCache{
		cache: cache,
		ttl:   ttl,
	}
}

func (c *CampaignCache) GetCampaign(ctx context.Context, id string) (*campaign.Campaign, error) {
	data, err := c.cache.Get(ctx, campaignKey(id))
	if err != nil {
		return nil, fmt.Errorf("cache get error: %w", err)
	}
	if data == nil {
		return nil, nil
	}

	var camp campaign.Campaign
	if err := json.Unmarshal(data, &camp); err != nil {
		return nil, fmt.Errorf("cache unmarshal error: %w", err)
	}
	return &camp, nil
}

func (c *CampaignCache) SetCampaign(ctx context.Context, camp *campaign.Campaign) error {
	data, err := json.Marshal(camp)
	if err != nil {
		return fmt.Errorf("cache marshal error: %w", err)
	}
	return c.cache.Set(ctx, campaignKey(camp.ID), data, c.ttl)
}

func (c *CampaignCache) GetActiveCampaigns(ctx context.Context) ([]campaign.Campaign, error) {
	data, err := c.cache.Get(ctx, "active_campaigns")
	if err != nil {
		return nil, fmt.Errorf("cache get error: %w", err)
	}
	if data == nil {
		return nil, nil
	}

	var campaigns []campaign.Campaign
	if err := json.Unmarshal(data, &campaigns); err != nil {
		return nil, fmt.Errorf("cache unmarshal error: %w", err)
	}
	return campaigns, nil
}

func (c *CampaignCache) SetActiveCampaigns(ctx context.Context, campaigns []campaign.Campaign) error {
	data, err := json.Marshal(campaigns)
	if err != nil {
		return fmt.Errorf("cache marshal error: %w", err)
	}
	return c.cache.Set(ctx, "active_campaigns", data, c.ttl)
}

func (c *CampaignCache) InvalidateCampaign(ctx context.Context, id string) error {
	if err := c.cache.Delete(ctx, campaignKey(id)); err != nil {
		return fmt.Errorf("cache delete error: %w", err)
	}
	return c.cache.Delete(ctx, "active_campaigns")
}

func campaignKey(id string) string {
	return fmt.Sprintf("campaign:%s", id)
}
