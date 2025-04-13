package health

import "context"

type Status string

const (
	StatusUp   Status = "UP"
	StatusDown Status = "DOWN"
)

type HealthResponse struct {
	Status  Status            `json:"status"`
	Details map[string]Status `json:"details,omitempty"`
}

type Checker interface {
	Check(ctx context.Context) error
}
