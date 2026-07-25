# go-changelog

`go-changelog` renders release sections from a `CHANGELOG.md` and provides a
ready-to-use Bootstrap 5 display with a version link, new-release badge, modal,
and `localStorage` reminder.

## Installation

```sh
go get github.com/realholgi/go-changelog
```

## Usage

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
	content := changelog.Prepare(changelogMarkdown, version)
	// Pass Content to html/template as part of the template data.
	_ = content
}
```

Place the trigger and modal where the application needs them:

```gotemplate
<nav>
  {{.Changelog.Trigger}}
</nav>

<script src="/static/bootstrap.bundle.min.js"></script>
{{.Changelog.Modal}}
```

The embedding application must load Bootstrap 5.

Customize the display with options:

```go
content := changelog.PrepareWithOptions(markdown, version, changelog.Options{
	ModalID:    "releaseNotes",
	StorageKey: "my_app_changelog_seen",
	Title:      "Release notes",
	NewIn:      "New in",
	Close:      "Close",
	Confirm:    "Got it",
})
```

Release headings must use this format:

```markdown
## [Unreleased]

## [1.2.0] - 2026-07-24
```

`Content.HTML` contains every non-empty release section. `Content.Date`
contains the date of the current version when that version exists in the
changelog.

## License

MIT
