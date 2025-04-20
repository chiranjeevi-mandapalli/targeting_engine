package targeting

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type Evaluator struct {
	ruleRepo RuleRepository
}

func NewEvaluator(ruleRepo RuleRepository) *Evaluator {
	return &Evaluator{ruleRepo: ruleRepo}
}

func (e *Evaluator) Evaluate(ctx context.Context, app, country, os string, campaignIDs []string) ([]string, error) {
	rules, err := e.ruleRepo.GetByCampaignIDs(ctx, campaignIDs)
	if err != nil {
		return nil, fmt.Errorf("error getting rules: %w", err)
	}

	rulesByCampaign := groupRulesByCampaign(rules)
	var matchingCampaigns []string

	for _, campaignID := range campaignIDs {
		if matchesAllRules(rulesByCampaign[campaignID], app, country, os) {
			matchingCampaigns = append(matchingCampaigns, campaignID)
		}
	}

	return matchingCampaigns, nil
}

func groupRulesByCampaign(rules []Rule) map[string][]Rule {
	result := make(map[string][]Rule)
	for _, rule := range rules {
		result[rule.CampaignID] = append(result[rule.CampaignID], rule)
	}
	return result
}

func matchesAllRules(rules []Rule, app, country, os string) bool {
	if len(rules) == 0 {
		return true
	}

	for _, rule := range rules {
		if !matchesRule(rule, app, country, os) {
			return false
		}
	}
	return true
}

func matchesRule(rule Rule, app, country, os string) bool {
	var value string
	switch rule.Dimension {
	case DimensionApp:
		value = app
	case DimensionCountry:
		value = country
	case DimensionOS:
		value = os
	default:
		return true
	}

	var values []string
	if err := json.Unmarshal(rule.Values, &values); err != nil {
		// Fail open: if values can't be decoded, don't filter out the campaign
		return true
	}

	for _, v := range values {
		if strings.EqualFold(v, value) {
			return rule.Operation == OperationInclude
		}
	}
	return rule.Operation == OperationExclude
}
