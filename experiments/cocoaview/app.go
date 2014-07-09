package cocoaview

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#include "start.h"
*/
import "C"

func StartApp() {
	C.StartApp()
}
