package utils

// SeverityToString converts numeric severity to text
func SeverityToString(severity int) string {
	switch severity {
	case 1:
		return "Low"
	case 2:
		return "Medium"
	case 3:
		return "High"
	case 4:
		return "Critical"
	default:
		return "Unknown"
	}
}

func Prepend[T any](slice []T, elems ...T) []T {
	return append(elems, slice...)
}
