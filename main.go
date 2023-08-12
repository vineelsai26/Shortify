package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Html struct {
	URL         string
	RedirectURL string
	Error       string
}

func getRedirectToURL(path string) string {
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
			return ""
		}
		panic(err)
	}

	return result["protocol"].(string) + "://" + result["redirectUrl"].(string)
}

func getRedirectFromURL(redirectUrl string) string {
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

func formattedURL(url string) (string, string, error) {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return strings.Split(url, "://")[1], strings.Split(url, "://")[0], nil
	} else if !strings.Contains(url, "://") {
		return url, "https", nil
	} else {
		return "", "", fmt.Errorf("protocol is not supported")
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateURLID(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	if getRedirectToURL(string(b)) != "" {
		return generateURLID(n)
	}
	return string(b)
}

func generateURL(id, url, protocol, createdAt string) error {
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

	item := bson.D{{Key: "url", Value: id}}

	_, err = coll.InsertOne(context.TODO(), item)
	if err != nil {
		return err
	}

	return nil
}

func getProtocol(r *http.Request) string {
	if r.TLS == nil {
		return "http"
	} else {
		return "https"
	}
}

func redirect(path string, w http.ResponseWriter, r *http.Request, t *template.Template) {
	if r.Method == "POST" {
		url := r.PostFormValue("url")

		redirectToURL, protocol, err := formattedURL(url)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		getRedirectFromURL := getRedirectFromURL(redirectToURL)
		isURLExists := getRedirectFromURL != ""
		if !isURLExists {
			id := generateURLID(6)
			createdAt := time.Now().String()
			err = generateURL(id, redirectToURL, protocol, createdAt)
			if err != nil {
				t.Execute(w, Html{
					URL:   url,
					Error: err.Error(),
				})
				return
			}
			t.Execute(w, Html{
				URL:         url,
				RedirectURL: getProtocol(r) + "://" + r.Host + "/" + id,
			})
		} else {
			t.Execute(w, Html{
				URL:         url,
				RedirectURL: getProtocol(r) + "://" + r.Host + "/" + getRedirectFromURL,
			})
		}
	} else if r.Method == "GET" {
		redirectToURL := getRedirectToURL(path)
		http.Redirect(w, r, redirectToURL, http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		t, err := template.ParseFiles("static/index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if path == "/" && r.Method == "GET" {
			t.Execute(w, nil)
		} else if path == "/style.css" && r.Method == "GET" {
			http.ServeFile(w, r, "static/style.css")
		} else if len(strings.Split(path, "/")) == 2 {
			redirect(strings.Split(path, "/")[1], w, r, t)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	http.ListenAndServe(":3000", nil)
}
