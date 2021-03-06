package gfx

import "log"

type Layer struct {
	sheet       *Sheet
	program     *Program
	buffer      *Buffer
	count       int
	matrixLoc   int
	samplerLoc  int
	positionLoc int
	uvLoc       int
	matrix      Matrix
}

func NewLayer(sheet *Sheet) *Layer {
	program, err := NewProgram(layerVertexSource, layerFragmentSource)
	if err != nil {
		log.Fatalln(err)
	}
	layer := Layer{}
	layer.sheet = sheet
	layer.program = program
	layer.buffer = NewBuffer()
	layer.count = 0
	layer.matrixLoc = program.UniformLocation("matrix")
	layer.samplerLoc = program.UniformLocation("sampler")
	layer.positionLoc = program.AttributeLocation("position")
	layer.uvLoc = program.AttributeLocation("uv")
	return &layer
}

func (layer *Layer) SetTiles(tiles []Tile) {
	layer.buffer.SetItems(tiles)
	layer.count = len(tiles)
}

func (layer *Layer) SetSprites(sprites []*Sprite) {
	tiles := make([]Tile, len(sprites))
	for i, sprite := range sprites {
		tiles[i] = sprite.Tile()
	}
	layer.SetTiles(tiles)
}

func (layer *Layer) SetMatrix(matrix Matrix) {
	layer.matrix = matrix
}

func (layer *Layer) Draw() {
	if layer.count == 0 {
		return
	}
	program := layer.program
	program.Use()
	program.SetMatrix(layer.matrixLoc, layer.matrix)
	program.SetInt(layer.samplerLoc, int32(layer.sheet.Texture.Unit))
	program.SetBuffer(layer.positionLoc, 2, 0, 16, layer.buffer)
	program.SetBuffer(layer.uvLoc, 2, 8, 16, layer.buffer)
	program.DrawTriangles(0, layer.count*6)
}

const layerVertexSource = `
#version 120

uniform mat4 matrix;

attribute vec4 position;
attribute vec2 uv;

varying vec2 fragment_uv;

void main() {
    gl_Position = matrix * position;
    fragment_uv = uv;
}
`

const layerFragmentSource = `
#version 120

uniform sampler2D sampler;

varying vec2 fragment_uv;

void main() {
    gl_FragColor = texture2D(sampler, fragment_uv);
}
`
