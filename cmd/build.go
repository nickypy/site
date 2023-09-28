package cmd

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/nickypy/site/buko"
	"github.com/spf13/cobra"
)

var ListUnpublished bool

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Renders all templates and markdown.",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		outputPath := path.Join(cwd, OutputPath)

		opts := make([]server.BlogOption, 0)

		if ListUnpublished {
			opts = append(opts, server.RenderUnpublished())
		}

		log.Default().Println("Cleaning out previous build...")
		err = os.RemoveAll(outputPath)
		if err != nil {
			if !os.IsNotExist(err) {
				panic(err)
			}
		}

		log.Default().Println("Building assets...")
		blog := server.NewBlogRenderCache(InputPath, opts...)

		start := time.Now().UnixMilli()
		blog.Render()
		elapsed := time.Now().UnixMilli() - start

		log.Default().Printf("Done. Took %dms.", elapsed)
	},
}

func init() {
	flags := buildCmd.Flags()
	flags.BoolVar(&ListUnpublished, "include-unpublished", false, "List unpublished posts")
}
