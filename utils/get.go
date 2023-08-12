package utils

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	filter := bson.D{{Key: "url", Value: path}}

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

	filter := bson.D{{Key: "redirectUrl", Value: redirectUrl}}

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
