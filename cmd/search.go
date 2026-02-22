package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/samarthgupta/brave-cli/internal/brave"
	"github.com/samarthgupta/brave-cli/internal/config"
	"github.com/samarthgupta/brave-cli/internal/output"
)

var (
	searchCount         int
	searchFreshness     string
	searchCountry       string
	searchExtraSnippets bool
	searchContent       bool
	searchSafesearch    string
	searchOutput        string
	searchQuiet         bool
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search the web",
	Long: `Search the web using Brave Search API.

Results go to stdout. Progress messages (when --content is used) go to stderr.

Freshness values: day, week, month, year
Output formats:  compact (default), json`,
	Args: cobra.ExactArgs(1),
	RunE: runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)

	f := searchCmd.Flags()
	f.IntVarP(&searchCount, "count", "n", 10, "Number of results (1-100; paginates automatically past 20)")
	f.StringVar(&searchFreshness, "freshness", "", "Date filter: day, week, month, year")
	f.StringVar(&searchCountry, "country", "", "2-char country code (e.g. US, GB)")
	f.BoolVar(&searchExtraSnippets, "extra-snippets", false, "Fetch additional text excerpts per result")
	f.BoolVar(&searchContent, "content", false, "Fetch and extract full page text for each result")
	f.StringVar(&searchSafesearch, "safesearch", "moderate", "Content filter: off, moderate, strict")
	f.StringVarP(&searchOutput, "output", "o", "compact", "Output format: compact or json")
	f.BoolVarP(&searchQuiet, "quiet", "q", false, "Suppress progress lines on stderr")
}

func runSearch(_ *cobra.Command, args []string) error {
	apiKey, err := config.LoadAPIKey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\nGet a key at: https://api-dashboard.search.brave.com/\nThen add to ~/.secrets:  export BRAVE_API_KEY=your_key\n", err)
		os.Exit(3)
	}

	client := brave.NewClient(apiKey)

	opts := brave.SearchOpts{
		Count:         searchCount,
		Freshness:     searchFreshness,
		Country:       searchCountry,
		ExtraSnippets: searchExtraSnippets,
		Safesearch:    searchSafesearch,
	}

	resp, err := client.Search(args[0], opts)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if searchContent {
		for i := range resp.Items {
			if !searchQuiet {
				fmt.Fprintf(os.Stderr, "fetching %s...\n", resp.Items[i].URL)
			}
			text, cerr := client.FetchContent(resp.Items[i].URL)
			if cerr == nil {
				resp.Items[i].Content = text
			}
		}
	}

	switch searchOutput {
	case "json":
		return output.PrintJSON(resp)
	default:
		output.PrintCompact(resp.Items)
	}
	return nil
}
