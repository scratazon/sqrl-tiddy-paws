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

type E6Generated struct {
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

type le struct {
	Result string
}

var l = le{
	Result: "asdf",
}

func docRoot(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	ec(err)

	ua := "Never gonna give you up / I don't know how to use JSON sorry"

	client := &http.Client{}

	rand.Seed(time.Now().UnixNano())
	var bj E6Generated
	//	var boobTags = []string{
	//		"boobie",
	//		"breasts",
	//		//"boob_fuck",
	//		//"boob_fondling",
	//		//"boob_size_difference",
	//	}
	//
	//randomBoobs := boobTags[0] //boobTags[rand.Intn(len(boobTags))]

	boobURL := "https://e621.net/posts.json?tags=squirrel+paws+boobie" // + randomBoobs

	req, e := http.NewRequest("GET", boobURL, nil)
	ec(e)

	req.Header.Add("User-Agent", ua)

	resp, e := client.Do(req)
	ec(e)

	defer resp.Body.Close()

	result, e := io.ReadAll(resp.Body)
	ec(e)
	json.Unmarshal(result, &bj)
	l.Result = bj.Posts[rand.Intn(len(bj.Posts))].File.Url
	tpl.ExecuteTemplate(w, IndexHtml, l)

}

func main() {
	http.HandleFunc("/", docRoot)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	log.Println("Server starting...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
