package window

import (
	"code.google.com/a/google.com/p/gojira/content"
	"image"
	"log"
)

// These are "event listeners": functionality that really
// ought to belong in the JavaScript code of the browser.
// As of yet, this code is busily running here though written in
// Go. Writing these handlers will help define the interfaces that
// we need to schlep events from browser to the v8 thread.

// Corresponds to JS registered with onmousedown
func mousedown(f *content.Frame, pt image.Point, button, buttons uint32) {
	log.Printf("OnMouseDown")	

}

// Corresponds to JS registered with onmouseup
func mouseup(f *content.Frame, pt image.Point, button, buttons  uint32) {
	if button == 0 {
		f.AddElement(pt)	
	}
}

//func mousewheel(f *content.Frame, pt image.Point, ) {
//
//}