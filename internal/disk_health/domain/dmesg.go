package domain

type DmesgEvent struct {
	Severity Severity `json:"severity"`
	Device   string   `json:"device,omitempty"`
	Message  string   `json:"message"`
	Line     string   `json:"line,omitempty"`
}

type DmesgAssessment struct {
	Status   Status       `json:"status"`
	Score    int          `json:"score"`
	Summary  string       `json:"summary"`
	Findings []DmesgEvent `json:"findings,omitempty"`
}
