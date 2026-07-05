package domain

type Status string

const (
	StatusHealthy  Status = "healthy"
	StatusWarning  Status = "warning"
	StatusCritical Status = "critical"
	StatusUnknown  Status = "unknown"
)

type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

type SmartAttribute struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Value      int    `json:"value"`
	Worst      int    `json:"worst"`
	Threshold  int    `json:"threshold"`
	Type       string `json:"type"`
	Updated    string `json:"updated"`
	WhenFailed string `json:"when_failed"`
	RawValue   int64  `json:"raw_value"`
	RawText    string `json:"raw_text"`
}

type SmartReport struct {
	Model         string           `json:"model,omitempty"`
	SerialNumber  string           `json:"serial_number,omitempty"`
	OverallHealth string           `json:"overall_health,omitempty"`
	Attributes    []SmartAttribute `json:"attributes,omitempty"`
}

type Finding struct {
	Severity      Severity `json:"severity"`
	AttributeID   int      `json:"attribute_id,omitempty"`
	AttributeName string   `json:"attribute_name,omitempty"`
	RawValue      int64    `json:"raw_value,omitempty"`
	Message       string   `json:"message"`
}

type Assessment struct {
	Status        Status           `json:"status"`
	Score         int              `json:"score"`
	Summary       string           `json:"summary"`
	OverallHealth string           `json:"overall_health,omitempty"`
	Model         string           `json:"model,omitempty"`
	SerialNumber  string           `json:"serial_number,omitempty"`
	Findings      []Finding        `json:"findings,omitempty"`
	Attributes    []SmartAttribute `json:"attributes,omitempty"`
}
