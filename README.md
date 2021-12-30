# tiletool

Tiletool is a command line interface utility for working with tilesets.

## Commands

### Parse

The parse command processes an image and identifies the set of unique tiles that compose it, which are then output as a tileset. Verbose output will list a frequency count for all tiles, their first location in the image and whether it was necessary to transform them by flipping or rotation.

Usage:

```
    tiletool parse <filename> [flags]
```

Flags:

```
    -h, --help              help for parse
    -s, --size uint16       tile size to parse. Tiles are square (default 16)
    -t, --transform         allow tiles to be flipped and rotated (default false)
    -x, --x-offset uint16   start at this x coordinate (default 0)
    -y, --y-offset uint16   start at this y coordinate (default 0)
```

## Global Flags

```
    -h, --help            help for tiletool
    -o, --output string   file name and format to output to. Valid extensions are: "jpg" (or "jpeg"), "png", "gif", "tif" (or "tiff"), and "bmp". (default "tileset.png")
    -v, --verbose         verbose output
```
