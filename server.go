package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type PageData struct {
	Body template.HTML
}

func renderTemplate(w http.ResponseWriter, bodyTemplate string) {
	// Parse the base template
	baseTemplate, err := template.ParseFiles("views/base.html")
	if err != nil {
		log.Printf("Error loading base template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Parse the body template
	body, err := template.New("body").Parse(bodyTemplate)
	if err != nil {
		log.Printf("Error parsing body template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the body template to produce the dynamic content
	var bodyContent bytes.Buffer
	if err := body.Execute(&bodyContent, nil); err != nil {
		log.Printf("Error executing body template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render the final page using the base template with the dynamic body
	data := PageData{Body: template.HTML(bodyContent.String())}
	if err := baseTemplate.Execute(w, data); err != nil {
		log.Printf("Error rendering base template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func bodyHandler(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Path[1:]
	fmt.Println("User requested path: " + request)
	if !strings.HasPrefix(request,"views") {
		fmt.Println("Invalid request")
		return
	}

	renderTemplate(w, strings.Split(request, "/")[1])
}

func main() {
	http.HandleFunc("/", bodyHandler)
	log.Println("Server is running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
