package rule

import (
	"fmt"

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
		cel.Variable("severity", cel.StringType),
		cel.Variable("tags", cel.ListType(cel.StringType)),
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

func (cr *CelRule) Match(data map[string]any) (bool, error) {

	out, _, err := cr.program.Eval(data)
	if err != nil {
		return false, fmt.Errorf("failed to eval program: %v", err)
	}

	if val, ok := out.Value().(bool); ok && val {
		return true, nil
	}

	// The rule dit not match
	return false, nil
}
