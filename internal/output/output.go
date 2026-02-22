package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/samarthgupta/brave-cli/internal/brave"
)

// NoColor disables ANSI escape codes when set to true.
var NoColor bool

func color(code, s string) string {
	if NoColor || os.Getenv("NO_COLOR") != "" {
		return s
	}
	return "\033[" + code + "m" + s + "\033[0m"
}

func bold(s string) string   { return color("1", s) }
func dim(s string) string    { return color("2", s) }
func cyan(s string) string   { return color("36", s) }
func yellow(s string) string { return color("33", s) }

// PrintCompact writes human-readable search results to stdout.
func PrintCompact(items []brave.SearchResult) {
	for i, r := range items {
		num := fmt.Sprintf("[%d]", i+1)
		fmt.Printf("%s %s\n", yellow(num), bold(r.Title))
		fmt.Printf("    %s\n", cyan(r.URL))
		if r.Description != "" {
			fmt.Printf("    %s\n", dim(r.Description))
		}
		if len(r.ExtraSnippets) > 0 {
			for _, s := range r.ExtraSnippets {
				fmt.Printf("    › %s\n", dim(s))
			}
		}
		if r.Content != "" {
			preview := r.Content
			if len(preview) > 400 {
				preview = preview[:400] + "…"
			}
			// indent content preview
			lines := strings.Split(preview, "\n")
			for _, l := range lines {
				if strings.TrimSpace(l) != "" {
					fmt.Printf("    | %s\n", l)
				}
			}
		}
		if i < len(items)-1 {
			fmt.Println()
		}
	}
}

// PrintJSON writes a SearchResponse as JSON to stdout.
func PrintJSON(resp *brave.SearchResponse) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(resp)
}

// PrintFetchText writes extracted page content to stdout.
func PrintFetchText(content string) {
	fmt.Println(content)
}

// PrintFetchJSON writes {url, content} as JSON to stdout.
func PrintFetchJSON(pageURL, content string) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(map[string]string{
		"url":     pageURL,
		"content": content,
	})
}
