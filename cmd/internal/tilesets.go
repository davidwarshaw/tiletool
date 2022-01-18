package internal

import (
	"fmt"
	"image"

	"github.com/disintegration/imaging"
)

func CropTiles(img *image.NRGBA, ps ParseConfig) []*image.NRGBA {

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

func GetRowsAndCols(img *image.NRGBA, ts TilesetConfig) (*int, *int, error) {
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

func ReadTileset(img *image.NRGBA, ts *TilesetConfig) error {
	numRows, numCols, err := GetRowsAndCols(img, *ts)
	if err != nil {
		return fmt.Errorf("error reading tileset: %s", err.Error())
	}

	ts.Columns = *numCols

	for row := 0; row < *numRows; row++ {
		for column := 0; column < *numCols; column++ {
			x, y := xyFromRowColumn(row, column, *ts)
			min := image.Point{x, y}
			max := image.Point{x + ts.TileWidth, y + ts.TileHeight}
			rectangle := image.Rectangle{min, max}
			tileImage := img.SubImage(rectangle).(*image.NRGBA)
			ts.TileImages = append(ts.TileImages, tileImage)
		}
	}

	return nil
}

func WriteTileset(ts *TilesetConfig) *image.NRGBA {
	rows := rows(*ts)

	width, height := widthHeightFromRowsColumns(rows, ts.Columns, *ts)
	tilesetImage := imaging.New(width, height, ts.Color)

	for i, tileImage := range ts.TileImages {
		column := i % ts.Columns
		row := i / ts.Columns
		x, y := xyFromRowColumn(row, column, *ts)
		position := image.Point{x, y}
		opacity := 1.0
		tilesetImage = imaging.Overlay(tilesetImage, tileImage, position, opacity)
	}

	return tilesetImage
}

func xyFromRowColumn(row, column int, ts TilesetConfig) (int, int) {
	x := (column * (ts.TileWidth + ts.Spacing)) + ts.Margin
	y := (row * (ts.TileHeight + ts.Spacing)) + ts.Margin
	return x, y
}

func widthHeightFromRowsColumns(numRows, numColumns int, ts TilesetConfig) (int, int) {
	width := (numColumns * (ts.TileWidth + ts.Spacing)) - ts.Spacing + (2 * ts.Margin)
	height := (numRows * (ts.TileHeight + ts.Spacing)) - ts.Spacing + (2 * ts.Margin)
	return width, height
}

func rows(ts TilesetConfig) (rows int) {
	extra := 0
	if len(ts.TileImages)%ts.Columns > 0 {
		extra = 1
	}
	rows = (len(ts.TileImages) / ts.Columns) + extra
	return
}
