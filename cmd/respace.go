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
		Long:  "The respace command outputs the tileset with the background replaced and the specified margin and spacing. NOTE: this command replaces the background, so it will remove tile extrusions.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				fmt.Fprintln(os.Stderr, "One arg required: <filename>")
				fmt.Fprintln(os.Stderr, "Use \"tiletool respace --help\" for more information.")
				os.Exit(1)
			}
			return nil
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := i.ValidatePixelValue(outMargin); err != nil {
				fmt.Fprintf(os.Stderr, "Invalid out-margin: %s\n", err.Error())
				os.Exit(1)
			}
			if err := i.ValidatePixelValue(outSpacing); err != nil {
				fmt.Fprintf(os.Stderr, "Invalid out-spacing: %s\n", err.Error())
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			filename := args[0]

			img := i.Open(filename, Verbose)
			tc.ReadImage(img)

			outTc := tc
			outTc.Margin = outMargin
			outTc.Spacing = outSpacing
			outTc.Color = BgColor

			if Verbose {
				fmt.Printf("Margin: reading: %d writing: %d\n", tc.Margin, outTc.Margin)
				fmt.Printf("Spacing: reading: %d writing: %d\n", tc.Spacing, outTc.Spacing)
				fmt.Printf("Background Color: writing: %s\n", BgColorHex)
			}

			tilesetImage := outTc.ToImage()
			i.Save(tilesetImage, Output, Verbose)
		},
	}
	respaceCmd.Flags().IntVar(&outMargin, "out-margin", 0, "output tileset margin in pixels (default 0)")
	respaceCmd.Flags().IntVar(&outSpacing, "out-spacing", 0, "output tile spacing in pixels (defualt 0)")
}
