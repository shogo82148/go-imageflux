package imageflux

import (
	"fmt"
	"image"
	"image/color"
	"net/url"
	"strconv"
)

// Config is configure of image.
type Config struct {
	// Width is width in pixel of the scaled image.
	Width int

	// Height is height in pixel of the scaled image.
	Height int

	// DisableEnlarge disables enlarge.
	DisableEnlarge bool

	// AspectMode is aspect mode.
	AspectMode AspectMode

	// DevicePixelRatio is a scale factor of device pixel ratio.
	// If DevicePixelRatio is 0, it is ignored.
	DevicePixelRatio float64

	// InputClip is a position in pixel of clipping area.
	// This is used for the input image.
	InputClip image.Rectangle

	// InputClipRatio is a position in ratio of clipping area.
	// The coordinates of the rectangle are divided by ClipMax.X or ClipMax.Y.
	// This is used for the input image.
	InputClipRatio image.Rectangle

	// InputOrigin is the position of the input image origin.
	InputOrigin Origin

	// OutputClip is a position in pixel of clipping area.
	// This is used for the output image.
	OutputClip image.Rectangle

	// Clip is an alias of OutputClip.
	// If both Clip and OutputClip are set, OutputClip is used.
	//
	// Deprecated: Use OutputClip instead.
	Clip image.Rectangle

	// OutputClipRatio is a position in ratio of clipping area.
	// The coordinates of the rectangle are divided by ClipMax.X or ClipMax.Y.
	OutputClipRatio image.Rectangle

	// ClipRatio is an alias of OutputClipRatio.
	// If both ClipRatio and OutputClipRatio are set, OutputClipRatio is used.
	//
	// Deprecated: Use OutputClipRatio instead.
	ClipRatio image.Rectangle

	// OutputOrigin is the position of the output image origin.
	OutputOrigin Origin

	// ClipMax is the denominators of ClipRatio.
	ClipMax image.Point

	// Origin is the position of the image origin.
	Origin Origin

	// Background is background color.
	Background color.Color

	// InputRotate rotates the image before processing.
	InputRotate Rotate

	// OutputRotate rotates the image after processing.
	OutputRotate Rotate

	// OutputRotate rotates the image after processing.
	// This is an alias of OutputRotate.
	// If both Rotate and OutputRotate are set, OutputRotate is used.
	//
	// Deprecated: Use OutputRotate instead.
	Rotate Rotate

	// Through is a format to pass through.
	Through Through

	// Overlay Parameters.
	Overlays []Overlay

	// Output Parameters.
	Format              Format
	Quality             int
	DisableOptimization bool

	// Unsharp configures unsharp mask.
	Unsharp Unsharp

	// Blur configures blur.
	Blur Blur

	// GrayScale converts to gray scale.
	// 0 means no conversion and 100 means full conversion.
	GrayScale int

	// Sepia converts to sepia.
	// 0 means no conversion and 100 means full conversion.
	Sepia int

	// Brightness adjusts brightness.
	// The value set in Brightness plus 100 is actually used.
	Brightness int

	// Contrast adjusts contrast.
	// The value set in Contrast plus 100 is actually used.
	Contrast int

	// Invert inverts the image if it is true.
	Invert bool
}

