package window

import (
	"github.com/rjkroege/wikitools/testhelpers"
	"testing"
)

func Test_windowCreation(t *testing.T) {
	w := new(Window)
	testhelpers.AssertInt(t, 0, int(w.width))
	testhelpers.AssertInt(t, 0, int(w.height))
	testhelpers.AssertInt(t, 0, int(w.width))
	testhelpers.AssertInt(t, 0, int(w.pointer.x))
	testhelpers.AssertInt(t, 0, int(w.pointer.y))
	testhelpers.AssertInt(t, 0, int(w.pointer.buttonmask))
}

func Test_onMouseBtn(t *testing.T) {
	w := new(Window)

	testhelpers.AssertInt(t, 1, MOUSE_BUTTON_LEFT)
	testhelpers.AssertInt(t, 2, MOUSE_BUTTON_MIDDLE)
	testhelpers.AssertInt(t, 4, MOUSE_BUTTON_RIGHT)

	w.onMouseBtn(0, 1)
	testhelpers.AssertInt(t, int(w.pointer.buttonmask), MOUSE_BUTTON_LEFT)
	w.onMouseBtn(0, 0)
	testhelpers.AssertInt(t, int(w.pointer.buttonmask), MOUSE_BUTTON_NONE)

	w.onMouseBtn(1, 1)
	testhelpers.AssertInt(t, int(w.pointer.buttonmask), MOUSE_BUTTON_MIDDLE)
	w.onMouseBtn(1, 0)
	testhelpers.AssertInt(t, int(w.pointer.buttonmask), MOUSE_BUTTON_NONE)

	w.onMouseBtn(2, 1)
	testhelpers.AssertInt(t, int(w.pointer.buttonmask), MOUSE_BUTTON_RIGHT)
	w.onMouseBtn(2, 0)
	testhelpers.AssertInt(t, int(w.pointer.buttonmask), MOUSE_BUTTON_NONE)
}

func Test_onMousePos(t *testing.T) {
	w := new(Window)
	w.onMousePos(1, 3)
	testhelpers.AssertInt(t, 1, int(w.pointer.x))
	testhelpers.AssertInt(t, 3, int(w.pointer.y))

}
