package main

import (
	"embed"
	"log"
	"net/http"
	"strings"
)

//go:embed template.html home.html
var contentFS embed.FS

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" || path == "template.html" {
			// Serve the main template for root or template path
			file, _ := contentFS.ReadFile("template.html")
			w.Header().Set("Content-Type", "text/html")
			w.Write(file)
			return
		}

		// Serve other content files for dynamic content
		content, err := contentFS.ReadFile(path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(content)
	})

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
