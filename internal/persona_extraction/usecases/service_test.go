package usecases

import (
	"strings"
	"testing"
)

func TestPersonaServiceHandleExtractionSuccess(t *testing.T) {
	t.Parallel()

	service := NewPersonaService()
	args := map[string]any{
		"context": "事件現場で出会った人物たち",
		"characters": []map[string]any{
			{
				"name":        "Mika",
				"age":         27,
				"gender":      "female",
				"spieces":     "human",
				"job":         "investigative journalist",
				"hobbies":     []string{"photography", "urban exploration"},
				"personality": "好奇心旺盛で慎重",
				"relationships": []map[string]any{
					{
						"with":  "Ken",
						"type":  "brother",
						"notes": "現場で情報共有",
					},
				},
				"evidence": []map[string]any{
					{
						"quote":  "真実を逃さない",
						"source": "chapter_02",
					},
				},
			},
		},
	}

	got, err := service.HandleExtraction(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(got, "Mika") || !strings.Contains(got, "human") || !strings.Contains(got, "photography") {
		t.Fatalf("formatted summary missing expected fields: %s", got)
	}
	if !strings.Contains(got, "Ken") || !strings.Contains(got, "chapter_02") {
		t.Fatalf("relationships/evidence not rendered: %s", got)
	}
}

func TestPersonaServiceHandleExtractionMissingArgs(t *testing.T) {
	service := NewPersonaService()
	if _, err := service.HandleExtraction(nil); err == nil || !strings.Contains(err.Error(), "characters引数") {
		t.Fatalf("expected error for missing args, got %v", err)
	}
}

func TestPersonaServiceHandleExtractionValidationError(t *testing.T) {
	service := NewPersonaService()
	args := map[string]any{
		"characters": []map[string]any{},
	}

	if _, err := service.HandleExtraction(args); err == nil || !strings.Contains(err.Error(), "charactersは最低1件") {
		t.Fatalf("expected validation error, got %v", err)
	}
}
