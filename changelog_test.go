package changelog

import (
	"strings"
	"testing"
)

func TestAbschnittGrenztReleasesVoneinanderAb(t *testing.T) {
	markdown := `# Changelog

## [Unreleased]

- Noch offen

## [1.2.0] - 2026-07-24

- Veröffentlicht

## [1.1.0]

- Älter`

	if got, want := abschnitt(markdown, "Unreleased"), "- Noch offen"; got != want {
		t.Fatalf("abschnitt(Unreleased) = %q, want %q", got, want)
	}
	if got, want := abschnitt(markdown, "1.2.0"), "- Veröffentlicht"; got != want {
		t.Fatalf("abschnitt(1.2.0) = %q, want %q", got, want)
	}
	if got := abschnitt(markdown, "9.9.9"); got != "" {
		t.Fatalf("abschnitt(unbekannt) = %q, want leer", got)
	}
}

func TestLeeresReleaseÜbernimmtKeinenFolgeinhalt(t *testing.T) {
	markdown := "## [Leer]\n\n## [Gefüllt]\n\n- Inhalt\n"
	if got := abschnitt(markdown, "Leer"); got != "" {
		t.Fatalf("abschnitt(Leer) = %q, want leer", got)
	}
}

func TestAufbereitenRendertAlleReleasesUndDatum(t *testing.T) {
	markdown := "## [Unreleased]\n\n- Offen\n\n## [1.2.0] - 2026-07-24\n\n- Fertig\n\n## [1.0.0]\n\n- Alt\n"
	got := Aufbereiten(markdown, "1.2.0")

	if got.Datum != "2026-07-24" {
		t.Fatalf("Datum = %q", got.Datum)
	}
	for _, version := range []string{"Unreleased", "1.2.0", "1.0.0"} {
		if !strings.Contains(got.HTML, `data-version="`+version+`"`) {
			t.Fatalf("Version %q fehlt:\n%s", version, got.HTML)
		}
	}
}

func TestFehlendeAktuelleVersionHatKeinDatum(t *testing.T) {
	got := Aufbereiten("## [1.0.0] - 2026-01-01\n\n- Erste Version\n", "dev")
	if got.Datum != "" {
		t.Fatalf("Datum = %q, want leer", got.Datum)
	}
}

func TestMarkdownZuHTML(t *testing.T) {
	got := markdownZuHTML("### Hinzugefügt\n\n- **Neue** Funktion\n")
	want := "<h3>Hinzugefügt</h3>\n<ul>\n<li><strong>Neue</strong> Funktion</li>\n</ul>\n"
	if got != want {
		t.Fatalf("markdownZuHTML() = %q, want %q", got, want)
	}
}

func TestReleaseMetadatenWerdenEscaped(t *testing.T) {
	got := Aufbereiten("## [<script>]\n\n- Inhalt\n", "dev").HTML
	if strings.Contains(got, "<script>") {
		t.Fatalf("HTML enthält ungeescapte Version:\n%s", got)
	}
}
