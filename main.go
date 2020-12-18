package main

import (
	"fmt"
	"html/template"
	"net/http"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"encoding/json"
)

type Page struct {
	Name string
	Title string
	DBStatus bool
}

type SearchResult struct {
	Title string
	Author string
	Year string
	ID string
}

func main() {
	templates := template.Must(template.ParseFiles("templates/index.html"))

	db, _ := sql.Open("sqlite3","dev.db")

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		p := Page{Name: "Gopher",Title: "Home Page"}
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
		results := []SearchResult{
			{"Saeed-RM6","Saeed Rahimi Manesh","2020","22222"},
			{"Ahmad-Kay","Ahmad Kaya","1995","0618"},
		}

		encoder := json.NewEncoder(writer)
		if err := encoder.Encode(results); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	})
	fmt.Println(http.ListenAndServe(":8080", nil))
}
