package arraysum

import "testing"

func TestService_Execute_Normal(t *testing.T) {
	service := NewService()

	result, err := service.Execute("sum", []float64{1, 2, 3})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if result != "sum([1 2 3]) = 6.00\n" {
		t.Fatalf("unexpected result: %q", result)
	}
}

func TestService_Execute_UnsupportedOperation(t *testing.T) {
	service := NewService()

	_, err := service.Execute("unknown", []float64{1, 2, 3})
	if err == nil {
		t.Fatal("expected error")
	}
}
