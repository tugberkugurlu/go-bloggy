package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Page struct {
	Title string
}

func handler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[1:]
	t, err  := template.ParseFiles("../../web/templates/home.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	t.Execute(w, &Page{
		Title: title,
	})
}