package action

import "fmt"

// PrintAction creates a simple action that prints a message
func PrintAction(message string) Action {
	return DoFunc(func() error {
		fmt.Println(message)
		return nil
	})
}
