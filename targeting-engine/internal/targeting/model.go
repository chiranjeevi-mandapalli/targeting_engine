package targeting

import (
	"encoding/json"
	"errors"
)

type Dimension string

const (
	DimensionApp     Dimension = "app"
	DimensionCountry Dimension = "country"
	DimensionOS      Dimension = "os"
)

type Operation string

const (
	OperationInclude Operation = "include"
	OperationExclude Operation = "exclude"
)

type Rule struct {
	CampaignID string          `json:"campaign_id" db:"campaign_id"`
	Dimension  Dimension       `json:"dimension" db:"dimension"`
	Operation  Operation       `json:"operation" db:"operation"`
	Values     json.RawMessage `json:"values" db:"values"`
}

var (
	ErrInvalidRule = errors.New("invalid targeting rule")
)

func (r Rule) Validate() error {
	if r.CampaignID == "" {
		return ErrInvalidRule
	}
	if r.Dimension == "" {
		return ErrInvalidRule
	}
	if r.Operation == "" {
		return ErrInvalidRule
	}
	if len(r.Values) == 0 {
		return ErrInvalidRule
	}
	return nil
}
