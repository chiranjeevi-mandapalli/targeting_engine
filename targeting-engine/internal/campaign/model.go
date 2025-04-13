package campaign

import (
	"time"
)

type Status string

const (
	StatusActive   Status = "ACTIVE"
	StatusInactive Status = "INACTIVE"
)

type Campaign struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	ImageURL  string    `json:"image_url" db:"image_url"`
	CTA       string    `json:"cta" db:"cta"`
	Status    Status    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (c Campaign) IsActive() bool {
	return c.Status == StatusActive
}
