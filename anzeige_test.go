package changelog

import (
	"strings"
	"testing"
)

func TestAnzeigeMitStandardoptionen(t *testing.T) {
	inhalt := Aufbereiten("## [1.2.0]\n\n- Neu\n", "v1.2.0")
	ausloeser := string(inhalt.Ausloeser())
	modal := string(inhalt.Modal())

	for _, want := range []string{
		`data-bs-target="#changelogModal"`,
		`id="changelogModal-badge"`,
		`Neu in v1.2.0`,
	} {
		if !strings.Contains(ausloeser, want) {
			t.Fatalf("Ausloeser enthält %q nicht:\n%s", want, ausloeser)
		}
	}
	for _, want := range []string{
		`id="changelogModal"`,
		`localStorage.getItem('changelog_seen')`,
		`localStorage.setItem('changelog_seen', v)`,
		`var v = 'v1.2.0'.replace(/^v/, '');`,
	} {
		if !strings.Contains(modal, want) {
			t.Fatalf("Modal enthält %q nicht:\n%s", want, modal)
		}
	}
}

func TestAnzeigeMitEigenenOptionen(t *testing.T) {
	inhalt := AufbereitenMitOptionen(
		"## [2.0.0]\n\n- New\n",
		"2.0.0",
		Optionen{
			ModalID:     "releaseNotes",
			StorageKey:  "my_app_seen",
			Titel:       "Release notes",
			NeuIn:       "New in",
			Schliessen:  "Close",
			Bestaetigen: "Got it",
		},
	)
	got := string(inhalt.Ausloeser()) + string(inhalt.Modal())

	for _, want := range []string{
		`id="releaseNotes"`,
		`id="releaseNotes-badge"`,
		`New in 2.0.0`,
		`Release notes`,
		`aria-label="Close"`,
		`>Got it</button>`,
		`localStorage.getItem('my_app_seen')`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("Anzeige enthält %q nicht:\n%s", want, got)
		}
	}
}

func TestAnzeigeOhneChangelog(t *testing.T) {
	inhalt := Aufbereiten("", "dev")
	if got := string(inhalt.Ausloeser()); got != "dev" {
		t.Fatalf("Ausloeser = %q, want dev", got)
	}
	if got := inhalt.Modal(); got != "" {
		t.Fatalf("Modal = %q, want leer", got)
	}
}

func TestAnzeigeEscapedVersionUndOptionen(t *testing.T) {
	inhalt := AufbereitenMitOptionen(
		"## [1.0.0]\n\n- Neu\n",
		`<script>alert("x")</script>`,
		Optionen{Titel: `<img src=x onerror=alert(1)>`},
	)
	got := string(inhalt.Ausloeser()) + string(inhalt.Modal())
	if strings.Contains(got, `<script>alert("x")</script>`) || strings.Contains(got, "<img src=x") {
		t.Fatalf("Anzeige enthält ungeescapte Werte:\n%s", got)
	}
}
