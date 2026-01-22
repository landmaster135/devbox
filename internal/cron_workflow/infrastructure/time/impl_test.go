package time

import "testing"

func TestStdRepositoryNowReturnsLocationTime(t *testing.T) {
	repo := NewRepository()
	tz := "Asia/Tokyo"

	tm, err := repo.Now(tz)
	if err != nil {
		t.Fatalf("Now() error = %v", err)
	}

	if loc := tm.Location().String(); loc != tz {
		t.Fatalf("Now() location = %s, want %s", loc, tz)
	}
}

func TestStdRepositoryNowRejectsEmptyTimezone(t *testing.T) {
	repo := NewRepository()
	if _, err := repo.Now(" "); err == nil {
		t.Fatalf("Now() error = nil, want error for empty timezone")
	}
}

func TestStdRepositoryNowRejectsInvalidTimezone(t *testing.T) {
	repo := NewRepository()
	if _, err := repo.Now("invalid/zone"); err == nil {
		t.Fatalf("Now() error = nil, want error for invalid timezone")
	}
}
