package action

import "fmt"

type PrintAction struct {
	Type    string
	Message string
}

// PrintAction creates a simple action that prints a message
func (pa *PrintAction) Do() error {
	fmt.Println(pa.Message)
	return nil
}

func NewPrintAction(message string) *PrintAction {
	return &PrintAction{
		Type:    "print",
		Message: message,
	}
}
