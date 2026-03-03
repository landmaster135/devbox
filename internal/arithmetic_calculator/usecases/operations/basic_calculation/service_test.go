package basiccalculation

import "testing"

func TestService_Execute_Normal(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		x         float64
		y         float64
		expected  string
		err       bool
	}{
		{name: "add", operation: "add", x: 10, y: 5, expected: "10.00 + 5.00 = 15.00\n"},
		{name: "subtract", operation: "subtract", x: 10, y: 5, expected: "10.00 - 5.00 = 5.00\n"},
		{name: "multiply", operation: "multiply", x: 10, y: 5, expected: "10.00 * 5.00 = 50.00\n"},
		{name: "divide", operation: "divide", x: 10, y: 5, expected: "10.00 / 5.00 = 2.00\n"},
		{name: "divide by zero", operation: "divide", x: 10, y: 0, err: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService()
			result, err := service.Execute(tt.operation, tt.x, tt.y)
			if tt.err {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("Execute returned error: %v", err)
			}
			if result != tt.expected {
				t.Fatalf("unexpected result: %q", result)
			}
		})
	}
}
