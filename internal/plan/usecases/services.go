package usecases

import (
	"encoding/json"
	"fmt"
	"strings"
)

// StepStatus represents the allowed state for each plan step.
type StepStatus string

const (
	StepStatusPending    StepStatus = "pending"
	StepStatusInProgress StepStatus = "in_progress"
	StepStatusCompleted  StepStatus = "completed"
)

var allowedStepStatuses = map[StepStatus]struct{}{
	StepStatusPending:    {},
	StepStatusInProgress: {},
	StepStatusCompleted:  {},
}

// PlanItem describes a single item inside a plan update request.
type PlanItem struct {
	Step   string     `json:"step"`
	Status StepStatus `json:"status"`
}

// UpdatePlanRequest is the payload accepted by the PlanService.
type UpdatePlanRequest struct {
	Explanation string     `json:"explanation,omitempty"`
	Plan        []PlanItem `json:"plan"`
}

// PlanItemJSONSchema returns a JSON Schema map describing PlanItem based on its struct definition.
func PlanItemJSONSchema() map[string]any {
	statusEnum := []string{
		string(StepStatusPending),
		string(StepStatusInProgress),
		string(StepStatusCompleted),
	}

	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"step": map[string]any{
				"type":        "string",
				"description": "Human-readable action description",
			},
			"status": map[string]any{
				"type":        "string",
				"description": "Progress state for the step",
				"enum":        statusEnum,
			},
		},
		"required":             []string{"step", "status"},
		"additionalProperties": false,
	}
}

func decodePlanArguments(args map[string]any) (UpdatePlanRequest, error) {
	var payload UpdatePlanRequest
	raw, err := json.Marshal(args)
	if err != nil {
		return payload, fmt.Errorf("引数のシリアライズに失敗しました: %w", err)
	}
	if len(raw) == 0 || string(raw) == "null" {
		return payload, fmt.Errorf("plan引数が指定されていません")
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return payload, fmt.Errorf("引数のデシリアライズに失敗しました: %w", err)
	}
	return payload, nil
}

// PlanService validates and renders incoming plan updates.
type PlanService struct{}

// NewPlanService creates a new PlanService instance ready for use.
func NewPlanService() *PlanService {
	return &PlanService{}
}

func (s *PlanService) HandleUpdatePlan(args map[string]any) (string, error) {
	payload, err := decodePlanArguments(args)
	if err != nil {
		return "", fmt.Errorf("引数の解析に失敗しました: %w", err)
	}
	summary, err := s.UpdatePlan(payload)
	if err != nil {
		return "", fmt.Errorf("プランの作成に失敗しました: %w", err)
	}

	return summary, nil
}

// UpdatePlan validates the request and returns a formatted summary string.
func (s *PlanService) UpdatePlan(request UpdatePlanRequest) (string, error) {
	if err := validatePlanRequest(request); err != nil {
		return "", err
	}
	return formatPlanSummary(request), nil
}

func validatePlanRequest(request UpdatePlanRequest) error {
	if request.Plan == nil {
		return fmt.Errorf("plan field is required")
	}
	var inProgressCount int
	for idx, item := range request.Plan {
		trimmed := strings.TrimSpace(item.Step)
		if trimmed == "" {
			return fmt.Errorf("plan[%d].step must not be empty", idx)
		}
		if _, ok := allowedStepStatuses[item.Status]; !ok {
			return fmt.Errorf("plan[%d].status must be one of pending, in_progress, completed", idx)
		}
		if item.Status == StepStatusInProgress {
			inProgressCount++
			if inProgressCount > 1 {
				return fmt.Errorf("plan must not contain more than one in_progress step")
			}
		}
	}
	return nil
}

func formatPlanSummary(request UpdatePlanRequest) string {
	var builder strings.Builder
	explanation := strings.TrimSpace(request.Explanation)
	if explanation != "" {
		builder.WriteString("Explanation:\n")
		builder.WriteString(explanation)
		builder.WriteString("\n\n")
	}

	if len(request.Plan) == 0 {
		builder.WriteString("(No plan steps were provided.)")
		return builder.String()
	}

	builder.WriteString("Plan Steps:\n")
	for i, item := range request.Plan {
		builder.WriteString(fmt.Sprintf("%d. [%s] %s\n", i+1, strings.ToUpper(string(item.Status)), item.Step))
	}
	return strings.TrimRight(builder.String(), "\n")
}
