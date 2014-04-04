package dom

import (
	"code.google.com/a/google.com/p/gojira/graphics"
	"image/color"
	"log"
)

const (
	VERTEX_NON   = iota // Draw the default vertex handle.
	VERTEX_HOVER        // Draw the hover vertex handle.
	VERTEX_PRESS        // Draw the mouse down vertex handle
)

const (
	QUAD_ELEMENT_DX = 45.
	QUAD_ELEMENT_DY = 45.
	QUAD_ELEMENT_DH = 4.
)

func (qe *QuadElement) vertexColor(mode int) color.RGBA {
	switch mode {
	case VERTEX_NON:
		return color.RGBA{0, 0, 0, 0xf0}
	case VERTEX_HOVER:
		return color.RGBA{0xff, 0, 0, 0xff}
	case VERTEX_PRESS:
		return color.RGBA{0, 0, 0xff, 0xff}
	}
	return color.RGBA{0, 0, 0, 0xf0}
}

// Element is a placeholder for an element. Will want, not a tree, but a true
// database format for efficent queries (including spatial).
type QuadElement struct {
	vertices     [4]graphics.Pointf
	color        color.RGBA
	hoverMode    int
	activeVertex int
}

func (qe *QuadElement) Init(pt graphics.Pointf) {
	qe.vertices[0] = graphics.Pointf{pt.X - QUAD_ELEMENT_DX, pt.Y - QUAD_ELEMENT_DY}
	qe.vertices[1] = graphics.Pointf{pt.X + QUAD_ELEMENT_DX, pt.Y - QUAD_ELEMENT_DY}
	qe.vertices[2] = graphics.Pointf{pt.X + QUAD_ELEMENT_DX, pt.Y + QUAD_ELEMENT_DY}
	qe.vertices[3] = graphics.Pointf{pt.X - QUAD_ELEMENT_DX, pt.Y + QUAD_ELEMENT_DY}
	qe.color = color.RGBA{uint8(0), uint8(0), uint8(0), uint8(25)}
	qe.hoverMode = VERTEX_NON
	qe.activeVertex = -1
}

func (qe *QuadElement) ActivateVertex(i int) graphics.Pointf {
	qe.hoverMode = VERTEX_PRESS
	qe.activeVertex = i
	return qe.vertices[i]
}

func (qe *QuadElement) Deactivate() {
	qe.hoverMode = VERTEX_HOVER
}

func (qe *QuadElement) SetActiveVertex(v graphics.Pointf) {
	qe.vertices[qe.activeVertex] = v
}

func (qe *QuadElement) drawHandle(dl *graphics.DisplayList) {
	dl.SetPointSize(2.)
	dl.SetColor(qe.vertexColor(VERTEX_NON))

	var ps [4]graphics.Pointf
	count := 0
	for i, p := range qe.vertices {
		if qe.hoverMode == VERTEX_NON || i != qe.activeVertex {
			ps[count] = p
			count++
		}
	}

	dl.DrawPoints(ps[:])
	if qe.hoverMode != VERTEX_NON {
		log.Printf("drawing hover vertex")
		dl.SetPointSize(2 * QUAD_ELEMENT_DH)
		if qe.hoverMode == VERTEX_HOVER {
			dl.SetColor(qe.vertexColor(VERTEX_HOVER))
		} else {
			dl.SetColor(qe.vertexColor(VERTEX_PRESS))
		}
		dl.DrawPoints([]graphics.Pointf{qe.vertices[qe.activeVertex]})
	}
}

// TODO(vollick): split this out into an interface.
func (qe *QuadElement) Draw(dl *graphics.DisplayList) {
	dl.SetColor(qe.color)
	dl.DrawQuads([][4]graphics.Pointf{qe.vertices})
	qe.drawHandle(dl)
}

func (qe *QuadElement) FindVertex(p graphics.Pointf) int {
	o := graphics.Pointf{QUAD_ELEMENT_DH, QUAD_ELEMENT_DH}
	for i, v := range qe.vertices {
		r := graphics.Rectanglef{v.Sub(o), v.Add(o)}
		if p.In(r) {
			return i
		}
	}
	return -1
}

func (qe *QuadElement) HoverOn(v int) {
	log.Printf("HoverOn %d", v)
	qe.hoverMode = VERTEX_HOVER
	qe.activeVertex = v
}

func (qe *QuadElement) HoverOff() {
	log.Printf("HoverOff")
	qe.hoverMode = VERTEX_NON
}
