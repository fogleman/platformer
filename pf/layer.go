package pf

import (
	"log"

	"github.com/fogleman/platformer/gfx"
)

type Layer struct {
	sheet       *Sheet
	program     *gfx.Program
	buffer      *gfx.Buffer
	matrixLoc   int
	samplerLoc  int
	positionLoc int
	uvLoc       int
	matrix      gfx.Matrix
}

func NewLayer(sheet *Sheet) *Layer {
	program, err := gfx.NewProgram(layerVertexSource, layerFragmentSource)
	if err != nil {
		log.Fatalln(err)
	}
	buffer := gfx.NewBuffer()
	buffer.Bind()
	tiles := make([]Tile, 17)
	for i := 0; i < 16; i++ {
		tiles[i] = sheet.Tile("GrassMid", i*128, 0)
	}
	tiles[16] = sheet.Tile("AlienBlueStand", 200, 128)
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

func (layer *Layer) Draw() {
	program := layer.program
	program.Use()
	program.SetMatrix(layer.matrixLoc, layer.matrix)
	program.SetInt(layer.samplerLoc, 0)
	program.SetBuffer(layer.positionLoc, 2, 0, 16, layer.buffer)
	program.SetBuffer(layer.uvLoc, 2, 8, 16, layer.buffer)
	program.DrawTriangles(0, 6*17)
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
