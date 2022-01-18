package cmd

import (
	"fmt"
	"image/color"
	"os"

	"github.com/spf13/cobra"

	i "github.com/davidwarshaw/tiletool/cmd/internal"
)

var Verbose bool
var Output string

var tileSize int
var margin int
var spacing int
var BgColorHex string

var BgColor color.Color

var rootCmd = &cobra.Command{
	Use:               "tiletool",
	Short:             "Command line interface utility for tilesets",
	CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if Verbose {
			fmt.Printf("Verbose output\n")
			fmt.Printf("Outputting to %s\n", Output)
		}
		var err error
		BgColor, err = i.ColorFromHex(BgColorHex)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing hex color: %s\n", err.Error())
			os.Exit(1)
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
	rootCmd.PersistentFlags().StringVarP(&Output, "output", "o", "tileset.png", fmt.Sprintf("file name and format to output to. %s", i.ValidOutputExtensionsMessage))
	rootCmd.PersistentFlags().IntVarP(&tileSize, "size", "s", 16, "tile size to parse. Tiles are square")
	rootCmd.PersistentFlags().IntVarP(&margin, "margin", "m", 0, "the tileset margin (default 0)")
	rootCmd.PersistentFlags().IntVarP(&spacing, "spacing", "p", 0, "the tile spacing (default 0)")
	rootCmd.PersistentFlags().StringVarP(&BgColorHex, "color", "c", "#00000000", "the 8 digit hex tileset background color to write (default #00000000 (transparent black))")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(parseCmd)
	rootCmd.AddCommand(respaceCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
