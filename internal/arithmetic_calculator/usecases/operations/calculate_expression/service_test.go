package calculateexpression

import (
	"strings"
	"testing"
)

func TestService_Execute_Normal(t *testing.T) {
	service := NewService()

	result, err := service.Execute("2+3*4")
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if result != "2+3*4 = 14.00\n" {
		t.Fatalf("unexpected result: %q", result)
	}
}

func TestService_Execute_Error(t *testing.T) {
	service := NewService()

	_, err := service.Execute("import os")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestService_CheckOSPattern(t *testing.T) {
	service := NewService()

	if err := service.CheckOSPattern("cos(0)"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := service.CheckOSPattern("os.system('ls')"); err == nil {
		t.Fatal("expected error")
	}

	indices := service.GetAllIndices(strings.ToLower("cos(cos(0))"), "os")
	if len(indices) == 0 {
		t.Fatal("expected indices")
	}
}
