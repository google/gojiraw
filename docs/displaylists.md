# Enso
`Enso` is a single simple drawable entity like a circle. The word is
from Zen Buddhist calligraphy and means (literally) circle. An Enzo
combines geometry and paint into an opaque object that has membership
in an Inkings.

```
type Enso uint32
```

As with a calligraphic strokes, an Enso is immutable once stroked into
an Inkings. Enzo immutability is intended to permit the library to
take any combination of deferred and eager rasterization of the Enzo
recorded into a specific Inkings.

For example, a ideal Gorjiraw implementation would convert a sequence
of `Enzo` into an array of floating point numbers and a shader
program. Redrawing the `Enzo` would proceed in the best case
GPU-accelerated fashion: call GPU library with pointer to float array
and pre-compiled shader program.

# Inkings
Inkings are a mutable record of multiple Enzo and other Inkings.
Inkings form a DAG. Inkings can be rasterized to an Image and hence
presented to the Display.

My intent with Inkings is to provide a plug-compatible substitute for
an Image except that, per the non-blocking principle, the pixel
contents are available only as a promise. 9p Images are frequently the
source and destination for the same draw operation. To preserve
this capability, lazy rasterization and the mutability of the Inkings,
Inkings offer a `Clone` function.

A naive algorithm for rasterizing an Inkings is to postfix depth-first
traverse the Inkings graph, allocating an Image for each Inkings and
rasterizing its Enzo before combining Inkings into their parent Image
via texture draw.

```golang
type Inkings Interface {
	Repl() bool
	Bound() Rectangle
	Clipr() Rectangle

	Clone() <-channel Inkings
	Zero()

	Render(i Image, opts ...interface{}) <-channel Image
	Show(i Image, opts ...interface{}) <-channel Image

	// r generalizes to a Path. Worry about that later.
	DrawInkings(r Rectangle, src, mask *Inkings, m *Matrix, d Drawop) Enzo
	ReDrawInkings(eid Enzo, r Rectangle, src, mask *Inkings, m *Matrix, d Drawop) Enzo

	GlyphString(glyphs []Glyph, styles []byte, stylefonts []StyledFont, p Point, m Matrix) Point
}

func NewInkings() (Inkings, error) {
}
```

## NewInkings
`NewInkings` creates a new Inkings. Each Inkings is (conceivably)
infinite in extent and has its own coordinate system. A new empty
Inkings is transparent (alpha = 0) and therefore has no effect when
added to another Inkings.

## Clone
`Clone` returns a channel that will contain a deep clone of the target
`Inkings`. Modifying `target` while `Clone` is in progress may have
undefined behaviour. Clone is central to how to use the Gojiraw API:
maintain a `Inkings` *model*, `Clone` it and render that.

```golang

// Many elided details...
model := NewInkings()
var display := NewDisplay()

wins := make(channel Image, 2)
wins <- currentDisplay.NewImage( /* ... */)
wins <- currentDisplay.NewImage(/* ... */),

rc := display.RefreshDeadlineChannel()
var snapshot Inkings

ni var <-channel Image

for {
	select {
	case <-rc:
		snapshot := <- ink.Clone()
		win := <- wins
		// TODO(rjk): The channel needs to be passed
		// in so that Image objects can be reused more easily.
		ni = shown = <-snapshot.Show(win)
	case im := <- ni
		win <- im
	case e := <-NextEvent
		appUpdatesModelForEvent(model, e)
	}
}
```

## Render
`Render` asynchronously rasterizes the target `Inking` to `Image i` and
triggers the returned channel with an updated `Image` object that contains
only rendered `Images`.

The returned `Image` may not be "flat". The `Inkings` may request use
of overlays by containing children `Image` objects that have the
`OverlayCandidate` attribute set. Then, the children `Image` objects
of the returned `Image` will correspond to hardware planes.

Between invoking this call and before the channel is triggered, the caller
should not modify the `Image`.

Some `Inkings` in the DAG may need longer than others. Intermediate
`Image` objects can be multi-buffered to accomodate this. {>>This 
would seem to require that `Clone` be synchronous.<<}

## Show
`Show` combines rendering with presentation. It asynchronously performs the following steps

* rasterize the target `Inking` to `Image i`.
* on the next vblank after rasterization, scan from `Image i`
* on next vblank after that, trigger the channel with `Image i`

## Zero
`Zero` removes all the Enzos from the `Inkings`. 

## Bound
`Bound` returns the smallest `Rectangle` enclosing the entire
`Inkings`

## DrawInkings
The most basic drawing operation in Gojiraw. It records an Enzo
wrapping `src`, `mask`, `m`, `op` in Inkings `target` and returns an
Enzo identifier.

DrawInkings has the following semantics for rasterization:

* create Rectangle `r` in the coordinate system of `target`.

* interesect `r` with `target.Clipr` (and other clips that might exist.)

