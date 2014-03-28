// TODO: make a header
package content

import (
	"github.com/go-gl/gl"
	"image"
	"image/color"
	"log"
	"code.google.com/a/google.com/p/gojira/geometry"
)

// Element is a placeholder for an element.
type Element struct {
	quad [4]geometry.Pointf
	color   color.RGBA
	hoveron bool
	vertex int
}

const (
	dX = 45.
	dY = 45.
)

// TODO(rjkroege): choose demaraction layer where we switch
// to floating point coordinates. (above here though.)
func (e *Element) init(tl image.Point) {
	pt := geometry.Ptfi(tl)

	e.quad[0] = geometry.Pointf{pt.X - dX, pt.Y - dY}
	e.quad[1] = geometry.Pointf{pt.X + dX, pt.Y - dY}
	e.quad[2] = geometry.Pointf{pt.X + dX, pt.Y + dY}
	e.quad[3] = geometry.Pointf{pt.X - dX, pt.Y + dY}

	e.color = color.RGBA{uint8(0), uint8(0), uint8(0), uint8(25)}
	e.hoveron = false
}

func (e *Element) DrawHandle() {
	gl.PointSize(2.)
	gl.Color4ub(0x0, 0, 0, 0xf0)
	gl.Begin(gl.POINTS)
	for i, p := range(e.quad) {
		if !e.hoveron || i != e.vertex {
			// does this make another copy of the point?
			gl.Vertex2f(p.X, p.Y)
		}
	}
	gl.End()

	if e.hoveron {
		gl.PointSize(8.)
		gl.Color4ub(0xff, 0, 0, 0xff)
		gl.Begin(gl.POINTS)
		gl.Vertex2f(e.quad[e.vertex].X, e.quad[e.vertex].Y)
		gl.End()
	}
}

// Draw renders a single quad using OpenGL returning the maximum point
// in the Element needed to contain the element.
// TODO(rjkroege): Return the tightest bounding box.
func (e *Element) Draw() (ow, oh float32) {
	c := e.color
	ow = 0.
	oh = 0.

	gl.Color4ub(c.R, c.G, c.R, c.A)
	gl.Begin(gl.QUADS)
	for _, p := range(e.quad) {
		// My presumption is that this structure does not copy.
		ow = Max(ow, p.X)
		oh = Max(oh, p.Y)
		gl.Vertex2f(p.X, p.Y)
	}
	gl.End()

	e.DrawHandle()

	return
}

func (e *Element) HoverOn(v int) {
	log.Printf("HoverOn")
	e.hoveron = true
	e.vertex = v
}

func (e *Element) HoverOff() {
	log.Printf("HoverOff")
	e.hoveron = false
}

// TODO(rjkroege): move this up in the file
// TODO(rjkroege): split apart Element and Frame code.
func (f *Element) FindVertex(p geometry.Pointf) int {
	o := geometry.Pointf{4., 4.}
	for i, v := range(f.quad) {
		r := geometry.Rectanglef{v.Sub(o), v.Add(o)}
		if p.In(r) {
			return i
		}
	}
	return -1
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
	// A list of elements. Note that Elements should really form a tree.
	// Perhaps also, we want to split the concept of the display list from
	// the model. At first however: the implicit Z order is from most recent.
	displaylist []Element

	// The most recently mouse-overed element or nil.
	overelement *Element
}

// AddElement extends the display list slice and fills in the new element
// with a quad.
func (f *Frame) AddElement(p image.Point) {
	ne := len(f.displaylist)
	ndl := f.displaylist[0 : ne+1]
	(&ndl[ne]).init(p)
	f.displaylist = ndl
}

// Find the control point, if any, under Point p. Return nil, 0 if there is no
// control point for an element under p. The returned int is the index
// of the vertex.
// TODO(rjkroege): The model could be a tree eventually. :-)
func (f *Frame) FindElementAtPoint(p image.Point) (*Element, int) {
	pf := geometry.Ptfi(p)
	log.Printf("FindElementAtPoint %v", p)
	for i := len(f.displaylist) - 1; i >= 0; i-- {
		e := &f.displaylist[i]
		v := e.FindVertex(pf)
		
		// TODO(rjkroege): Replace this appropriately.
		// r := geometry.Rectanglef{e.quad[0], e.quad[2]}
		if v > -1 {
			return e, v
		}
	}
	return nil, -1
}

// Adjusts visual style for elements that are under the
// mouse pointer.
func (f *Frame) MouseOver(e *Element,  v int) {
	log.Printf("MouseOver: %+v", e)
	if f.overelement != nil {
		f.overelement.HoverOff()
	}
	f.overelement = e
	if e != nil {
		e.HoverOn(v)
	}
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
	f.x = Min(f.x+dx, f.w)
	f.y += Min(f.y+dy, f.h)
	log.Printf("translation: %d %d", f.x, f.y)
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

	for _, e := range frame.displaylist {
		ow, oh := e.Draw()
		fw = Max(fw, ow)
		fh = Max(fh, oh)
	}
	gl.PopMatrix()
	return
}
