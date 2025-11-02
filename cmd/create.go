package cmd

import (
	"log"
	"os"
	"path"

	"github.com/nickypy/site/buko"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Renders all templates and markdown.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		title := args[0]

		outputPath := path.Join(cwd, "markdown/posts")
		server.CreateNewMarkdownFile(title, outputPath)
	},
}
