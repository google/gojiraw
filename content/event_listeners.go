// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package content

import (
	"image"
	"log"
)

// | These together to record which button is up or down.
// In particular: the buttons value will be some or of the following.
const (
	MOUSE_BUTTON_NONE = 0
)

// Reset iota.
const (
	MOUSE_BUTTON_LEFT = 1 << iota
	MOUSE_BUTTON_MIDDLE
	MOUSE_BUTTON_RIGHT
)

const (
	EVD_NON     = iota // There is no handler.
	EVD_PREVDEF        // A handler exists and it wishes the default action for the event to be suppressed.
	EVD_DEF            // A handler exists and it wants the default action.
)

// Anything Frame-like entity capable of receiving events
// needs to implement this interface.
type EventHandler interface {
	// Corresponds to JS registered with onmousedown
	Mousedown(pt image.Point, button, buttons uint32) uint32

	// Corresponds to JS registered with onmouseup
	Mouseup(pt image.Point, button, buttons uint32) uint32

	// Corresponds to JS registered with onmousemove
	Mousemove(pt image.Point, buttons uint32) uint32

	// Corresponds to JS registered with onmousewheel
	Wheel(pt image.Point, buttons uint32, dx, dy, dz float32) uint32
}

// These are "event listeners": functionality that really
// ought to belong in the JavaScript code of the browser.
// As of yet, this code is busily running here though written in
// Go. Writing these handlers will help define the interfaces that
// we need to schlep events from browser to the v8 thread.
//
// General scheme: in a future with JavaScript, the content side (i.e. main thread)
// for these handlers would receive the event bundle and call into v8
// On the browser side, the messaging proxy would handle this.
//
func (f *Frame) Mousedown(pt image.Point, button, buttons uint32) uint32 {
	log.Printf("OnMouseDown")

	e, v := f.FindElementAtPoint(pt)
	if e != nil && v > -1 {
		f.StartMouseDownMode(pt, e, v)
	}
	return EVD_PREVDEF
}

func (f *Frame) Mouseup(pt image.Point, button, buttons uint32) uint32 {
	if button == 0 && f.mouseDown {
		f.EndMouseDownMode()
	} else if button == 0 {
		f.AddElement(pt)
	}
	return EVD_PREVDEF
}

func (f *Frame) Mousemove(pt image.Point, buttons uint32) uint32 {
	if buttons&MOUSE_BUTTON_LEFT == MOUSE_BUTTON_LEFT {
		f.InMouseDownMode(pt)
	} else {
		e, v := f.FindElementAtPoint(pt)
		f.MouseOver(e, v)
	}
	return EVD_PREVDEF
}

func (f *Frame) Wheel(pt image.Point, buttons uint32, dx, dy, dz float32) uint32 {
	// The container for the frame scrolls.
	return EVD_NON
}
