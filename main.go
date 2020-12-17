package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Page struct {
	Name string
}

func main() {
	templates := template.Must(template.ParseFiles("templates/index.html"))

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		p := Page{Name: "Gopher"}
		if name := request.FormValue("name"); name != "" {
			p.Name = name
		}
		if err := templates.ExecuteTemplate(writer, "index.html", p); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	})
	fmt.Println(http.ListenAndServe(":8080", nil))
}