// Overlay is the configure of an overlay image.
type Overlay struct {
	// URL is an url for overlay image.
	URL string

	// Width is width in pixel of the scaled image.
	Width int

	// Height is height in pixel of the scaled image.
	Height int

	// DisableEnlarge disables enlarge.
	DisableEnlarge bool

	// AspectMode is aspect mode.
	AspectMode AspectMode

	// InputClip is a position in pixel of clipping area.
	// This is used for the input image.
	InputClip image.Rectangle

	// InputClipRatio is a position in ratio of clipping area.
	// The coordinates of the rectangle are divided by ClipMax.X or ClipMax.Y.
	// This is used for the input image.
	InputClipRatio image.Rectangle

	// InputOrigin is the position of the input image origin.
	InputOrigin Origin

	// OutputClip is a position in pixel of clipping area.
	// This is used for the output image.
	OutputClip image.Rectangle

	// Clip is an alias of OutputClip.
	// If both Clip and OutputClip are set, OutputClip is used.
	//
	// Deprecated: Use OutputClip instead.
	Clip image.Rectangle

	// OutputClipRatio is a position in ratio of clipping area.
	// The coordinates of the rectangle are divided by ClipMax.X or ClipMax.Y.
	OutputClipRatio image.Rectangle

	// ClipRatio is an alias of OutputClipRatio.
	// If both ClipRatio and OutputClipRatio are set, OutputClipRatio is used.
	//
	// Deprecated: Use OutputClipRatio instead.
	ClipRatio image.Rectangle

	// OutputOrigin is the position of the output image origin.
	OutputOrigin Origin

	// ClipMax is the denominators of ClipRatio.
	ClipMax image.Point

	// Origin is the position of the image origin.
	Origin Origin

	// Background is background color.
	Background color.Color

	// InputRotate rotates the image before processing.
	InputRotate Rotate

	// OutputRotate rotates the image after processing.
	OutputRotate Rotate

	// OutputRotate rotates the image after processing.
	// This is an alias of OutputRotate.
	// If both Rotate and OutputRotate are set, OutputRotate is used.
	//
	// Deprecated: Use OutputRotate instead.
	Rotate Rotate

	// Offset is an offset in pixel of overlay image.
	Offset image.Point

	// OffsetRatio is an offset in ratio of overlay image.
	// The coordinates of the rectangle are divided by OffsetMax.X or OffsetMax.Y.
	OffsetRatio image.Point

	// OffsetMax is the denominators of OffsetRatio.
	OffsetMax image.Point

	// OverlayOrigin is the position of the overlay image origin.
	OverlayOrigin Origin
}

// Unsharp is an unsharp filter config.
type Unsharp struct {
	Radius    int
	Sigma     float64
	Gain      float64
	Threshold float64
}

func (u Unsharp) append(buf []byte) []byte {
	buf = strconv.AppendInt(buf, int64(u.Radius), 10)
	buf = append(buf, 'x')
	buf = strconv.AppendFloat(buf, u.Sigma, 'f', -1, 64)
	if u.Gain != 0 && u.Threshold != 0 {
		buf = append(buf, '+')
		buf = strconv.AppendFloat(buf, u.Gain, 'f', -1, 64)
		buf = append(buf, '+')
		buf = strconv.AppendFloat(buf, u.Threshold, 'f', -1, 64)
	}
	return buf
}

// Blur is a blur config.
type Blur struct {
	Radius int
	Sigma  float64
}

func (b Blur) append(buf []byte) []byte {
	buf = strconv.AppendInt(buf, int64(b.Radius), 10)
	buf = append(buf, 'x')
	buf = strconv.AppendFloat(buf, b.Sigma, 'f', -1, 64)
	return buf
}

// AspectMode is aspect mode.
type AspectMode int

const (
	// AspectModeDefault is the default value of aspect mode.
	AspectModeDefault AspectMode = iota

	// AspectModeScale holds the the aspect ratio of the input image,
	// and scales to fit in the specified size.
	AspectModeScale

	// AspectModeForceScale ignores the aspect ratio of the input image.
	AspectModeForceScale

	// AspectModeCrop holds the the aspect ratio of the input image,
	// and crops the image.
	AspectModeCrop

	// AspectModePad holds the the aspect ratio of the input image,
	// and fills the unfilled portion with the specified background color.
	AspectModePad
)

// Origin is the origin.
type Origin int

const (
	// OriginDefault is default origin.
	OriginDefault Origin = iota

	// OriginTopLeft is top-left
	OriginTopLeft

	// OriginTopCenter is top-center
	OriginTopCenter

	// OriginTopRight is top-right
	OriginTopRight

	// OriginMiddleLeft is middle-left
	OriginMiddleLeft

	// OriginMiddleCenter is middle-center
	OriginMiddleCenter

	// OriginMiddleRight is middle-right
	OriginMiddleRight

	// OriginBottomLeft is bottom-left
	OriginBottomLeft

	// OriginBottomCenter is bottom-center
	OriginBottomCenter

	// OriginBottomRight is bottom-right
	OriginBottomRight
)

