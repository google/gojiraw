// TODO: make a header
package geometry

import (
  "math"
)

// Returns sin(2 * atan(d))
func Sin2Atan(d float64) float64 {
  return 2.0 * d / (1.0 + d * d)
}

// Returns cos(2 * atan(d))
func Cos2Atan(d float64) float64 {
  return (1.0 - d * d) / (1.0 + d * d)
}

// Represents both projective points and lines.
type Dual struct {
  x, y, w float64
}

type SignedVector struct {
  x, y float64
  negated bool
}

type Arc struct {
  p0, p1 *Dual
  // The two endpoints of the chord and the furthest point from the line
  // coincident with the chord form a isosceles triangle. |d| is the tangent
  // of the doubled angle. Roughly speaking, it's the "depth" of the arc.
  //
  // NB: It MUST be the case that -1 < d < 1. In fact, it would do well to be
  // smaller than 1/2 in magnitude.
  d float64
}

func (arc *Arc) Normals() (l0, l1 *Dual) {
  // (dx, dy) is the vector from p0 to p1. We will rotate it to get the normal
  // to the lines l0 and l1 by +/- 2 * atan(d) respectively.
  //
  // Why 2 * atan(d)? Great question. Let's do some geometry.
  //
  // Consider the following diagram. 
  //
  //                                  c
  //                                  |
  //                                  |
  //                    a_____________m_____________b
  //                                  | 
  //                                  |
  //                                  |
  //                                  o
  //                                  |
  //                                  |
  //                                  |
  //                                  |
  //                                  |
  //                                  |
  //                                  d
  //
  // Duals a, b, c and d are circumscribed by the same circle. Duals a and b
  // are p0 and p1 in our arc parlance. o is the center of the circle, m is the
  // midpoint of the chord subtended by a and b and c is the furthest point on
  // our arc from the chord. d is not on our arc, but is the other point of
  // intection on the circle of the bisector of a and b.
  //
  // We want angle(aoc). That is the angle we would have to rotate vector(co)
  // to get the normal at a or b.
  //
  // Now we know that atan(d) = angle(cam). Since cam forms a right triangle
  // angle(acm) = pi/2 - angle(cam). Since cd passes through the center of the
  // circle, we also that angle(dac) = pi/2. This means that the triangles cam
  // and cad are similar, which is handy because it means that angle(cad) =
  // angle(cam) = atan(d). Next, notice that triangles aod and aoc (which form
  // triangle cad) are both isosceles. This implies that angle(dao) = atan(d) as
  // well. Since angle(dac) = pi/2, we know that angle(oac) = angle(aco) =
  // pi/2 - atan(d). Finally, since the angles in triangle aoc must sum to pi,
  // we know that angle(aoc) = pi - angle(oac) - angle(aco)
  //                         = pi - 2 * (p/2 - atan(d))
  //                         = pi - pi + 2 * atan(d)
  //                         = 2 * atan(d).
  dx := arc.p1.x - arc.p0.x
  dy := arc.p1.y - arc.p0.y

  // We happen to want sin and cos of this angle for rotating. There happen to
  // be cute closed forms for getting these values.
  sin := Sin2Atan(arc.d)
  cos := Cos2Atan(arc.d)

  // Initialize the lines with the correct orientation, but passing through the
  // origin.
  l0 = &Dual{ dx * cos - dy * sin, dx * sin + dy * cos, 0.0 }
  l1 = &Dual{ -dx * cos - dy * sin, dx * sin - dy * cos, 0.0 }

  // Adjust the lines to pass through p0, and p1.
  l0.w = -l0.AngularDistanceTo(arc.p0)
  l1.w = -l1.AngularDistanceTo(arc.p1)

  return l0, l1
}

func IsInWedge(point, n0, n1 *Dual) bool {
  return n0.ProjectiveDistanceTo(point) > 0 && n1.ProjectiveDistanceTo(point) > 0
}

func (arc *Arc) IsInWedge(point *Dual) bool {
  // Perhaps we should cache these?
  n0, n1 := arc.Normals()
  return IsInWedge(point, n0, n1)
}

func (arc *Arc) Apex() *Dual {
  dx := arc.p1.x - arc.p0.x
  dy := arc.p1.y - arc.p0.y

  return &Dual{ arc.p0.x + 0.5 * (dx - arc.d * dy),
                arc.p0.y + 0.5 * (dy + arc.d * dx),
                1.0 }
}

