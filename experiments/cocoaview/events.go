package cocoaview

// Example: bridging Objective-C into Go.

import (
	"log"
	"unsafe"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>
#include <Foundation/Foundation.h>

// Make into void* to pass around.
// It ought to be possible to persuade the Go compiler to build
// this code dynamically somehow.
inline static NSPoint
getLocation(void* event) {
	return [(NSEvent*)event locationInWindow];
}

*/
import "C"

// Return true if the event is handled. Callable from C.
//export mouseUp
func mouseUp(event unsafe.Pointer) bool {
	pt := C.getLocation(event)
	log.Printf("mouseUp at %f %f\n", pt.x, pt.y)
	return false
}

// Return true if the event is handled. Callable from C.
//export mouseDown
func mouseDown(event unsafe.Pointer) bool {
	pt := C.getLocation(event)
	log.Printf("mouseDown at %f %f\n", pt.x, pt.y)
	return false
}

