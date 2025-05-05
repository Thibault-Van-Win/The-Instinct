package rule

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type RuleConfig struct {
	Type   string         `yaml:"type" json:"type" bson:"type"`
	Params map[string]any `yaml:"params" json:"params" bson:"params"`
}

func (rc *RuleConfig) Validate() error {
	ruleReg := NewRuleRegistry(
		WithStandardRules(),
	)

	_, err := ruleReg.Create(*rc)

	return err
}

func ConvertToConfig(data any) (*RuleConfig, error) {
	var config RuleConfig
	err := mapstructure.Decode(data, &config)
	if err != nil {
		return nil, fmt.Errorf("child config has an unexpected structure: %v", err)
	}

	return &config, nil
}
