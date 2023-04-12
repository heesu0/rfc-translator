package rfc

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Regex pattern.
const (
	newLinesPattern         = `\r?\n{2,}\s*`
	newLinesPattern2        = `\r?\n{2,}`
	whiteSpacesPattern      = `(?:\r?\n|\s)+`
	firstPageHeaderPattern  = `^.*\s+(?:Standards Track|Informational|Experimental|Best Current Practice|Historic)\s+\[Page [0-9]+\]$`
	secondPageHeaderPattern = `^RFC\s[0-9]+\s+[\w\s]+\s+\w+\s\d{4}$`
	tableOfContentsPattern  = `Table\sof\sContents((?s).*)1\.\s{2}Introduction`
)

const (
	newLine             = "\n"
	newLines            = "\n\n"
	newLinesPlaceholder = "NEWLINES_PLACEHOLDER"
	tablePlaceholder    = "TABLE_PLACEHOLDER"
)

// Document represents an RFC document and contains the methods to process and fetch the document text.
type Document struct {
	client *http.Client
	regex  *regexPatterns
}

// regexPatterns contains all the compiled regex patterns for cleaning the document text.
type regexPatterns struct {
	newLines         *regexp.Regexp
	newLines2        *regexp.Regexp
	whiteSpaces      *regexp.Regexp
	firstPageHeader  *regexp.Regexp
	secondPageHeader *regexp.Regexp
	tableOfContents  *regexp.Regexp
}

// NewDocument initializes a new Document struct and sets up the regex patterns and HTTP client.
func NewDocument() *Document {
	d := &Document{
		client: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: false,
					MinVersion:         tls.VersionTLS12,
				},
				TLSHandshakeTimeout: 10 * time.Second,
			},
		},
		regex: &regexPatterns{
			newLines:         regexp.MustCompile(newLinesPattern),
			newLines2:        regexp.MustCompile(newLinesPattern2),
			whiteSpaces:      regexp.MustCompile(whiteSpacesPattern),
			firstPageHeader:  regexp.MustCompile(firstPageHeaderPattern),
			secondPageHeader: regexp.MustCompile(secondPageHeaderPattern),
			tableOfContents:  regexp.MustCompile(tableOfContentsPattern),
		},
	}

	return d
}

// GetText fetches the plain-text content of an RFC document with the given number and returns the cleaned text.
func (d *Document) GetText(number int) (string, error) {
	url := fmt.Sprintf("https://www.rfc-editor.org/rfc/rfc%d.txt", number)

	response, err := d.client.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var text string
	if response.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return text, err
		}

		text = d.removePageHeader(string(bodyBytes))
		text = d.cleanWhitespace(text)

		return text, nil
	}

	return text, fmt.Errorf("failed to fetch RFC document, status code: %d", response.StatusCode)
}

// removePageHeader removes the first and second page headers from the input text.
func (d *Document) removePageHeader(text string) string {
	lines := strings.Split(text, "\n")
	filteredLines := []string{}
	for _, line := range lines {
		if !d.regex.firstPageHeader.MatchString(line) && !d.regex.secondPageHeader.MatchString(line) {
			filteredLines = append(filteredLines, line)
		}
	}

	return strings.Join(filteredLines, newLine)
}

// cleanWhitespace normalizes newlines and spaces in the input string by removing extra whitespace while preserving the desired formatting.
func (d *Document) cleanWhitespace(text string) string {
	// Extract the table of contents section and clean it separately
	table := d.regex.tableOfContents.FindString(text)
	table = d.regex.newLines2.ReplaceAllString(table, newLine)

	// Replace the table of contents with a placeholder in the original text.
	text = d.regex.tableOfContents.ReplaceAllString(text, tablePlaceholder)

	// Replace consecutive newlines (with optional trailing spaces) with a placeholder.
	text = d.regex.newLines.ReplaceAllString(text, newLinesPlaceholder)

	// Replace newline and consecutive spaces with a single space.
	text = d.regex.whiteSpaces.ReplaceAllString(text, " ")

	// Replace the placeholder with consecutive newline characters.
	text = strings.ReplaceAll(text, newLinesPlaceholder, newLines)

	// Replace the table of contents placeholder with the cleaned table of contents.
	text = strings.ReplaceAll(text, tablePlaceholder, table)

	return text
}
