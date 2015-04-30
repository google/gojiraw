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
