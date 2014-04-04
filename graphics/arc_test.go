// TODO: make a header
package graphics

import (
	"math"
	"testing"
)

const (
	FLOAT_TOL = 1.0e-5
)

func AssertFloatEqual(t *testing.T, expected float32, actual float32) {
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
	ATAN_TEST_CASES = []float32{-.99, 0, 0.0001, 0.99}
)

func TestSin2Atan(t *testing.T) {
	for _, theta := range ATAN_TEST_CASES {
		AssertFloatEqual(t, Sin2Atan(theta), float32(math.Sin(2.0*math.Atan(float64(theta)))))
	}
}

func TestCos2Atan(t *testing.T) {
	for _, theta := range ATAN_TEST_CASES {
		AssertFloatEqual(t, Cos2Atan(theta), float32(math.Cos(2.0*math.Atan(float64(theta)))))
	}
}

// Returns tan(0.5 * acos(d))
func TanHalfAcos(d float32) float32 {
	return float32(math.Tan(math.Acos(float64(d)) * 0.25))
}

func TestIsInWedge(t *testing.T) {
	p0 := &Dual{1.0, 0.0, 1.0}
	p1 := &Dual{2.0, 0.0, 1.0}
	center := &Dual{1.5, -0.5, 1.0}

	c0 := p0.Intersection(center)
	c1 := p1.Intersection(center)

	cos_d := c0.AngularDistanceTo(c1)
	a := Arc{p0, p1, float32(math.Tan(math.Acos(float64(cos_d)) * 0.25))}

	type TestCase struct {
		p        *Dual
		expected bool
	}

	testCases := []TestCase{
		{&Dual{1.5, 0.5, 1.0}, true},
		{&Dual{1.5, -0.2, 1.0}, true},
		{&Dual{1.5, -1.5, 1.0}, false},
		{&Dual{3.5, 0.5, 1.0}, false},
		{&Dual{-1.5, 0.5, 1.0}, false}}

	for i, test := range testCases {
		if a.IsInWedge(test.p) != test.expected {
			t.Errorf("Case %d Failed: p(%f, %f), is in = %t\n",
				i, test.p.x, test.p.y, a.IsInWedge(test.p))
		}
	}
}

func TestEuclideanDistanceTo(t *testing.T) {
	p0 := &Dual{1.0, 0.0, 1.0}
	p1 := &Dual{2.0, 0.0, 1.0}
	center := &Dual{1.5, -0.5, 1.0}

	c0 := p0.Intersection(center)
	c1 := p1.Intersection(center)

	cos_d := c0.AngularDistanceTo(c1)
	radius := float32(math.Hypot(float64(c0.x), float64(c0.y)))
	cos_d /= radius * radius
	arc := Arc{p0, p1, float32(math.Tan(math.Acos(float64(cos_d)) * 0.25))}
	apex := arc.Apex()

	type TestCase struct {
		p                *Dual
		expectedDistance float32
	}

	exterior := Dual{1.7872, 5.22, 1.0}
	interior := Dual{1.26, 0.0115, 1.0}

	testCases := []TestCase{
		{apex, 0.0},
		{&Dual{apex.x, apex.y + 1.0, 1.0}, 1.0},
		{&Dual{apex.x, apex.y - 0.2, 1.0}, -0.2},
		{&Dual{apex.x, apex.y + 10.0, 1.0}, 10.0},
		{&exterior, float32(math.Hypot(float64(exterior.x-center.x), float64(exterior.y-center.y))) - radius},
		{&interior, float32(math.Hypot(float64(interior.x-center.x), float64(interior.y-center.y))) - radius},
		{p0, 0.0},
		{p1, 0.0}}

	// Crap is not normalized. Is of limited use..
	for _, test := range testCases {
		AssertFloatEqual(t, test.expectedDistance, arc.EuclideanDistanceTo(test.p))
	}
}