func (arc *Arc) SignedVectorToClosestArcPoint(point *Dual) *SignedVector {
  n0, n1 := arc.Normals()
  if (!IsInWedge(point, n0, n1)) {
    dx0 := point.x - arc.p0.x
    dy0 := point.y - arc.p0.y
    dx1 := point.x - arc.p1.x
    dy1 := point.y - arc.p1.y
    if (dx0 * dx0 + dy0 * dy0 < dx1 * dx1 + dy1 * dy1) {
      return &SignedVector{ dx0, dy0, false }
    }
    return &SignedVector{ dx1, dy1, false }
  }

  d0 := n0.ProjectiveDistanceTo(point)
  d1 := n1.ProjectiveDistanceTo(point)

  // (px, py) is a vector in the dirction of the tangent at the closest point on
  // the circle to |point|.
  px := d1 * n0.x - d0 * n1.x
  py := d1 * n0.y - d0 * n1.y

  // This is the line that passes through the center of the circle and p.
  l0 := Dual{ px, py, -(px * point.x + py * point.y) }

  // FIXME: the only reason for the following sqrt's is to get the angular
  // bisector of n0 and (px, py). Can that be obtained more cheaply?
  length0 := math.Hypot(n0.x, n0.y)
  length1 := math.Hypot(px, py)

  // This is is a line (through the origin) coincident with the angular bisector
  // of n0 and (px, py).
  l1 := Dual{ -py / length1 - n0.y / length0,
               px / length1 + n0.x / length0,
               0.0 }

  // Make it pass through p0.
  l1.w = -(l1.x * arc.p0.x + l1.y * arc.p0.y)

  // This is the closest point on the arc.
  arcPoint := l1.Intersection(&l0).Normalize()

  v := SignedVector{ arcPoint.x - point.x, arcPoint.y - point.y, true }

  v.negated = v.x * py - v.y * px < 0
  if (arc.d < 0) {
    v.negated = !v.negated
  }

  return &v
}

func (arc *Arc) EuclideanDistanceTo(point *Dual) float64 {
  v := arc.SignedVectorToClosestArcPoint(point)
  distance := math.Hypot(v.x, v.y)
  if (v.negated) {
    distance *= -1.0
  }
  return distance
}

func (d *Dual) Bisector(other *Dual) *Dual {
  sx := d.x + other.x
  sy := d.y + other.y
  dx := other.x - d.x
  dy := other.y - d.y
  return &Dual{ dx, dy, -0.5 * (dx * sx + dy * sy) }
}

func (d *Dual) Normalize() *Dual {
  return &Dual{ d.x / d.w, d.y / d.w, 1 }
}

func (a *Dual) ProjectiveDistanceTo(b *Dual) float64 {
  return a.x * b.x + a.y * b.y + a.w * b.w
}

func (a *Dual) AngularDistanceTo(b *Dual) float64 {
  return a.x * b.x + a.y * b.y
}

func (d0 *Dual) Intersection(d1 *Dual) *Dual {
  return &Dual{ d0.y * d1.w - d0.w * d1.y,
                d0.w * d1.x - d0.x * d1.w,
                d0.x * d1.y - d0.y * d1.x }
}

////////////////////////////////////////////////////////////////////////////////
//
// The following methods are "interesting" to me, and I thought they might be
// useful in the future, but aren't yet used.
// 
// FIXME: If these turn out to be dead methods, kill them.
//
////////////////////////////////////////////////////////////////////////////////

// Returns the closest point to |p| on |line|
func (line *Dual) ClosestPoint(p *Dual) *Dual {
  return line.Intersection(line.Intersection(p))
}

func (arc *Arc) Tangents() (t0, t1 *Dual) {
  // See |Normals| for an explanation of the following.
  dx := arc.p1.x - arc.p0.x
  dy := arc.p1.y - arc.p0.y
  sin := Sin2Atan(arc.d)
  cos := Cos2Atan(arc.d)
  t0 = &Dual{ -dx * sin - dy * cos, dx * cos - dy * sin, 0.0 }
  t1 = &Dual{ dy * cos - dx * sin, -dx * cos - dy * sin, 0.0 }
  t0.w = -t0.AngularDistanceTo(arc.p0)
  t1.w = -t1.AngularDistanceTo(arc.p1)

  return t0, t1
}

// Note the return value of this function is NOT normalized, and is highly non-
// linear. It should only be used to determine inside vs outside-ness.
// NB: Must be true that arc.IsInWedge(point).
func (arc *Arc) ProjectiveDistanceTo(point *Dual) float64 {
  // NB: most of the values computed on the fly in this function could be
  // precomputed.
  t0, t1 := arc.Tangents()
  c := arc.p0.Intersection(arc.p1)
  apex := arc.Apex()

  // The equation for a circle is implicitly formed by k^2 - l*m = 0, where k
  // is the distance to the line connecting p0 and p1 and l and m are the signed
  // distances to the tangents at those points. However, this does not imply a
  // circle; we may have to scale the term l * m. We choose a scale factor such
  // that the apex of the arc is on the circle.
  k := c.ProjectiveDistanceTo(apex)
  l := t0.ProjectiveDistanceTo(apex)
  m := t1.ProjectiveDistanceTo(apex)
  scale := k * k / (l * m)

  // FIXME: are there simplifications we could acheive, knowing that the
  // homogenous components of our inputs will always be 0?
  k = c.ProjectiveDistanceTo(point)
  l = t0.ProjectiveDistanceTo(point)
  m = t1.ProjectiveDistanceTo(point)
  return k * k - scale * l * m
}

// Measures the euclidean distance between the two in the plane.
func (d *Dual) EuclideanDistanceTo(other *Dual) float64 {
  n0 := d.Normalize()
  n1 := other.Normalize()
  return math.Hypot(n0.x - n1.x, n0.y - n1.y)
}

