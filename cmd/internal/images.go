package internal

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"strings"

	"github.com/disintegration/imaging"
)

const ValidOutputExtensionsMessage = "Valid extensions are: \"jpg\" (or \"jpeg\"), \"png\", \"gif\", \"tif\" (or \"tiff\") and \"bmp\"."

func ColorFromHex(hex string) (c color.RGBA, err error) {
	_, err = fmt.Sscanf(hex, "#%02x%02x%02x%02x", &c.R, &c.G, &c.B, &c.A)
	return
}
func HexFromColor(c color.Color) (hex string) {
	r, g, b, a := c.RGBA()
	hex = fmt.Sprintf("#%02x%02x%02x%02x", r, g, b, a)
	return
}

func Open(filename string, verbose bool) *image.NRGBA {
	if verbose {
		fmt.Printf("Opening %s\n", filename)
	}
	rawImg, err := imaging.Open(filename, imaging.AutoOrientation(true))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %s\n", err.Error())
		os.Exit(1)
	}

	img := ImageToNRGBA(rawImg)
	return img
}

func Save(tilesetImage *image.NRGBA, filename string, verbose bool) {
	if verbose {
		fmt.Printf("Saving to %s\n", filename)
	}
	err := imaging.Save(tilesetImage, filename)
	if err != nil {
		if strings.Contains(err.Error(), "unsupported image format") {
			fmt.Fprintf(os.Stderr, "Error: the tileset could not be saved because the output extension is invalid. %s\n", ValidOutputExtensionsMessage)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Error saving file: %s\n", err.Error())
		os.Exit(1)
	}
}

func GetContigousSubPixels(img *image.NRGBA) []byte {
	rect := img.Bounds()
	nrgba := image.NewNRGBA(rect)
	draw.Draw(nrgba, rect, img, rect.Min, draw.Src)
	return nrgba.Pix
}

func ImageToNRGBA(src image.Image) *image.NRGBA {
	if dst, ok := src.(*image.NRGBA); ok {
		return dst
	}
	b := src.Bounds()
	dst := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	return dst
}
