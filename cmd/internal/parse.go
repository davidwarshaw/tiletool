package internal

import (
	"image"
	"image/color"
)

type ParseConfig struct {
	TileWidth  int
	TileHeight int
	XOffset    int
	YOffset    int
}

func (ps ParseConfig) NewTilesetConfigFromParseConfig() TilesetConfig {
	return TilesetConfig{
		TileWidth:  ps.TileWidth,
		TileHeight: ps.TileHeight,
		Margin:     0,
		Spacing:    0,
		Columns:    10,
		Color:      color.Transparent,
		TileImages: []*image.NRGBA{},
	}
}

func (ps ParseConfig) CropTiles(img *image.NRGBA) []*image.NRGBA {

	columns := (img.Bounds().Dx() - int(ps.XOffset)) / ps.TileWidth
	rows := (img.Bounds().Dy() - int(ps.YOffset)) / ps.TileHeight

	crops := []*image.NRGBA{}
	for column := 0; column < columns; column++ {
		for row := 0; row < rows; row++ {
			x := (column * ps.TileWidth) + ps.XOffset
			y := (row * ps.TileHeight) + ps.YOffset
			min := image.Point{x, y}
			max := image.Point{x + ps.TileWidth, y + ps.TileHeight}
			rectangle := image.Rectangle{min, max}

			crop := img.SubImage(rectangle).(*image.NRGBA)
			crops = append(crops, crop)
		}
	}

	return crops
}

type FrequencyTile struct {
	Hash            string
	Image           *image.NRGBA
	Count           int
	FirstLocation   image.Point
	Transformations bool
}
