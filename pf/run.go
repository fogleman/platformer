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
	tiles = append(tiles, sheet.Tile("SignRight", 128, 128))
	tiles = append(tiles, sheet.Tile("Bush", 512, 128))
	tiles = append(tiles, sheet.Tile("Cactus", 768, 128))
	tiles = append(tiles, sheet.Tile("LadderMid", 1024, 128))
	tiles = append(tiles, sheet.Tile("LadderMid", 1024, 256))
	tiles = append(tiles, sheet.Tile("LadderMid", 1024, 384))
	tiles = append(tiles, sheet.Tile("LadderTop", 1024, 512))
	tiles = append(tiles, sheet.Tile("LeverLeft", 128, 512))
	// tiles = append(tiles, sheet.Tile("BoxItem", 256, 384))
	for i := 0; i < 32; i++ {
		tiles = append(tiles, sheet.Tile("GrassMid", i*128, 0))
		if i%2 == 1 {
			tiles = append(tiles, sheet.Tile("GrassHalf", i*128, 384))
		}
	}
	backgroundLayer.SetTiles(tiles)

	var sprites []*gfx.Sprite
	sprite := sheet.Sprite("AlienBlueStand")
	sprite.SetAnchor(0.5, 0)
	sprite.SetPosition(0, 128)
	sprites = append(sprites, sprite)
	spriteLayer.SetSprites(sprites)

	names := []string{"AlienBlueWalk1", "AlienBlueWalk2"}

	// run loop
	for !window.ShouldClose() {
		t := glfw.GetTime()
		w, h := window.GetFramebufferSize()
		matrix := gfx.Orthographic(0, float64(w), 0, float64(h), -1, 1)
		sprite.X = float64(w/2) + math.Sin(t)*float64(w/3)
		sprite.FlipX = math.Cos(t+0.25) < 0
		sprite.Name = names[int((math.Sin(t)+1)*16)%len(names)]
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
