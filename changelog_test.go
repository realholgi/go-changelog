package changelog

import (
	"strings"
	"testing"
)

func TestReleaseSectionStopsAtNextRelease(t *testing.T) {
	markdown := `# Changelog

## [Unreleased]

- Not released yet

## [1.2.0] - 2026-07-24

- Released

## [1.1.0]

- Older`

	document := parseReleaseDocument(markdown)
	if got, want := document.first("Unreleased").section, "- Not released yet"; got != want {
		t.Fatalf("first(Unreleased).section = %q, want %q", got, want)
	}
	if got, want := document.first("1.2.0").section, "- Released"; got != want {
		t.Fatalf("first(1.2.0).section = %q, want %q", got, want)
	}
	if got := document.first("9.9.9").section; got != "" {
		t.Fatalf("first(unknown).section = %q, want empty", got)
	}
}

func TestEmptyReleaseDoesNotIncludeNextRelease(t *testing.T) {
	markdown := "## [Empty]\n\n## [Populated]\n\n- Content\n"
	if got := parseReleaseDocument(markdown).first("Empty").section; got != "" {
		t.Fatalf("first(Empty).section = %q, want empty", got)
	}
}

func TestDuplicateVersionsUseTheFirstSection(t *testing.T) {
	markdown := "## [1.0.0] - 2026-01-01\n\n- First\n\n## [1.0.0] - 2026-02-01\n\n- Second\n"
	got := Prepare(markdown, "1.0.0")

	if got.Date != "2026-01-01" {
		t.Fatalf("Date = %q, want first heading date", got.Date)
	}
	if strings.Count(got.HTML, ">First</li>") != 2 || strings.Contains(got.HTML, ">Second</li>") {
		t.Fatalf("HTML does not preserve first-section lookup:\n%s", got.HTML)
	}
}

func TestReleaseHeadingRemainsPermissive(t *testing.T) {
	document := parseReleaseDocument("## [1.0.0] notes\n\n- Content\n")
	if got := document.first("1.0.0").section; got != "- Content" {
		t.Fatalf("first(1.0.0).section = %q, want content", got)
	}
}

func TestPrepareRendersEveryReleaseAndDate(t *testing.T) {
	markdown := "## [Unreleased]\n\n- Open\n\n## [1.2.0] - 2026-07-24\n\n- Done\n\n## [1.0.0]\n\n- Old\n"
	got := Prepare(markdown, "1.2.0")

	if got.Date != "2026-07-24" {
		t.Fatalf("Date = %q", got.Date)
	}
	for _, version := range []string{"Unreleased", "1.2.0", "1.0.0"} {
		if !strings.Contains(got.HTML, `data-version="`+version+`"`) {
			t.Fatalf("version %q missing:\n%s", version, got.HTML)
		}
	}
}

func TestMissingCurrentVersionHasNoDate(t *testing.T) {
	got := Prepare("## [1.0.0] - 2026-01-01\n\n- First version\n", "dev")
	if got.Date != "" {
		t.Fatalf("Date = %q, want empty", got.Date)
	}
}

func TestMarkdownToHTML(t *testing.T) {
	got := markdownToHTML("### Added\n\n- **New** feature\n")
	want := "<h3>Added</h3>\n<ul>\n<li><strong>New</strong> feature</li>\n</ul>\n"
	if got != want {
		t.Fatalf("markdownToHTML() = %q, want %q", got, want)
	}
}

func TestReleaseMetadataIsEscaped(t *testing.T) {
	got := Prepare("## [<script>]\n\n- Content\n", "dev").HTML
	if strings.Contains(got, "<script>") {
		t.Fatalf("HTML contains unescaped version:\n%s", got)
	}
}
