package utils

import (
	"strings"
	"testing"
)

func TestGenerateURLID(t *testing.T) {
	n := 6
	id := GenerateURLID(n)
	if id == "" {
		t.Errorf("Error while generating URL ID: %v", id)
	}

	if len(strings.TrimSpace(id)) != n {
		t.Errorf("Error while generating URL ID: %v", id)
	}
}
