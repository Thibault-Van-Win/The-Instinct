package rule

type Rule interface {
	Match(data map[string]any) (bool, error)
}
