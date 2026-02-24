package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/gupsammy/brave-cli/internal/brave"
	"github.com/gupsammy/brave-cli/internal/config"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Verify BRAVE_API_KEY is set and valid",
	Long:  `Makes a minimal test search to confirm your API key is accepted.`,
	Args:         cobra.NoArgs,
	RunE:         runAuth,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(authCmd)
}

func runAuth(_ *cobra.Command, _ []string) error {
	apiKey, err := config.LoadAPIKey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\nGet a key at: https://api-dashboard.search.brave.com/\nThen add to ~/.secrets:  export BRAVE_API_KEY=your_key\n", err)
		os.Exit(3)
	}

	client := brave.NewClient(apiKey)
	resp, err := client.Search("test", brave.SearchOpts{Count: 1})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Auth failed: %v\n", err)
		os.Exit(3)
	}

	fmt.Printf("OK — key accepted. Test query returned %d result(s).\n", len(resp.Items))
	if len(resp.Items) > 0 {
		fmt.Printf("First result: %s\n", resp.Items[0].URL)
	}
	os.Exit(0)
	return nil
}
