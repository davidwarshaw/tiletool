package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Verbose bool
var Output string

const ValidOutputExtensionsMessage = "Valid extensions are: \"jpg\" (or \"jpeg\"), \"png\", \"gif\", \"tif\" (or \"tiff\") and \"bmp\"."

var rootCmd = &cobra.Command{
	Use:               "tiletool",
	Short:             "Command line interface utility for tilesets",
	CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if Verbose {
			fmt.Printf("Verbose output\n")
			fmt.Printf("Outputting to %s\n", Output)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

func init() {
	cobra.OnInitialize()

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&Output, "output", "o", "tileset.png", fmt.Sprintf("file name and format to output to. %s", ValidOutputExtensionsMessage))

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(parseCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
