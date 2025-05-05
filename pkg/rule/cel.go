package rule

import (
	"encoding/json"
	"fmt"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/security_context"
	"github.com/google/cel-go/cel"
)

type CelRule struct {
	Expression string
	program    cel.Program
}

func NewCelRule(expression string) (*CelRule, error) {
	instance := &CelRule{
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
	// Cannot give the ctx directly as the eval expect a map[string]any
	// Use the json tags to marshall this into a map
	// Another option would be to create a new map, benefits:
	// 	- Faster
	//	- Type safety
	// Negatives:
	//	- Need to add a new field for each extension
	jsonBytes, err := json.Marshal(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to marshal context: %v", err)
	}

	var activation map[string]any
	if err := json.Unmarshal(jsonBytes, &activation); err != nil {
		return false, fmt.Errorf("failed to unmarshal into map: %v", err)
	}

	out, _, err := cr.program.Eval(activation)
	if err != nil {
		return false, fmt.Errorf("failed to eval program: %v", err)
	}

	if val, ok := out.Value().(bool); ok && val {
		return true, nil
	}

	// The rule dit not match
	return false, nil
}
