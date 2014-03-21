// TODO make a header.
package window

import (
	"code.google.com/a/google.com/p/gojira/content"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"image"
	"log"
)

// | These together to record which button is up or down.
// TODO(rjkroege): Bring this into alignment with webkit events.
const (
	MOUSE_BUTTON_NONE = 0
)

// Reset iota.
const (
	MOUSE_BUTTON_LEFT = 1 << iota
	MOUSE_BUTTON_MIDDLE
	MOUSE_BUTTON_RIGHT
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

func (window *Window) RunMessageLoop() {
	should_continue := true
	for should_continue {
		// TODO(rjkroege): full generality: provide the transform to bring the Frame into
		// Window coordinates and the width and height.
		window.fw, window.fh = window.frame.Draw(window.x, window.y, float32(window.width), float32(window.height))
		glfw.SwapBuffers()
		should_continue = glfw.WindowParam(glfw.Opened) == 1
	}
}

// Based on https://raw.github.com/go-gl/examples/master/glfw/simplewindow
func (window *Window) Open() {
	red_bits := 8
	green_bits := 8
	blue_bits := 8
	alpha_bits := 8
	depth_bits := 0
	stencil_bits := 0
	mode := glfw.Windowed

	// TODO: what's go style here? How do you get clang-format for go?
	// TODO(rjkroege): sizes should be uints
	err := glfw.OpenWindow(int(window.width), int(window.height),
		red_bits, green_bits, blue_bits, alpha_bits,
		depth_bits, stencil_bits,
		mode)

	if err != nil {
		// TODO: error?
	}

	defer glfw.CloseWindow()

	// Apparantly, this enables vsync?
	glfw.SetSwapInterval(1)
	glfw.SetWindowTitle("Hello, World!")
	glfw.SetWindowSizeCallback(func(w, h int) {
		window.onResize(w, h)
	})

	glfw.SetMouseButtonCallback(func(b, s int) {
		window.onMouseBtn(b, s)
	})

	glfw.SetMouseWheelCallback(func(d int) {
		window.onMouseWheel(d)
	})

	glfw.SetKeyCallback(func(k, s int) {
		window.onKey(k, s)
	})

	glfw.SetCharCallback(func(k, s int) {
		window.onChar(k, s)
	})

	glfw.SetMousePosCallback(func(x, y int) {
		window.onMousePos(x, y)
	})

	window.RunMessageLoop()
}

func (window *Window) onResize(w, h int) {
	gl.DrawBuffer(gl.FRONT_AND_BACK)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Viewport(0, 0, w, h)
	gl.Ortho(0, float64(w), float64(h), 0, -1.0, 1.0)
	gl.ClearColor(1, 1, 1, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	log.Printf("Resize %d %d", window.width, window.height)
}

func (window *Window) onMouseBtn(button, state int) {
	if button < 0 || button > 31 || state < 0 || state > 1 {
		log.Fatal("button/state values from glfw are silly: ", button, state)
	}

	b := uint32(button)
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
func (w *Window) mousePositionInFrame() image.Point {
	p := image.Pt(w.pointer.x, w.pointer.y)
	p.Add(image.Pt(int(w.x), int(w.y)))
	return p
}

// We have a problem. Mouse wheel events have an X and a Y and a
// whole bunch of other things. But here, mousewheels are one-dimensional.
// TODO(rjkroege): plumb a richer wheel event type to here.
// It would appear that delta isn't a delta..
func (window *Window) onMouseWheel(absolute_displacement int) {
	log.Printf("mouse wheel abs: %d\n", absolute_displacement)

	// delta is not a delta. I am instead interested in the difference.
	delta := absolute_displacement - window.previous_absolute_displacement
	window.previous_absolute_displacement = absolute_displacement
	log.Printf("mouse wheel delta: %d\n", delta)

	// TODO(rjkroege): Is this still signed correctly?
	//  delta = -delta

	// Scrolling is currently only in one dimension because of glfw limitation.
	// TODO(rjkroege): enable two-dimensional scrolling.
	// window.x = content.Max(window.width - window.fw, window.x + float32(delta.x))
	window.y = content.Max(float32(window.height)-window.fh, window.y+float32(delta))
	// Can't scroll before the start of the content area.
	if window.y > 0.0 {
		window.y = 0.0
	}

	log.Printf("window.y is max of %f %f\n", float32(window.height)-window.fh, window.y+float32(delta))

	if window.y >= 0.0 {
		log.Print("i have the scrolling sign backwards")
	}
}

func (window *Window) onKey(key, state int) {
	log.Printf("key: %d, %d\n", key, state)
}

func (window *Window) onChar(key, state int) {
	log.Printf("char: %d, %d\n", key, state)
}

func (window *Window) onMousePos(x, y int) {
	// log.Printf("mouse motion %d %d\n", x, y)
	window.pointer.x = x
	window.pointer.y = y

	p := window.mousePositionInFrame()

	// TODO(rjkroege): filter/collapse/schedule the events as desirable.
	window.ev.Mousemove(p, window.pointer.buttonmask)
}
