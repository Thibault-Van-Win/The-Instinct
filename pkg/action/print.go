package action

import "fmt"

// Used to cast the output function to the Action interface
func PrintAction() Action {
	return DoFunc(output)
}

func output() error {
	fmt.Println("A very important action was carried out!")
	return nil
}
