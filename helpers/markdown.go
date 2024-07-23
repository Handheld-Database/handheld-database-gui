package helpers

import (
	"regexp"
	"strings"
)

func MarkdownToPlaintext(markdownText string) string {
	// Remove HTML tags
	reHTML := regexp.MustCompile(`<[^>]*>`)
	noHTML := reHTML.ReplaceAllString(markdownText, "")

	// Remove Markdown emphasis (bold and italic)
	// Primeiro, remove negrito (**) e depois itálico (*)
	reBold := regexp.MustCompile(`(?s)\*\*(.*?)\*\*`)
	noBold := reBold.ReplaceAllString(noHTML, `$1`)

	reItalic := regexp.MustCompile(`(?s)_(.*?)_|\*(.*?)\*`)
	noItalic := reItalic.ReplaceAllString(noBold, `$1$2`)

	// Remove Markdown headers, lists, and other elements that are not plain text
	reList := regexp.MustCompile(`(?m)^\s*[-*+]\s+`)
	noLists := reList.ReplaceAllString(noItalic, "")

	reHeader := regexp.MustCompile(`(?m)^#+\s+`)
	noHeaders := reHeader.ReplaceAllString(noLists, "")

	// Remove additional Markdown elements if needed
	reImage := regexp.MustCompile(`!\[.*?\]\(.*?\)`)
	noImages := reImage.ReplaceAllString(noHeaders, "") // Images

	reLink := regexp.MustCompile(`\[.*?\]\(.*?\)`)
	noLinks := reLink.ReplaceAllString(noImages, "") // Links

	// Ensure no extra newlines are introduced
	lines := strings.Split(noLinks, "\n")
	var result []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}
	plainText := strings.Join(result, "\n")

	return plainText
}
