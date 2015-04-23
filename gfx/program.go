package gfx

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
)

type Program struct {
	Handle             uint32
	attributeLocations []int
}

func NewProgram(vertexSource, fragmentSource string) (*Program, error) {
	handle, err := CompileProgram(vertexSource, fragmentSource)
	if err != nil {
		return nil, err
	}
	program := Program{Handle: handle}
	program.attributeLocations = program.AttributeLocations()
	return &program, nil
}

func NewProgramFromFile(vertexPath, fragmentPath string) (*Program, error) {
	vertexSource, err := ioutil.ReadFile(vertexPath)
	if err != nil {
		return nil, err
	}
	fragmentSource, err := ioutil.ReadFile(fragmentPath)
	if err != nil {
		return nil, err
	}
	return NewProgram(string(vertexSource), string(fragmentSource))
}

func (p *Program) Delete() {
	gl.DeleteProgram(p.Handle)
}

func (p *Program) Use() {
	gl.UseProgram(p.Handle)
	p.AttributeNames()
}

func (p *Program) AttributeNames() []string {
	var count int32
	gl.GetProgramiv(p.Handle, gl.ACTIVE_ATTRIBUTES, &count)
	result := make([]string, count)
	for i := 0; i < int(count); i++ {
		var size int32
		var dataType uint32
		name := strings.Repeat("\x00", 256)
		gl.GetActiveAttrib(
			p.Handle, uint32(i), 256, nil,
			&size, &dataType, gl.Str(name))
		result[i] = name
	}
	return result
}

func (p *Program) AttributeLocations() []int {
	names := p.AttributeNames()
	result := make([]int, len(names))
	for i, name := range names {
		result[i] = p.AttributeLocation(name)
	}
	return result
}

func (p *Program) AttributeLocation(name string) int {
	return int(gl.GetAttribLocation(p.Handle, gl.Str(name+"\x00")))
}

func (p *Program) UniformLocation(name string) int {
	return int(gl.GetUniformLocation(p.Handle, gl.Str(name+"\x00")))
}

func (p *Program) SetBuffer(location, size, offset, stride int, buffer *Buffer) {
	buffer.Bind()
	gl.VertexAttribPointer(
		uint32(location), int32(size), gl.FLOAT, false, int32(stride),
		gl.PtrOffset(offset))
	buffer.Unbind()
}

func (p *Program) SetMatrix(location int, value Matrix) {
	data := value.ColMajor()
	gl.UniformMatrix4fv(int32(location), 1, false, &data[0])
}

func (p *Program) SetInt(location int, value int32) {
	gl.Uniform1i(int32(location), value)
}

func (p *Program) SetFloat(location int, value float32) {
	gl.Uniform1f(int32(location), value)
}

func (p *Program) DrawTriangles(offset, count int) {
	for _, location := range p.attributeLocations {
		gl.EnableVertexAttribArray(uint32(location))
	}
	gl.DrawArrays(gl.TRIANGLES, int32(offset), int32(count))
	for _, location := range p.attributeLocations {
		gl.DisableVertexAttribArray(uint32(location))
	}
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
