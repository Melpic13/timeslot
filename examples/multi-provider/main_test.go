package main

import (
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	main()
}

func TestCreateProvider(t *testing.T) {
	if p := createProvider("x", time.Monday); p == nil {
		t.Fatalf("expected provider")
	}
}