func (o Origin) String() string {
	switch o {
	case OriginDefault:
		return "default"
	case OriginTopLeft:
		return "top-left"
	case OriginTopCenter:
		return "top-center"
	case OriginTopRight:
		return "top-right"
	case OriginMiddleLeft:
		return "middle-left"
	case OriginMiddleCenter:
		return "middle-center"
	case OriginMiddleRight:
		return "middle-right"
	case OriginBottomLeft:
		return "bottom-left"
	case OriginBottomCenter:
		return "bottom-center"
	case OriginBottomRight:
		return "bottom-right"
	}
	return ""
}

// Format is the format of the output image.
type Format string

const (
	// FormatAuto encodes the image by the same format with the input image.
	FormatAuto Format = "auto"

	// FormatJPEG encodes the image as a JPEG.
	FormatJPEG Format = "jpg"

	// FormatPNG encodes the image as a PNG.
	FormatPNG Format = "png"

	// FormatGIF encodes the image as a GIF.
	FormatGIF Format = "gif"

	// FormatWebPFromJPEG encodes the image as a WebP.
	// The input image should be a JPEG.
	FormatWebPFromJPEG Format = "webp:jpeg"

	// FormatWebPFromPNG encodes the image as a WebP.
	// The input image should be a PNG.
	FormatWebPFromPNG Format = "webp:png"
)

func (f Format) String() string {
	return string(f)
}

// Rotate rotates the image.
type Rotate int

const (
	// RotateDefault is the default value of Rotate. It is same as RotateTopLeft.
	RotateDefault Rotate = iota

	// RotateTopLeft does not anything.
	RotateTopLeft

	// RotateTopRight flips the image left and right.
	RotateTopRight

	// RotateBottomRight rotates the image 180 degrees.
	RotateBottomRight

	// RotateBottomLeft flips the image upside down.
	RotateBottomLeft

	// RotateLeftTop mirrors the image around the diagonal axis.
	RotateLeftTop

	// RotateRightTop rotates the image left 90 degrees.
	RotateRightTop

	// RotateRightBottom rotates the image 180 degrees and mirrors the image around the diagonal axis.
	RotateRightBottom

	// RotateLeftBottom rotates the image right 90 degrees.
	RotateLeftBottom

	// RotateAuto parses the Orientation of the Exif information and rotates the image.
	RotateAuto Rotate = -1
)

func (r Rotate) String() string {
	switch r {
	case RotateDefault:
		return "default"
	case RotateTopLeft:
		return "top-left"
	case RotateTopRight:
		return "top-right"
	case RotateBottomRight:
		return "bottom-right"
	case RotateBottomLeft:
		return "bottom-left"
	case RotateLeftTop:
		return "left-top"
	case RotateRightTop:
		return "right-top"
	case RotateRightBottom:
		return "right-bottom"
	case RotateLeftBottom:
		return "left-bottom"
	case RotateAuto:
		return "auto"
	}
	return ""
}

// Through is an image format list for skipping converting.
type Through int

const (
	// ThroughJPEG skips converting JPEG images.
	ThroughJPEG Through = 1 << iota

	// ThroughPNG skips converting PNG images.
	ThroughPNG

	// ThroughGIF skips converting GIF images.
	ThroughGIF
)

func (t Through) String() string {
	var buf [12]byte
	return string(t.append(buf[:]))
}

func (t Through) append(buf []byte) []byte {
	if t == 0 {
		return buf
	}
	if (t & ThroughJPEG) != 0 {
		buf = append(buf, "jpg:"...)
	}
	if (t & ThroughPNG) != 0 {
		buf = append(buf, "png:"...)
	}
	if (t & ThroughGIF) != 0 {
		buf = append(buf, "gif:"...)
	}
	return buf[:len(buf)-1]
}

