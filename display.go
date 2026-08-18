package changelog

import (
	"bytes"
	"embed"
	"html/template"
	"strings"
	"unicode"
)

// Options customizes display IDs, browser storage, and labels.
// Empty IDs and storage keys are derived from the current version; empty labels
// are replaced by English defaults.
type Options struct {
	ModalID    string
	StorageKey string
	Title      string
	NewIn      string
	Close      string
	Confirm    string
}

func (o Options) withDefaults(version string) Options {
	identity := displayIdentity(version)
	if o.ModalID == "" {
		o.ModalID = identity.modalID
	} else {
		o.ModalID = normalizeIdentifier(o.ModalID)
	}
	if o.StorageKey == "" {
		o.StorageKey = identity.storageKey
	}
	if o.Title == "" {
		o.Title = "Changelog"
	}
	if o.NewIn == "" {
		o.NewIn = "New in"
	}
	if o.Close == "" {
		o.Close = "Close"
	}
	if o.Confirm == "" {
		o.Confirm = "OK"
	}
	return o
}

type displayIDs struct {
	modalID    string
	storageKey string
}

func displayIdentity(version string) displayIDs {
	suffix := normalizeIdentifier(version)
	return displayIDs{
		modalID:    "changelogModal-" + suffix,
		storageKey: "changelog_seen_" + suffix,
	}
}

func normalizeIdentifier(value string) string {
	var output strings.Builder
	previousSeparator := false
	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
			output.WriteRune(r)
			previousSeparator = false
			continue
		}
		if !previousSeparator {
			output.WriteByte('-')
			previousSeparator = true
		}
	}
	normalized := strings.Trim(output.String(), "-")
	if normalized == "" {
		return "changelog"
	}
	first := normalized[0]
	if first >= '0' && first <= '9' {
		return "changelog-" + normalized
	}
	return normalized
}

//go:embed templates/*.html
var displayTemplates embed.FS

var triggerTemplate = template.Must(template.ParseFS(displayTemplates, "templates/trigger.html"))

var modalTemplate = template.Must(template.ParseFS(displayTemplates, "templates/modal.html"))

type view struct {
	Version    string
	HTML       template.HTML
	Available  bool
	ModalID    string
	StorageKey string
	Title      string
	NewIn      string
	Close      string
	Confirm    string
}

// Trigger renders the version display and, when changelog content is
// available, the trigger and badge for the Bootstrap modal.
func (c Content) Trigger() template.HTML {
	return execute(triggerTemplate, c.view())
}

// Modal renders the styles, Bootstrap modal, and localStorage behavior.
// The embedding application must load Bootstrap 5.
func (c Content) Modal() template.HTML {
	return execute(modalTemplate, c.view())
}

func (c Content) view() view {
	return view{
		Version:    c.currentVersion,
		HTML:       template.HTML(c.HTML),
		Available:  c.HTML != "",
		ModalID:    c.options.ModalID,
		StorageKey: c.options.StorageKey,
		Title:      c.options.Title,
		NewIn:      c.options.NewIn,
		Close:      c.options.Close,
		Confirm:    c.options.Confirm,
	}
}

func execute(tmpl *template.Template, data view) template.HTML {
	var output bytes.Buffer
	if err := tmpl.Execute(&output, data); err != nil {
		return ""
	}
	return template.HTML(output.String())
}
