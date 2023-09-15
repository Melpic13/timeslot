package timezone

import (
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	loc, err := Load("UTC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.String() != "UTC" {
		t.Fatalf("got %s", loc.String())
	}
}

func TestConvert(t *testing.T) {
	loc := MustLoad("America/New_York")
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	converted := Convert(base, loc)
	if converted.Location().String() != "America/New_York" {
		t.Fatalf("unexpected location %s", converted.Location())
	}
}
