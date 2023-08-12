package utils

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Removes the protocol from the URL
func GenerateFormattedURL(url string) (string, string, error) {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return strings.Split(url, "://")[1], strings.Split(url, "://")[0], nil
	} else if !strings.Contains(url, "://") {
		return url, "https", nil
	} else {
		return "", "", fmt.Errorf("protocol is not supported")
	}
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

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
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		panic("MONGODB_URI is not set")
	}

	opts := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}

	defer func() {
		client.Disconnect(context.TODO())
	}()

	coll := client.Database("URLS").Collection("urls")

	item := bson.D{{Key: "url", Value: id}, {Key: "redirectUrl", Value: url}, {Key: "protocol", Value: protocol}, {Key: "createdAt", Value: createdAt}}

	_, err = coll.InsertOne(context.TODO(), item)
	if err != nil {
		return err
	}

	return nil
}
