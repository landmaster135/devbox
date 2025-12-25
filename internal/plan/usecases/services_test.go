package usecases

import (
	"strings"
	"testing"
)

func TestPlanService_UpdatePlan_Summary(t *testing.T) {
	service := NewPlanService()
	req := UpdatePlanRequest{
		Explanation: "MVPに集中する",
		Plan: []PlanItem{
			{Step: "仕様調査", Status: StepStatusCompleted},
			{Step: "PoCの作成", Status: StepStatusInProgress},
			{Step: "最終レビュー", Status: StepStatusPending},
		},
	}

	got, err := service.UpdatePlan(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "Explanation:\nMVPに集中する\n\nPlan Steps:\n1. [COMPLETED] 仕様調査\n2. [IN_PROGRESS] PoCの作成\n3. [PENDING] 最終レビュー"
	if got != want {
		t.Fatalf("summary mismatch\nwant: %q\ngot:  %q", want, got)
	}
}

func TestPlanService_UpdatePlan_AllowsEmptyPlan(t *testing.T) {
	service := NewPlanService()
	req := UpdatePlanRequest{Explanation: "やることは後で追加", Plan: []PlanItem{}}

	got, err := service.UpdatePlan(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "Explanation:\nやることは後で追加\n\n(No plan steps were provided.)"
	if got != want {
		t.Fatalf("summary mismatch\nwant: %q\ngot:  %q", want, got)
	}
}

func TestPlanService_UpdatePlan_Errors(t *testing.T) {
	t.Run("missing plan", func(t *testing.T) {
		service := NewPlanService()
		_, err := service.UpdatePlan(UpdatePlanRequest{})
		if err == nil {
			t.Fatalf("expected error for missing plan")
		}
	})

	t.Run("empty step", func(t *testing.T) {
		service := NewPlanService()
		_, err := service.UpdatePlan(UpdatePlanRequest{
			Plan: []PlanItem{{Step: "   ", Status: StepStatusPending}},
		})
		if err == nil {
			t.Fatalf("expected error for empty step")
		}
	})

	t.Run("invalid status", func(t *testing.T) {
		service := NewPlanService()
		_, err := service.UpdatePlan(UpdatePlanRequest{
			Plan: []PlanItem{{Step: "setup", Status: StepStatus("unknown")}},
		})
		if err == nil {
			t.Fatalf("expected error for invalid status")
		}
	})

	t.Run("multiple in_progress", func(t *testing.T) {
		service := NewPlanService()
		_, err := service.UpdatePlan(UpdatePlanRequest{
			Plan: []PlanItem{
				{Step: "A", Status: StepStatusInProgress},
				{Step: "B", Status: StepStatusInProgress},
			},
		})
		if err == nil {
			t.Fatalf("expected error for multiple in_progress")
		}
	})
}

func TestPlanService_HandleUpdatePlan(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		args    map[string]any
		want    string
		wantErr string
	}{
		{
			name: "valid request",
			args: map[string]any{
				"explanation": "Focus on MVP",
				"plan": []map[string]any{
					{"step": "Setup env", "status": string(StepStatusCompleted)},
					{"step": "Implement feature", "status": string(StepStatusInProgress)},
				},
			},
			want: "Explanation:\nFocus on MVP\n\nPlan Steps:\n1. [COMPLETED] Setup env\n2. [IN_PROGRESS] Implement feature",
		},
		{
			name:    "missing arguments",
			args:    nil,
			wantErr: "plan引数が指定されていません",
		},
		{
			name: "validation error",
			args: map[string]any{
				"plan": []map[string]any{
					{"step": "A", "status": string(StepStatusInProgress)},
					{"step": "B", "status": string(StepStatusInProgress)},
				},
			},
			wantErr: "plan must not contain more than one in_progress step",
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			service := NewPlanService()
			got, err := service.HandleUpdatePlan(tt.args)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got != tt.want {
					t.Fatalf("summary mismatch\nwant: %q\ngot:  %q", tt.want, got)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error but got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error mismatch\nwant substring: %q\ngot: %v", tt.wantErr, err)
			}
		})
	}
}
