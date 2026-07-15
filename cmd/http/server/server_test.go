package server

import "testing"

func TestGetPortUsesCronURLPort(t *testing.T) {
	t.Setenv(PortEnvKey, "30000")

	got := getPort(nil)

	if got != 30000 {
		t.Fatalf("getPort() = %d, want 30000", got)
	}
}

func TestGetPortDefaultsWhenCronURLPortIsMissing(t *testing.T) {
	t.Setenv(PortEnvKey, "")

	got := getPort(nil)

	if got != DefaultPort {
		t.Fatalf("getPort() = %d, want %d", got, DefaultPort)
	}
}

func TestGetPortDefaultsWhenCronURLPortIsInvalid(t *testing.T) {
	t.Setenv(PortEnvKey, "invalid")

	got := getPort(nil)

	if got != DefaultPort {
		t.Fatalf("getPort() = %d, want %d", got, DefaultPort)
	}
}
