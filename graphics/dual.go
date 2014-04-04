// TODO: make a header
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
