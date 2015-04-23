// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphics

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
