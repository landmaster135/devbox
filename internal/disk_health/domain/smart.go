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

type DiskInfo struct {
	RotationRateRPM    *int64 `json:"rotation_rate_rpm,omitempty"`
	PowerOnHours       *int64 `json:"power_on_hours,omitempty"`
	PowerCycleCount    *int64 `json:"power_cycle_count,omitempty"`
	TemperatureCelsius *int64 `json:"temperature_celsius,omitempty"`
	TotalLBAsWritten   *int64 `json:"total_lbas_written,omitempty"`
	TotalBytesWritten  *int64 `json:"total_bytes_written,omitempty"`
	TotalLBAsRead      *int64 `json:"total_lbas_read,omitempty"`
	TotalBytesRead     *int64 `json:"total_bytes_read,omitempty"`
}

type SmartReport struct {
	Model                  string           `json:"model,omitempty"`
	SerialNumber           string           `json:"serial_number,omitempty"`
	OverallHealth          string           `json:"overall_health,omitempty"`
	LogicalSectorSizeBytes *int64           `json:"-"`
	DiskInfo               *DiskInfo        `json:"disk_info,omitempty"`
	Attributes             []SmartAttribute `json:"attributes,omitempty"`
}

type Finding struct {
	Severity      Severity `json:"severity"`
	AttributeID   int      `json:"attribute_id,omitempty"`
	AttributeName string   `json:"attribute_name,omitempty"`
	RawValue      int64    `json:"raw_value"`
	Message       string   `json:"message"`
}

type Assessment struct {
	Status        Status           `json:"status"`
	Score         int              `json:"score"`
	Summary       string           `json:"summary"`
	OverallHealth string           `json:"overall_health,omitempty"`
	Model         string           `json:"model,omitempty"`
	SerialNumber  string           `json:"serial_number,omitempty"`
	DiskInfo      *DiskInfo        `json:"disk_info,omitempty"`
	Findings      []Finding        `json:"findings,omitempty"`
	Attributes    []SmartAttribute `json:"attributes,omitempty"`
}
