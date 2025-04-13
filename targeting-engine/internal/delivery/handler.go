package delivery

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeHTTPHandler(s Service) http.Handler {
	r := mux.NewRouter()

	deliveryHandler := httptransport.NewServer(
		makeDeliveryEndpoint(s),
		decodeDeliveryRequest,
		encodeDeliveryResponse,
	)

	r.Handle("/v1/delivery", deliveryHandler).Methods("GET")
	return r
}

func makeDeliveryEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(Request)
		return s.GetMatchingCampaigns(ctx, req)
	}
}

func decodeDeliveryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	app := r.URL.Query().Get("app")
	country := r.URL.Query().Get("country")
	os := r.URL.Query().Get("os")

	return Request{
		App:     app,
		Country: country,
		OS:      os,
	}, nil
}

func encodeDeliveryResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(Response)

	switch {
	case resp.Error == ErrMissingApp.Error() ||
		resp.Error == ErrMissingCountry.Error() ||
		resp.Error == ErrMissingOS.Error():
		w.WriteHeader(http.StatusBadRequest)
	case resp.Error == ErrNoCampaigns.Error():
		w.WriteHeader(http.StatusNoContent)
		return nil
	case resp.Error != "":
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(resp)
}
