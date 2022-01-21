package cmd

import (
	"os"
	"testing"

	"github.com/disintegration/imaging"

	i "github.com/davidwarshaw/tiletool/cmd/internal"
)

type imageTest struct {
	filename        string
	transformations []string
}

func TestImage(t *testing.T) {

	tests := []imageTest{
		{filename: "../fixtures/test_01.png",
			transformations: []string{"flipH-none", "flipV-none", "none-rotate90", "none-rotate180", "none-rotate270"}},
		{filename: "../fixtures/test_02.png",
			transformations: []string{"flipH-none", "flipV-none", "none-rotate90", "none-rotate180", "none-rotate270"}},
		{filename: "../fixtures/test_03.png",
			transformations: []string{"flipH-none", "flipV-none", "none-rotate90", "none-rotate180", "none-rotate270"}},
	}

	for _, tc := range tests {
		t.Run(tc.filename, func(t *testing.T) {
			img, err := imaging.Open(tc.filename, imaging.AutoOrientation(true))
			if err != nil {
				t.Fatalf("Error opening file: %s\n", err.Error())
				os.Exit(1)
			}

			nrgba := i.ImageToNRGBA(img)
			parseConfig := i.ParseConfig{TileWidth: 16, TileHeight: 16, XOffset: 0, YOffset: 0}

			for _, transformation := range tc.transformations {
				crops := parseConfig.CropTiles(nrgba)
				for j, crop := range crops {
					cropHash := hashNrgba(crop)
					transformed := transformCrop(transformation, crop)
					transformedHash := hashNrgba(transformed)

					if cropHash == transformedHash {
						t.Errorf(
							"transformation: %s crop %d: %v :expected %s to be different than %s",
							transformation, j, crop.Bounds().Min, cropHash, transformedHash)
					}
				}
			}
		})
	}
}
