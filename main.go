package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"vineelsai.com/shortify/src/utils"
)

type Html struct {
	URL         string
	RedirectURL string
	Error       string
}

func render(res http.ResponseWriter, req *http.Request, template *template.Template) {
	url := req.PostFormValue("url")

	redirectToURL, protocol, err := utils.GetFormattedURL(url)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	getRedirectFromURL := utils.GetRedirectFromURL(redirectToURL)
	isURLExists := getRedirectFromURL != ""
	if !isURLExists {
		id := utils.GenerateURLID(6)
		createdAt := time.Now().String()
		err = utils.GenerateURL(id, redirectToURL, protocol, createdAt)
		if err != nil {
			template.Execute(res, Html{
				URL:   url,
				Error: err.Error(),
			})
			return
		}
		template.Execute(res, Html{
			URL:         url,
			RedirectURL: "https://" + req.Host + "/" + id,
		})
	} else {
		template.Execute(res, Html{
			URL:         url,
			RedirectURL: "https://" + req.Host + "/" + getRedirectFromURL,
		})
	}
}

// redirects to the URL if it exists in the database
func redirect(path string, res http.ResponseWriter, req *http.Request) {
	redirectToURL, err := utils.GetRedirectToURL(path)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Println("Redirecting to URl - " + redirectToURL)

	http.Redirect(res, req, redirectToURL, http.StatusTemporaryRedirect)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file, using default values")
	}

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		template, err := template.ParseFiles("static/index.html")
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println("req - " + req.Method + " - " + path + " FROM " + req.Host)

		if path == "/" && req.Method == "GET" {
			template.Execute(res, nil)
		} else if path == "/style.css" && req.Method == "GET" {
			http.ServeFile(res, req, "static/style.css")
		} else if path == "/icon.png" && req.Method == "GET" {
			http.ServeFile(res, req, "static/icon.png")
		} else if len(strings.Split(path, "/")) == 2 && req.Method == "GET" {
			redirect(strings.Split(path, "/")[1], res, req)
		} else if len(strings.Split(path, "/")) == 2 && req.Method == "POST" {
			render(res, req, template)
		} else {
			res.WriteHeader(http.StatusNotFound)
		}
	})

	fmt.Println("Server started on port 3000")
	http.ListenAndServe(":3000", nil)
}
