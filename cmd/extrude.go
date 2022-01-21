package cmd

import (
	"fmt"
	"image"
	"image/draw"
	"os"

	"github.com/spf13/cobra"

	i "github.com/davidwarshaw/tiletool/cmd/internal"
)

var extrudeCmd *cobra.Command
var thickness int

func extrudeTile(tileImage *image.NRGBA, thickness int) (extruded *image.NRGBA) {
	extrudedRect := tileImage.Bounds().Inset(-thickness)
	extruded = image.NewNRGBA(extrudedRect)
	draw.Draw(extruded, tileImage.Bounds(), tileImage, tileImage.Bounds().Min, draw.Src)

	// Horizontal
	for x := tileImage.Bounds().Min.X; x < tileImage.Bounds().Max.X; x++ {
		// Top
		for y := extruded.Bounds().Min.Y; y < tileImage.Bounds().Min.Y; y++ {
			extruded.Set(x, y, tileImage.At(x, tileImage.Bounds().Min.Y))
		}
		// Bottom
		for y := extruded.Bounds().Max.Y - 1; y > tileImage.Bounds().Max.Y-1; y-- {
			extruded.Set(x, y, tileImage.At(x, tileImage.Bounds().Max.Y-1))
		}
	}
	// Vertical
	for y := tileImage.Bounds().Min.Y; y < tileImage.Bounds().Max.Y; y++ {
		// Left
		for x := extruded.Bounds().Min.X; x < tileImage.Bounds().Min.X; x++ {
			extruded.Set(x, y, tileImage.At(tileImage.Bounds().Min.X, y))
		}
		// Right
		for x := extruded.Bounds().Max.X - 1; x > tileImage.Bounds().Max.X-1; x-- {
			extruded.Set(x, y, tileImage.At(tileImage.Bounds().Max.X-1, y))
		}
	}

	// The four corners
	// Left
	for x := extruded.Bounds().Min.X; x < extruded.Bounds().Min.X+thickness; x++ {
		// Top
		for y := extruded.Bounds().Min.Y; y < extruded.Bounds().Min.Y+thickness; y++ {
			extruded.Set(x, y, tileImage.At(tileImage.Bounds().Min.X, tileImage.Bounds().Min.Y))
		}
		// Bottom
		for y := extruded.Bounds().Max.Y - 1; y > extruded.Bounds().Max.Y-thickness-1; y-- {
			extruded.Set(x, y, tileImage.At(tileImage.Bounds().Min.X, tileImage.Bounds().Max.Y-1))
		}
	}
	// Right
	for x := extruded.Bounds().Max.X - 1; x > extruded.Bounds().Max.X-thickness-1; x-- {
		// Top
		for y := extruded.Bounds().Min.Y; y < extruded.Bounds().Min.Y+thickness; y++ {
			extruded.Set(x, y, tileImage.At(tileImage.Bounds().Min.X, tileImage.Bounds().Min.Y))
		}
		// Bottom
		for y := extruded.Bounds().Max.Y - 1; y > extruded.Bounds().Max.Y-thickness-1; y-- {
			extruded.Set(x, y, tileImage.At(tileImage.Bounds().Min.X, tileImage.Bounds().Max.Y-1))
		}
	}

	return
}

func init() {

	extrudeCmd = &cobra.Command{
		Use:   "extrude <filename>",
		Short: "Extrude the tiles of a tileset.",
		Long:  "The extrude command copies tile content into the margin around, and spacing between, tiles. Extrusion mitigates texture bleeding or tearing during tileset map scrolling. The extrude command will increase the tileset margin by the amount of extrusion thickness and increase the tileset spacing by twice the extrusion thickness.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				fmt.Fprintln(os.Stderr, "One arg required: <filename>")
				fmt.Fprintln(os.Stderr, "Use \"tiletool extrude --help\" for more information.")
				os.Exit(1)
			}
			return nil
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := i.ValidatePixelValue(thickness); err != nil {
				fmt.Fprintf(os.Stderr, "Invalid thickness: %s\n", err.Error())
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			filename := args[0]

			img := i.Open(filename, Verbose)
			tc.ReadImage(img)

			outTc := tc
			outTc.TileWidth += 2 * thickness
			outTc.TileHeight += 2 * thickness

			etis := make([]*image.NRGBA, len(outTc.TileImages))
			for i, tileImage := range outTc.TileImages {
				extruded := extrudeTile(tileImage, thickness)
				etis[i] = extruded
			}
			outTc.TileImages = etis

			if Verbose {
				fmt.Printf("Extruding with thickness: %d\n", thickness)
			}

			fmt.Printf("Extruded tileset has margin: %d and spacing: %d\n", thickness, 2*thickness)

			tilesetImage := outTc.ToImage()
			i.Save(tilesetImage, Output, Verbose)
		},
	}
	extrudeCmd.Flags().IntVar(&thickness, "thickness", 1, "extrusion thickness in pixels (default 1)")
}
