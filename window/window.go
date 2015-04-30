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


package window

import (
	"image"
	"log"

	"code.google.com/a/google.com/p/gojiraw/content"
	"code.google.com/a/google.com/p/gojiraw/graphics"
	"github.com/go-gl/gl"

	glfw "github.com/go-gl/glfw3"
)

// Stores the current mouse pointer position.
type Mousepointer struct {
	x, y       int
	buttonmask uint32
}

// TODO: when this gets big, it might want to be a package.
type Window struct {
	width, height uint32
	frame         *content.Frame
	pointer       Mousepointer

	// Position of the Frame's origin in the viewport. Starts as 0, 0 where the
	// Window (aka viewport) will show the upper left corner of the content.Frame.
	// rjk presumes that positive y extends down the page and positive x increase
	// towards the right. As a result, x, y are <= 0.
	//
	// rjk:  T S cs = ws for Translation and Scale and cs is in content coordinates. Verify.
	x, y float32

	// Width, height of the frame. Must start the same size as the Window.
	fw, fh float32

	// TODO(rjkroege): When the mouse wheel handling is rational and actually
	// ships us deltas, then we can remove this code.
	previous_absolute_displacement int

	// Event handling interface
	ev content.EventHandler
}

func NewWindow(width int, height int) *Window {
	c := content.NewFrame()
	return &Window{uint32(width), uint32(height), c, Mousepointer{0, 0, 0}, 0.0, 0.0,
		float32(width), float32(height), 0, c}
}

func (window *Window) RunMessageLoop(w *glfw.Window, program *gl.Program) {
	gl.GetError()
	for !w.ShouldClose() {
		// TODO(rjkroege): full generality: provide the transform to bring the Frame into
		// Window coordinates and the width and height.
		window.fw, window.fh = window.frame.Draw(window.x, window.y, float32(window.width), float32(window.height), program)
		w.SwapBuffers()
		glfw.PollEvents()
	}
}

// Based on https://raw.github.com/go-gl/examples/master/glfw/simplewindow
func (window *Window) Open() {
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenglForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)

	// TODO: what's go style here? How do you get clang-format for go?
	// TODO(rjkroege): sizes should be uints
	glfwWindow, err := glfw.CreateWindow(int(window.width), int(window.height), "Testing", nil, nil)

	if err != nil {
		log.Panic(err)
	}

	defer glfwWindow.Destroy()

	glfwWindow.MakeContextCurrent()

	// Apparantly, this enables vsync?
	glfw.SwapInterval(1)

	gl.Init()

	glfwWindow.SetSizeCallback(func(_ *glfw.Window, w, h int) {
		window.onResize(w, h)
	})

	glfwWindow.SetMouseButtonCallback(func(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, _ glfw.ModifierKey) {
		window.onMouseBtn(button, action)
	})

	// TODO(rjkroege): Figure out how these events should be dealt with.
	// glfw.SetMouseWheelCallback(func(d int) {
	// 	window.onMouseWheel(d)
	// })

	// glfw.SetKeyCallback(func(k, s int) {
	// 	window.onKey(k, s)
	// })

	// glfw.SetCharCallback(func(k, s int) {
	// 	window.onChar(k, s)
	// })

	// TODO(vollick): Passing around one program like this is a stopgap. We
	// should really be initializing our shader library here.
	program := graphics.CreateDefaultShaders()
	defer program.Delete()

	glfwWindow.SetCursorPositionCallback(func(_ *glfw.Window, x, y float64) {
		window.onMousePos(int(x), int(y))
	})

	window.RunMessageLoop(glfwWindow, &program)
}

func (window *Window) onResize(w, h int) {
	window.width = uint32(w)
	window.height = uint32(h)
	log.Printf("Resize %d %d", window.width, window.height)
}

func (window *Window) onMouseBtn(button glfw.MouseButton, action glfw.Action) {
	state := uint32(action)
	if button < 0 || button > 31 || state < 0 || state > 1 {
		log.Fatal("button/state values from glfw are silly: ", button, state)
	}

	b := uint32(button)

	// log.Printf("onMouseButton. state = %d, button = %d", state, b)
	p := window.mousePositionInFrame()
	if state == 1 {
		window.pointer.buttonmask |= 1 << b
		window.ev.Mousedown(p, b, window.pointer.buttonmask)
	} else {
		window.pointer.buttonmask &= ^(1 << b)
		window.ev.Mouseup(p, b, window.pointer.buttonmask)
	}
}

// Get the position of the mouse in the coordinates of the frame.
// The sub is right for y-axis but since we can't scroll x-axis.
func (w *Window) mousePositionInFrame() image.Point {
	p := image.Pt(w.pointer.x, w.pointer.y)
	p = p.Sub(image.Pt(int(w.x), int(w.y)))
	return p
}

// Mini-essay on the way of scrolling. In scrolling, there are two
// conceptual entities: the viewport and the contents.
//
// I maintain that scrolling should be handled by the viewport. The
// contents should provide the capability of drawing an arbitrary
// subrectangle of itself. But the choice of subrectangle is managed by
// the viewport.
//
// The principle here is to handle the event as "close" to where it
// arrives as possible. In particular: scrolling will happen in the
// browser.
func (window *Window) onMouseWheel(absolute_displacement int) {
	log.Printf("mouse wheel abs: %d\n", absolute_displacement)

	// TODO(rjkroege): displacement is stupid because of glfw issue. Fix.
	// delta is not a delta. I am instead interested in the difference.
	delta := absolute_displacement - window.previous_absolute_displacement
	window.previous_absolute_displacement = absolute_displacement
	// log.Printf("mouse wheel delta: %d\n", delta)

	dx := float32(delta)
	// Scrolling is currently only in one dimension because of glfw limitation.
	// TODO(rjkroege): enable two-dimensional scrolling.

	p := window.mousePositionInFrame()
	if window.ev.Wheel(p, window.pointer.buttonmask, dx, 0, 0) != content.EVD_PREVDEF {
		// window.x = content.Max(window.width - window.fw, window.x + float32(delta.x))
		// Consider putting Max in some kind of base-like class.
		window.y = graphics.MaxF(float32(window.height)-window.fh, window.y+dx)

		// Can't scroll before the start of the content area (need x limit in the future.)
		if window.y > 0.0 {
			window.y = 0.0
		}

		log.Printf("window.y is max of %f %f\n", float32(window.height)-window.fh, window.y+float32(delta))

		if window.y >= 0.0 {
			log.Print("i have the scrolling sign backwards")
		}
	}
}

// TODO(rjkroege): Add support for delivering of key events.
func (window *Window) onKey(key, state int) {
	log.Printf("key: %d, %d\n", key, state)
}

func (window *Window) onChar(key, state int) {
	log.Printf("char: %d, %d\n", key, state)
}

func (window *Window) onMousePos(x, y int) {
	window.pointer.x = x
	window.pointer.y = y

	p := window.mousePositionInFrame()

	// TODO(rjkroege): filter/collapse/schedule the events as desirable.
	window.ev.Mousemove(p, window.pointer.buttonmask)
}
