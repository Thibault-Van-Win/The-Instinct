package action

// ActionConfig represents the structure of an action configuration
type ActionConfig struct {
	Type   string         `yaml:"type" json:"type"`
	Params map[string]any `yaml:"params" json:"params"`
}

func (ac *ActionConfig) Validate() error {
	actionReg := NewActionRegistry(
		WithStandardActions(),
	)

	_, err := actionReg.Create(*ac)

	return err
}
