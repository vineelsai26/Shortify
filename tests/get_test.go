package tests

import (
	"testing"

	"vineelsai.com/shortify/src/utils"
)

func TestGetFormattedURL(t *testing.T) {
	valid_urls := []string{
		"https://vineelsai.com",
		"http://vineelsai.com",
		"vineelsai.com",
	}

	for _, url := range valid_urls {
		_, _, err := utils.GetFormattedURL(url)
		if err != nil {
			t.Errorf("Error while generating formatted URL: %v", err)
		}
	}

	invalid_urls := []string{
		"ftp://vineelsai.com",
		"ssh://vineelsai.com",
	}

	for _, url := range invalid_urls {
		_, _, err := utils.GetFormattedURL(url)
		if err == nil {
			t.Errorf("Error while generating formatted URL: %v", err)
		}
	}
}

func TestGetRedirectToURL(t *testing.T) {
	url, err := utils.GetRedirectToURL("LkO8rE1")
	if err != nil {
		t.Errorf("Error while fetching URL: %v", err)
	}

	if url == "" {
		t.Errorf("Error while fetching URL: %v", err)
	}

	if url != "https://vineelsai.com" {
		t.Errorf("Error while fetching URL: %v", err)
	}
}

func TestGetRedirectFromURL(t *testing.T) {
	url := utils.GetRedirectFromURL("vineelsai.com")
	if url == "" {
		t.Errorf("Error while fetching URL: %v", url)
	}

	if url != "LkO8rE1" {
		t.Errorf("Error while fetching URL: %v", url)
	}
}
