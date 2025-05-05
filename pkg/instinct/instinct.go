package instinct

import (
	"context"
	"fmt"
	"sync"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/reflex"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/rule"
)

type Instinct struct {
	Reflexes []reflex.Reflex
	mu       sync.RWMutex
}

// This is implemented by both the loader and repository interface
type ReflexLister interface {
	ListReflexes(ctx context.Context) ([]*reflex.Reflex, error)
}

// New creates and returns a new Instinct instance
func New(ruleRegistry *rule.RuleRegistry, actionRegistry *action.ActionRegistry) *Instinct {
	return &Instinct{
		Reflexes: []reflex.Reflex{},
	}
}

// AddReflex adds a new reflex to the instinct
func (i *Instinct) AddReflex(r reflex.Reflex) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.Reflexes = append(i.Reflexes, r)
}

// LoadReflexes loads reflexes using the specified loader
func (i *Instinct) LoadReflexes(reflexLister ReflexLister) error {

	reflexes, err := reflexLister.ListReflexes(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list reflexes: %w", err)
	}

	// Dereference pointers
	reflexVals := []reflex.Reflex{}
	for _, reflexPtr := range reflexes {
		reflexVals = append(reflexVals, *reflexPtr)
	}

	// Add the reflexes
	i.mu.Lock()
	defer i.mu.Unlock()
	i.Reflexes = append(i.Reflexes, reflexVals...)

	return nil
}

// ProcessEvent processes an event through all reflexes
func (i *Instinct) ProcessEvent(data map[string]any) error {
	i.mu.RLock()
	defer i.mu.RUnlock()

	// Create a context for the event
	ctx := &action.SecurityContext{
		Event:           data,
		Variables:       make(map[string]any),
		ExecutionStatus: make(map[string]action.Status),
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(i.Reflexes))

	for _, r := range i.Reflexes {
		wg.Add(1)
		// Process each reflex in its own goroutine
		go func(reflex reflex.Reflex) {
			defer wg.Done()
			match, err := reflex.Match(data)
			if err != nil {
				errChan <- fmt.Errorf("error matching reflex %s: %w", reflex.Name, err)
				return
			}
			if match {
				if err := reflex.Execute(ctx); err != nil {
					errChan <- fmt.Errorf("error executing reflex %s: %w", reflex.Name, err)
				}
			}
		}(r)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Collect errors
	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("encountered %d errors while processing event", len(errs))
	}

	return nil
}
