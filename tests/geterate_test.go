package tests

import (
	"strings"
	"testing"

	"vineelsai.com/shortify/src/common"
	"vineelsai.com/shortify/src/utils"
)

func TestGenerateURLID(t *testing.T) {
	n := 6
	id := common.SanitizeString(utils.GenerateURLID(n))
	if id == "" {
		t.Errorf("Error while generating URL ID: %v", id)
	}

	if len(strings.TrimSpace(id)) != n {
		t.Errorf("Error while generating URL ID: %v", id)
	}
}
