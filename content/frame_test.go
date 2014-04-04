package content

import (
	"code.google.com/a/google.com/p/gojira/content/dom"
	"github.com/rjkroege/wikitools/testhelpers"
	"image"
	"testing"
)

func Test_FindElementAtPoint(t *testing.T) {
	f := NewFrame()
	testhelpers.AssertInt(t, 0, len(f.document))

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

	e, v = f.FindElementAtPoint(image.Pt(30-dom.QUAD_ELEMENT_DX, 10+dom.QUAD_ELEMENT_DY))
	if e == nil || *e != f.document[0] {
		t.Errorf("point 30,10 in %+v but found element is %+v", f.document[0], e)
	}
	if v != 3 {
		t.Errorf("desired vertex 3 but didn't get it %d", v)
	}
}
