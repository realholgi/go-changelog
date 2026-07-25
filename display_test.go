package changelog

import (
	"strings"
	"testing"
)

func TestDisplayWithDefaultOptions(t *testing.T) {
	content := Prepare("## [1.2.0]\n\n- New\n", "v1.2.0")
	trigger := string(content.Trigger())
	modal := string(content.Modal())

	for _, want := range []string{
		`data-bs-target="#changelogModal-v1-2-0"`,
		`id="changelogModal-v1-2-0-badge"`,
		`New in v1.2.0`,
	} {
		if !strings.Contains(trigger, want) {
			t.Fatalf("Trigger does not contain %q:\n%s", want, trigger)
		}
	}
	for _, want := range []string{
		`id="changelogModal-v1-2-0"`,
		`aria-label="Close"`,
		`localStorage.getItem('changelog_seen_v1-2-0')`,
		`localStorage.setItem('changelog_seen_v1-2-0', v)`,
		`var v = 'v1.2.0'.replace(/^v/, '');`,
	} {
		if !strings.Contains(modal, want) {
			t.Fatalf("Modal does not contain %q:\n%s", want, modal)
		}
	}
}

func TestDisplayNormalizesCustomModalID(t *testing.T) {
	content := PrepareWithOptions(
		"## [1.0.0]\n\n- New\n",
		"1.0.0",
		Options{ModalID: "release notes/#1"},
	)
	got := string(content.Trigger()) + string(content.Modal())

	for _, want := range []string{
		`id="release-notes-1"`,
		`data-bs-target="#release-notes-1"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("display does not contain normalized identifier %q:\n%s", want, got)
		}
	}
}

func TestDefaultDisplayIdentityUsesVersion(t *testing.T) {
	first := Prepare("## [1.0.0]\n\n- First\n", "1.0.0")
	second := Prepare("## [2.0.0]\n\n- Second\n", "2.0.0")

	if string(first.Modal()) == string(second.Modal()) {
		t.Fatal("default displays share an identity")
	}
}

func TestDisplayPrefixesNumericCustomModalID(t *testing.T) {
	content := PrepareWithOptions(
		"## [1.0.0]\n\n- New\n",
		"1.0.0",
		Options{ModalID: "1 release"},
	)
	if got := string(content.Modal()); !strings.Contains(got, `id="changelog-1-release"`) {
		t.Fatalf("display does not contain a CSS-safe identifier:\n%s", got)
	}
}

func TestDisplayWithCustomOptions(t *testing.T) {
	content := PrepareWithOptions(
		"## [2.0.0]\n\n- New\n",
		"2.0.0",
		Options{
			ModalID:    "releaseNotes",
			StorageKey: "my_app_seen",
			Title:      "Release notes",
			NewIn:      "Fresh in",
			Close:      "Dismiss",
			Confirm:    "Got it",
		},
	)
	got := string(content.Trigger()) + string(content.Modal())

	for _, want := range []string{
		`id="releaseNotes"`,
		`id="releaseNotes-badge"`,
		`Fresh in 2.0.0`,
		`Release notes`,
		`aria-label="Dismiss"`,
		`>Got it</button>`,
		`localStorage.getItem('my_app_seen')`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("display does not contain %q:\n%s", want, got)
		}
	}
}

func TestDisplayWithoutChangelog(t *testing.T) {
	content := Prepare("", "dev")
	if got := string(content.Trigger()); got != "dev" {
		t.Fatalf("Trigger = %q, want dev", got)
	}
	if got := content.Modal(); got != "" {
		t.Fatalf("Modal = %q, want empty", got)
	}
}

func TestDisplayEscapesVersionAndOptions(t *testing.T) {
	content := PrepareWithOptions(
		"## [1.0.0]\n\n- New\n",
		`<script>alert("x")</script>`,
		Options{Title: `<img src=x onerror=alert(1)>`},
	)
	got := string(content.Trigger()) + string(content.Modal())
	if strings.Contains(got, `<script>alert("x")</script>`) || strings.Contains(got, "<img src=x") {
		t.Fatalf("display contains unescaped values:\n%s", got)
	}
}
