package internal

import (
	"fmt"
	"image"
	"image/color"

	"github.com/disintegration/imaging"
)

type TilesetConfig struct {
	TileWidth  int
	TileHeight int
	Margin     int
	Spacing    int
	Columns    int
	Color      color.Color
	TileImages []*image.NRGBA
}

func NewDefaultTilesetConfig() TilesetConfig {
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

func NewTilesetConfig(tileSize, margin, spacing int, bgColor color.Color) TilesetConfig {
	return TilesetConfig{
		TileWidth:  tileSize,
		TileHeight: tileSize,
		Margin:     margin,
		Spacing:    spacing,
		Columns:    10,
		Color:      bgColor,
		TileImages: []*image.NRGBA{},
	}
}

func (ts *TilesetConfig) TilePosition(row, column int) (pos image.Point) {
	pos.X = (column * (ts.TileWidth + ts.Spacing)) + ts.Margin
	pos.Y = (row * (ts.TileHeight + ts.Spacing)) + ts.Margin
	return
}

func (ts *TilesetConfig) TileRectangle(row, column int) (rect image.Rectangle) {
	min := ts.TilePosition(row, column)
	rect = image.Rect(min.X, min.Y, min.X+ts.TileWidth, min.Y+ts.TileHeight)
	return
}

func (ts *TilesetConfig) Dims() (width, height int) {
	rows := ts.Rows()
	width = (ts.Columns * (ts.TileWidth + ts.Spacing)) - ts.Spacing + (2 * ts.Margin)
	height = (rows * (ts.TileHeight + ts.Spacing)) - ts.Spacing + (2 * ts.Margin)
	return
}

func (ts *TilesetConfig) Rows() (rows int) {
	extra := 0
	if len(ts.TileImages)%ts.Columns > 0 {
		extra = 1
	}
	rows = (len(ts.TileImages) / ts.Columns) + extra
	return
}

func (ts *TilesetConfig) GetRowsAndCols(img *image.NRGBA) (*int, *int, error) {
	imageWidth := img.Bounds().Dx()
	imageHeight := img.Bounds().Dy()

	tileableWidth := imageWidth + ts.Spacing - (2 * ts.Margin)
	tileableHeight := imageHeight + ts.Spacing - (2 * ts.Margin)

	tileingWidth := ts.TileWidth + ts.Spacing
	tileingHeight := ts.TileHeight + ts.Spacing

	if tileableWidth%tileingWidth != 0 {
		return nil, nil, fmt.Errorf(
			"bad margin, spacing, or tile size for image width: margin: %d, spacing: %d, tile size: %d, image width: %d",
			ts.Margin, ts.Spacing, ts.TileWidth, imageWidth)
	}
	numCols := tileableWidth / tileingWidth

	if tileableHeight%tileingHeight != 0 {
		return nil, nil, fmt.Errorf(
			"bad margin, spacing, or tile size for image height: margin: %d, spacing: %d, tile size: %d, image height: %d",
			ts.Margin, ts.Spacing, ts.TileHeight, imageHeight)
	}
	numRows := tileableHeight / tileingHeight

	return &numRows, &numCols, nil
}

func (ts *TilesetConfig) ReadImage(img *image.NRGBA) error {
	numRows, numCols, err := ts.GetRowsAndCols(img)
	if err != nil {
		return fmt.Errorf("error reading tileset: %s", err.Error())
	}

	ts.Columns = *numCols

	for row := 0; row < *numRows; row++ {
		for column := 0; column < *numCols; column++ {
			rect := ts.TileRectangle(row, column)
			tileImage := img.SubImage(rect).(*image.NRGBA)
			ts.TileImages = append(ts.TileImages, tileImage)
		}
	}

	return nil
}

func (ts *TilesetConfig) ToImage() *image.NRGBA {
	width, height := ts.Dims()
	tilesetImage := imaging.New(width, height, ts.Color)

	for i, tileImage := range ts.TileImages {
		column := i % ts.Columns
		row := i / ts.Columns
		pos := ts.TilePosition(row, column)
		tilesetImage = imaging.Paste(tilesetImage, tileImage, pos)
	}

	return tilesetImage
}
