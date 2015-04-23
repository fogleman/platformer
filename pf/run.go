package pf

import (
	"log"
	"math"
	"runtime"

	"github.com/fogleman/platformer/gfx"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const (
	width  = 1280
	height = 720
	title  = "Platformer"
)

func init() {
	runtime.LockOSThread()
}

func Run() {
	// initialize glfw
	if err := glfw.Init(); err != nil {
		log.Fatalln(err)
	}
	defer glfw.Terminate()

	// create window
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		log.Fatalln(err)
	}
	window.MakeContextCurrent()

	// initialize glew
	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}

	// gl setup
	gl.Enable(gl.TEXTURE_2D)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)

	sheet, err := gfx.NewSheetFromFile(0, "textures/sprites.png", "textures/sprites.csv")
	if err != nil {
		log.Fatalln(err)
	}

	backgroundLayer := gfx.NewLayer(sheet)
	spriteLayer := gfx.NewLayer(sheet)

	var tiles []gfx.Tile
	for i := 0; i < 32; i++ {
		tiles = append(tiles, sheet.Tile("GrassMid", i*128, 0))
	}
	backgroundLayer.SetTiles(tiles)

	var sprites []*gfx.Sprite
	sprite := sheet.Sprite("AlienGreenStand")
	sprite.SetAnchor(0.5, 0)
	sprite.SetPosition(0, 128)
	sprites = append(sprites, sprite)
	spriteLayer.SetSprites(sprites)

	// run loop
	for !window.ShouldClose() {
		t := glfw.GetTime()
		w, h := window.GetFramebufferSize()
		matrix := gfx.Orthographic(0, float64(w), 0, float64(h), -1, 1)
		sprite.X = math.Mod(t*100, float64(w))
		spriteLayer.SetSprites(sprites)
		gl.ClearColor(0.78, 0.95, 0.96, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		backgroundLayer.SetMatrix(matrix)
		backgroundLayer.Draw()
		spriteLayer.SetMatrix(matrix)
		spriteLayer.Draw()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
