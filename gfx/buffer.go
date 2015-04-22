package gfx

import (
	"reflect"

	"github.com/go-gl/gl/v2.1/gl"
)

type Buffer struct {
	Handle uint32
}

func NewBuffer() *Buffer {
	var handle uint32
	gl.GenBuffers(1, &handle)
	return &Buffer{handle}
}

func (b *Buffer) Delete() {
	gl.DeleteBuffers(1, &b.Handle)
}

func (b *Buffer) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, b.Handle)
}

func (b *Buffer) Unbind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (b *Buffer) SetItems(data interface{}) {
	v := reflect.ValueOf(data)
	size := int(v.Type().Elem().Size()) * v.Len()
	gl.BufferData(gl.ARRAY_BUFFER, size, gl.Ptr(data), gl.STATIC_DRAW)
}

func (b *Buffer) SetItem(index int, data interface{}) {
	v := reflect.ValueOf(data)
	size := int(v.Type().Size())
	slice := reflect.Append(
		reflect.MakeSlice(reflect.SliceOf(v.Type()), 0, 1), v).Interface()
	gl.BufferSubData(gl.ARRAY_BUFFER, size*index, size, gl.Ptr(slice))
}
