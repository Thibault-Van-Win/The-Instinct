package rule

type CelRule struct {
	Expression string
}

func (cr *CelRule) Match(data map[string]any) (bool, error) {

	return true, nil
}