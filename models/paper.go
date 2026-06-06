package models

// Paper represents a research paper
type Paper struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	Abstract        string   `json:"abstract"`
	Authors         []Author `json:"authors"`
	Categories      []string `json:"categories"`
	PrimaryCategory string   `json:"primary_category"`
	Published       string   `json:"published"`
	Updated         string   `json:"updated"`
	AbstractURL     string   `json:"abstract_url"`
	PDFURL          string   `json:"pdf_url"`
	DOI             string   `json:"doi,omitempty"`
	JournalRef      string   `json:"journal_ref,omitempty"`
	Comment         string   `json:"comment,omitempty"`
	Source          string   `json:"source"`
}

// Author represents a paper author
type Author struct {
	Name        string `json:"name"`
	Affiliation string `json:"affiliation,omitempty"`
}

// SearchResult represents results from a search query
type SearchResult struct {
	Query      string  `json:"query"`
	TotalFound int     `json:"total_found"`
	Papers     []Paper `json:"papers"`
	Source     string  `json:"source"`
}

// PaperAnalysis represents the full analysis output
type PaperAnalysis struct {
	Paper          Paper    `json:"paper"`
	Summary        string   `json:"summary"`
	KeyFindings    []string `json:"key_findings"`
	RelatedPapers  []Paper  `json:"related_papers"`
	IsWorthReading bool     `json:"is_worth_reading"`
}
