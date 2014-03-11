package geometry

// Assorted primitives inspired by the image.Point and image.Rect
// code from Go.

import (
	"image"
	"strconv"
)

// A Pointf is an X, Y coordinate pair. The axes increase right and down.
// TODO(rjkroege|ianvollick): Write conversion Pointf <-> Dual as needed.
type Pointf struct {
	X, Y float32
}

// String returns a string representation of p like "(3,4)".
func (p Pointf) String() string {
	return "(" + strconv.FormatFloat(float64(p.X), 'f', -1, 32) + "," + strconv.FormatFloat(float64(p.Y), 'f', -1, 32) + ")"
}

// Add returns the vector p+q.
func (p Pointf) Add(q Pointf) Pointf {
	return Pointf{p.X + q.X, p.Y + q.Y}
}

// Sub returns the vector p-q.
func (p Pointf) Sub(q Pointf) Pointf {
	return Pointf{p.X - q.X, p.Y - q.Y}
}

// Mul returns the vector p*k.
func (p Pointf) Mul(k float32) Pointf {
	return Pointf{p.X * k, p.Y * k}
}

// Div returns the vector p/k.
func (p Pointf) Div(k float32) Pointf {
	return Pointf{p.X / k, p.Y / k}
}

// In reports whether p is in r.
func (p Pointf) In(r Rectanglef) bool {
	return r.Min.X <= p.X && p.X < r.Max.X &&
		r.Min.Y <= p.Y && p.Y < r.Max.Y
}

// Eq reports whether p and q are equal.
func (p Pointf) Eq(q Pointf) bool {
	return p.X == q.X && p.Y == q.Y
}

// ZP is the zero Pointf.
var ZP Pointf

// Ptf is shorthand for Pointf{X, Y}.
func Ptf(X, Y float32) Pointf {
	return Pointf{X, Y}
}

// Make a Pointf from an integer point
func Ptfi(p image.Point) Pointf {
	return Pointf{float32(p.X), float32(p.Y)}
}

// A Rectanglef contains the points with Min.X <= X < Max.X, Min.Y <= Y < Max.Y.
// It is well-formed if Min.X <= Max.X and likewise for Y. Pointfs are always
// well-formed. A rectangle's methods always return well-formed outputs for
// well-formed inputs.
type Rectanglef struct {
	Min, Max Pointf
}

// String returns a string representation of r like "(3,4)-(6,5)".
func (r Rectanglef) String() string {
	return r.Min.String() + "-" + r.Max.String()
}

// Dx returns r's width.
func (r Rectanglef) Dx() float32 {
	return r.Max.X - r.Min.X
}

// Dy returns r's height.
func (r Rectanglef) Dy() float32 {
	return r.Max.Y - r.Min.Y
}

// Size returns r's width and height.
func (r Rectanglef) Size() Pointf {
	return Pointf{
		r.Max.X - r.Min.X,
		r.Max.Y - r.Min.Y,
	}
}

// Add returns the rectangle r translated by p.
func (r Rectanglef) Add(p Pointf) Rectanglef {
	return Rectanglef{
		Pointf{r.Min.X + p.X, r.Min.Y + p.Y},
		Pointf{r.Max.X + p.X, r.Max.Y + p.Y},
	}
}

// Sub returns the rectangle r translated by -p.
func (r Rectanglef) Sub(p Pointf) Rectanglef {
	return Rectanglef{
		Pointf{r.Min.X - p.X, r.Min.Y - p.Y},
		Pointf{r.Max.X - p.X, r.Max.Y - p.Y},
	}
}

// Intersect returns the largest rectangle contained by both r and s. If the
// two rectangles do not overlap then the zero rectangle will be returned.
func (r Rectanglef) Intersect(s Rectanglef) Rectanglef {
	if r.Min.X < s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y < s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X > s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y > s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	if r.Min.X > r.Max.X || r.Min.Y > r.Max.Y {
		return ZR
	}
	return r
}

// Union returns the smallest rectangle that contains both r and s.
func (r Rectanglef) Union(s Rectanglef) Rectanglef {
	if r.Min.X > s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y > s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X < s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y < s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	return r
}

// Empty reports whether the rectangle contains no points.
func (r Rectanglef) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

// Eq reports whether r and s are equal.
func (r Rectanglef) Eq(s Rectanglef) bool {
	return r.Min.X == s.Min.X && r.Min.Y == s.Min.Y &&
		r.Max.X == s.Max.X && r.Max.Y == s.Max.Y
}

// Overlaps reports whether r and s have a non-empty intersection.
func (r Rectanglef) Overlaps(s Rectanglef) bool {
	return r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

// In reports whether every point in r is in s.
func (r Rectanglef) In(s Rectanglef) bool {
	if r.Empty() {
		return true
	}
	// Note that r.Max is an exclusive bound for r, so that r.In(s)
	// does not require that r.Max.In(s).
	return s.Min.X <= r.Min.X && r.Max.X <= s.Max.X &&
		s.Min.Y <= r.Min.Y && r.Max.Y <= s.Max.Y
}

// Canon returns the canonical version of r. The returned rectangle has minimum
// and maximum coordinates swapped if necessary so that it is well-formed.
func (r Rectanglef) Canon() Rectanglef {
	if r.Max.X < r.Min.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Max.Y < r.Min.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}

// ZR is the zero Rectanglef.
var ZR Rectanglef

// Rect is shorthand for Rectanglef{Pt(x0, y0), Pt(x1, y1)}.
func Rect(x0, y0, x1, y1 float32) Rectanglef {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return Rectanglef{Pointf{x0, y0}, Pointf{x1, y1}}
}
