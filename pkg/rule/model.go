package rule

import "github.com/Thibault-Van-Win/The-Instinct/pkg/security_context"


type Rule interface {
	Match(ctx *security_context.SecurityContext) (bool, error)
}
