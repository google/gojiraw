// Copyright 2014 The Gojiraw Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package graphics

import (
	"fmt"
	"image/color"
	"log"

	"github.com/go-gl/gl"
	"github.com/go-gl/glu"
)

const (
	// Sort these alphabetically or vollick will hunt you down.
	DRAW_OP_COLOR = iota
	DRAW_OP_QUADS
)

func CheckForGLErrors() {
	errcode := gl.GetError()
	if errcode != gl.NO_ERROR {
		// The error is non-nil here if we can't get an error string.
		s, err := glu.ErrorString(errcode)
		if err != nil {
			log.Panic("GLError(string): ", s)
		} else {
			log.Panicf("GLError(code): %x", errcode)
		}
	}
}

const (
	// The default vertex shader takes 2D points and scales them to fit in
	// the viewport.
	defaultVertexShader = `
#version 400

// The x and y components represent the reciprocal of the width and height of the
// viewport. We've used the reciprocal to avoid unnecessary division.
uniform vec2 u_Viewport;

// The default shader only supports 2D points.
in vec2 in_Position;

void main()
{
    gl_Position = vec4(2.0 * in_Position.x * u_Viewport.x - 1.0,
                       -(2.0 * in_Position.y * u_Viewport.y - 1.0),
                       0.0, 1.0);
}` // defaultVertexShader

	// The default fragment shader simply passes along a uniform color.
	defaultFragmentShader = `
#version 400

uniform vec4 u_Color;
out vec4 out_Color;

void main()
{
    out_Color = u_Color;
}` // defaultFragmentShader

)

func CreateDefaultShaders() (program gl.Program) {
	vertex_shader := gl.CreateShader(gl.VERTEX_SHADER)
	vertex_shader.Source(defaultVertexShader)
	vertex_shader.Compile()
	fmt.Println(vertex_shader.GetInfoLog())
	defer vertex_shader.Delete()

	fragment_shader := gl.CreateShader(gl.FRAGMENT_SHADER)
	fragment_shader.Source(defaultFragmentShader)
	fragment_shader.Compile()
	fmt.Println(fragment_shader.GetInfoLog())
	defer fragment_shader.Delete()

	program = gl.CreateProgram()
	program.AttachShader(vertex_shader)
	program.AttachShader(fragment_shader)

	program.BindFragDataLocation(0, "out_Color")
	program.Link()
	program.Use()
	return
}

// TODO(vollick): We might just want to store a gob encoder here? Do we need to
// have tighter control of the binary rep so we can pass data directly to card.

// TODO(vollick): We need to consider spatial queries, mutability and display
// list optimization. I am not at all convinced that this representation is
// ideal for all of these purposes.
type DisplayList struct {
	opCodes                          []uint8
	integers                         []uint32
	floats                           []float32
	bytes                            []uint8
	cur_integer, cur_float, cur_byte int
	W, H                             float32
	cur_point_size                   float32
}

func (dl *DisplayList) SetColor(c color.RGBA) {
	dl.opCodes = append(dl.opCodes, DRAW_OP_COLOR)
	dl.bytes = append(dl.bytes, c.R, c.G, c.B, c.A)
}

func (dl *DisplayList) SetPointSize(s float32) {
	dl.cur_point_size = s
}

func (dl *DisplayList) DrawPoints(ps []Pointf) {
	// TODO(vollick): We could, if we wanted, use the "unsafe" package here to
	// add the raw bytes.
	quads := make([][4]Pointf, len(ps))
	delta := 0.5 * dl.cur_point_size
	for i, p := range ps {
		quads[i] = [...]Pointf{
			Pointf{p.X - delta, p.Y - delta},
			Pointf{p.X + delta, p.Y - delta},
			Pointf{p.X + delta, p.Y + delta},
			Pointf{p.X - delta, p.Y + delta}}
	}
	dl.DrawQuads(quads)
}

func (dl *DisplayList) DrawQuads(qs [][4]Pointf) {
	dl.opCodes = append(dl.opCodes, DRAW_OP_QUADS)
	dl.integers = append(dl.integers, uint32(len(qs)))
	for _, q := range qs {
		for _, p := range q {
			dl.W = MaxF(dl.W, p.X)
			dl.H = MaxF(dl.H, p.Y)
			dl.floats = append(dl.floats, p.X, p.Y)
		}
	}
}

func (dl *DisplayList) Draw(program *gl.Program, width, height float32) {
	viewportUniform := program.GetUniformLocation("u_Viewport")
	viewportUniform.Uniform2f(1.0/width, 1.0/height)

	// TODO(vollick): Can we do something like this in parallel?
	dl.cur_integer = 0
	dl.cur_float = 0
	dl.cur_byte = 0
	for _, op := range dl.opCodes {
		switch op {
		case DRAW_OP_COLOR:
			dl.DoColor(program)
		case DRAW_OP_QUADS:
			dl.DoQuads(program)
		}
		CheckForGLErrors()
	}
}

func (dl *DisplayList) DoColor(program *gl.Program) {
	colorLocation := program.GetUniformLocation("u_Color")
	colorLocation.Uniform4f(
		float32(dl.bytes[dl.cur_byte])/255,
		float32(dl.bytes[dl.cur_byte+1])/255,
		float32(dl.bytes[dl.cur_byte+2])/255,
		float32(dl.bytes[dl.cur_byte+3])/255)
	dl.cur_byte += 4
	CheckForGLErrors()
}

func (dl *DisplayList) DoQuads(program *gl.Program) {
	num_quads := dl.integers[dl.cur_integer]
	dl.cur_integer++
	quads := []float32{}

	// FIXME: we shouldn't recreate this every time the display list is
	// drawn.
	for i := uint32(0); i < num_quads; i++ {
		quads = append(quads, dl.floats[dl.cur_float:dl.cur_float+6]...)
		quads = append(quads, dl.floats[dl.cur_float:dl.cur_float+2]...)
		quads = append(quads, dl.floats[dl.cur_float+4:dl.cur_float+6]...)
		quads = append(quads, dl.floats[dl.cur_float+6:dl.cur_float+8]...)
		dl.cur_float += 8
	}

	// FIXME: this is atrocious. We need to retain these objects rather than
	// creating and destroying them constantly.
	vao := gl.GenVertexArray()
	vao.Bind()

	vbo := gl.GenBuffer()
	vbo.Bind(gl.ARRAY_BUFFER)

	gl.BufferData(gl.ARRAY_BUFFER, len(quads)*4, quads, gl.STATIC_DRAW)

	positionAttrib := program.GetAttribLocation("in_Position")
	positionAttrib.AttribPointer(2, gl.FLOAT, false, 0, nil)
	positionAttrib.EnableArray()
	defer positionAttrib.DisableArray()

	gl.DrawArrays(gl.TRIANGLES, 0, len(quads)/2)

	CheckForGLErrors()
}
