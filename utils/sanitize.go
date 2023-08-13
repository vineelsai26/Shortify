package utils

import "regexp"

func sanitizeString(urlId string) string {
	urlIdRegex := regexp.MustCompile(`[^a-zA-Z0-9]`)
	return string(urlIdRegex.ReplaceAll([]byte(urlId), []byte("")))
}

func sanitizeUrl(url string) string {
	urlRegex := regexp.MustCompile(`[^a-zA-Z0-9-_/:.,?&@=#%]`)
	return string(urlRegex.ReplaceAll([]byte(url), []byte("")))
}
