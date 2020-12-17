package main

import (
	"fmt"
	"net/http"
)

func main()  {
	fmt.Println("Hello , Go Web Development")

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer,"Hello Go Web Development\n")
	})
	fmt.Println(http.ListenAndServe(":8080",nil))
}
