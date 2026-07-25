// Paket changelog bereitet ein Keep-a-Changelog-ähnliches Markdown-Dokument
// für eine wiederverwendbare Bootstrap-Anzeige auf.
package changelog

import (
	"bytes"
	"html"
	"strings"

	"github.com/yuin/goldmark"
)

// Inhalt enthält das gerenderte Changelog und das Datum der aktuellen Version.
type Inhalt struct {
	HTML  string
	Datum string

	aktuelleVersion string
	optionen        Optionen
}

// Aufbereiten rendert alle nicht leeren Release-Abschnitte mit den
// Standardoptionen.
func Aufbereiten(markdown, aktuelleVersion string) Inhalt {
	return AufbereitenMitOptionen(markdown, aktuelleVersion, Optionen{})
}

// AufbereitenMitOptionen rendert alle nicht leeren Release-Abschnitte in ihrer
// ursprünglichen Reihenfolge und konfiguriert die Anzeige.
func AufbereitenMitOptionen(markdown, aktuelleVersion string, optionen Optionen) Inhalt {
	optionen = optionen.mitStandardwerten()

	var ausgabe strings.Builder
	for _, version := range versionen(markdown) {
		abschnitt := abschnitt(markdown, version)
		if abschnitt == "" {
			continue
		}

		datum := releaseDatum(markdown, version)
		ausgabe.WriteString(`<div class="changelog-section" data-version="`)
		ausgabe.WriteString(html.EscapeString(version))
		ausgabe.WriteString(`"><h3>Version `)
		ausgabe.WriteString(html.EscapeString(version))
		if datum != "" {
			ausgabe.WriteString(` — `)
			ausgabe.WriteString(html.EscapeString(datum))
		}
		ausgabe.WriteString(`</h3>`)
		ausgabe.WriteString(markdownZuHTML(abschnitt))
		ausgabe.WriteString(`</div>`)
	}

	return Inhalt{
		HTML:            ausgabe.String(),
		Datum:           releaseDatum(markdown, aktuelleVersion),
		aktuelleVersion: aktuelleVersion,
		optionen:        optionen,
	}
}

func versionen(markdown string) []string {
	var ergebnis []string
	for _, zeile := range strings.Split(markdown, "\n") {
		if version, _, ok := releaseKopf(zeile); ok {
			ergebnis = append(ergebnis, version)
		}
	}
	return ergebnis
}

func abschnitt(markdown, gesuchteVersion string) string {
	zeilen := strings.Split(markdown, "\n")
	start := -1
	for i, zeile := range zeilen {
		version, _, ok := releaseKopf(zeile)
		if !ok {
			continue
		}
		if start != -1 {
			return strings.TrimSpace(strings.Join(zeilen[start:i], "\n"))
		}
		if version == gesuchteVersion {
			start = i + 1
		}
	}
	if start == -1 {
		return ""
	}
	return strings.TrimSpace(strings.Join(zeilen[start:], "\n"))
}

func releaseDatum(markdown, gesuchteVersion string) string {
	for _, zeile := range strings.Split(markdown, "\n") {
		version, datum, ok := releaseKopf(zeile)
		if ok && version == gesuchteVersion {
			return datum
		}
	}
	return ""
}

func releaseKopf(zeile string) (version, datum string, ok bool) {
	zeile = strings.TrimSuffix(zeile, "\r")
	if !strings.HasPrefix(zeile, "## [") {
		return "", "", false
	}
	ende := strings.Index(zeile[len("## ["):], "]")
	if ende == -1 {
		return "", "", false
	}
	ende += len("## [")
	version = zeile[len("## ["):ende]
	if version == "" {
		return "", "", false
	}

	rest := strings.TrimSpace(zeile[ende+1:])
	if strings.HasPrefix(rest, "- ") {
		datum = strings.TrimSpace(strings.TrimPrefix(rest, "- "))
	}
	return version, datum, true
}

func markdownZuHTML(markdown string) string {
	var ausgabe bytes.Buffer
	if err := goldmark.Convert([]byte(markdown), &ausgabe); err != nil {
		return markdown
	}
	return ausgabe.String()
}
