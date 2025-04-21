package delivery

import (
	"context"
	"targeting-engine/internal/campaign"
	"targeting-engine/internal/targeting"
)

type Service interface {
	GetMatchingCampaigns(ctx context.Context, req Request) (Response, error)
}

type serviceImpl struct {
	campaignSvc *campaign.Service
	targeting   *targeting.RuleEngine
}

func NewService(campaignSvc *campaign.Service, targeting *targeting.RuleEngine) Service {
	return &serviceImpl{
		campaignSvc: campaignSvc,
		targeting:   targeting,
	}
}

func (s *serviceImpl) GetMatchingCampaigns(ctx context.Context, req Request) (Response, error) {
	if err := req.Validate(); err != nil {
		return Response{Error: err.Error()}, err
	}

	activeCampaigns, err := s.campaignSvc.GetActiveCampaigns(ctx)
	if err != nil {
		return Response{Error: err.Error()}, err
	}

	if len(activeCampaigns) == 0 {
		return Response{Error: ErrNoCampaigns.Error()}, ErrNoCampaigns
	}

	campaignIDs := make([]string, len(activeCampaigns))
	for i, c := range activeCampaigns {
		campaignIDs[i] = c.ID
	}

	matchedIDs, err := s.targeting.Evaluate(ctx, req.App, req.Country, req.OS, campaignIDs)
	if err != nil {
		return Response{Error: err.Error()}, err
	}

	matchedCampaigns, err := s.campaignSvc.GetCampaignsByIDs(ctx, matchedIDs)
	if err != nil {
		return Response{Error: err.Error()}, err
	}

	response := Response{
		Campaigns: make([]CampaignResponse, len(matchedCampaigns)),
	}
	for i, c := range matchedCampaigns {
		response.Campaigns[i] = CampaignResponse{
			ID:    c.ID,
			Image: c.ImageURL,
			CTA:   c.CTA,
		}
	}

	return response, nil
}
