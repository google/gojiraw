// Copyright 2014 The Gojiraw Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

