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

var version = "1.1.0"

func main() {
	inhalt := changelog.Aufbereiten(changelogMarkdown, version)
	vorlage := template.Must(template.ParseFiles("base.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		daten := struct {
			Changelog changelog.Inhalt
		}{
			Changelog: inhalt,
		}
		if err := vorlage.Execute(w, daten); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
