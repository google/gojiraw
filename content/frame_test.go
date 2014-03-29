package content

import (
	"github.com/rjkroege/wikitools/testhelpers"
	"code.google.com/a/google.com/p/gojira/geometry"
	"image"
	"testing"
)

func Test_FindVertex(t *testing.T) {
	e := new(Element)
	e.init(image.Pt(10, 15))

	v := e.FindVertex(geometry.Pointf{5, 5})
	testhelpers.AssertInt(t, -1, v)

	v = e.FindVertex(geometry.Pointf{10 - dX, 15 - dY})
	testhelpers.AssertInt(t, 0, v)

	v = e.FindVertex(geometry.Pointf{10 - dX - dH, 15 - dY - dH})
	testhelpers.AssertInt(t, 0, v)

	v = e.FindVertex(geometry.Pointf{10 - dX - dH - .01, 15 - dY - dH})
	testhelpers.AssertInt(t, -1, v)

}

func Test_FindElementAtPoint(t *testing.T) {
	f := NewFrame()
	testhelpers.AssertInt(t, 0, len(f.displaylist))

	p := image.Pt(10, 10)

	e, _ := f.FindElementAtPoint(p)
	if e != nil {
		t.Errorf("empty display list can't have element in it")
	}

	f.AddElement(image.Pt(30, 10))
	f.AddElement(image.Pt(10, 30))

	e, v := f.FindElementAtPoint(p)
	if e != nil {
		t.Errorf("despite a point outside the elements, we hit anyway: %+v", e)
	}

	e, v = f.FindElementAtPoint(image.Pt(30 - dX, 10 + dY))
	if e == nil || *e != f.displaylist[0] {
		t.Errorf("point 30,10 in %+v but found element is %+v", f.displaylist[0], e)
	}
	if v != 3 {
		t.Errorf("desired vertex 3 but didn't get it %d", v)
	}
}
