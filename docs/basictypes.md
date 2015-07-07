This section describes the primtive types largely to flesh out some
of the other more interesting sections. See the source for the
actual definition.

# Geometry
Standard geometry definitions.

```go
package geometry

type Point struct {
	x float
	y float
}

type Rectangle struct {
	min Point
	max Point
}

type Matrix struct {
	// numbers...
}

func (m Matrix) TransformRect(r Rectangle) Rectangle
// Produces a Matrix that translates by p
func Translate(p Point) Matrix

```

# Drawop
Each of the compositing operations in Porter-Duff compositing.
Names and definitions per Plan9.

```go
package drawop

const (
	Clear 	= 0
	DoutS 	= 1 << iota
	SoutD 	= 2
	DinS	=  4
	SinD	= 8
	S		= SinD|SoutD
	SoverD	= SinD|SoutD|DoutS
	SatopD	= SinD|DoutS
	SxorD	= SoutD|DoutS
	D		= DinS|DoutS
	DoverS	= DinS|DoutS|SoutD
	DatopS	= DinS|SoutD
	DxorS	= DoutS|SoutD     /* == SxorD */
	Ncomp	= 12
)
```

# Glyph
Gojiraw deliberately excludes font shaping from its scope. Instead, it deals
in *glyphs*. 


```go
type Glyph int32
```

NB: a `glyph` is *not* a UTF `rune`. Gojiraw will attempt to create an
equivalence between valid runes and the letter shape associated with
the identical glyph. For example, rune 0x41 is a valid UTF code point
and also a valid glyph to reference the letter shape for `A` in the
current style. This convenience makes sense when font shaping is
unused as it permits treating most rune arrays as an identical glyph
array. The `utf8.ValidRune(rune)` method will determine if a given
glyph corresponds to a actual UTF code point.

