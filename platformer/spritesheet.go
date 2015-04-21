package platformer

import (
	"image"

	"github.com/go-gl/gl/v2.1/gl"
)

type SpriteSheetItem struct {
	x, y, w, h int
}

type SpriteSheet struct {
	texture uint32
	items   map[string]SpriteSheetItem
}

func NewSpriteSheet(pngPath, csvPath string) (*SpriteSheet, error) {
	// load csv
	rows, err := LoadCSV(csvPath)
	if err != nil {
		return nil, err
	}
	items := make(map[string]SpriteSheetItem)
	for _, row := range rows {
		x := ParseInts(row[1:])
		items[row[0]] = SpriteSheetItem{x[0], x[1], x[2], x[3]}
	}

	// load png
	im, err := LoadPNG(pngPath)
	if err != nil {
		return nil, err
	}
	texture := createTexture()
	setTexture(ImageToRGBA(im))

	// create sprite sheet
	sheet := SpriteSheet{}
	sheet.texture = texture
	sheet.items = items
	return &sheet, nil
}

func createTexture() uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
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
