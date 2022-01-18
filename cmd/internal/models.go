package internal

import (
	"image"
	"image/color"
)

func NewTilesetConfig() TilesetConfig {
	return TilesetConfig{
		TileWidth:  16,
		TileHeight: 16,
		Margin:     0,
		Spacing:    0,
		Columns:    10,
		Color:      color.Transparent,
		TileImages: []*image.NRGBA{},
	}
}

type TilesetConfig struct {
	TileWidth  int
	TileHeight int
	Margin     int
	Spacing    int
	Columns    int
	Color      color.Color
	TileImages []*image.NRGBA
}

func NewTilesetConfigFromParseConfig(parseConfig ParseConfig) TilesetConfig {
	return TilesetConfig{
		TileWidth:  parseConfig.TileWidth,
		TileHeight: parseConfig.TileHeight,
		Margin:     0,
		Spacing:    0,
		Columns:    10,
		Color:      color.Transparent,
		TileImages: []*image.NRGBA{},
	}
}

type ParseConfig struct {
	TileWidth  int
	TileHeight int
	XOffset    int
	YOffset    int
}

type FrequencyTile struct {
	Hash            string
	Image           *image.NRGBA
	Count           int
	FirstLocation   image.Point
	Transformations bool
}
