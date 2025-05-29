package procedure

type ProcedureRequest struct {
	PatternID   string `json:"patternId"`
	OccurDate   int64  `json:"occurDate"`
	Tactic      string `json:"tactic"`
	Description string `json:"description"`
}
