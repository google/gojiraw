// TODO make a header.
package main

import (
	"fmt"
	"github.com/go-gl/glfw"
	"github.com/ianvollick/gojira/window"
	"os"
)

func main() {
	// Initialize glfw.
	err := glfw.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return
	}

	defer glfw.Terminate()

	width := 256
	height := 256
	window := window.NewWindow(width, height)

	// This will block until the window is closed.
	window.Open()
}
