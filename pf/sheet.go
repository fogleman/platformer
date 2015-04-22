package pf

import (
	"image"

	"github.com/go-gl/gl/v2.1/gl"
)

type TileVertex struct {
	x, y, u, v float32
}

type Tile struct {
	vertices [6]TileVertex
}

type SheetItem struct {
	x, y, w, h int
}

type Sheet struct {
	w, h    int
	texture uint32
	items   map[string]SheetItem
}

func NewSheet(pngPath, csvPath string) (*Sheet, error) {
	// load csv
	rows, err := LoadCSV(csvPath)
	if err != nil {
		return nil, err
	}
	items := make(map[string]SheetItem)
	for _, row := range rows {
		x := ParseInts(row[1:])
		items[row[0]] = SheetItem{x[0], x[1], x[2], x[3]}
	}

	// load png
	im, err := LoadPNG(pngPath)
	if err != nil {
		return nil, err
	}
	texture := createTexture()
	setTexture(ImageToRGBA(im))

	size := im.Bounds().Size()
	return &Sheet{size.X, size.Y, texture, items}, nil
}

// func (sheet *Sheet) Bind() {
// 	gl.ActiveTexture(gl.TEXTURE0)
// 	gl.BindTexture(gl.TEXTURE_2D, sheet.texture)
// }

func (sheet *Sheet) Tile(name string, x, y int) Tile {
	item := sheet.items[name]
	x0 := float32(x)
	y0 := float32(y)
	x1 := float32(x + item.w)
	y1 := float32(y + item.h)
	u0 := float32(item.x) / float32(sheet.w)
	v1 := float32(item.y) / float32(sheet.h)
	u1 := float32(item.x+item.w) / float32(sheet.w)
	v0 := float32(item.y+item.h) / float32(sheet.h)
	return Tile{[6]TileVertex{
		TileVertex{x0, y0, u0, v0},
		TileVertex{x1, y0, u1, v0},
		TileVertex{x0, y1, u0, v1},
		TileVertex{x1, y0, u1, v0},
		TileVertex{x1, y1, u1, v1},
		TileVertex{x0, y1, u0, v1},
	}}
}

func createTexture() uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	// gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	return texture
}

func setTexture(im *image.RGBA) {
	size := im.Rect.Size()
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA, int32(size.X), int32(size.Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(im.Pix))
}
