package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var tpl *template.Template

const IndexHtml string = "index.gohtml"

//const Error string = "err.gohtml"

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

func init() {
	tpl = template.Must(tpl.ParseGlob("templates/*.gohtml"))
}

func ec(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func docRoot(w http.ResponseWriter, req *http.Request) {

	//err := req.ParseForm()
	//ec(err)

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

	randomBoobs := boobTags[rand.Intn(len(boobTags))]

	boobUrl := "https://e621.net/posts.json?tags=squirrel+paws+" + randomBoobs

	request, err := http.NewRequest("GET", boobUrl, nil)
	ec(err)

	request.Header.Add("User-Agent", userAgent)

	response, err := client.Do(request)
	ec(err)

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	ec(err)

	json.Unmarshal(responseBody, &e6Json)

	tpl.ExecuteTemplate(w, IndexHtml, e6Json.Posts[rand.Intn(len(e6Json.Posts))].File.Url)

}

func main() {
	http.HandleFunc("/", docRoot)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	log.Println("Server starting...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
