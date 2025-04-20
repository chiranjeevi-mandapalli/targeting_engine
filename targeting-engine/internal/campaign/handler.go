package campaign

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeGetHTTPHandler(s *Service) http.Handler {
	r := mux.NewRouter()

	campaignHandler := httptransport.NewServer(
		makeGetCampaignEndpoint(s),
		decodeGetCampaignRequest,
		encodeResponse,
	)
	r.Handle("/v1/campaigns/{id}", campaignHandler).Methods("GET")
	return r
}

func MakeGetAllHTTPHandler(s *Service) http.Handler {
	r := mux.NewRouter()

	campaignGetAllHandler := httptransport.NewServer(
		makeGetActiveCampaignsEndpoint(s),
		decodeGetActiveCampaignsRequest,
		encodeResponse,
	)

	r.Handle("/v1/campaigns", campaignGetAllHandler).Methods("GET")
	return r
}

func makeGetCampaignEndpoint(s *Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getCampaignRequest)
		return s.GetCampaignByID(ctx, req.ID)
	}
}

func makeGetActiveCampaignsEndpoint(s *Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return s.GetActiveCampaigns(ctx)
	}
}

type getCampaignRequest struct {
	ID string
}

func decodeGetCampaignRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, fmt.Errorf("missing campaign ID")
	}
	return getCampaignRequest{ID: id}, nil
}

func decodeGetActiveCampaignsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
