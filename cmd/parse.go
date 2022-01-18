package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"os"
	"sort"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	i "github.com/davidwarshaw/tiletool/cmd/internal"
)

var parseCmd *cobra.Command

// Flags
var xOffset int
var yOffset int

var transform bool

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
	pixels := i.GetContigousSubPixels(nrgba)
	shaBytes := md5.Sum(pixels)
	hash := hex.EncodeToString(shaBytes[:])
	return hash
}

func computeFreq(img *image.NRGBA, crops []*image.NRGBA, transformations []string) []i.FrequencyTile {
	frequencyTiles := []i.FrequencyTile{}
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
			frequencyTile := i.FrequencyTile{
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

func CreateTransformations() []string {
	var transformations []string
	for _, flip := range []string{"flipH", "flipV", "none"} {
		for _, rotation := range []string{"rotate90", "rotate180", "rotate270", "none"} {
			transformations = append(transformations, fmt.Sprintf("%s-%s", flip, rotation))
		}
	}
	// Drop the identity transformation from the array. We'll search it separately
	transformations = transformations[:len(transformations)-1]
	return transformations
}

func parse(img *image.NRGBA, parseConfig i.ParseConfig, transformations []string, verbose bool) ([]*image.NRGBA, []i.FrequencyTile, error) {
	if verbose {
		bounds := img.Bounds()
		imgSize := fmt.Sprintf("%dx%d", bounds.Dx(), bounds.Dy())
		tileSize := fmt.Sprintf("%dx%d", parseConfig.TileWidth, parseConfig.TileHeight)
		leftOverSize := fmt.Sprintf("%dx%d", bounds.Dx()%parseConfig.TileWidth, bounds.Dy()%parseConfig.TileHeight)
		offsetSize := fmt.Sprintf("%dx%d", xOffset, yOffset)
		fmt.Printf("Parsing %s image (offset by %s) for %s tiles with %s remainder\n", imgSize, offsetSize, tileSize, leftOverSize)
	}
	tiles := i.CropTiles(img, parseConfig)

	frequencyTiles := computeFreq(img, tiles, transformations)

	return tiles, frequencyTiles, nil
}

func outputTable(frequencyTiles []i.FrequencyTile) {
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
		Run: func(cmd *cobra.Command, args []string) {
			filename := args[0]

			parseConfig := i.ParseConfig{
				TileWidth:  tileSize,
				TileHeight: tileSize,
				XOffset:    xOffset,
				YOffset:    yOffset,
			}
			var transformations []string
			if transform {
				transformations = CreateTransformations()
			}

			img := i.Open(filename, Verbose)

			tiles, frequencyTiles, err := parse(img, parseConfig, transformations, Verbose)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening file: %s\n", err.Error())
				os.Exit(1)
			}
			if Verbose {
				fmt.Printf("Parsed %d total tiles, %d unique\n", len(tiles), len(frequencyTiles))
				outputTable(frequencyTiles)
			}

			ts := i.NewTilesetConfigFromParseConfig(parseConfig)
			for _, frequencyTile := range frequencyTiles {
				ts.TileImages = append(ts.TileImages, frequencyTile.Image)
			}

			tilesetImage := i.WriteTileset(&ts)
			i.Save(tilesetImage, Output, Verbose)
		},
	}
	parseCmd.Flags().IntVarP(&xOffset, "x-offset", "x", 0, "start at this x coordinate (default 0)")
	parseCmd.Flags().IntVarP(&yOffset, "y-offset", "y", 0, "start at this y coordinate (default 0)")
	parseCmd.Flags().BoolVarP(&transform, "transform", "t", false, "allow tiles to be flipped and rotated (default false)")

}
