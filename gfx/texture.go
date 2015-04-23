package gfx

import (
	"image"

	"github.com/go-gl/gl/v2.1/gl"
)

type Texture struct {
	Handle uint32
	Unit   int
	Width  int
	Height int
}

func NewTexture(unit int, im *image.RGBA) *Texture {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0 + uint32(unit))
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	size := im.Rect.Size()
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA, int32(size.X), int32(size.Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(im.Pix))
	return &Texture{texture, unit, size.X, size.Y}
}

func NewTextureFromFile(unit int, path string) (*Texture, error) {
	im, err := LoadPNG(path)
	if err != nil {
		return nil, err
	}
	return NewTexture(unit, ImageToRGBA(im)), nil
}

func (t *Texture) Delete() {
	gl.DeleteTextures(1, &t.Handle)
}

func (t *Texture) Bind() {
	gl.ActiveTexture(gl.TEXTURE0 + uint32(t.Unit))
	gl.BindTexture(gl.TEXTURE_2D, t.Handle)
}
