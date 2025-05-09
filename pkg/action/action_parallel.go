package action

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/security_context"
	"github.com/mitchellh/mapstructure"
)

const (
	ActionTypeParallel = "parallel"
)

type ParallelAction struct {
	BaseAction
	Children []Action `json:"children"`
}

func NewParallelAction(params map[string]any, reg *ActionRegistry) (*ParallelAction, error) {
	name, ok := params["name"].(string)
	if !ok {
		return nil, fmt.Errorf("a parallel action requires a name")
	}

	childConfigs, ok := params["children"].([]any)
	if !ok {
		return nil, fmt.Errorf("a parallel action requires children action configs, got %T", params["children"])
	}

	children := make([]Action, 0, len(childConfigs))
	for _, rawChildConfig := range childConfigs {
		// Need to convert to an ActionConfig
		var childConfig ActionConfig
		err := mapstructure.Decode(rawChildConfig, &childConfig)
		if err != nil {
			return nil, fmt.Errorf("child config has an unexpected structure: %v", err)
		}

		child, err := reg.Create(childConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create child: %v", err)
		}
		children = append(children, child)
	}

	instance := &ParallelAction{
		BaseAction: BaseAction{
			Type: ActionTypeParallel,
			Name: name,
		},
		Children: children,
	}

	if err := instance.Validate(); err != nil {
		return nil, fmt.Errorf("parallel validation failed: %v", err)
	}

	return instance, nil
}

func (pa *ParallelAction) Execute(ctx *security_context.SecurityContext) error {
	ctx.ExecutionStatus[pa.Name] = security_context.StatusRunning

	errs := make(chan error, len(pa.Children))
	var wg sync.WaitGroup

	// Mutex for thread safe context updates
	var ctxMutex sync.Mutex

	for _, child := range pa.Children {
		wg.Add(1)

		go func(action Action) {
			defer wg.Done()

			localCtx := &security_context.SecurityContext{
				Event:           ctx.Event,
				Variables:       make(map[string]any),
				ExecutionStatus: make(map[string]security_context.Status),
			}

			// Copy over values to avoid race conditions
			ctxMutex.Lock()
			for k, v := range ctx.Variables {
				localCtx.Variables[k] = v
			}
			ctxMutex.Unlock()

			err := action.Execute(localCtx)

			// Merge results back to main context
			ctxMutex.Lock()
			for k, v := range localCtx.Variables {
				ctx.Variables[k] = v
			}

			for k, v := range localCtx.ExecutionStatus {
				ctx.ExecutionStatus[k] = v
			}
			ctxMutex.Unlock()

			if err != nil {
				errs <- err
			}
		}(child)
	}

	wg.Wait()
	close(errs)

	var errList []error
	for err := range errs {
		errList = append(errList, err)
	}

	if len(errList) > 0 {
		ctx.ExecutionStatus[pa.Name] = "failed"
		return fmt.Errorf("parallel action %s had %d failures: %w", pa.Name, len(errList), errors.Join(errList...))
	}

	ctx.ExecutionStatus[pa.Name] = security_context.StatusCompleted
	return nil
}

func (pa *ParallelAction) Validate() error {
	if len(pa.Children) == 0 {
		return fmt.Errorf("parallel action %s has no children", pa.Name)
	}

	var errs []error
	for _, child := range pa.Children {
		errs = append(errs, child.Validate())
	}

	if err := errors.Join(errs...); err != nil {
		return fmt.Errorf("parallel action %s has invalid children: %v", pa.Name, err)
	}

	return nil
}
