package trigonometry

import "testing"

func TestService_Execute_Normal(t *testing.T) {
	service := NewService()

	result, err := service.Execute("sin", 90, "degrees")
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if result != "sin(90.00 degrees) = 1.000000\n" {
		t.Fatalf("unexpected result: %q", result)
	}
}

func TestService_Execute_Error(t *testing.T) {
	service := NewService()

	_, err := service.Execute("bad", 90, "degrees")
	if err == nil {
		t.Fatal("expected error")
	}
}
