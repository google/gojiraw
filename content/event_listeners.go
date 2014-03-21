package content

import (
	"image"
	"log"
)

const (
	EVD_NON = iota		// There is no handler.
	EVD_PREVDEF	// A handler exists and it wishes the default action for the event to be suppressed.
	EVD_DEF // A handler exists and it wants the default action.
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
	return EVD_NON
}

func (f *Frame) Mouseup(pt image.Point, button, buttons  uint32) uint32 {
	if button == 0 {
		f.AddElement(pt)	
	}
	return EVD_PREVDEF
}

func (f *Frame) Mousemove(pt image.Point, buttons  uint32) uint32 {
	f.MouseOver(f.FindElementAtPoint(pt))
	return EVD_PREVDEF
}

//func mousewheel(f *content.Frame, pt image.Point, ) {
//
//}


