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
	options = options.withDefaults()

	var output strings.Builder
	for _, version := range versions(markdown) {
		section := releaseSection(markdown, version)
		if section == "" {
			continue
		}

		date := releaseDate(markdown, version)
		output.WriteString(`<div class="changelog-section" data-version="`)
		output.WriteString(html.EscapeString(version))
		output.WriteString(`"><h3>Version `)
		output.WriteString(html.EscapeString(version))
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
		Date:           releaseDate(markdown, currentVersion),
		currentVersion: currentVersion,
		options:        options,
	}
}

func versions(markdown string) []string {
	var result []string
	for _, line := range strings.Split(markdown, "\n") {
		if version, _, ok := releaseHeading(line); ok {
			result = append(result, version)
		}
	}
	return result
}

func releaseSection(markdown, wantedVersion string) string {
	lines := strings.Split(markdown, "\n")
	start := -1
	for i, line := range lines {
		version, _, ok := releaseHeading(line)
		if !ok {
			continue
		}
		if start != -1 {
			return strings.TrimSpace(strings.Join(lines[start:i], "\n"))
		}
		if version == wantedVersion {
			start = i + 1
		}
	}
	if start == -1 {
		return ""
	}
	return strings.TrimSpace(strings.Join(lines[start:], "\n"))
}

func releaseDate(markdown, wantedVersion string) string {
	for _, line := range strings.Split(markdown, "\n") {
		version, date, ok := releaseHeading(line)
		if ok && version == wantedVersion {
			return date
		}
	}
	return ""
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
