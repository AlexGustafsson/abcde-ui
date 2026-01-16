package abcde

import (
	"fmt"
	"os"
	"testing"
)

func TestParseLogInfo(t *testing.T) {
	file, err := os.Open("testdata/logs.txt")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	info := ParseLogInfo(file)
	fmt.Printf("%+v\n", info)
}
