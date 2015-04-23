// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#import "goview.h"

// To support the execution of foreign code (C, C++, ObjC etc.), Go
// uses the cgo tool to automatically generate interfaces for calling
// into foreign code and back into go. The interface for C -> Go is in 
// _cgo_export.h. It needs to be included here so that this ObjC file
// can call functions exported from Go.
#import "_cgo_export.h"

@implementation GoView

- (void)mouseDown:(NSEvent *)theEvent {
	if (!mouseDown(theEvent)) {
		[[self nextResponder] mouseDown:theEvent];
	}
}

- (void)mouseUp:(NSEvent *)theEvent {
	if (!mouseUp(theEvent)) {
		[[self nextResponder] mouseUp:theEvent];
	}
}

- (BOOL)acceptsFirstResponder {
    return YES;
}

- (id)initWithFrame:(NSRect)frame {
    self = [super initWithFrame:frame];
    if (self) {
	NSLog(@"set up the go structs");
       }
    return self;
}

@end

