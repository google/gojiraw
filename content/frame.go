// TODO: make a header
package content

import (
	"image"
	"log"

	"code.google.com/a/google.com/p/gojira/content/dom"
	"code.google.com/a/google.com/p/gojira/graphics"
	"github.com/go-gl/gl"
)

// Frame is the Gojira equivalent of a RenderFrame in Chrome?
// The basic skeleton of shipping the display list and RenderFrame
// sizes between processes will probably entail re-writing this.
// Note(vollick): Frames always have documents. Documents are currently
// "special" in blink in that they don't have a renderer. Let's just leave the
// document related code in Frame. If anything, "document" could just be an
// interface Frame implements.
type Frame struct {
	// Probably gonna own a page. :)
	// (Note(vollick): Pages are the things that own frame trees. I think
	// that's analogous to windows. Let's move the page code into window.)

	// Translation
	// TODO(vollick): Probably want a transform here to capture scaling.
	x, y, z float32

	// Viewport
	w, h float32

	// The most recently mouse-overed element or nil.
	overElement *dom.QuadElement

	// The mouse is down.
	mouseDown bool

	// A floating point offset from a mouseDown to the centroid of the handle.
	// TODO(vollick): It's fishy that Frame knows anything about "handles."
	offset graphics.Pointf

	// The root of the document.
	document []dom.QuadElement
}

// AddElement extends the document slice and fills in the new element with a
// quad.
func (f *Frame) AddElement(p image.Point) {
	ne := len(f.document)
	nd := f.document[0 : ne+1]
	pf := graphics.Ptfi(p)
	// Note: the following code initializes the new point into the allocated
	// memory in the document list. It does *not* resize the list; if we have
	// more than 1000 quads, we'll die.
	// TODO(vollick): make this dynamic.
	(&nd[ne]).Init(pf)
	f.document = nd
}

// Find the control point, if any, under Point p. Return nil, 0 if there is no
// control point for an element under p. The returned int is the index
// of the vertex.
// TODO(rjkroege): The model could be a tree eventually. :-)
func (f *Frame) FindElementAtPoint(p image.Point) (*dom.QuadElement, int) {
	pf := graphics.Ptfi(p)
	log.Printf("FindElementAtPoint %v", p)
	for i := len(f.document) - 1; i >= 0; i-- {
		if v := f.document[i].FindVertex(pf); v != -1 {
			return &f.document[i], v
		}
	}
	return nil, -1
}

// Adjusts visual style for elements that are under the
// mouse pointer.
func (f *Frame) MouseOver(qe *dom.QuadElement, v int) {
	log.Printf("MouseOver: %+v, %d", qe, v)
	if f.overElement != nil && f.overElement != qe {
		f.overElement.HoverOff()
	}
	f.overElement = qe
	if qe != nil {
		qe.HoverOn(v)
	}
}

// Pan sets a translation on the Frame to permit the Frame to move around within
// its viewport. Translates the viewport w.r.t. the Frame origin by p.
func (f *Frame) Pan(dx, dy float32) {
	f.x = graphics.MinF(f.x+dx, f.w)
	f.y += graphics.MinF(f.y+dy, f.h)
	log.Printf("translation: %d %d", f.x, f.y)
}

// Resize tells the frame what its size should be.
func (f *Frame) Resize(w, h float32) {
	f.w = graphics.MaxF(w, f.w)
	f.h = graphics.MaxF(h, f.h)
	log.Printf("current size: %d %d", f.w, f.h)
}

func NewFrame() *Frame {
	// TODO(vollick): allow more than 1000 things.
	d := make([]dom.QuadElement, 0, 1000)
	return &Frame{document: d}
}

// TODO(rjk): Tell the Frame to clip its drawing to a given viewport.
// x, y is the position of the Frame's origin in the containing port.
// w, h is the width and height of the port in the port's coordinates.
// Returns the enclosing boundary of the Frame.
// TODO(rjkroege): boundaries should admit objects outside [0, w), [0. h)?
// TODO(rjkroege): Provide and wire in types for stuff, boxes, etc.
func (frame *Frame) Draw(x, y, vw, vh float32) (fw, fh float32) {
	// Build the display list.
	dl := &graphics.DisplayList{}
	for _, e := range frame.document {
		e.Draw(dl)
	}

	gl.PushMatrix()
	gl.Translatef(x, y, 0)
	gl.ClearColor(1, 1, 1, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.Enable(gl.BLEND)
	gl.Enable(gl.POINT_SMOOTH)
	gl.Enable(gl.LINE_SMOOTH)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	dl.Draw()

	gl.PopMatrix()

	return dl.W, dl.H
}

func (f *Frame) StartMouseDownMode(pt image.Point, qe *dom.QuadElement, v int) {
	f.overElement = qe
	f.mouseDown = true
	pf := graphics.Ptfi(pt)
	f.offset = pf.Sub(qe.ActivateVertex(v))
}

func (f *Frame) InMouseDownMode(pt image.Point) {
	// Is this idiomatic?
	pf := graphics.Ptfi(pt)
	if qe := f.overElement; qe != nil {
		qe.SetActiveVertex(pf.Add(f.offset))
	}
}

func (f *Frame) EndMouseDownMode() {
	f.mouseDown = false
	f.overElement.Deactivate()
}
