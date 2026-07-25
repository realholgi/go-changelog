// Package changelog prepares a Keep a Changelog-style Markdown document for a
// reusable Bootstrap display.
package changelog

import (
	"bytes"
	"html"
	"strings"

	"github.com/yuin/goldmark"
)

// Content contains the rendered changelog and the current release date.
type Content struct {
	HTML string
	Date string

	currentVersion string
	options        Options
}

// Prepare renders every non-empty release section with the default options.
func Prepare(markdown, currentVersion string) Content {
	return PrepareWithOptions(markdown, currentVersion, Options{})
}

// PrepareWithOptions renders every non-empty release section in source order
// and configures the display.
func PrepareWithOptions(markdown, currentVersion string, options Options) Content {
	document := parseReleaseDocument(markdown)
	options = options.withDefaults(currentVersion)

	var output strings.Builder
	for _, release := range document.releases {
		first := document.first(release.version)
		section := first.section
		if section == "" {
			continue
		}

		date := first.date
		output.WriteString(`<div class="changelog-section" data-version="`)
		output.WriteString(html.EscapeString(first.version))
		output.WriteString(`"><h3>Version `)
		output.WriteString(html.EscapeString(first.version))
		if date != "" {
			output.WriteString(` — `)
			output.WriteString(html.EscapeString(date))
		}
		output.WriteString(`</h3>`)
		output.WriteString(markdownToHTML(section))
		output.WriteString(`</div>`)
	}

	return Content{
		HTML:           output.String(),
		Date:           document.first(currentVersion).date,
		currentVersion: currentVersion,
		options:        options,
	}
}

type release struct {
	version string
	date    string
	section string
}

type releaseDocument struct {
	releases       []release
	firstByVersion map[string]release
}

func parseReleaseDocument(markdown string) releaseDocument {
	lines := strings.Split(markdown, "\n")
	var releases []release
	start := -1
	var current release
	for i, line := range lines {
		version, date, ok := releaseHeading(line)
		if !ok {
			continue
		}
		if start != -1 {
			current.section = strings.TrimSpace(strings.Join(lines[start:i], "\n"))
			releases = append(releases, current)
		}
		current = release{version: version, date: date}
		start = i + 1
	}
	if start != -1 {
		current.section = strings.TrimSpace(strings.Join(lines[start:], "\n"))
		releases = append(releases, current)
	}

	firstByVersion := make(map[string]release, len(releases))
	for _, release := range releases {
		if _, exists := firstByVersion[release.version]; !exists {
			firstByVersion[release.version] = release
		}
	}
	return releaseDocument{releases: releases, firstByVersion: firstByVersion}
}

func (d releaseDocument) first(version string) release {
	return d.firstByVersion[version]
}

func releaseHeading(line string) (version, date string, ok bool) {
	line = strings.TrimSuffix(line, "\r")
	if !strings.HasPrefix(line, "## [") {
		return "", "", false
	}
	end := strings.Index(line[len("## ["):], "]")
	if end == -1 {
		return "", "", false
	}
	end += len("## [")
	version = line[len("## ["):end]
	if version == "" {
		return "", "", false
	}

	rest := strings.TrimSpace(line[end+1:])
	if strings.HasPrefix(rest, "- ") {
		date = strings.TrimSpace(strings.TrimPrefix(rest, "- "))
	}
	return version, date, true
}

func markdownToHTML(markdown string) string {
	var output bytes.Buffer
	if err := goldmark.Convert([]byte(markdown), &output); err != nil {
		return markdown
	}
	return output.String()
}
