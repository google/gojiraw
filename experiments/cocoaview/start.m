// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#import <Cocoa/Cocoa.h>
#import "goview.h"

int
StartApp(void) {
    [NSAutoreleasePool new];
    [NSApplication sharedApplication];
    [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];

    id menubar = [[NSMenu new] autorelease];
    id appMenuItem = [[NSMenuItem new] autorelease];
    [menubar addItem:appMenuItem];
    [NSApp setMainMenu:menubar];

    id appMenu = [[NSMenu new] autorelease];
    id appName = [[NSProcessInfo processInfo] processName];
    id quitTitle = [@"Quit " stringByAppendingString:appName];
    id quitMenuItem = [[[NSMenuItem alloc] initWithTitle:quitTitle
        action:@selector(terminate:) keyEquivalent:@"q"] autorelease];
    [appMenu addItem:quitMenuItem];
    [appMenuItem setSubmenu:appMenu];

    NSRect winFrame = NSMakeRect(0, 0, 200, 200);

    NSWindow* window = [[[NSWindow alloc]
        initWithContentRect:winFrame
        styleMask:NSTitledWindowMask
        backing:NSBackingStoreBuffered defer:NO]
        autorelease];
    [window cascadeTopLeftFromPoint:NSMakePoint(20,20)];
    [window setTitle:appName];
    [window makeKeyAndOrderFront:nil];

    NSView* view = [[GoView alloc]
        initWithFrame:[window contentRectForFrameRect: winFrame]];
    [window setContentView: view];

    [NSApp activateIgnoringOtherApps:YES];
    [NSApp run];
    return 0;
}
