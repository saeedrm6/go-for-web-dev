package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Page struct {
	Name     string
	Title    string
	DBStatus bool
}

type SearchResult struct {
	Title  string `xml:"title,attr"`
	Author string `xml:"author,attr"`
	Year   string `xml:"hyr,attr"`
	ID     string `xml:"owi,attr"`
}

func main() {
	templates := template.Must(template.ParseFiles("templates/index.html"))

	db, _ := sql.Open("sqlite3", "dev.db")

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		p := Page{Name: "Gopher", Title: "Home Page"}
		if name := request.FormValue("name"); name != "" {
			p.Name = name
		}
		p.DBStatus = db.Ping() == nil
		if err := templates.ExecuteTemplate(writer, "index.html", p); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		//db.Close()
	})

	http.HandleFunc("/search", func(writer http.ResponseWriter, request *http.Request) {
		/*
			[[[static mode :]]]
			results := []SearchResult{
					{"Saeed-RM6", "Saeed Rahimi Manesh", "2020", "22222"},
					{"Ahmad-Kay", "Ahmad Kaya", "1995", "0618"},
				}
		*/
		var results []SearchResult
		var err error
		if results, err = search(request.FormValue("search")); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		encoder := json.NewEncoder(writer)
		if err := encoder.Encode(results); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	})
	fmt.Println(http.ListenAndServe(":8080", nil))
}

type ClassifySearchResponse struct {
	Results []SearchResult `xml:"works>work"`
}

func search(query string) ([]SearchResult, error) {
	var resp *http.Response
	var err error

	if resp, err = http.Get("http://classify.oclc.org/classify2/Classify?&summary=true&title=" + url.QueryEscape(query)); err != nil {
		return []SearchResult{}, err
	}

	defer resp.Body.Close()
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return []SearchResult{}, err
	}

	var c ClassifySearchResponse
	err = xml.Unmarshal(body, &c)
	return c.Results, err
}
