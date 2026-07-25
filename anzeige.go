package changelog

import (
	"bytes"
	"html/template"
)

// Optionen passt IDs, Browser-Speicher und Beschriftungen der Anzeige an.
// Leere Felder werden durch deutsche Standardwerte ersetzt.
type Optionen struct {
	ModalID     string
	StorageKey  string
	Titel       string
	NeuIn       string
	Schliessen  string
	Bestaetigen string
}

func (o Optionen) mitStandardwerten() Optionen {
	if o.ModalID == "" {
		o.ModalID = "changelogModal"
	}
	if o.StorageKey == "" {
		o.StorageKey = "changelog_seen"
	}
	if o.Titel == "" {
		o.Titel = "Changelog"
	}
	if o.NeuIn == "" {
		o.NeuIn = "Neu in"
	}
	if o.Schliessen == "" {
		o.Schliessen = "Schließen"
	}
	if o.Bestaetigen == "" {
		o.Bestaetigen = "OK"
	}
	return o
}

var ausloeserVorlage = template.Must(template.New("changelog-ausloeser").Parse(
	`{{if .Vorhanden}}<a href="#" class="text-white text-decoration-none" style="cursor: pointer;" data-bs-toggle="modal" data-bs-target="#{{.ModalID}}" title="{{.Titel}} anzeigen">
  {{.Version}}
</a>
<a id="{{.ModalID}}-badge" href="#"
   class="badge bg-warning text-dark ms-2 text-decoration-none d-none"
   data-bs-toggle="modal" data-bs-target="#{{.ModalID}}">
  {{.NeuIn}} {{.Version}}
</a>{{else}}{{.Version}}{{end}}`,
))

var modalVorlage = template.Must(template.New("changelog-modal").Parse(
	`{{if .Vorhanden}}<style>
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
        <h5 class="modal-title" id="{{.ModalID}}Label">{{.Titel}}</h5>
        <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="{{.Schliessen}}"></button>
      </div>
      <div class="modal-body">
        <div class="changelog-container">{{.HTML}}</div>
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-primary" data-bs-dismiss="modal">{{.Bestaetigen}}</button>
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

type ansicht struct {
	Version     string
	HTML        template.HTML
	Vorhanden   bool
	ModalID     string
	StorageKey  string
	Titel       string
	NeuIn       string
	Schliessen  string
	Bestaetigen string
}

// Ausloeser rendert die Versionsanzeige und, falls ein Changelog vorhanden
// ist, den Auslöser samt Badge für das Bootstrap-Modal.
func (i Inhalt) Ausloeser() template.HTML {
	return ausführen(ausloeserVorlage, i.ansicht())
}

// Modal rendert Styles, Bootstrap-Modal und das localStorage-Verhalten.
// Bootstrap 5 muss von der einbindenden Anwendung geladen werden.
func (i Inhalt) Modal() template.HTML {
	return ausführen(modalVorlage, i.ansicht())
}

func (i Inhalt) ansicht() ansicht {
	return ansicht{
		Version:     i.aktuelleVersion,
		HTML:        template.HTML(i.HTML),
		Vorhanden:   i.HTML != "",
		ModalID:     i.optionen.ModalID,
		StorageKey:  i.optionen.StorageKey,
		Titel:       i.optionen.Titel,
		NeuIn:       i.optionen.NeuIn,
		Schliessen:  i.optionen.Schliessen,
		Bestaetigen: i.optionen.Bestaetigen,
	}
}

func ausführen(vorlage *template.Template, daten ansicht) template.HTML {
	var ausgabe bytes.Buffer
	if err := vorlage.Execute(&ausgabe, daten); err != nil {
		return ""
	}
	return template.HTML(ausgabe.String())
}