* `m` transforms `r` from the target's coordinate system into `src`.
Combining `r` and `m` permit picking out an arbitrary
rectangle from `src`. 

*  `m` also transforms `r` from the target's CS into `mask` in identical fashion.

*  if `src.Repl` is set, consider `src` to tile its plane and if `mask.Repl` is set, consider `mask` to tile its plane.

*  for each pixel location in `r`, combine `target` pixel with `src` pixel using alpha
provided by `mask` pixel using Porter-Duff composition specified by `op`. 

Note that `Image` objects are usable as `Inkings` so this method also
offers the capability of drawing `Image` objects.

## ReDrawInkings
Operates identically to `DrawInkings` except that rather than adding a new
Enzo to Inkings `target`, replaces the specified Enzo. 


# GlyphString
`GlyphString` efficiently draws runs of styled text.
The server does not line-break the provided run.

*  `glyphs` an array of glyphs. NB: glyphs are not necessarily runes in order to support
languages like Arabic and ligatures. It seems convenient that runes are 1-1 with glyphs
where possible. For example, glyph 0x41 could be `A`.
*  `styles` an array of 1 byte style identies, one per glyph or alternatively an array of
length 1 which selects the StyleFont used for every glyph.
*  'styledfonts' is an array of `StyleFont` objects. Each byte in `styles` chooses
the `StyleFont` used to draw the corresponding glyph in `glyphs`
*  `p` is a `Point` specifying where to position the origin corner of the glyph in
the coordinate system of the target `Inkings`.
*  'm' is a `Matrix` that specifies the how to advance to the next glyph position: next `p` 
is `p += m * glyph-bounds`.

Because no line breaking is done. `\n` is treated as white space. Per the
Plan9 `string` function, `GlyphString` "returns a `Point` that is the
position of the next character that would be drawn if the string were
longer."

Drawing proceeds as the following  "Go with vector operations" pseudo-code

```golang
// A conceptually workable but brutally slow definition of outline fonts.
// The outline origin has its natural meaning.
type StyleFont []Inkings

m1 := geometry.Translate(p)

for i, _ := range glyphs {
	im := styledfonts[styles[i]][glyphs[i]]
	DrawInkings(m1.TransformRect(im.Bound()), im, im, m1, drawop.DoverS)
	m1 = m1 + m * (im.Bound.Max - im.Bound.Min)
}

return /* translation in x and y from m1 */
```

*insert diagram*


## Discussion
We require:

*  efficient encoding of styled text runs
*  a way to control the advance so that runs can be drawn right to left or vertically
*  efficient font metrics access on both the client and the server.
*  mechanism for client to prepare and upload glyph representations. 

Lots of really tiny text runs need to be efficient. In particular,
this means not shipping the text font name, font description, style
name, size, etc. over again. For example, trying to use the existing
`<canvas>` `strokeText()` function to draw a styled text run like
"hello *there* gentle reader" entails something like this:

```javascript
x, y = ... // set to baseline corner
ctx.strokeText('hello ', x, y)
sz = ctx.measureText('hello ')
ctx.font = ... // string contstant for italic font
x += sz.width
ctx.strokeText('there', x, y)
sz = ctx.measureText('there')
x += sz.width
ctx.font = ... // default font string constant
ctx.strokeTest(' gentle reader', x, y, )
```
A naïve conversion of this into a series of display list commands results
in context adjustment data exceeding the size of the string payload.
Instead, Gojiraw provides a byte-per glyph picking from a current cache
of styles along with an optimized form for character spans of uniform style.

## String Measurement
We need a way to efficiently measure the size of a string where we do
layout. I want to leave layout *outside* of the Gojira server. So on
the client. One layout style does not fit all use cases. *Is this
worthy of discussion?*

Doing layout in the client requires that the client be able to
efficiently measure the size of glyphs. Efficiently advancing down a
text run to draw the next glyph in it requires the server to *also* be
able to efficiently measure the size of glyphs. The most efficient
measurement mechanism is a (hopefully) O(1) look up of the glyph and
style combination's dimensions in a map of some kind.

Font metrics can be large. It is desirable to not replicate this map
between server and client. Further, it is desirable to share font
metrics between clients. Both blink and the Plan9 font code presume
that the client is responsible for loading and preparing fonts.
Moreover, blink admits dynamically downloading webfonts into the
client.

I think we can satisfy all of these requirements like this:

*  Gojira server maintains a metric cache in shared memory
*  Gojira clients can obtain a read-only memory mapping of the metric cache
*  Gojira clients open font files, parse and validate them and convert them into metric and (TBD) arc or curve data. 
*  Client ships the prepared font data to Gojira
*  Server validates the presented data.
*  Server places *the metrics* in memory shared RO with clients.
*  The same library on client and server can read the metrics from shared memory.
*  Server preps the glyph rendering data in its font cache and prepares it for use.

*yet another picture here*
