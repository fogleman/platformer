package gg

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
)

type Program struct {
	Handle uint32
}

func NewProgram(vertexPath, fragmentPath string) (*Program, error) {
	vertexSource, err := ioutil.ReadFile(vertexPath)
	if err != nil {
		return nil, err
	}
	fragmentSource, err := ioutil.ReadFile(fragmentPath)
	if err != nil {
		return nil, err
	}
	handle, err := CompileProgram(string(vertexSource), string(fragmentSource))
	if err != nil {
		return nil, err
	}
	return &Program{handle}, nil
}

func (p *Program) Delete() {
	gl.DeleteProgram(p.Handle)
}

func (p *Program) Use() {
	gl.UseProgram(p.Handle)
}

func (p *Program) UniformLocation(name string) int32 {
	return gl.GetUniformLocation(p.Handle, gl.Str(name+"\x00"))
}

func (p *Program) AttributeLocation(name string) int32 {
	return gl.GetAttribLocation(p.Handle, gl.Str(name+"\x00"))
}

func (p *Program) UniformMatrix(location int32, value Matrix) {
	data := value.ColMajor()
	gl.UniformMatrix4fv(location, 1, false, &data[0])
}

func (p *Program) UniformInt(location int32, value int32) {
	gl.Uniform1i(location, value)
}

func (p *Program) UniformFloat(location int32, value float32) {
	gl.Uniform1f(location, value)
}

func CompileShader(shaderType uint32, shaderSource string) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	source := gl.Str(shaderSource + "\x00")
	gl.ShaderSource(shader, 1, &source, nil)
	gl.CompileShader(shader)
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var length int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &length)
		log := strings.Repeat("\x00", int(length+1))
		gl.GetShaderInfoLog(shader, length, nil, gl.Str(log))
		gl.DeleteShader(shader)
		return 0, fmt.Errorf(log)
	}
	return shader, nil
}

func CompileProgram(vertexSource, fragmentSource string) (uint32, error) {
	vs, err := CompileShader(gl.VERTEX_SHADER, vertexSource)
	if err != nil {
		return 0, err
	}
	fs, err := CompileShader(gl.FRAGMENT_SHADER, fragmentSource)
	if err != nil {
		gl.DeleteShader(vs)
		return 0, err
	}
	program := gl.CreateProgram()
	gl.AttachShader(program, vs)
	gl.AttachShader(program, fs)
	gl.LinkProgram(program)
	gl.DeleteShader(vs)
	gl.DeleteShader(fs)
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var length int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &length)
		log := strings.Repeat("\x00", int(length+1))
		gl.GetProgramInfoLog(program, length, nil, gl.Str(log))
		gl.DeleteProgram(program)
		return 0, fmt.Errorf(log)
	}
	return program, nil
}
