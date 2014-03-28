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
}

/*
func Test_FindElementAtPoint(t *testing.T) {
	f := NewFrame()
	testhelpers.AssertInt(t, 0, len(f.displaylist))

	p := image.Pt(10, 10)

	e := f.FindElementAtPoint(p)
	if e != nil {
		t.Errorf("empty display list can't have element in it")
	}

	f.AddElement(image.Pt(30, 10))
	f.AddElement(image.Pt(10, 30))

	e = f.FindElementAtPoint(p)
	if e != nil {
		t.Errorf("despite a point outside the elements, we hit anyway: %+v", e)
	}

	e = f.FindElementAtPoint(image.Pt(30, 10))
	if e == nil || *e != f.displaylist[0] {
		t.Errorf("point 30,10 in %+v but found element is %+v", f.displaylist[0], e)
	}
}
*/
