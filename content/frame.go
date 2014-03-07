// TODO: make a header
package content

import (
	"github.com/go-gl/gl"
	"image"
	"image/color"
	"log"
)

// Element is a placeholder for an element.
type Element struct {
	// TODO(rjkroege): Convert this to a float.
	image.Rectangle	
	color.RGBA
}

func (e *Element) init(tl image.Point) {
	e.Min = tl
	e.Max = tl.Add(image.Pt(90, 90))
	e.RGBA = color.RGBA{uint8(0xff), uint8(0), uint8(0), uint8(25)}
}

// Draw renders a single quad using OpenGL
func (e *Element) Draw() (ow, oh float32) {
	ow = float32(e.Max.X)
	oh = float32(e.Max.Y)

	gl.Color4b(int8(e.R), int8(e.G), int8(e.R), int8(e.A))
	gl.Begin(gl.QUADS)
	gl.Vertex2i(e.Min.X, e.Min.Y)
	gl.Vertex2i(e.Max.X, e.Min.Y)
	gl.Vertex2i(e.Max.X, e.Max.Y)
	gl.Vertex2i(e.Min.X, e.Max.Y)
	gl.End()

	return
}

// Frame is the Gojira equivalent of a RenderFrame in Chrome?
// The basic skeleton of shipping the display list and RenderFrame
// sizes between processes will probably entail re-writing this.
type Frame struct {
	// Probably gonna own a page. :)

	// Translation
	x, y, z float32

	// Viewport
	w, h float32

	// rjk's understanding of Go: the backing for this is big but this object
	// is a slice and therefore Frame is small.
	displaylist []Element
}

// AddElement extends the display list slice and fills in the new element
// with a quad.
func (f *Frame) AddElement(p image.Point) {
	ne := len(f.displaylist)
	ndl := f.displaylist[0:ne + 1]
	(&ndl[ne]).init(p)
	f.displaylist = ndl
}

func Max(x1, x2 float32) float32 {
	if x1 > x2 {
		return x1
	}
	return x2
}

func Min(x1, x2 float32) float32 {
	if x1 < x2 {
		return x1
	}
	return x2
}

// Pan sets a translation on the Frame to permit the Frame to move around within
// its viewport. Translates the viewport w.r.t. the Frame origin by p.
func (f *Frame) Pan(dx, dy float32) {
	f.x =Min(f.x + dx, f.w)
	f.y += Min(f.y + dy, f.h)
	log.Printf("translation: %d %d", f.x, f.y )
}

// Resize tells the frame what its size should be.
func (f *Frame) Resize(w, h float32) {
	f.w = Max(w, f.w)
	f.h = Max(h, f.h)
	log.Printf("current size: %d %d", f.w, f.h)
}

func NewFrame() *Frame {
	dl := make([]Element, 0, 1000)
	return &Frame{displaylist: dl}
}

// TODO(rjk): Tell the Frame to clip its drawing to a given viewport.
// x, y is the position of the Frame's origin in the containing port.
// w, h is the width and height of the port in the port's coordinates.
// Returns the enclosing boundary of the Frame.
// TODO(rjkroege): boundaries should admit objects outside [0, w), [0. h)?
// TODO(rjkroege): Provide and wire in types for stuff, boxes, etc.
func (frame *Frame) Draw(x, y, vw, vh float32) (fw, fh float32) {
	gl.PushMatrix()
	gl.Translatef(x, y, 0)
	gl.ClearColor(1, 1, 1, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.Enable(gl.BLEND)
	gl.Enable(gl.POINT_SMOOTH)
	gl.Enable(gl.LINE_SMOOTH)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	fw = 0.0
	fh = 0.0

	for _, e := range(frame.displaylist) {
		ow, oh := e.Draw()
		fw  = Max(fw, ow)
		fh = Max(fh, oh)
	}
	gl.PopMatrix()
	return
}
