// TODO make a header.
package main

import (
	"log"

	"code.google.com/a/google.com/p/gojira/window"
	glfw "github.com/go-gl/glfw3"
)

func main() {
	// Initialize glfw.
	if !glfw.Init() {
		log.Panic("Couldn't initialize glfw3")
		return
	}

	defer glfw.Terminate()

	width := 256
	height := 256
	window := window.NewWindow(width, height)

	// This will block until the window is closed.
	window.Open()
}
