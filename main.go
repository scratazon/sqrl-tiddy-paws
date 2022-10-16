package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var tpl *template.Template

const IndexHtml string = "index.gohtml"

type E621Json struct {
	Posts []struct {
		File struct {
			Url string
		}
		Tags struct {
			General []string
			Species []string
		}
	}
}

type FileExtension struct {
	Video bool
	Image bool
	Url   string
}

var e6Content FileExtension

func init() {
	tpl = template.Must(tpl.ParseGlob("templates/*.gohtml"))
}

func ec(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func findExtension(hornyUrl string) {
	if strings.Contains(hornyUrl, ".gif") {
		e6Content.Image = true
	}
	if strings.Contains(hornyUrl, ".jpeg") {
		e6Content.Image = true
	}
	if strings.Contains(hornyUrl, ".mp4") {
		e6Content.Video = true
	}
	if strings.Contains(hornyUrl, ".webm") {
		e6Content.Video = true
	}
	if strings.Contains(hornyUrl, ".png") {
		e6Content.Image = true
	}
	if strings.Contains(hornyUrl, ".jpg") {
		e6Content.Image = true
	}
}

func docRoot(w http.ResponseWriter, req *http.Request) {
	// e621 Says this is necessary to avoid them taking a cummy dump on your computer
	userAgent := "Never gonna give you up / I don't know how to use JSON sorry"
	client := &http.Client{}
	// Use a pseudo-random value for rand
	rand.Seed(time.Now().UnixNano())
	var e6Json E621Json
	var boobTags = []string{
		"boobie",
		"breasts",
		"boob_fuck",
		"boob_fondling",
		"boob_size_difference",
	}

	var hornyUrl string

	for {
		randomBoobs := boobTags[rand.Intn(len(boobTags))]
		boobUrl := "https://e621.net/posts.json?tags=squirrel+paws+" + randomBoobs
		request, err := http.NewRequest("GET", boobUrl, nil)
		ec(err)
		request.Header.Add("User-Agent", userAgent)

		// Initiate the hornt
		response, err := client.Do(request)
		ec(err)
		defer response.Body.Close()

		responseBody, err := io.ReadAll(response.Body)
		ec(err)
		json.Unmarshal(responseBody, &e6Json)
		// Start empty so that we never display a blank image
		hornyUrl = e6Json.Posts[rand.Intn(len(e6Json.Posts))].File.Url
		if hornyUrl != "" {
			break
		}
	}

	e6Content.Image = false
	e6Content.Video = false
	findExtension(hornyUrl)
	e6Content.Url = hornyUrl
	tpl.ExecuteTemplate(w, IndexHtml, e6Content)
}

func main() {
	http.HandleFunc("/", docRoot)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	log.Println("Server starting...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
