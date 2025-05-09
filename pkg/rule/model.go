package rule

import (
	"errors"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/security_context"
)

type Rule interface {
	Match(ctx *security_context.SecurityContext) (bool, error)
	GetType() string
	Validate() error
}

// BaseRule provides common functionality for all rules
type BaseRule struct {
	Type string `json:"type"`
}

func (br *BaseRule) GetType() string {
	return br.Type
}

func (br *BaseRule) Validate() error {
	if br.Type == "" {
		return errors.New("missing type: all rules must have a type")
	}

	return nil
}
