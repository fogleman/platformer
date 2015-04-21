package platformer

import (
	"fmt"
	"log"
	"runtime"

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
	gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)

	sheet, err := NewSpriteSheet("textures/sprites.png", "textures/sprites.csv")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(sheet.texture)

	// run loop
	for !window.ShouldClose() {
		gl.ClearColor(0.78, 0.95, 0.96, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		drawBuffer(window)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func drawBuffer(window *glfw.Window) {
	w, h := window.GetFramebufferSize()
	var x, y float32 = (float32(w) - 2048) / 2, float32(h) - 4096
	gl.PushMatrix()
	gl.Ortho(0, float64(w), 0, float64(h), -1, 1)
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(x, y)
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(x+2048, y)
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(x+2048, y+4096)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(x, y+4096)
	gl.End()
	gl.PopMatrix()
}
