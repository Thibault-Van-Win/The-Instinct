package rule

import (
	"fmt"

	"github.com/google/cel-go/cel"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/security_context"
)

const (
	RuleTypeCel string = "cel"
)

type CelRule struct {
	BaseRule
	Expression string `json:"expression"`
	program    cel.Program
}

func NewCelRule(params map[string]any) (*CelRule, error) {
	expression, ok := params["expression"].(string)
	if !ok {
		return nil, fmt.Errorf("cel rules requires an expression")
	}

	instance := &CelRule{
		BaseRule: BaseRule{
			Type: RuleTypeCel,
		},
		Expression: expression,
	}

	env, _ := cel.NewEnv(
		cel.Variable("event", cel.MapType(
			cel.StringType,
			cel.DynType,
		)),
		cel.Variable("variables", cel.MapType(
			cel.StringType,
			cel.DynType,
		)),
	)

	ast, iss := env.Compile(instance.Expression)
	if iss.Err() != nil {
		return nil, fmt.Errorf("type check error for expression %s: %s", instance.Expression, iss.Err())
	}

	// Create an evaluable instance of the AST
	prog, err := env.Program(ast)
	if err != nil {
		return nil, fmt.Errorf("failed to create Program: %v", err)
	}

	instance.program = prog
	return instance, nil
}

func (cr *CelRule) Match(ctx *security_context.SecurityContext) (bool, error) {
	activation, err := ctx.ToMap()
	if err != nil {
		return false, fmt.Errorf("failed to convert security context to a map: %v", err)
	}

	out, _, err := cr.program.Eval(activation)
	if err != nil {
		// Ignore matching errors
		// These include if the field does not exists
		return false, nil
	}

	if val, ok := out.Value().(bool); ok && val {
		return true, nil
	}

	// The rule dit not match
	return false, nil
}

func (cr *CelRule) Validate() error {
	if err := cr.BaseRule.Validate(); err != nil {
		return fmt.Errorf("base validation failed: %v", err)
	}

	if cr.Expression == "" {
		return fmt.Errorf("a cel rule must have an expression")
	}

	return nil
}
