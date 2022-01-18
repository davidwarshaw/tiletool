package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	i "github.com/davidwarshaw/tiletool/cmd/internal"
)

var respaceCmd *cobra.Command

var outMargin int
var outSpacing int

func init() {

	respaceCmd = &cobra.Command{
		Use:   "respace <filename>",
		Short: "Respace a tileset.",
		Long:  "The respace command outputs the tileset with specified margin and spacing.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				fmt.Fprintln(os.Stderr, "One arg required: <filename>")
				fmt.Fprintln(os.Stderr, "Use \"tiletool respace --help\" for more information.")
				os.Exit(1)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			filename := args[0]

			readTc := i.NewTilesetConfig()
			readTc.TileWidth = tileSize
			readTc.TileHeight = tileSize

			img := i.Open(filename, Verbose)
			i.ReadTileset(img, &readTc)

			writeTc := readTc
			writeTc.Margin = outMargin
			writeTc.Spacing = outSpacing
			writeTc.Color = BgColor

			if Verbose {
				fmt.Printf("Margin: reading: %d writing: %d\n", readTc.Margin, writeTc.Margin)
				fmt.Printf("Spacing: reading: %d writing: %d\n", readTc.Spacing, writeTc.Spacing)
				fmt.Printf("Background Color: writing: %s\n", BgColorHex)
			}

			tilesetImage := i.WriteTileset(&writeTc)
			i.Save(tilesetImage, Output, Verbose)
		},
	}
	respaceCmd.Flags().IntVar(&outMargin, "out-margin", 0, "the output tileset margin (default 0)")
	respaceCmd.Flags().IntVar(&outSpacing, "out-spacing", 0, "the output tile spacing (defualt 0)")
}
