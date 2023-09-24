package cmd

import (
	"github.com/nickypy/site/buko"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "experimental search index creation",
	Run: func(cmd *cobra.Command, args []string) {
		server.Search(args[0])
	},
}
