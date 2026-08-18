package main

import (
	_ "embed"
	"html/template"
	"log"
	"net/http"

	"github.com/realholgi/go-changelog"
)

//go:embed CHANGELOG.md
var changelogMarkdown string

// version is overridden for release builds with:
// go build -ldflags "-X main.version=v1.1.0" .
var version = "dev"

func main() {
	content := changelog.Prepare(changelogMarkdown, version)
	page := template.Must(template.ParseFiles("base.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		data := struct {
			Changelog changelog.Content
		}{
			Changelog: content,
		}
		if err := page.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
