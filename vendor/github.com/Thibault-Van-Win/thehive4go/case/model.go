package cases

import (
	"time"
)

// Case structure (simplified for TheHive)
type Case struct {
	ID          string    `json:"_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Summary     string    `json:"summary"`
	Severity    int       `json:"severity"`
	StartDate   int64     `json:"startDate"`
	Status      string    `json:"status"`
	Tags        []string  `json:"tags"`
	TLP         int       `json:"tlp"`
	CreatedAt   time.Time `json:"createdAt"`
}
