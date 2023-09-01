package common

import "regexp"

func SanitizeString(urlId string) string {
	urlIdRegex := regexp.MustCompile(`[^a-zA-Z0-9]`)
	return string(urlIdRegex.ReplaceAll([]byte(urlId), []byte("")))
}

func SanitizeUrl(url string) string {
	urlRegex := regexp.MustCompile(`[^a-zA-Z0-9-_/:.,?&@=#%]`)
	return string(urlRegex.ReplaceAll([]byte(url), []byte("")))
}
