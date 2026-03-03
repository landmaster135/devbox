package getconstants

import "testing"

func TestService_Execute_Normal(t *testing.T) {
	service := NewService()

	result, err := service.Execute()
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if result == "" {
		t.Fatal("expected non-empty result")
	}
}
