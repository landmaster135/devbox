package factorial

import "testing"

func TestService_Execute_Normal(t *testing.T) {
	service := NewService()

	result, err := service.Execute(5)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if result != "5! = 120\n" {
		t.Fatalf("unexpected result: %q", result)
	}
}

func TestService_Execute_Error(t *testing.T) {
	service := NewService()

	_, err := service.Execute(-1)
	if err == nil {
		t.Fatal("expected error")
	}
}
