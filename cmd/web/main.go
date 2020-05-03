package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("../../web/static"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Page struct {
	Title string
}

func handler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[1:]
	t, err  := template.ParseFiles("../../web/template/home.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	t.Execute(w, &Page{
		Title: title,
	})
}