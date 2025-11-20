package timezone

import (
	"testing"
	"time"
)

func TestTimezoneAdditional(t *testing.T) {
	if _, err := Load(""); err == nil {
		t.Fatalf("expected empty timezone error")
	}

	if got := NowIn(nil); got.Location().String() != "UTC" {
		t.Fatalf("NowIn nil should return UTC time")
	}
	if got := Convert(time.Now(), nil); got.IsZero() {
		t.Fatalf("convert nil should return original time")
	}

	if !EqualLocation(nil, nil) {
		t.Fatalf("nil locations should be equal")
	}
	if EqualLocation(time.UTC, nil) {
		t.Fatalf("one nil location should not be equal")
	}
	if !EqualLocation(time.UTC, time.FixedZone("UTC", 0)) {
		t.Fatalf("same location names should be equal")
	}

	_ = IsDST(time.Now())
}
