package pf

import (
	"log"
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

	sheet, err := NewSheet("textures/sprites.png", "textures/sprites.csv")
	if err != nil {
		log.Fatalln(err)
	}

	layer := NewLayer(sheet)

	// run loop
	for !window.ShouldClose() {
		gl.ClearColor(0.78, 0.95, 0.96, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		w, h := window.GetFramebufferSize()
		layer.matrix = gfx.Orthographic(0, float64(w), 0, float64(h), -1, 1)
		layer.Draw()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
