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

func MakeHTTPHandler(s *Service) http.Handler {
	r := mux.NewRouter()

	r.Methods("GET").Path("/campaigns/{id}").Handler(httptransport.NewServer(
		makeGetCampaignEndpoint(s),
		decodeGetCampaignRequest,
		encodeResponse,
	))

	r.Methods("GET").Path("/campaigns").Handler(httptransport.NewServer(
		makeGetActiveCampaignsEndpoint(s),
		decodeGetActiveCampaignsRequest,
		encodeResponse,
	))

	return r
}

// Endpoints
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

// Request/Response types
type getCampaignRequest struct {
	ID string
}

// Decoders
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

// Encoder
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
