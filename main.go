package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"vineelsai.com/shortify/src/utils"
)

type RedirectHTMLResponse struct {
	URL         string
	RedirectURL string
	Error       string
}

type RedirectResponse struct {
	RedirectToUrl   string   `json:"redirectToUrl"`
	RedirectFromUrl string   `json:"redirectFromUrl"`
	Errors          []string `json:"errors"`
}

func generateRedirectURL(url string) (string, error) {
	redirectToURL, protocol, err := utils.GetFormattedURL(url)
	if err != nil {
		return "", err
	}

	getRedirectFromURL := utils.GetRedirectFromURL(redirectToURL)
	isURLExists := getRedirectFromURL != ""
	if !isURLExists {
		id := utils.GenerateURLID(6)
		createdAt := time.Now().String()
		err = utils.GenerateURL(id, redirectToURL, protocol, createdAt)
		if err != nil {
			return "", err
		}
		return id, nil
	} else {
		return getRedirectFromURL, nil
	}
}

func render(res http.ResponseWriter, req *http.Request, template *template.Template) {
	url := req.PostFormValue("url")
	redirectToURL, err := generateRedirectURL(url)

	if err != nil {
		template.Execute(res, RedirectHTMLResponse{
			URL:   url,
			Error: err.Error(),
		})
		return
	} else {
		template.Execute(res, RedirectHTMLResponse{
			URL:         url,
			RedirectURL: "https://" + req.Host + "/" + redirectToURL,
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
		} else if path == "/api" && req.Method == "POST" {
			url := req.PostFormValue("url")
			redirectToURL, err := generateRedirectURL(url)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			redirectResponse := RedirectResponse{RedirectToUrl: url, RedirectFromUrl: "https://" + req.Host + "/" + redirectToURL, Errors: []string{}}
			json.NewEncoder(res).Encode(redirectResponse)
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
