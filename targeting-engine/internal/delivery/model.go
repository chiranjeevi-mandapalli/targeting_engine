package delivery

import "errors"

type Request struct {
	App     string `json:"app"`
	Country string `json:"country"`
	OS      string `json:"os"`
}

type Response struct {
	Campaigns []CampaignResponse `json:"campaigns,omitempty"`
	Error     string             `json:"error,omitempty"`
}

type CampaignResponse struct {
	ID    string `json:"cid"`
	Image string `json:"img"`
	CTA   string `json:"cta"`
}

var (
	ErrMissingApp     = errors.New("missing app parameter")
	ErrMissingCountry = errors.New("missing country parameter")
	ErrMissingOS      = errors.New("missing os parameter")
	ErrNoCampaigns    = errors.New("no matching campaigns found")
)

func (r Request) Validate() error {
	if r.App == "" {
		return ErrMissingApp
	}
	if r.Country == "" {
		return ErrMissingCountry
	}
	if r.OS == "" {
		return ErrMissingOS
	}
	return nil
}
