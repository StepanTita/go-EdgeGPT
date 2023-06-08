package chat_bot

import (
	"regexp"
	"strings"
)

var rLink = regexp.MustCompile(`(\[\d+\]):\s((https?):\/\/([\w.-]+)(\/[\w\-.&+]*)*(\?[^#\n]*)?)\s("(.+)")\n`)

func ExtractURLs(s string) (string, []ResponseLink) {
	linksList := make([]ResponseLink, 0, 10)

	referencesText := strings.Builder{}

	for _, match := range rLink.FindAllStringSubmatch(s, -1) {
		linksList = append(linksList, ResponseLink{
			ID:    strings.TrimSuffix(strings.TrimPrefix(match[1], "["), "]"),
			URL:   match[2],
			Title: match[7],
		})
		referencesText.WriteString(match[0])
	}

	return referencesText.String(), linksList
}
