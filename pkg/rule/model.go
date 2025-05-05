package rule

import "github.com/Thibault-Van-Win/The-Instinct/pkg/action"

type Rule interface {
	Match(ctx *action.SecurityContext) (bool, error)
}
