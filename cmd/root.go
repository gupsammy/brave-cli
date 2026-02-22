package cmd

import (
	"github.com/spf13/cobra"
	"github.com/samarthgupta/brave-cli/internal/output"
)

var rootCmd = &cobra.Command{
	Use:     "brave-cli",
	Short:   "Search the web via Brave Search API",
	Long:    "brave-cli: search the web and extract page content via the Brave Search API.\nNo browser or Node.js required. Requires BRAVE_API_KEY.",
	Version: "1.0.0",
}

// Execute is the entry point called from main.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&output.NoColor, "no-color", false, "Disable ANSI color output")
}
