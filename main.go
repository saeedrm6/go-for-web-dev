package main

import (
	"fmt"
	"html/template"
	"net/http"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Page struct {
	Name string
	Title string
	DBStatus bool
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
	})
	fmt.Println(http.ListenAndServe(":8080", nil))
}
