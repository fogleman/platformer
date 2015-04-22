package platformer

import (
	"log"

	"github.com/fogleman/platformer/gg"
	"github.com/go-gl/gl/v2.1/gl"
)

type Layer struct {
	sheet            *Sheet
	program          *gg.Program
	buffer           *gg.Buffer
	matrixLocation   int32
	samplerLocation  int32
	positionLocation uint32
	uvLocation       uint32
	matrix           gg.Matrix
}

func NewLayer(sheet *Sheet) *Layer {
	program, err := gg.NewProgram("shaders/vertex.glsl", "shaders/fragment.glsl")
	if err != nil {
		log.Fatalln(err)
	}
	buffer := gg.NewBuffer()
	buffer.Bind()
	tiles := make([]Tile, 17)
	for i := 0; i < 16; i++ {
		tiles[i] = sheet.Tile("GrassMid", i*128, 0)
	}
	tiles[16] = sheet.Tile("AlienBlueStand", 100, 128)
	buffer.SetItems(tiles)
	layer := Layer{}
	layer.sheet = sheet
	layer.program = program
	layer.buffer = buffer
	layer.matrixLocation = program.UniformLocation("matrix")
	layer.samplerLocation = program.UniformLocation("sampler")
	layer.positionLocation = uint32(program.AttributeLocation("position"))
	layer.uvLocation = uint32(program.AttributeLocation("uv"))
	layer.matrix = gg.Orthographic(0, 1280, 0, 720, -1, 1)
	return &layer
}

func (layer *Layer) Draw() {
	layer.program.Use()
	layer.program.UniformMatrix(layer.matrixLocation, layer.matrix)
	layer.program.UniformInt(layer.samplerLocation, 0)
	layer.buffer.Bind()
	gl.EnableVertexAttribArray(layer.positionLocation)
	gl.EnableVertexAttribArray(layer.uvLocation)
	gl.VertexAttribPointer(layer.positionLocation, 2, gl.FLOAT, false, 16, gl.PtrOffset(0))
	gl.VertexAttribPointer(layer.uvLocation, 2, gl.FLOAT, false, 16, gl.PtrOffset(8))
	gl.DrawArrays(gl.TRIANGLES, 0, 6*17)
	gl.DisableVertexAttribArray(layer.positionLocation)
	gl.DisableVertexAttribArray(layer.uvLocation)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
