package rule

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
