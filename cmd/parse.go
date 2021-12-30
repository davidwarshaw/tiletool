package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"sort"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var parseCmd *cobra.Command

// Flags
var xOffset uint16
var yOffset uint16
var tileSize uint16
var transform bool

// Internal vars
var tileWidth = 16
var tileHeight = 16

var tilesetColumns = 10

var transformations = []string{}

type FrequencyTile struct {
	Hash            string
	Image           *image.NRGBA
	Count           int
	FirstLocation   image.Point
	Transformations bool
}

func getContigousSubPixels(img *image.NRGBA) []byte {
	rect := img.Bounds()
	nrgba := image.NewNRGBA(rect)
	draw.Draw(nrgba, rect, img, rect.Min, draw.Src)
	return nrgba.Pix
}

func imageToNRGBA(src image.Image) *image.NRGBA {
	if dst, ok := src.(*image.NRGBA); ok {
		return dst
	}
	b := src.Bounds()
	dst := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	return dst
}

func cropTiles(img *image.NRGBA, width, height int) []*image.NRGBA {

	columns := (img.Bounds().Dx() - int(xOffset)) / width
	rows := (img.Bounds().Dy() - int(yOffset)) / height

	crops := []*image.NRGBA{}
	for column := 0; column < columns; column++ {
		for row := 0; row < rows; row++ {
			x := (column * width) + int(xOffset)
			y := (row * height) + int(yOffset)
			min := image.Point{x, y}
			max := image.Point{x + width, y + height}
			rectangle := image.Rectangle{min, max}

			crop := img.SubImage(rectangle).(*image.NRGBA)
			crops = append(crops, crop)
		}
	}

	return crops
}

func transformCrop(transformType string, crop *image.NRGBA) *image.NRGBA {
	transformTypes := strings.Split(transformType, "-")
	var first, second *image.NRGBA

	// First pass transofrmation
	switch transformTypes[0] {
	case "none":
		{
			first = crop
		}
	case "flipH":
		{
			first = imaging.FlipH(crop)
		}
	case "flipV":
		{
			first = imaging.FlipV(crop)
		}
	}

	// Second pass transofrmation
	switch transformTypes[1] {
	case "none":
		{
			second = first
		}
	case "rotate90":
		{
			second = imaging.Rotate90(first)
		}
	case "rotate180":
		{
			second = imaging.Rotate180(first)
		}
	case "rotate270":
		{
			second = imaging.Rotate270(first)
		}
	}

	return second
}

func hashNrgba(nrgba *image.NRGBA) string {
	pixels := getContigousSubPixels(nrgba)
	shaBytes := md5.Sum(pixels)
	hash := hex.EncodeToString(shaBytes[:])
	return hash
}

func computeFreq(img *image.NRGBA, crops []*image.NRGBA, transformations []string) []FrequencyTile {
	frequencyTiles := []FrequencyTile{}
	lookup := map[string]int{}
	tileIndex := 0
	fmt.Println()
	for _, crop := range crops {
		baseOrientationHash := hashNrgba(crop)

		foundTransformation := false
		for _, transformation := range transformations {
			transformedCrop := transformCrop(transformation, crop)
			hash := hashNrgba(transformedCrop)
			// fmt.Printf("%v: %s: %v\n", pixels, hash, transformedCrop.Bounds().Min)

			// If the hash is the same as the base orientation has, then don't bother
			// searching with it
			if hash == baseOrientationHash {
				continue
			}

			if index, ok := lookup[hash]; ok {
				frequencyTiles[index].Count++
				frequencyTiles[index].Transformations = true
				foundTransformation = true
				break
			}
		}

		// If we've already found this tile with a transformation, skip the base orientation
		if foundTransformation {
			continue
		}
		if index, ok := lookup[baseOrientationHash]; ok {
			frequencyTiles[index].Count++
		} else {
			frequencyTile := FrequencyTile{
				Hash:            baseOrientationHash,
				Image:           crop,
				Count:           1,
				FirstLocation:   crop.Bounds().Min,
				Transformations: false,
			}
			frequencyTiles = append(frequencyTiles, frequencyTile)
			lookup[baseOrientationHash] = tileIndex
			tileIndex++
		}
	}

	sort.Slice(frequencyTiles, func(i, j int) bool {
		return frequencyTiles[i].Count > frequencyTiles[j].Count
	})

	return frequencyTiles
}

func createTileset(frequencyTiles []FrequencyTile) *image.NRGBA {
	columns := tilesetColumns
	rows := (len(frequencyTiles) / columns) + 1

	width := tilesetColumns * tileWidth
	height := rows * tileHeight

	tileset := imaging.New(width, height, color.Transparent)

	for i, frequencyTile := range frequencyTiles {
		column := i % tilesetColumns
		row := i / tilesetColumns
		position := image.Point{X: column * tileWidth, Y: row * tileHeight}
		opacity := 1.0
		tileset = imaging.Overlay(tileset, frequencyTile.Image, position, opacity)
	}

	return tileset
}

