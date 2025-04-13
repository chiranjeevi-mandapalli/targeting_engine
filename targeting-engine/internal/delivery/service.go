package delivery

import (
	"context"
	"targeting-engine/internal/campaign"
	"targeting-engine/internal/targeting"
)

type Service struct {
	campaignSvc *campaign.Service
	targeting   *targeting.Evaluator
}

func NewService(campaignSvc *campaign.Service, targeting *targeting.Evaluator) *Service {
	return &Service{
		campaignSvc: campaignSvc,
		targeting:   targeting,
	}
}

func (s *Service) GetMatchingCampaigns(ctx context.Context, req Request) (Response, error) {
	if err := req.Validate(); err != nil {
		return Response{Error: err.Error()}, err
	}

	// Get all active campaigns
	activeCampaigns, err := s.campaignSvc.GetActiveCampaigns(ctx)
	if err != nil {
		return Response{Error: err.Error()}, err
	}

	if len(activeCampaigns) == 0 {
		return Response{Error: ErrNoCampaigns.Error()}, ErrNoCampaigns
	}

	// Get campaign IDs for targeting evaluation
	campaignIDs := make([]string, len(activeCampaigns))
	for i, c := range activeCampaigns {
		campaignIDs[i] = c.ID
	}

	// Evaluate which campaigns match the targeting rules
	matchedIDs, err := s.targeting.Evaluate(ctx, req.App, req.Country, req.OS, campaignIDs)
	if err != nil {
		return Response{Error: err.Error()}, err
	}

	// Get full campaign details for matched IDs
	matchedCampaigns, err := s.campaignSvc.GetCampaignsByIDs(ctx, matchedIDs)
	if err != nil {
		return Response{Error: err.Error()}, err
	}

	// Convert to API response format
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
