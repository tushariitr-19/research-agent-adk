package util

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/tushariitr-19/research-agent-adk/models"
)

// ArxivFeed represents the top-level XML response from ArXiv
type ArxivFeed struct {
	XMLName      xml.Name     `xml:"feed"`
	TotalResults int          `xml:"totalResults"`
	Entries      []ArxivEntry `xml:"entry"`
}

// ArxivEntry represents a single paper in the ArXiv response
type ArxivEntry struct {
	ID              string          `xml:"id"`
	Title           string          `xml:"title"`
	Summary         string          `xml:"summary"`
	Published       string          `xml:"published"`
	Updated         string          `xml:"updated"`
	Authors         []ArxivAuthor   `xml:"author"`
	Categories      []ArxivCategory `xml:"category"`
	PrimaryCategory ArxivCategory   `xml:"primary_category"`
	Links           []ArxivLink     `xml:"link"`
	Comment         string          `xml:"comment"`
	JournalRef      string          `xml:"journal_ref"`
	DOI             string          `xml:"doi"`
}

// ArxivAuthor represents a paper author in XML
type ArxivAuthor struct {
	Name        string `xml:"name"`
	Affiliation string `xml:"affiliation"`
}

// ArxivCategory represents a paper category in XML
type ArxivCategory struct {
	Term string `xml:"term,attr"`
}

// ArxivLink represents a link in XML
type ArxivLink struct {
	Href  string `xml:"href,attr"`
	Rel   string `xml:"rel,attr"`
	Type  string `xml:"type,attr"`
	Title string `xml:"title,attr"`
}

// ParseArxivFeed parses raw XML bytes into an ArxivFeed
func ParseArxivFeed(data []byte) (*ArxivFeed, error) {
	var feed ArxivFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, fmt.Errorf("failed to parse ArXiv XML: %w", err)
	}
	return &feed, nil
}

// EntryToPaper converts an ArXiv XML entry to our Paper model
func EntryToPaper(entry ArxivEntry) *models.Paper {
	var abstractURL, pdfURL string
	for _, link := range entry.Links {
		if link.Type == "text/html" {
			abstractURL = link.Href
		}
		if link.Title == "pdf" {
			pdfURL = link.Href
		}
	}

	categories := make([]string, 0, len(entry.Categories))
	for _, cat := range entry.Categories {
		categories = append(categories, cat.Term)
	}

	authors := make([]models.Author, 0, len(entry.Authors))
	for _, a := range entry.Authors {
		authors = append(authors, models.Author{
			Name:        strings.TrimSpace(a.Name),
			Affiliation: strings.TrimSpace(a.Affiliation),
		})
	}

	// Extract clean ID from URL and remove version suffix
	id := entry.ID
	if idx := strings.LastIndex(id, "/"); idx != -1 {
		id = id[idx+1:]
	}
	if idx := strings.LastIndex(id, "v"); idx != -1 {
		id = id[:idx]
	}

	return &models.Paper{
		ID:              id,
		Title:           strings.TrimSpace(entry.Title),
		Abstract:        strings.TrimSpace(entry.Summary),
		Authors:         authors,
		Categories:      categories,
		PrimaryCategory: entry.PrimaryCategory.Term,
		Published:       entry.Published,
		Updated:         entry.Updated,
		AbstractURL:     abstractURL,
		PDFURL:          pdfURL,
		DOI:             entry.DOI,
		JournalRef:      entry.JournalRef,
		Comment:         entry.Comment,
		Source:          ArxivSource,
	}
}

func BuildArxivURL(params map[string]string) string {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return ArxivBaseURL + "?" + values.Encode()
}

// ParseDate formats an ISO date string
func ParseDate(raw string) string {
	raw = strings.TrimSpace(raw)
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return raw
	}
	return t.Format("2006-01-02")
}
