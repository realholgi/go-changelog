package changelog

import (
	"bytes"
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

var triggerTemplate = template.Must(template.New("changelog-trigger").Parse(
	`{{if .Available}}<a href="#" class="text-white text-decoration-none" style="cursor: pointer;" data-bs-toggle="modal" data-bs-target="#{{.ModalID}}" title="{{.Title}}">
  {{.Version}}
</a>
<a id="{{.ModalID}}-badge" href="#"
   class="badge bg-warning text-dark ms-2 text-decoration-none d-none"
   data-bs-toggle="modal" data-bs-target="#{{.ModalID}}">
  {{.NewIn}} {{.Version}}
</a>{{else}}{{.Version}}{{end}}`,
))

var modalTemplate = template.Must(template.New("changelog-modal").Parse(
	`{{if .Available}}<style>
#{{.ModalID}} .modal-body h3 { margin-top: 1.2rem; margin-bottom: 0.6rem; font-size: 1.1rem; }
#{{.ModalID}} .modal-body h3:first-child { margin-top: 0; }
#{{.ModalID}} .modal-body ul { margin-bottom: 0.8rem; padding-left: 1.5rem; }
#{{.ModalID}} .modal-body li { margin-bottom: 0.4rem; }
#{{.ModalID}} .changelog-section { margin-bottom: 1.5rem; padding-bottom: 1.5rem; border-bottom: 1px solid #e9ecef; }
#{{.ModalID}} .changelog-section h3 { margin-top: 0; margin-bottom: 0.8rem; }
#{{.ModalID}} .changelog-section:last-child { border-bottom: none; }
</style>
<div class="modal fade" id="{{.ModalID}}" tabindex="-1" aria-labelledby="{{.ModalID}}Label" aria-hidden="true">
  <div class="modal-dialog modal-lg">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title" id="{{.ModalID}}Label">{{.Title}}</h5>
        <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="{{.Close}}"></button>
      </div>
      <div class="modal-body">
        <div class="changelog-container">{{.HTML}}</div>
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-primary" data-bs-dismiss="modal">{{.Confirm}}</button>
      </div>
    </div>
  </div>
</div>
<script>
(function () {
  var v = '{{.Version}}'.replace(/^v/, '');
  if (v === 'dev') return;
  var badge = document.getElementById('{{.ModalID}}-badge');
  if (!badge) return;
  var lastSeen = localStorage.getItem('{{.StorageKey}}');
  if (lastSeen !== v) {
    badge.classList.remove('d-none');
  }
  document.getElementById('{{.ModalID}}')
    .addEventListener('show.bs.modal', function () {
      localStorage.setItem('{{.StorageKey}}', v);
      badge.classList.add('d-none');
    });
}());
</script>{{end}}`,
))

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
