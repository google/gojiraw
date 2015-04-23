// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"code.google.com/a/google.com/p/gojiraw/window"
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
