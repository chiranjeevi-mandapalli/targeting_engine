package campaign

import (
	"context"
	"fmt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetCampaignByID(ctx context.Context, id string) (*Campaign, error) {
	if id == "" {
		return nil, fmt.Errorf("campaign ID cannot be empty")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetActiveCampaigns(ctx context.Context) ([]Campaign, error) {
	return s.repo.GetActive(ctx)
}

func (s *Service) GetCampaignsByIDs(ctx context.Context, ids []string) ([]Campaign, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	return s.repo.GetByIDs(ctx, ids)
}
