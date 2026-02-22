package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/samarthgupta/brave-cli/internal/brave"
	"github.com/samarthgupta/brave-cli/internal/config"
	"github.com/samarthgupta/brave-cli/internal/output"
)

var fetchOutput string

var fetchCmd = &cobra.Command{
	Use:   "fetch <url>",
	Short: "Extract readable text from a URL",
	Long: `Fetch a web page and extract clean readable text from the HTML.
Script/style blocks, navigation, and boilerplate are removed.

Output formats: text (default), json`,
	Args:         cobra.ExactArgs(1),
	RunE:         runFetch,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().StringVarP(&fetchOutput, "output", "o", "text", "Output format: text or json")
}

func runFetch(_ *cobra.Command, args []string) error {
	// Fetch does not hit the Brave API — no key needed — but we load it
	// anyway so missing-key errors surface consistently for all subcommands.
	// If the key is absent we proceed (content fetch is key-independent).
	apiKey, _ := config.LoadAPIKey()
	client := brave.NewClient(apiKey)

	pageURL := args[0]
	content, err := client.FetchContent(pageURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching %s: %v\n", pageURL, err)
		os.Exit(1)
	}

	switch fetchOutput {
	case "json":
		return output.PrintFetchJSON(pageURL, content)
	default:
		output.PrintFetchText(content)
	}
	return nil
}
