#import <Cocoa/Cocoa.h>

@interface  GoView : NSView

- (void)mouseDown:(NSEvent *)theEvent;
- (void)mouseUp:(NSEvent *)theEvent;
- (BOOL)acceptsFirstResponder;
- (id)initWithFrame:(NSRect)frame;

@end

