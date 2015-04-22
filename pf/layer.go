package pf

import (
	"log"

	"github.com/fogleman/platformer/gfx"
)

type Layer struct {
	sheet            *Sheet
	program          *gfx.Program
	buffer           *gfx.Buffer
	matrixLocation   int
	samplerLocation  int
	positionLocation int
	uvLocation       int
	matrix           gfx.Matrix
}

func NewLayer(sheet *Sheet) *Layer {
	program, err := gfx.NewProgram("shaders/vertex.glsl", "shaders/fragment.glsl")
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
	layer.matrixLocation = program.UniformLocation("matrix")
	layer.samplerLocation = program.UniformLocation("sampler")
	layer.positionLocation = program.AttributeLocation("position")
	layer.uvLocation = program.AttributeLocation("uv")
	return &layer
}

func (layer *Layer) Draw() {
	program := layer.program
	program.Use()
	program.UniformMatrix(layer.matrixLocation, layer.matrix)
	program.UniformInt(layer.samplerLocation, 0)
	program.AttributeBuffer(layer.buffer, layer.positionLocation, 2, 0, 16)
	program.AttributeBuffer(layer.buffer, layer.uvLocation, 2, 8, 16)
	program.Draw(6*17, layer.positionLocation, layer.uvLocation)
}
