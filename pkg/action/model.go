package action

type Action interface {
	Do() error
}

// Function as interface pattern
type DoFunc func() error

func (df DoFunc) Do() error {
	return df()
}
