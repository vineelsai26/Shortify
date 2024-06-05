package utils

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"vineelsai.com/shortify/src/common"
)

type Redirect struct {
	url         string
	redirectUrl string
	protocol    string
}

// Removes the protocol from the URL
func GetFormattedURL(url string) (string, string, error) {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return strings.Split(url, "://")[1], strings.Split(url, "://")[0], nil
	} else if !strings.Contains(url, "://") {
		return url, "https", nil
	} else {
		return "", "", fmt.Errorf("protocol is not supported")
	}
}

// Fetches the URL to redirect to from the database
func GetRedirectToURL(path string) (string, error) {
	database_url := os.Getenv("TURSO_DATABASE")
	auth_token := os.Getenv("TURSO_AUTH_TOKEN")

	turso_url := "libsql://" + database_url + ".turso.io?authToken=" + auth_token
	db, err := sql.Open("libsql", turso_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", turso_url, err)
		os.Exit(1)
	}
	defer db.Close()

	select_query := fmt.Sprintf("SELECT protocol, redirectUrl FROM urls WHERE url='%s'", common.SanitizeString(path))
	res, err := db.Query(select_query)
	if err != nil {
		return "", err
	}
	defer res.Close()

	fmt.Println(res.Next())

	var redirect Redirect
	res.Scan(&redirect.protocol, &redirect.redirectUrl)
	if redirect.redirectUrl == "" {
		return "", fmt.Errorf("URL Not found")
	}

	return redirect.protocol + "://" + redirect.redirectUrl, nil
}

// Fetches Short URL Path Name from the database
func GetRedirectFromURL(redirectUrl string) string {
	database_url := os.Getenv("TURSO_DATABASE")
	auth_token := os.Getenv("TURSO_AUTH_TOKEN")

	turso_url := "libsql://" + database_url + ".turso.io?authToken=" + auth_token
	db, err := sql.Open("libsql", turso_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", turso_url, err)
		os.Exit(1)
	}
	defer db.Close()

	select_query := fmt.Sprintf("SELECT url FROM urls WHERE redirectUrl='%s'", redirectUrl)
	res, err := db.Query(select_query)
	if err != nil {
		return ""
	}
	defer res.Close()

	fmt.Println(res.Next())

	var redirect Redirect
	res.Scan(&redirect.url)

	return redirect.url
}
