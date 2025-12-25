package usecases

import "testing"

func TestPlanService_HandleUpdatePlan_Summary(t *testing.T) {
	service := NewPlanService()
	req := UpdatePlanRequest{
		Explanation: "MVPに集中する",
		Plan: []PlanItem{
			{Step: "仕様調査", Status: StepStatusCompleted},
			{Step: "PoCの作成", Status: StepStatusInProgress},
			{Step: "最終レビュー", Status: StepStatusPending},
		},
	}

	got, err := service.HandleUpdatePlan(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "Explanation:\nMVPに集中する\n\nPlan Steps:\n1. [COMPLETED] 仕様調査\n2. [IN_PROGRESS] PoCの作成\n3. [PENDING] 最終レビュー"
	if got != want {
		t.Fatalf("summary mismatch\nwant: %q\ngot:  %q", want, got)
	}
}

func TestPlanService_HandleUpdatePlan_AllowsEmptyPlan(t *testing.T) {
	service := NewPlanService()
	req := UpdatePlanRequest{Explanation: "やることは後で追加", Plan: []PlanItem{}}

	got, err := service.HandleUpdatePlan(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "Explanation:\nやることは後で追加\n\n(No plan steps were provided.)"
	if got != want {
		t.Fatalf("summary mismatch\nwant: %q\ngot:  %q", want, got)
	}
}

func TestPlanService_HandleUpdatePlan_Errors(t *testing.T) {
	t.Run("missing plan", func(t *testing.T) {
		service := NewPlanService()
		_, err := service.HandleUpdatePlan(UpdatePlanRequest{})
		if err == nil {
			t.Fatalf("expected error for missing plan")
		}
	})

	t.Run("empty step", func(t *testing.T) {
		service := NewPlanService()
		_, err := service.HandleUpdatePlan(UpdatePlanRequest{
			Plan: []PlanItem{{Step: "   ", Status: StepStatusPending}},
		})
		if err == nil {
			t.Fatalf("expected error for empty step")
		}
	})

	t.Run("invalid status", func(t *testing.T) {
		service := NewPlanService()
		_, err := service.HandleUpdatePlan(UpdatePlanRequest{
			Plan: []PlanItem{{Step: "setup", Status: StepStatus("unknown")}},
		})
		if err == nil {
			t.Fatalf("expected error for invalid status")
		}
	})

	t.Run("multiple in_progress", func(t *testing.T) {
		service := NewPlanService()
		_, err := service.HandleUpdatePlan(UpdatePlanRequest{
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
