package rule

import (
	"fmt"
	"log"

	"github.com/google/cel-go/cel"
)

type CelRule struct {
	Expression string
	program    cel.Program
}

func NewCelRule(expression string) *CelRule {
	instance := &CelRule{
		Expression: expression,
	}

	env, _ := cel.NewEnv(
		cel.Variable("severity", cel.StringType),
		cel.Variable("tags", cel.ListType(cel.StringType)),
	)

	ast, iss := env.Compile(instance.Expression)
	if iss.Err() != nil {
		log.Fatalf("Failed to compile cel rule expression: %s", instance.Expression)
	}

	// Create an evaluable instance of the AST
	prog, err := env.Program(ast)
	if err != nil {
		log.Fatalf("Failed to create Program: %v", err)
	}

	instance.program = prog
	return instance
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
