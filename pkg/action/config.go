package action

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

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

func convertToConfig(data any) (*ActionConfig, error) {
	var config ActionConfig
	err := mapstructure.Decode(data, &config)
	if err != nil {
		return nil, fmt.Errorf("child config has an unexpected structure: %v", err)
	}

	return &config, nil
}
