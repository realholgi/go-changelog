# go-changelog

`go-changelog` rendert Release-Abschnitte aus einer `CHANGELOG.md` und liefert
eine direkt einsetzbare Bootstrap-5-Anzeige mit Versionslink, Neu-Badge, Modal
und `localStorage`-Erinnerung.

## Installation

```sh
go get github.com/realholgi/go-changelog
```

## Verwendung

```go
package main

import (
	_ "embed"

	"github.com/realholgi/go-changelog"
)

//go:embed CHANGELOG.md
var changelogMarkdown string

var version = "dev"

func main() {
	inhalt := changelog.Aufbereiten(changelogMarkdown, version)
	// Inhalt als Teil der Daten an html/template übergeben.
	_ = inhalt
}
```

Im Go-Template werden Auslöser und Modal dort platziert, wo die Anwendung sie
benötigt:

```gotemplate
<nav>
  {{.Changelog.Ausloeser}}
</nav>

<script src="/static/bootstrap.bundle.min.js"></script>
{{.Changelog.Modal}}
```

Bootstrap 5 muss von der Anwendung geladen werden.

Die Anzeige kann angepasst werden:

```go
inhalt := changelog.AufbereitenMitOptionen(markdown, version, changelog.Optionen{
	ModalID:     "releaseNotes",
	StorageKey:  "my_app_changelog_seen",
	Titel:       "Release notes",
	NeuIn:       "New in",
	Schliessen:  "Close",
	Bestaetigen: "Got it",
})
```

Erwartet werden Release-Überschriften in dieser Form:

```markdown
## [Unreleased]

## [1.2.0] - 2026-07-24
```

`Inhalt.HTML` enthält alle nicht leeren Release-Abschnitte. `Inhalt.Datum`
enthält das Datum der aktuellen Version, sofern diese im Changelog vorhanden
ist.

## Lizenz

MIT
