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

type ClassifySearchResponse struct {
	Results []SearchResult `xml:"works>work"`
}

type ClassifyBookResponse struct {
	BookData struct {
		Title  string `xml:"title,attr"`
		Author string `xml:"author,attr"`
		ID     string `xml:"owi,attr"`
	} `xml:"work"`
	Classification struct {
		MostPopular string `xml:"sfa,attr"`
	} `xml:"recommendations>ddc>mostPopular"`
}

func main() {
	templates := template.Must(template.ParseFiles("templates/index.html"))

	db, _ := sql.Open("sqlite3", "dev.db")

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		p := Page{Name: "Book store", Title: "Home Page"}
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

	http.HandleFunc("/books/add", func(writer http.ResponseWriter, request *http.Request) {
		var book ClassifyBookResponse
		var err error
		var target_book_id = request.FormValue("id")
		if book, err = find(target_book_id); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		if err = db.Ping(); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		_, err = db.Exec("insert into books (pk, title, author, id, classification) values (?,?,?,?,?)", nil, book.BookData.Title, book.BookData.Author, book.BookData.ID, book.Classification.MostPopular)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	})
	fmt.Println(http.ListenAndServe(":8080", nil))
}

func find(id string) (ClassifyBookResponse, error) {
	var c ClassifyBookResponse
	body, err := classifyAPI("http://classify.oclc.org/classify2/Classify?&summary=true&owi=" + url.QueryEscape(id))
	if err != nil {
		return ClassifyBookResponse{}, err
	}
	err = xml.Unmarshal(body, &c)
	return c, err
}

func search(query string) ([]SearchResult, error) {
	var c ClassifySearchResponse
	body, err := classifyAPI("http://classify.oclc.org/classify2/Classify?&summary=true&title=" + url.QueryEscape(query))
	if err != nil {
		return []SearchResult{}, err
	}
	err = xml.Unmarshal(body, &c)
	return c.Results, err
}

func classifyAPI(url string) ([]byte, error) {
	var resp *http.Response
	var err error

	if resp, err = http.Get(url); err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
