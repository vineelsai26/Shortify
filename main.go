package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"vineelsai.com/shortify/utils"
)

type Html struct {
	URL         string
	RedirectURL string
	Error       string
}

func getProtocol(r *http.Request) string {
	if r.TLS == nil {
		return "http"
	} else {
		return "https"
	}
}

func render(path string, w http.ResponseWriter, r *http.Request, t *template.Template) {
	url := r.PostFormValue("url")

	redirectToURL, protocol, err := utils.GetFormattedURL(url)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	getRedirectFromURL := utils.GetRedirectFromURL(redirectToURL)
	isURLExists := getRedirectFromURL != ""
	if !isURLExists {
		id := utils.GenerateURLID(6)
		createdAt := time.Now().String()
		err = utils.GenerateURL(id, redirectToURL, protocol, createdAt)
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
}

// redirects to the URL if it exists in the database
func redirect(path string, w http.ResponseWriter, r *http.Request) {
	redirectToURL, err := utils.GetRedirectToURL(path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	http.Redirect(w, r, redirectToURL, http.StatusTemporaryRedirect)
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
		} else if len(strings.Split(path, "/")) == 2 && r.Method == "GET" {
			redirect(strings.Split(path, "/")[1], w, r)
		} else if len(strings.Split(path, "/")) == 2 && r.Method == "POST" {
			render(strings.Split(path, "/")[1], w, r, t)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	http.ListenAndServe(":3000", nil)
}