func CreateTransformations() []string {
	for _, flip := range []string{"flipH", "flipV", "none"} {
		for _, rotation := range []string{"rotate90", "rotate180", "rotate270", "none"} {
			transformations = append(transformations, fmt.Sprintf("%s-%s", flip, rotation))
		}
	}
	// Drop the identity transformation from the array. We'll search it separately
	transformations = transformations[:len(transformations)-1]
	return transformations
}

func Parse(filename string, transformations []string, verbose bool) ([]*image.NRGBA, []FrequencyTile, error) {
	img, err := imaging.Open(filename, imaging.AutoOrientation(true))
	if err != nil {
		return nil, nil, err
	}

	nrgba := imageToNRGBA(img)

	if verbose {
		bounds := nrgba.Bounds()
		imgSize := fmt.Sprintf("%dx%d", bounds.Dx(), bounds.Dy())
		tileSize := fmt.Sprintf("%dx%d", tileWidth, tileHeight)
		leftOverSize := fmt.Sprintf("%dx%d", bounds.Dx()%tileWidth, bounds.Dy()%tileHeight)
		offsetSize := fmt.Sprintf("%dx%d", xOffset, yOffset)
		fmt.Printf("Parsing %s: %s image (offset by %s) for %s tiles with %s remainder\n", filename, imgSize, offsetSize, tileSize, leftOverSize)
	}
	tiles := cropTiles(nrgba, tileWidth, tileHeight)

	frequencyTiles := computeFreq(nrgba, tiles, transformations)

	return tiles, frequencyTiles, nil
}

func outputTable(frequencyTiles []FrequencyTile) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	if transform {
		t.AppendHeader(table.Row{"Tileset Index", "Count", "First Location", "Transformation Required"})
	} else {
		t.AppendHeader(table.Row{"Tileset Index", "Count", "First Location"})
	}

	for i, frequencyTile := range frequencyTiles {
		if transform {
			t.AppendRow(table.Row{
				fmt.Sprintf("%d", i),
				fmt.Sprintf("%d", frequencyTile.Count),
				fmt.Sprintf("%v", frequencyTile.FirstLocation),
				fmt.Sprintf("%t", frequencyTile.Transformations),
			})
		} else {
			t.AppendRow(table.Row{
				fmt.Sprintf("%d", i),
				fmt.Sprintf("%d", frequencyTile.Count),
				fmt.Sprintf("%v", frequencyTile.FirstLocation),
			})
		}
	}
	t.Render()
}

func init() {

	parseCmd = &cobra.Command{
		Use:   "parse <filename>",
		Short: "Parse a tileset from an image.",
		Long:  "The parse command processes an image and identifies the set of unique tiles that compose it, which are then output as a tileset. Verbose output will list a frequency count for all tiles, their first location in the image and whether it was necessary to transform them by flipping or rotation.",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				fmt.Fprintln(os.Stderr, "One arg required: <filename>")
				fmt.Fprintln(os.Stderr, "Use \"tiletool parse --help\" for more information.")
				os.Exit(1)
			}
			return nil
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			tileWidth = int(tileSize)
			tileHeight = int(tileSize)
			if transform {
				transformations = CreateTransformations()
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			// verbose, _ := cmd.Flags().GetBool("verbose")
			// output, _ := cmd.Flags().GetString("output")
			filename := args[0]

			tiles, frequencyTiles, err := Parse(filename, transformations, Verbose)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening file: %s\n", err.Error())
				os.Exit(1)
			}
			if Verbose {
				fmt.Printf("Parsed %d total tiles, %d unique\n", len(tiles), len(frequencyTiles))
				outputTable(frequencyTiles)
			}

			tileset := createTileset(frequencyTiles)

			err = imaging.Save(tileset, Output)
			if err != nil {
				if strings.Contains(err.Error(), "unsupported image format") {
					fmt.Fprintf(os.Stderr, "Error: the tileset could not be saved because the output extension is invalid. %s\n", ValidOutputExtensionsMessage)
					os.Exit(1)
				}
				fmt.Fprintf(os.Stderr, "Error saving file: %s\n", err.Error())
				os.Exit(1)
			}
		},
	}
	parseCmd.Flags().Uint16VarP(&xOffset, "x-offset", "x", 0, "start at this x coordinate (default 0)")
	parseCmd.Flags().Uint16VarP(&yOffset, "y-offset", "y", 0, "start at this y coordinate (default 0)")
	parseCmd.Flags().Uint16VarP(&tileSize, "size", "s", 16, "tile size to parse. Tiles are square")
	parseCmd.Flags().BoolVarP(&transform, "transform", "t", false, "allow tiles to be flipped and rotated (default false)")

}
