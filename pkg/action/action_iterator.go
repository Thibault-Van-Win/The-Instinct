package action

import (
	"fmt"

	"github.com/google/cel-go/cel"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/security_context"
)

const (
	ActionTypeIterator = "iterator"
)

type IteratorAction struct {
	BaseAction
	Expression  string `json:"expression"`
	ItemVarName string `json:"item_var_name"`
	InnerAction Action `json:"inner_action"`
	StopOnError bool   `json:"stop_on_error"`

	compiledProgram cel.Program
}

func NewIteratorAction(params map[string]any, reg *ActionRegistry) (*IteratorAction, error) {
	name, ok := params["name"].(string)
	if !ok {
		return nil, fmt.Errorf("iterator action requires a name")
	}

	expression, ok := params["expression"].(string)
	if !ok {
		return nil, fmt.Errorf("iterator action requires an expression")
	}

	stopOnError, ok := params["stop_on_error"].(bool)
	if !ok {
		stopOnError = true
	}

	program, err := compileCelExpression(expression)
	if err != nil {
		return nil, fmt.Errorf("iterator has an invalid expression: %v", err)
	}

	itemVarName, ok := params["item_var_name"].(string)
	if !ok {
		return nil, fmt.Errorf("iterator action requires an item_var_name")
	}

	rawInnerAction, ok := params["inner_action"]
	if !ok {
		return nil, fmt.Errorf("iterator action requires an inner_action")
	}

	innerActionConfig, err := convertToConfig(rawInnerAction)
	if err != nil {
		return nil, fmt.Errorf("inner action config has an unexpected structure: %v", err)
	}

	innerAction, err := reg.Create(*innerActionConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create inner action: %v", err)
	}

	instance := &IteratorAction{
		BaseAction: BaseAction{
			Type: ActionTypeIterator,
			Name: name,
		},
		Expression:      expression,
		ItemVarName:     itemVarName,
		InnerAction:     innerAction,
		StopOnError:     stopOnError,
		compiledProgram: program,
	}

	if err := instance.Validate(); err != nil {
		return nil, fmt.Errorf("iterator validation failed: %v", err)
	}

	return instance, nil
}

func (ia *IteratorAction) Execute(ctx *security_context.SecurityContext) error {
	ctx.ExecutionStatus[ia.Name] = security_context.StatusRunning

	// Fetch the data	
	activation, err := ctx.ToMap()
	if err != nil {
		return fmt.Errorf("failed to convert security context to a map: %v", err)
	}

	data, _, err := ia.compiledProgram.Eval(activation)
	if err != nil {
		return fmt.Errorf("failed to evaluate expression: %v", err)
	}

	values, ok := data.Value().([]any)
	if !ok {
		return fmt.Errorf("retrieved value need to be a slice, got %T", values)
	}

	// Make sure we can restore the original value
	original, hadOriginal := ctx.Variables[ia.ItemVarName]

	// Start Loop
	for index, value := range values {
		ctx.ExecutionStatus[fmt.Sprintf("%s:iteration %d", ia.Name, index)] = security_context.StatusRunning
		// Place in ctx
		ctx.Variables[ia.ItemVarName] = value

		// Execute Action
		if err := ia.InnerAction.Execute(ctx); err != nil {
			ctx.ExecutionStatus[fmt.Sprintf("%s:iteration %d", ia.Name, index)] = security_context.StatusFailed
			ctx.ExecutionStatus[ia.Name] = security_context.StatusFailed
			if ia.StopOnError {
				break
			}

			continue
		}

		ctx.ExecutionStatus[fmt.Sprintf("%s:iteration %d", ia.Name, index)] = security_context.StatusCompleted
	}

	// Replace shadowed variable
	if hadOriginal {
		ctx.Variables[ia.ItemVarName] = original
	} else {
		delete(ctx.Variables, ia.ItemVarName)
	}

	if ctx.ExecutionStatus[ia.Name] == security_context.StatusFailed {
		return fmt.Errorf("iterator action %s failed ", ia.Name)
	}

	ctx.ExecutionStatus[ia.Name] = security_context.StatusCompleted
	return nil
}

func (ia *IteratorAction) Validate() error {
	if err := ia.BaseAction.Validate(); err != nil {
		return fmt.Errorf("basic validation failed: %v", err)
	}

	if ia.Expression == "" {
		return fmt.Errorf("iterator action %s has an empty expression", ia.Name)
	}

	if ia.ItemVarName == "" {
		return fmt.Errorf("iterator action %s has no item_var_name set", ia.Name)
	}

	if err := ia.InnerAction.Validate(); err != nil {
		return fmt.Errorf("iterator action %s has an invalid inner child: %v", ia.Name, err)
	}

	return nil
}

func compileCelExpression(expression string) (cel.Program, error) {
	// Variables are found in the security context
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

	ast, iss := env.Compile(expression)
	if iss.Err() != nil {
		return nil, fmt.Errorf("type check error for expression %s: %s", expression, iss.Err())
	}

	// Create an evaluable instance of the AST
	return env.Program(ast)
}
