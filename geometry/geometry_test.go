// TODO: make a header
package geometry

import (
	"math"
	"testing"
	//  "fmt"
)

const (
	FLOAT_TOL = 1.0e-5
)

func AssertFloatEqual(t *testing.T, expected float64, actual float64) {
	delta := expected - actual
	if delta > FLOAT_TOL || delta < -FLOAT_TOL {
		t.Errorf("Expected near %f, received %f", expected, actual)
	}
}

func AssertTrue(t *testing.T, value bool) {
	if !value {
		t.Error("Expected true, got false.")
	}
}

func AssertFalse(t *testing.T, value bool) {
	if value {
		t.Error("Expected false, got true.")
	}
}

var (
	ATAN_TEST_CASES = []float64{-.99, 0, 0.0001, 0.99}
)

func TestTan2Atan(t *testing.T) {
	for _, theta := range ATAN_TEST_CASES {
		AssertFloatEqual(t, Tan2Atan(theta), math.Tan(2.0*math.Atan(theta)))
	}
}

func TestSin2Atan(t *testing.T) {
	for _, theta := range ATAN_TEST_CASES {
		AssertFloatEqual(t, Sin2Atan(theta), math.Sin(2.0*math.Atan(theta)))
	}
}

func TestCos2Atan(t *testing.T) {
	for _, theta := range ATAN_TEST_CASES {
		AssertFloatEqual(t, Cos2Atan(theta), math.Cos(2.0*math.Atan(theta)))
	}
}

/*
func TestPointEqual(t *testing.T) {
  p1 := new(Point)
  p2 := new(Point)
  AssertTrue(t, p1.Equal(p2))
  p1.x = 1
  AssertFalse(t, p1.Equal(p2))
  p2.x = 1
  AssertTrue(t, p1.Equal(p2))
  p2.y = 1
  AssertFalse(t, p1.Equal(p2))
}

func TestPointAdd(t *testing.T) {
  p := Point{2, 2}
  v := Vector{3, 4}
  AssertTrue(t, p.Add(&v).Equal(&Point{5, 6}))
}

func TestPointSubtract(t *testing.T) {
  p := Point{2, 2}
  v := Vector{3, 4}
  AssertTrue(t, p.Subtract(&v).Equal(&Point{-1, -2}))
}

func TestPointMidpoint(t *testing.T) {
  type TestCase struct {
    p Point
    expected Point
  }

  cases := []TestCase {
    TestCase{Point{1, 1}, Point{0.5, 0.5}},
    TestCase{Point{0, 1}, Point{0, 0.5}},
    TestCase{Point{0, 0}, Point{0, 0}},
    TestCase{Point{-1, -1}, Point{-0.5, -0.5}}}

  o := new(Point)
  for _, c := range cases {
    AssertTrue(t, o.Midpoint(&c.p).Equal(&c.expected))
  }
}

func TestPointBisector(t *testing.T) {
  type TestCase struct {
    p Point
    expected *Line
  }
  cases := []TestCase {
    TestCase{Point{1, 1}, new(Line)},
    TestCase{Point{0, 1}, new(Line)},
    TestCase{Point{-1, -1}, new(Line)},
  }
  o := new(Point)
  for _, c := range cases {
    b := o.Bisector(&c.p)
    m := o.Midpoint(&c.p)
    cp1 := b.ClosestPoint(&HomogenousPoint{o.x, o.y, 1}).Normalize()
    cp2 := b.ClosestPoint(&HomogenousPoint{c.p.x, c.p.y, 1}).Normalize()
    AssertTrue(t, m.Equal(cp1))
    AssertTrue(t, m.Equal(cp2))
  }
}

func TestPointBadBisector(t *testing.T) {
  o := new(Point)
  p := new(Point)
  b := o.Bisector(p)
  AssertFloatEqual(t, 0, b.n.x)
  AssertFloatEqual(t, 0, b.n.y)
  AssertFloatEqual(t, 0, b.c)
}
*/

