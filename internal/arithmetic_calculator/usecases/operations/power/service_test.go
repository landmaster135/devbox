package power

import "testing"

func TestService_Execute_Normal(t *testing.T) {
	service := NewService()

	result, err := service.Execute(2, 8)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if result != "2.00^8.00 = 256.00\n" {
		t.Fatalf("unexpected result: %q", result)
	}
}
