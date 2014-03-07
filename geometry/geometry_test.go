// TODO: make a header
package geometry

import (
	"math"
	"testing"
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
