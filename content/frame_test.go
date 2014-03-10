package content

import (
	"github.com/rjkroege/wikitools/testhelpers"
	"image"
	"testing"
)

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
