package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type Page struct {
	Title string
	Body  []byte
}

func handler(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Path[1:]
	title := ""

	if strings.HasPrefix(request, "views") {
		title = strings.Split(request, "/")[1]
	} else {
		fmt.Println("invalid request: " + request)
		return
	}

	filename := "views/" + title + ".html"
	fmt.Println("Getting page: " + filename)
	body, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Got body: " + string(body))

  base,err := template.New("base").ParseFiles("views/base.html", filename)
	if err != nil {
		log.Fatal(err)
	}

	pageData := &Page{Title: title, Body: body}

	//tp.ExecuteTemplate(w, pg, page)
	base.Execute(w, pageData)
  buf := bytes.Buffer
}

func main() {
	http.Handle("/views", http.FileServer(http.Dir("./views")))
	http.HandleFunc("/", handler)

	fmt.Println("Serving https://localhost:443")
	log.Fatal(http.ListenAndServeTLS(":443", "server.crt", "server.key", nil))
}
