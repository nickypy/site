package cmd

import (
	"path"

	"github.com/nickypy/site/buko"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search from the index, generating one as needed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		searchPath := path.Join(InputPath, "blog")
		expectedPath := path.Join(OutputPath, "search_index.json")

		server.Search(args[0])
	},
}
