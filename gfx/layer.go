package gfx

import "log"

type Layer struct {
	sheet       *Sheet
	program     *Program
	buffer      *Buffer
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
	buffer := NewBuffer()
	tiles := make([]Tile, 33)
	for i := 0; i < 32; i++ {
		tiles[i] = sheet.Tile("GrassMid", i*128, 0)
	}
	tiles[32] = sheet.Tile("AlienGreenStand", 200, 128)
	buffer.SetItems(tiles)
	layer := Layer{}
	layer.sheet = sheet
	layer.program = program
	layer.buffer = buffer
	layer.matrixLoc = program.UniformLocation("matrix")
	layer.samplerLoc = program.UniformLocation("sampler")
	layer.positionLoc = program.AttributeLocation("position")
	layer.uvLoc = program.AttributeLocation("uv")
	return &layer
}

func (layer *Layer) SetMatrix(matrix Matrix) {
	layer.matrix = matrix
}

func (layer *Layer) Draw() {
	program := layer.program
	program.Use()
	program.SetMatrix(layer.matrixLoc, layer.matrix)
	program.SetInt(layer.samplerLoc, 0)
	program.SetBuffer(layer.positionLoc, 2, 0, 16, layer.buffer)
	program.SetBuffer(layer.uvLoc, 2, 8, 16, layer.buffer)
	program.DrawTriangles(0, 6*33)
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
