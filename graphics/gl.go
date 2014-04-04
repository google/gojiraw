package graphics

import (
	"github.com/go-gl/gl"
	"image/color"
)

const (
	// Sort this alphabetically or vollick will hunt you down.
	DRAW_OP_COLOR = iota
	DRAW_OP_POINTS
	DRAW_OP_POINT_SIZE
	DRAW_OP_QUADS
)

// TODO(vollick): We might just want to store a gob encoder here? Do we need to
// have tighter control of the binary rep so we can pass data directly to card?

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
}

func (dl *DisplayList) SetColor(c color.RGBA) {
	dl.opCodes = append(dl.opCodes, DRAW_OP_COLOR)
	dl.bytes = append(dl.bytes, c.R, c.G, c.B, c.A)
}

// TODO(vollick): It would be nice if we know the current "state" so that we
// could reject redundant sets immediately.
func (dl *DisplayList) SetPointSize(s float32) {
	dl.opCodes = append(dl.opCodes, DRAW_OP_POINT_SIZE)
	dl.floats = append(dl.floats, s)
}

func (dl *DisplayList) DrawPoints(ps []Pointf) {
	dl.opCodes = append(dl.opCodes, DRAW_OP_POINTS)
	dl.integers = append(dl.integers, uint32(len(ps)))
	// TODO(vollick): We could, if we wanted, use the "unsafe" package here to
	// add the raw bytes.
	for _, p := range ps {
		dl.W = MaxF(dl.W, p.X)
		dl.H = MaxF(dl.H, p.Y)
		dl.floats = append(dl.floats, p.X, p.Y)
	}
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

func (dl *DisplayList) Do() {
	// TODO(vollick): Can we do something like this in parallel?
	dl.cur_integer = 0
	dl.cur_float = 0
	dl.cur_byte = 0
	for _, op := range dl.opCodes {
		switch op {
		case DRAW_OP_COLOR:
			dl.DoColor()
		case DRAW_OP_POINTS:
			dl.DoPoints()
		case DRAW_OP_POINT_SIZE:
			dl.DoPointSize()
		case DRAW_OP_QUADS:
			dl.DoQuads()
		}
	}
}

func (dl *DisplayList) DoColor() {
	gl.Color4ub(dl.bytes[dl.cur_byte],
		dl.bytes[dl.cur_byte+1],
		dl.bytes[dl.cur_byte+2],
		dl.bytes[dl.cur_byte+3])
	dl.cur_byte += 4
}

func (dl *DisplayList) DoPoints() {
	gl.Begin(gl.POINTS)
	num_points := dl.integers[dl.cur_integer]
	dl.cur_integer++
	for i := uint32(0); i < num_points; i++ {
		gl.Vertex2f(dl.floats[dl.cur_float], dl.floats[dl.cur_float+1])
		dl.cur_float += 2
	}
	gl.End()
}

func (dl *DisplayList) DoPointSize() {
	gl.PointSize(dl.floats[dl.cur_float])
	dl.cur_float++
}

func (dl *DisplayList) DoQuads() {
	gl.Begin(gl.QUADS)
	num_points := dl.integers[dl.cur_integer] * 4
	dl.cur_integer++
	for i := uint32(0); i < num_points; i++ {
		gl.Vertex2f(dl.floats[dl.cur_float], dl.floats[dl.cur_float+1])
		dl.cur_float += 2
	}
	gl.End()
}