func (c *Config) String() string {
	if c == nil {
		return ""
	}
	buf := bufPool.Get().(*[]byte)
	*buf = c.append((*buf)[:0])
	str := string(*buf)
	bufPool.Put(buf)
	return str
}

func (c *Config) append(buf []byte) []byte {
	var zr image.Rectangle
	var zp image.Point
	if c == nil {
		return buf
	}

	l := len(buf)
	if c.Width != 0 {
		buf = append(buf, 'w', '=')
		buf = strconv.AppendInt(buf, int64(c.Width), 10)
		buf = append(buf, ',')
	}
	if c.Height != 0 {
		buf = append(buf, 'h', '=')
		buf = strconv.AppendInt(buf, int64(c.Height), 10)
		buf = append(buf, ',')
	}
	if c.DisableEnlarge {
		buf = append(buf, 'u', '=', '0', ',')
	}
	if c.AspectMode != AspectModeDefault {
		buf = append(buf, 'a', '=')
		buf = strconv.AppendInt(buf, int64(c.AspectMode-1), 10)
		buf = append(buf, ',')
	}
	if c.DevicePixelRatio != 0 {
		buf = append(buf, 'd', 'p', 'r', '=')
		buf = strconv.AppendFloat(buf, c.DevicePixelRatio, 'f', -1, 64)
		buf = append(buf, ',')
	}

	// clipping parameters
	if ic := c.InputClip; ic != zr {
		buf = append(buf, "ic="...)
		buf = strconv.AppendInt(buf, int64(ic.Min.X), 10)
		buf = append(buf, ':')
		buf = strconv.AppendInt(buf, int64(ic.Min.Y), 10)
		buf = append(buf, ':')
		buf = strconv.AppendInt(buf, int64(ic.Max.X), 10)
		buf = append(buf, ':')
		buf = strconv.AppendInt(buf, int64(ic.Max.Y), 10)
		buf = append(buf, ',')
	}
	if cm, ic := c.ClipMax, c.InputClipRatio; cm != zp && ic != zr {
		x1 := float64(ic.Min.X) / float64(cm.X)
		y1 := float64(ic.Min.Y) / float64(cm.Y)
		x2 := float64(ic.Max.X) / float64(cm.X)
		y2 := float64(ic.Max.Y) / float64(cm.Y)
		buf = append(buf, "icr="...)
		buf = strconv.AppendFloat(buf, x1, 'f', -1, 64)
		buf = append(buf, ':')
		buf = strconv.AppendFloat(buf, y1, 'f', -1, 64)
		buf = append(buf, ':')
		buf = strconv.AppendFloat(buf, x2, 'f', -1, 64)
		buf = append(buf, ':')
		buf = strconv.AppendFloat(buf, y2, 'f', -1, 64)
		buf = append(buf, ',')
	}
	if ig := c.InputOrigin; ig != OriginDefault {
		buf = append(buf, "ig="...)
		buf = strconv.AppendInt(buf, int64(ig), 10)
		buf = append(buf, ',')
	}
	if c, oc := c.Clip, c.OutputClip; c != zr || oc != zr {
		if oc == zr {
			oc = c
		}
		buf = append(buf, "oc="...)
		buf = strconv.AppendInt(buf, int64(oc.Min.X), 10)
		buf = append(buf, ':')
		buf = strconv.AppendInt(buf, int64(oc.Min.Y), 10)
		buf = append(buf, ':')
		buf = strconv.AppendInt(buf, int64(oc.Max.X), 10)
		buf = append(buf, ':')
		buf = strconv.AppendInt(buf, int64(oc.Max.Y), 10)
		buf = append(buf, ',')
	}
	if c, oc, cm := c.ClipRatio, c.OutputClipRatio, c.ClipMax; cm != zp && (c != zr || oc != zr) {
		if oc == zr {
			oc = c
		}
		x1 := float64(oc.Min.X) / float64(cm.X)
		y1 := float64(oc.Min.Y) / float64(cm.Y)
		x2 := float64(oc.Max.X) / float64(cm.X)
		y2 := float64(oc.Max.Y) / float64(cm.Y)
		buf = append(buf, "ocr="...)
		buf = strconv.AppendFloat(buf, x1, 'f', -1, 64)
		buf = append(buf, ':')
		buf = strconv.AppendFloat(buf, y1, 'f', -1, 64)
		buf = append(buf, ':')
		buf = strconv.AppendFloat(buf, x2, 'f', -1, 64)
		buf = append(buf, ':')
		buf = strconv.AppendFloat(buf, y2, 'f', -1, 64)
		buf = append(buf, ',')
	}
	if og := c.OutputOrigin; og != OriginDefault {
		buf = append(buf, "og="...)
		buf = strconv.AppendInt(buf, int64(og), 10)
		buf = append(buf, ',')
	}

	if c.Origin != OriginDefault {
		buf = append(buf, 'g', '=')
		buf = strconv.AppendInt(buf, int64(c.Origin), 10)
		buf = append(buf, ',')
	}
	if c.Background != nil {
		r, g, b, a := c.Background.RGBA()
		if a == 0xffff {
			buf = append(buf, 'b', '=')
			buf = appendByte(buf, byte(r>>8))
			buf = appendByte(buf, byte(g>>8))
			buf = appendByte(buf, byte(b>>8))
			buf = append(buf, ',')
		} else if a == 0 {
			buf = append(buf, "b=000000,"...)
		} else {
			c := fmt.Sprintf("b=%02x%02x%02x%02x,", r>>8, g>>8, b>>8, a>>8)
			buf = append(buf, c...)
		}
	}

	// rotation
	if ir := c.InputRotate; ir != RotateDefault {
		if ir == RotateAuto {
			buf = append(buf, "ir=auto,"...)
		} else {
			buf = append(buf, "ir="...)
			buf = strconv.AppendInt(buf, int64(ir), 10)
			buf = append(buf, ',')
		}
	}
	if r, or := c.Rotate, c.OutputRotate; r != RotateDefault || or != RotateDefault {
		if or == RotateDefault {
			or = r
		}
		if or == RotateAuto {
			buf = append(buf, "or=auto,"...)
		} else {
			buf = append(buf, "or="...)
			buf = strconv.AppendInt(buf, int64(or), 10)
			buf = append(buf, ',')
		}
	}

	if c.Through != 0 {
		buf = append(buf, "through="...)
		buf = c.Through.append(buf)
		buf = append(buf, ',')
	}

	if len(c.Overlays) > 0 {
		for _, overlay := range c.Overlays {
			buf = append(buf, "l=("...)
			buf = overlay.append(buf)
			buf = append(buf, "),"...)
		}
	}

	// output formats
	if c.Format != "" {
		buf = append(buf, "f="...)
		buf = append(buf, c.Format...)
		buf = append(buf, ',')
	}
	if c.Quality != 0 {
		buf = append(buf, "q="...)
		buf = strconv.AppendInt(buf, int64(c.Quality), 10)
		buf = append(buf, ',')
	}
	if c.DisableOptimization {
		buf = append(buf, "o=0,"...)
	}

	// image filters
	if c.Unsharp.Radius != 0 {
		buf = append(buf, "unsharp="...)
		buf = c.Unsharp.append(buf)
		buf = append(buf, ',')
	}
	if c.Blur.Radius != 0 {
		buf = append(buf, "blur="...)
		buf = c.Blur.append(buf)
		buf = append(buf, ',')
	}
	if c.GrayScale != 0 {
		buf = append(buf, "grayscale="...)
		buf = strconv.AppendInt(buf, int64(c.GrayScale), 10)
		buf = append(buf, ',')
	}
	if c.Sepia != 0 {
		buf = append(buf, "sepia="...)
		buf = strconv.AppendInt(buf, int64(c.Sepia), 10)
		buf = append(buf, ',')
	}
	if c.Brightness != 0 {
		buf = append(buf, "brightness="...)
		buf = strconv.AppendInt(buf, int64(c.Brightness+100), 10)
		buf = append(buf, ',')
	}
	if c.Contrast != 0 {
		buf = append(buf, "contrast="...)
		buf = strconv.AppendInt(buf, int64(c.Contrast+100), 10)
		buf = append(buf, ',')
	}
	if c.Invert {
		buf = append(buf, "invert=1,"...)
	}

	if len(buf) != l {
		buf = buf[:len(buf)-1]
	}
	return buf
}

