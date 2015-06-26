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

package content

import (
	"github.com/google/gojiraw/content/dom"
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
