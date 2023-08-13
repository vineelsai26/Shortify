package utils

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		panic("MONGODB_URI is not set")
	}

	opts := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("URLS").Collection("urls")

	filter := bson.D{{Key: "url", Value: sanitizeString(path)}}

	var result bson.M
	doc := coll.FindOne(context.TODO(), filter)

	err = doc.Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", fmt.Errorf("URL not found")
		}
		return "", err
	}

	return result["protocol"].(string) + "://" + result["redirectUrl"].(string), nil
}

// Fetches Short URL Path Name from the database
func GetRedirectFromURL(redirectUrl string) string {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		panic("MONGODB_URI is not set")
	}

	opts := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("URLS").Collection("urls")

	filter := bson.D{{Key: "redirectUrl", Value: sanitizeUrl(redirectUrl)}}

	var result bson.M
	doc := coll.FindOne(context.TODO(), filter)

	err = doc.Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ""
		}
		panic(err)
	}

	return result["url"].(string)
}
