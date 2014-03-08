// TODO make a header.
package main

import (
	"code.google.com/a/google.com/p/gojira/window"
	"fmt"
	"github.com/go-gl/glfw"
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
