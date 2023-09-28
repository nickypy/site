package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const DEFAULT_INPUT_PATH = "assets"
const DEFAULT_OUTPUT_PATH = "dist"

var InputPath string
var OutputPath string

var rootCmd = &cobra.Command{
	Use:   "site",
	Short: "Utilities for building and serving nickypy-site.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

func Execute() {
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(searchCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	flags := rootCmd.Flags()
	flags.StringVarP(&InputPath, "input", "i", DEFAULT_INPUT_PATH, "Source directory for markdown")
	flags.StringVarP(&OutputPath, "output", "o", DEFAULT_OUTPUT_PATH, "Output directory for rendered assets")
}
