package server

import (
	"os"
	"testing"
)

func TestRender(t *testing.T) {
	logs, err := os.ReadFile("testdata/logs.txt")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	err = Render(os.Stdout, true, string(logs), nil)
	if err != nil {
		t.Fatalf("failed to render template: %v", err)
	}
}
