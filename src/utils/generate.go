package utils

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"time"

	_ "github.com/tursodatabase/go-libsql"
	"vineelsai.com/shortify/src/common"
)

const letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Generates a random string of length n
func GenerateURLID(n int) string {
	b := make([]byte, n)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	fmt.Printf("Generated URL ID: %s\n", string(b))
	if _, err := GetRedirectToURL(string(b)); err == nil {
		return GenerateURLID(n)
	}
	return string(b)
}

// Create a new URL in the database
func GenerateURL(id, url, protocol, createdAt string) error {
	database_url := os.Getenv("TURSO_DATABASE")
	auth_token := os.Getenv("TURSO_AUTH_TOKEN")

	turso_url := "libsql://" + database_url + ".turso.io?authToken=" + auth_token

	db, err := sql.Open("libsql", turso_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", turso_url, err)
		os.Exit(1)
	}
	defer db.Close()

	insert_query := fmt.Sprintf("INSERT INTO urls (url, redirectUrl, protocol, createdAt) VALUES ('%s', '%s', '%s', '%s')", id, common.SanitizeString(url), common.SanitizeString(protocol), "-1")

	res, err := db.Exec(insert_query)
	if err != nil {
		return err
	}

	fmt.Println(res.LastInsertId())

	return nil
}
