package tlp

import "fmt"

type TLP int

const (
	WHITE TLP = iota
	GREEN
	AMBER
	AMBER_STRICT
	RED
)

func (t TLP) String() string {
	level := [...]string{"WHITE", "GREEN", "AMBER", "AMBER+STRICT", "RED"}[t]
	return fmt.Sprintf("TLP:%s", level)
}

func Validate(level TLP) error {
	if level < WHITE || level > RED {
		return fmt.Errorf("invalid TLP value, must be between 0 (WHITE) and 4 (RED): %d", level)
	}

	return nil
}
