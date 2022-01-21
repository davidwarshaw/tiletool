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

var tc i.TilesetConfig

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
			fmt.Fprintf(os.Stderr, "Invalid hex color: %s\n", err.Error())
			os.Exit(1)
		}

		if err := i.ValidatePositivePixelValue(tileSize); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid size: %s\n", err.Error())
			os.Exit(1)
		}
		if err := i.ValidatePixelValue(margin); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid margin: %s\n", err.Error())
			os.Exit(1)
		}
		if err := i.ValidatePixelValue(spacing); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid spacing: %s\n", err.Error())
			os.Exit(1)
		}

		tc = i.NewTilesetConfig(tileSize, margin, spacing, BgColor)
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
	rootCmd.PersistentFlags().IntVarP(&tileSize, "size", "s", 16, "input tile size in pixels. Tiles are square")
	rootCmd.PersistentFlags().IntVarP(&margin, "margin", "m", 0, "input tileset margin in pixels (default 0)")
	rootCmd.PersistentFlags().IntVarP(&spacing, "spacing", "p", 0, "input tile spacing in pixels (default 0)")
	rootCmd.PersistentFlags().StringVarP(&BgColorHex, "color", "c", "#00000000", "output tileset background color in 8 digit hex format (RGBA)")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(parseCmd)
	rootCmd.AddCommand(respaceCmd)
	rootCmd.AddCommand(extrudeCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
