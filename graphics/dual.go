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

// Represents both projective points and lines.
type Dual struct {
	x, y, w float32
}

func (d *Dual) Bisector(other *Dual) *Dual {
	sx := d.x + other.x
	sy := d.y + other.y
	dx := other.x - d.x
	dy := other.y - d.y
	return &Dual{dx, dy, -0.5 * (dx*sx + dy*sy)}
}

func (d *Dual) Normalize() *Dual {
	return &Dual{d.x / d.w, d.y / d.w, 1}
}

func (a *Dual) ProjectiveDistanceTo(b *Dual) float32 {
	return a.x*b.x + a.y*b.y + a.w*b.w
}

func (a *Dual) AngularDistanceTo(b *Dual) float32 {
	return a.x*b.x + a.y*b.y
}

func (d0 *Dual) Intersection(d1 *Dual) *Dual {
	return &Dual{d0.y*d1.w - d0.w*d1.y,
		d0.w*d1.x - d0.x*d1.w,
		d0.x*d1.y - d0.y*d1.x}
}
