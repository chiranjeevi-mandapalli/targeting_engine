package health

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func MakeHandler(service *HealthService) http.Handler {
	r := chi.NewRouter()

	r.Get("/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		response := service.Check(r.Context())

		w.Header().Set("Content-Type", "application/json")
		if response.Status == StatusDown {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		json.NewEncoder(w).Encode(response)
	})

	return r
}