func appendByte(buf []byte, b byte) []byte {
	const digits = "0123456789abcdef"
	return append(buf, digits[b>>4], digits[b&0x0F])
}

func (a AspectMode) String() string {
	switch a {
	case AspectModeDefault:
		return "default"
	case AspectModeScale:
		return "scale"
	case AspectModeForceScale:
		return "force-scale"
	case AspectModePad:
		return "pad"
	}
	return ""
}

func (o Overlay) String() string {
	return string(o.append([]byte{}))
}

func (o Overlay) append(buf []byte) []byte {
	var zr image.Rectangle
	var zp image.Point

	l := len(buf)
	if o.Width != 0 {
		buf = append(buf, 'w', '=')
		buf = strconv.AppendInt(buf, int64(o.Width), 10)
		buf = append(buf, ',')
	}
	if o.Height != 0 {
		buf = append(buf, 'h', '=')
		buf = strconv.AppendInt(buf, int64(o.Height), 10)
		buf = append(buf, ',')
	}
	if o.DisableEnlarge {
		buf = append(buf, 'u', '=', '0', ',')
	}
	if o.AspectMode != AspectModeDefault {
		buf = append(buf, 'a', '=')
		buf = strconv.AppendInt(buf, int64(o.AspectMode-1), 10)
		buf = append(buf, ',')
	}

	// clipping parameters
	if ic := o.InputClip; ic != zr {
		buf = append(buf, "ic="...)
		buf = strconv.AppendInt(buf, int64(ic.Min.X), 10)
		buf = append(buf, ':')
		buf = strconv.AppendInt(buf, int64(ic.Min.Y), 10)
		buf = append(buf, ':')
		buf = strconv.AppendInt(buf, int64(ic.Max.X), 10)
		buf = append(buf, ':')
		buf = strconv.AppendInt(buf, int64(ic.Max.Y), 10)
		buf = append(buf, ',')
	}
	if cm, ic := o.ClipMax, o.InputClipRatio; cm != zp && ic != zr {
		x1 := float64(ic.Min.X) / float64(cm.X)
		y1 := float64(ic.Min.Y) / float64(cm.Y)
		x2 := float64(ic.Max.X) / float64(cm.X)
		y2 := float64(ic.Max.Y) / float64(cm.Y)
		buf = append(buf, "icr="...)
		buf = strconv.AppendFloat(buf, x1, 'f', -1, 64)
		buf = append(buf, ':')
		buf = strconv.AppendFloat(buf, y1, 'f', -1, 64)
		buf = append(buf, ':')
		buf = strconv.AppendFloat(buf, x2, 'f', -1, 64)
		buf = append(buf, ':')
		buf = strconv.AppendFloat(buf, y2, 'f', -1, 64)
		buf = append(buf, ',')
	}
	if ig := o.InputOrigin; ig != OriginDefault {
		buf = append(buf, "ig="...)
		buf = strconv.AppendInt(buf, int64(ig), 10)
		buf = append(buf, ',')
	}
	if c, oc := o.Clip, o.OutputClip; c != zr || oc != zr {
		if oc == zr {
			oc = c
		}
		buf = append(buf, "oc="...)
		buf = strconv.AppendInt(buf, int64(oc.Min.X), 10)
		buf = append(buf, ':')
		buf = strconv.AppendInt(buf, int64(oc.Min.Y), 10)
		buf = append(buf, ':')
		buf = strconv.AppendInt(buf, int64(oc.Max.X), 10)
		buf = append(buf, ':')
		buf = strconv.AppendInt(buf, int64(oc.Max.Y), 10)
		buf = append(buf, ',')
	}
	if c, oc, cm := o.ClipRatio, o.OutputClipRatio, o.ClipMax; cm != zp && (c != zr || oc != zr) {
		if oc == zr {
			oc = c
		}
		x1 := float64(oc.Min.X) / float64(cm.X)
		y1 := float64(oc.Min.Y) / float64(cm.Y)
		x2 := float64(oc.Max.X) / float64(cm.X)
		y2 := float64(oc.Max.Y) / float64(cm.Y)
		buf = append(buf, "ocr="...)
		buf = strconv.AppendFloat(buf, x1, 'f', -1, 64)
		buf = append(buf, ':')
		buf = strconv.AppendFloat(buf, y1, 'f', -1, 64)
		buf = append(buf, ':')
		buf = strconv.AppendFloat(buf, x2, 'f', -1, 64)
		buf = append(buf, ':')
		buf = strconv.AppendFloat(buf, y2, 'f', -1, 64)
		buf = append(buf, ',')
	}
	if og := o.OutputOrigin; og != OriginDefault {
		buf = append(buf, "og="...)
		buf = strconv.AppendInt(buf, int64(og), 10)
		buf = append(buf, ',')
	}

	if o.Origin != OriginDefault {
		buf = append(buf, 'g', '=')
		buf = strconv.AppendInt(buf, int64(o.Origin), 10)
		buf = append(buf, ',')
	}
	if o.Background != nil {
		r, g, b, a := o.Background.RGBA()
		if a == 0xffff {
			buf = append(buf, 'b', '=')
			buf = appendByte(buf, byte(r>>8))
			buf = appendByte(buf, byte(g>>8))
			buf = appendByte(buf, byte(b>>8))
			buf = append(buf, ',')
		} else if a == 0 {
			buf = append(buf, "b=000000,"...)
		} else {
			c := fmt.Sprintf("b=%02x%02x%02x%02x,", r>>8, g>>8, b>>8, a>>8)
			buf = append(buf, c...)
		}
	}

	// rotation
	if ir := o.InputRotate; ir != RotateDefault {
		if ir == RotateAuto {
			buf = append(buf, "ir=auto,"...)
		} else {
			buf = append(buf, "ir="...)
			buf = strconv.AppendInt(buf, int64(ir), 10)
			buf = append(buf, ',')
		}
	}
	if r, or := o.Rotate, o.OutputRotate; r != RotateDefault || or != RotateDefault {
		if or == RotateDefault {
			or = r
		}
		if or == RotateAuto {
			buf = append(buf, "or=auto,"...)
		} else {
			buf = append(buf, "or="...)
			buf = strconv.AppendInt(buf, int64(or), 10)
			buf = append(buf, ',')
		}
	}

	if o.Offset != zp {
		buf = append(buf, "x="...)
		buf = strconv.AppendInt(buf, int64(o.Offset.X), 10)
		buf = append(buf, ",y="...)
		buf = strconv.AppendInt(buf, int64(o.Offset.Y), 10)
		buf = append(buf, ',')
	}
	if o.OffsetRatio != zp && o.OffsetMax != zp {
		x := float64(o.OffsetRatio.X) / float64(o.OffsetMax.X)
		y := float64(o.OffsetRatio.Y) / float64(o.OffsetMax.Y)
		buf = append(buf, "xr="...)
		buf = strconv.AppendFloat(buf, x, 'f', -1, 64)
		buf = append(buf, ",yr="...)
		buf = strconv.AppendFloat(buf, y, 'f', -1, 64)
		buf = append(buf, ',')
	}
	if o.OverlayOrigin != OriginDefault {
		buf = append(buf, "lg="...)
		buf = strconv.AppendInt(buf, int64(o.OverlayOrigin), 10)
		buf = append(buf, ',')
	}

	if len(buf) > l && buf[len(buf)-1] == ',' {
		buf = buf[:len(buf)-1]
	}
	buf = append(buf, "%2f"...)
	buf = append(buf, url.QueryEscape(o.URL)...)
	return buf
}