// Returns tan(0.5 * acos(d))
func TanHalfAcos(d float64) float64 {
	return math.Tan(math.Acos(d) * 0.25)
}

func TestIsInWedge(t *testing.T) {
	p0 := &Dual{1.0, 0.0, 1.0}
	p1 := &Dual{2.0, 0.0, 1.0}
	center := &Dual{1.5, -0.5, 1.0}

	c0 := p0.Intersection(center)
	c1 := p1.Intersection(center)

	cos_d := c0.AngularDistanceTo(c1)
	a := Arc{p0, p1, math.Tan(math.Acos(cos_d) * 0.25)}

	type TestCase struct {
		p        *Dual
		expected bool
	}

	testCases := []TestCase{
		TestCase{&Dual{1.5, 0.5, 1.0}, true},
		TestCase{&Dual{1.5, -0.2, 1.0}, true},
		TestCase{&Dual{1.5, -1.5, 1.0}, false},
		TestCase{&Dual{3.5, 0.5, 1.0}, false},
		TestCase{&Dual{-1.5, 0.5, 1.0}, false}}

	for i, test := range testCases {
		if a.IsInWedge(test.p) != test.expected {
			t.Errorf("Case %d Failed: p(%f, %f), is in = %t\n",
				i, test.p.x, test.p.y, a.IsInWedge(test.p))
		}
	}
}

func TestProjectiveDistanceTo(t *testing.T) {
	p0 := &Dual{1.0, 0.0, 1.0}
	p1 := &Dual{2.0, 0.0, 1.0}
	center := &Dual{1.5, -0.5, 1.0}

	c0 := p0.Intersection(center)
	c1 := p1.Intersection(center)

	cos_d := c0.AngularDistanceTo(c1)
	arc := Arc{p0, p1, math.Tan(math.Acos(cos_d) * 0.25)}
	apex := arc.Apex()

	type TestCase struct {
		p                *Dual
		expectedDistance float64
	}

	testCases := []TestCase{
		TestCase{apex, 0.0},
		TestCase{p0, 0.0},
		TestCase{p1, 0.0}}

	// Crap is not normalized. Is of limited use..
	for _, test := range testCases {
		AssertFloatEqual(t, test.expectedDistance, arc.ProjectiveDistanceTo(test.p))
	}
}

func TestEuclideanDistanceTo(t *testing.T) {
	p0 := &Dual{1.0, 0.0, 1.0}
	p1 := &Dual{2.0, 0.0, 1.0}
	center := &Dual{1.5, -0.5, 1.0}

	c0 := p0.Intersection(center)
	c1 := p1.Intersection(center)

	cos_d := c0.AngularDistanceTo(c1)
	radius := math.Hypot(c0.x, c0.y)
	cos_d /= radius * radius
	arc := Arc{p0, p1, math.Tan(math.Acos(cos_d) * 0.25)}
	apex := arc.Apex()

	type TestCase struct {
		p                *Dual
		expectedDistance float64
	}

	exterior := Dual{1.7872, 5.22, 1.0}
	interior := Dual{1.26, 0.0115, 1.0}

	testCases := []TestCase{
		TestCase{apex, 0.0},
		TestCase{&Dual{apex.x, apex.y + 1.0, 1.0}, 1.0},
		TestCase{&Dual{apex.x, apex.y - 0.2, 1.0}, -0.2},
		TestCase{&Dual{apex.x, apex.y + 10.0, 1.0}, 10.0},
		TestCase{&exterior, math.Hypot(exterior.x-center.x, exterior.y-center.y) - radius},
		TestCase{&interior, math.Hypot(interior.x-center.x, interior.y-center.y) - radius},
		TestCase{p0, 0.0},
		TestCase{p1, 0.0}}

	// Crap is not normalized. Is of limited use..
	for _, test := range testCases {
		AssertFloatEqual(t, test.expectedDistance, arc.EuclideanDistanceTo(test.p))
	}
}
