package imageflux

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const rectangleScale = 100

// nowFunc is for testing.
var nowFunc = time.Now

// ErrExpired is returned when the image is expired.
var ErrExpired = errors.New("imageflux: expired")

// ErrInvalidSignature is returned when the signature is invalid.
var ErrInvalidSignature = errors.New("imageflux: invalid signature")

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
	Format Format

	// Quality is quality of the output image.
	// It is used when the output format is JPEG or WebP.
	Quality int

	// DisableOptimization disables optimization of the Huffman coding table
	// of the output image when the output format is JPEG.
	DisableOptimization bool

	// Lossless enables lossless compression when the output format is WebP.
	Lossless bool

	// ExifOption specifies the Exif information to be included in the output image.
	ExifOption ExifOption

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

	// MaskType specifies the area to be treated as a mask.
	MaskType MaskType

	// PaddingMode specifies processing when the specified image is smaller than the input image.
	PaddingMode PaddingMode
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
	if u.Threshold != 0 {
		buf = append(buf, '+')
		buf = strconv.AppendFloat(buf, u.Gain, 'f', -1, 64)
		buf = append(buf, '+')
		buf = strconv.AppendFloat(buf, u.Threshold, 'f', -1, 64)
	}
	return buf
}

func parseUnsharp(s string) (Unsharp, error) {
	var u Unsharp

	// radius
	idx := strings.IndexByte(s, 'x')
	if idx < 0 {
		return Unsharp{}, errors.New("imageflux: invalid unsharp format")
	}
	r, err := strconv.ParseInt(s[:idx], 10, 0)
	if err != nil {
		return Unsharp{}, fmt.Errorf("imageflux: invalid unsharp format: %w", err)
	}
	if r <= 0 {
		return Unsharp{}, errors.New("imageflux: invalid unsharp format")
	}
	u.Radius = int(r)
	s = s[idx+1:]

	// sigma
	idx = strings.IndexByte(s, '+')
	if idx < 0 {
		sigma, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return u, fmt.Errorf("imageflux: invalid unsharp format: %w", err)
		}
		if sigma <= 0 {
			return u, errors.New("imageflux: invalid unsharp format")
		}
		u.Sigma = sigma
		return u, nil
	}
	sigma, err := strconv.ParseFloat(s[:idx], 64)
	if err != nil {
		return Unsharp{}, fmt.Errorf("imageflux: invalid unsharp format: %w", err)
	}
	if sigma == 0 {
		return u, errors.New("imageflux: invalid unsharp format")
	}
	u.Sigma = sigma
	s = s[idx+1:]

	// gain
	idx = strings.IndexByte(s, '+')
	if idx < 0 {
		return Unsharp{}, errors.New("imageflux: invalid unsharp format")
	}
	gain, err := strconv.ParseFloat(s[:idx], 64)
	if err != nil {
		return Unsharp{}, fmt.Errorf("imageflux: invalid unsharp format: %w", err)
	}
	u.Gain = gain
	s = s[idx+1:]

	// threshold
	threshold, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return Unsharp{}, fmt.Errorf("imageflux: invalid unsharp format: %w", err)
	}
	if threshold <= 0 || threshold >= 1 {
		return Unsharp{}, errors.New("imageflux: invalid unsharp format")
	}
	u.Threshold = threshold

	return u, nil
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

func parseBlur(s string) (Blur, error) {
	idx := strings.IndexByte(s, 'x')
	if idx < 0 {
		return Blur{}, errors.New("imageflux: invalid blur format")
	}

	// radius
	r, err := strconv.ParseInt(s[:idx], 10, 0)
	if err != nil {
		return Blur{}, fmt.Errorf("imageflux: invalid blur format: %w", err)
	}
	if r <= 0 {
		return Blur{}, errors.New("imageflux: invalid blur format")
	}

	// sigma
	sigma, err := strconv.ParseFloat(s[idx+1:], 64)
	if err != nil {
		return Blur{}, fmt.Errorf("imageflux: invalid blur format: %w", err)
	}
	if sigma <= 0 {
		return Blur{}, errors.New("imageflux: invalid blur format")
	}

	return Blur{
		Radius: int(r),
		Sigma:  sigma,
	}, nil
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

	aspectModeMax
)

// Origin is the origin.
type Origin int

const (
	// OriginDefault is default origin.
	OriginDefault Origin = 0

	// OriginTopLeft is top-left
	OriginTopLeft Origin = 1

	// OriginTopCenter is top-center
	OriginTopCenter Origin = 2

	// OriginTopRight is top-right
	OriginTopRight Origin = 3

	// OriginMiddleLeft is middle-left
	OriginMiddleLeft Origin = 4

	// OriginMiddleCenter is middle-center
	OriginMiddleCenter Origin = 5

	// OriginMiddleRight is middle-right
	OriginMiddleRight Origin = 6

	// OriginBottomLeft is bottom-left
	OriginBottomLeft Origin = 7

	// OriginBottomCenter is bottom-center
	OriginBottomCenter Origin = 8

	// OriginBottomRight is bottom-right
	OriginBottomRight Origin = 9

	originMax Origin = 10
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

	// FormatJPEG encodes the image as JPEG.
	FormatJPEG Format = "jpg"

	// FormatPNG encodes the image as PNG.
	FormatPNG Format = "png"

	// FormatGIF encodes the image as GIF.
	FormatGIF Format = "gif"

	// FormatWebP encodes the image as WebP.
	FormatWebP Format = "webp"

	// FormatWebPAuto encodes the image as a WebP if the client supports WebP.
	// Otherwise, the image is encoded as the same format with the input image.
	FormatWebPAuto Format = "webp:auto"

	// FormatWebPJPEG encodes the image as a WebP if the client supports WebP.
	// Otherwise, the image is encoded as JPEG.
	FormatWebPJPEG Format = "webp:jpg"

	// FormatWebPPNG encodes the image as a WebP if the client supports WebP.
	// Otherwise, the image is encoded as PNG.
	FormatWebPPNG Format = "webp:png"

	// FormatWebPGIF encodes the image as a WebP if the client supports WebP.
	// Otherwise, the image is encoded as GIF.
	FormatWebPGIF Format = "webp:gif"

	// FormatWebPFromJPEG encodes the image as a WebP.
	//
	// Deprecated: use FormatWebPJPEG instead.
	FormatWebPFromJPEG Format = "webp:jpeg"

	// FormatWebPFromPNG encodes the image as a WebP.
	//
	// Deprecated: use FormatWebPPNG instead.
	FormatWebPFromPNG Format = "webp:png"
)

func (f Format) String() string {
	return string(f)
}

// Rotate rotates the image.
type Rotate int

const (
	// RotateDefault is the default value of Rotate.
	// It is same effect as RotateTopLeft.
	RotateDefault Rotate = 0

	rotateMin Rotate = 1

	// RotateTopLeft does not anything.
	RotateTopLeft Rotate = 1

	// RotateTopRight flips the image left and right.
	RotateTopRight Rotate = 2

	// RotateBottomRight rotates the image 180 degrees.
	RotateBottomRight Rotate = 3

	// RotateBottomLeft flips the image upside down.
	RotateBottomLeft Rotate = 4

	// RotateLeftTop mirrors the image around the diagonal axis.
	RotateLeftTop Rotate = 5

	// RotateRightTop rotates the image left 90 degrees.
	RotateRightTop Rotate = 6

	// RotateRightBottom rotates the image 180 degrees and mirrors the image around the diagonal axis.
	RotateRightBottom Rotate = 7

	// RotateLeftBottom rotates the image right 90 degrees.
	RotateLeftBottom Rotate = 8

	rotateMax Rotate = 9

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

	// ThroughWebP skips converting WebP images.
	ThroughWebP
)

func (t Through) String() string {
	var buf [32]byte
	return string(t.append(buf[:]))
}

func (t Through) append(buf []byte) []byte {
	if (t & ThroughJPEG) != 0 {
		buf = append(buf, "jpg:"...)
	}
	if (t & ThroughPNG) != 0 {
		buf = append(buf, "png:"...)
	}
	if (t & ThroughGIF) != 0 {
		buf = append(buf, "gif:"...)
	}
	if (t & ThroughWebP) != 0 {
		buf = append(buf, "webp:"...)
	}
	if len(buf) == 0 {
		return buf
	}
	return buf[:len(buf)-1]
}

func parseThrough(s string) (Through, error) {
	var t Through
	for s != "" {
		var v string
		if idx := strings.IndexByte(s, ':'); idx >= 0 {
			v = s[:idx]
			s = s[idx+1:]
		} else {
			v = s
			s = ""
		}
		switch v {
		case "jpg":
			t |= ThroughJPEG
		case "png":
			t |= ThroughPNG
		case "gif":
			t |= ThroughGIF
		case "webp":
			t |= ThroughWebP
		}
	}
	return t, nil
}

// MaskType specifies the area to be treated as a mask.
type MaskType string

const (
	// MaskTypeWhite clips the mask image leaving the white parts.
	MaskTypeWhite MaskType = "white"

	// MaskTypeBlack clips the mask image leaving the black parts.
	MaskTypeBlack MaskType = "black"

	// MaskTypeAlpha clips the mask image leaving the opaque parts.
	MaskTypeAlpha MaskType = "alpha"
)

// PaddingMode specifies processing when the specified image is smaller than the input image.
type PaddingMode int

const (
	// PaddingModeDefault makes the part of the image that protrudes from the specified image transparent.
	PaddingModeDefault PaddingMode = 0

	// PaddingModeLeave leaves the overflow area of the specified image as it is.
	PaddingModeLeave PaddingMode = 1
)

// ExifOption specifies the Exif information to be included in the output image.
type ExifOption int

const (
	// ExifOptionDefault is the default value of ExifOption.
	ExifOptionDefault ExifOption = 0

	exifOptionMin ExifOption = 1

	// ExifOptionStrip removes all Exif information from the output image.
	ExifOptionStrip ExifOption = 1

	// ExifOptionKeepOrientation removes all Exif information
	// except Orientation from the output image.
	ExifOptionKeepOrientation ExifOption = 2

	exifOptionMax ExifOption = 3
)

// String returns a string representing the Config.
// If c is nil or zero value, it returns "f=auto".
func (c *Config) String() string {
	if c == nil {
		return "f=auto"
	}
	buf := bufPool.Get().(*[]byte)
	*buf = c.append((*buf)[:0], false)
	str := string(*buf)
	bufPool.Put(buf)
	return str
}

func (c *Config) append(buf []byte, escapeComma bool) []byte {
	var zr image.Rectangle
	var zp image.Point
	if c == nil {
		buf = append(buf, "f=auto"...)
		return buf
	}

	l := len(buf)
	if c.Width != 0 {
		buf = append(buf, "w="...)
		buf = strconv.AppendInt(buf, int64(c.Width), 10)
		buf = appendComma(buf, escapeComma)
	}
	if c.Height != 0 {
		buf = append(buf, "h="...)
		buf = strconv.AppendInt(buf, int64(c.Height), 10)
		buf = appendComma(buf, escapeComma)
	}
	if c.DisableEnlarge {
		buf = append(buf, "u=0"...)
		buf = appendComma(buf, escapeComma)
	}
	if c.AspectMode != AspectModeDefault {
		buf = append(buf, "a="...)
		buf = strconv.AppendInt(buf, int64(c.AspectMode-1), 10)
		buf = appendComma(buf, escapeComma)
	}
	if c.DevicePixelRatio != 0 {
		buf = append(buf, "dpr="...)
		buf = strconv.AppendFloat(buf, c.DevicePixelRatio, 'f', -1, 64)
		buf = appendComma(buf, escapeComma)
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
		buf = appendComma(buf, escapeComma)
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
		buf = appendComma(buf, escapeComma)
	}
	if ig := c.InputOrigin; ig != OriginDefault {
		buf = append(buf, "ig="...)
		buf = strconv.AppendInt(buf, int64(ig), 10)
		buf = appendComma(buf, escapeComma)
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
		buf = appendComma(buf, escapeComma)
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
		buf = appendComma(buf, escapeComma)
	}
	if og := c.OutputOrigin; og != OriginDefault {
		buf = append(buf, "og="...)
		buf = strconv.AppendInt(buf, int64(og), 10)
		buf = appendComma(buf, escapeComma)
	}

	if c.Origin != OriginDefault {
		buf = append(buf, "g="...)
		buf = strconv.AppendInt(buf, int64(c.Origin), 10)
		buf = appendComma(buf, escapeComma)
	}
	if c.Background != nil {
		b := color.NRGBAModel.Convert(c.Background).(color.NRGBA)
		if b.A == 0xff {
			// opaque background
			buf = append(buf, "b="...)
			buf = appendByte(buf, b.R)
			buf = appendByte(buf, b.G)
			buf = appendByte(buf, b.B)
			buf = appendComma(buf, escapeComma)
		} else {
			buf = append(buf, "b="...)
			buf = appendByte(buf, b.R)
			buf = appendByte(buf, b.G)
			buf = appendByte(buf, b.B)
			buf = appendByte(buf, b.A)
			buf = appendComma(buf, escapeComma)
		}
	}

	// rotation
	if ir := c.InputRotate; ir != RotateDefault {
		if ir == RotateAuto {
			buf = append(buf, "ir=auto"...)
			buf = appendComma(buf, escapeComma)
		} else {
			buf = append(buf, "ir="...)
			buf = strconv.AppendInt(buf, int64(ir), 10)
			buf = appendComma(buf, escapeComma)
		}
	}
	if r, or := c.Rotate, c.OutputRotate; r != RotateDefault || or != RotateDefault {
		if or == RotateDefault {
			or = r
		}
		if or == RotateAuto {
			buf = append(buf, "or=auto"...)
			buf = appendComma(buf, escapeComma)
		} else {
			buf = append(buf, "or="...)
			buf = strconv.AppendInt(buf, int64(or), 10)
			buf = appendComma(buf, escapeComma)
		}
	}

	if c.Through != 0 {
		buf = append(buf, "through="...)
		buf = c.Through.append(buf)
		buf = appendComma(buf, escapeComma)
	}

	if len(c.Overlays) > 0 {
		for _, overlay := range c.Overlays {
			buf = append(buf, "l=("...)
			buf = overlay.append(buf, escapeComma)
			buf = append(buf, ')')
			buf = appendComma(buf, escapeComma)
		}
	}

	// output formats
	if c.Format != "" {
		buf = append(buf, "f="...)
		buf = append(buf, c.Format...)
		buf = appendComma(buf, escapeComma)
	}
	if c.Quality != 0 {
		buf = append(buf, "q="...)
		buf = strconv.AppendInt(buf, int64(c.Quality), 10)
		buf = appendComma(buf, escapeComma)
	}
	if c.DisableOptimization {
		buf = append(buf, "o=0"...)
		buf = appendComma(buf, escapeComma)
	}
	if c.Lossless {
		buf = append(buf, "lossless=1"...)
		buf = appendComma(buf, escapeComma)
	}
	if c.ExifOption != ExifOptionDefault {
		buf = append(buf, "s="...)
		buf = strconv.AppendInt(buf, int64(c.ExifOption), 10)
		buf = appendComma(buf, escapeComma)
	}

	// image filters
	if c.Unsharp.Radius != 0 {
		buf = append(buf, "unsharp="...)
		buf = c.Unsharp.append(buf)
		buf = appendComma(buf, escapeComma)
	}
	if c.Blur.Radius != 0 {
		buf = append(buf, "blur="...)
		buf = c.Blur.append(buf)
		buf = appendComma(buf, escapeComma)
	}
	if c.GrayScale != 0 {
		buf = append(buf, "grayscale="...)
		buf = strconv.AppendInt(buf, int64(c.GrayScale), 10)
		buf = appendComma(buf, escapeComma)
	}
	if c.Sepia != 0 {
		buf = append(buf, "sepia="...)
		buf = strconv.AppendInt(buf, int64(c.Sepia), 10)
		buf = appendComma(buf, escapeComma)
	}
	if c.Brightness != 0 {
		buf = append(buf, "brightness="...)
		buf = strconv.AppendInt(buf, int64(c.Brightness+100), 10)
		buf = appendComma(buf, escapeComma)
	}
	if c.Contrast != 0 {
		buf = append(buf, "contrast="...)
		buf = strconv.AppendInt(buf, int64(c.Contrast+100), 10)
		buf = appendComma(buf, escapeComma)
	}
	if c.Invert {
		buf = append(buf, "invert=1"...)
		buf = appendComma(buf, escapeComma)
	}

	if len(buf) == l {
		buf = append(buf, "f=auto"...)
		buf = appendComma(buf, escapeComma)
	}
	if escapeComma {
		return buf[:len(buf)-3]
	}
	return buf[:len(buf)-1]
}

func appendByte(buf []byte, b byte) []byte {
	const digits = "0123456789abcdef"
	return append(buf, digits[b>>4], digits[b&0x0F])
}

func appendComma(buf []byte, escape bool) []byte {
	if escape {
		return append(buf, "%2C"...)
	}
	return append(buf, ',')
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
	return string(o.append([]byte{}, false))
}

func (o Overlay) append(buf []byte, escapeComma bool) []byte {
	var zr image.Rectangle
	var zp image.Point

	l := len(buf)
	if o.Width != 0 {
		buf = append(buf, "w="...)
		buf = strconv.AppendInt(buf, int64(o.Width), 10)
		buf = appendComma(buf, escapeComma)
	}
	if o.Height != 0 {
		buf = append(buf, "h="...)
		buf = strconv.AppendInt(buf, int64(o.Height), 10)
		buf = appendComma(buf, escapeComma)
	}
	if o.DisableEnlarge {
		buf = append(buf, "u=0"...)
		buf = appendComma(buf, escapeComma)
	}
	if o.AspectMode != AspectModeDefault {
		buf = append(buf, "a="...)
		buf = strconv.AppendInt(buf, int64(o.AspectMode-1), 10)
		buf = appendComma(buf, escapeComma)
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
		buf = appendComma(buf, escapeComma)
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
		buf = appendComma(buf, escapeComma)
	}
	if ig := o.InputOrigin; ig != OriginDefault {
		buf = append(buf, "ig="...)
		buf = strconv.AppendInt(buf, int64(ig), 10)
		buf = appendComma(buf, escapeComma)
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
		buf = appendComma(buf, escapeComma)
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
		buf = appendComma(buf, escapeComma)
	}
	if og := o.OutputOrigin; og != OriginDefault {
		buf = append(buf, "og="...)
		buf = strconv.AppendInt(buf, int64(og), 10)
		buf = appendComma(buf, escapeComma)
	}

	if o.Origin != OriginDefault {
		buf = append(buf, 'g', '=')
		buf = strconv.AppendInt(buf, int64(o.Origin), 10)
		buf = appendComma(buf, escapeComma)
	}
	if o.Background != nil {
		r, g, b, a := o.Background.RGBA()
		if a == 0xffff {
			buf = append(buf, 'b', '=')
			buf = appendByte(buf, byte(r>>8))
			buf = appendByte(buf, byte(g>>8))
			buf = appendByte(buf, byte(b>>8))
			buf = appendComma(buf, escapeComma)
		} else if a == 0 {
			buf = append(buf, "b=000000"...)
			buf = appendComma(buf, escapeComma)
		} else {
			c := fmt.Sprintf("b=%02x%02x%02x%02x,", r>>8, g>>8, b>>8, a>>8)
			buf = append(buf, c...)
		}
	}

	// rotation
	if ir := o.InputRotate; ir != RotateDefault {
		if ir == RotateAuto {
			buf = append(buf, "ir=auto"...)
			buf = appendComma(buf, escapeComma)
		} else {
			buf = append(buf, "ir="...)
			buf = strconv.AppendInt(buf, int64(ir), 10)
			buf = appendComma(buf, escapeComma)
		}
	}
	if r, or := o.Rotate, o.OutputRotate; r != RotateDefault || or != RotateDefault {
		if or == RotateDefault {
			or = r
		}
		if or == RotateAuto {
			buf = append(buf, "or=auto"...)
			buf = appendComma(buf, escapeComma)
		} else {
			buf = append(buf, "or="...)
			buf = strconv.AppendInt(buf, int64(or), 10)
			buf = appendComma(buf, escapeComma)
		}
	}

	if o.Offset != zp {
		buf = append(buf, "x="...)
		buf = strconv.AppendInt(buf, int64(o.Offset.X), 10)
		buf = appendComma(buf, escapeComma)
		buf = append(buf, "y="...)
		buf = strconv.AppendInt(buf, int64(o.Offset.Y), 10)
		buf = appendComma(buf, escapeComma)
	}
	if o.OffsetRatio != zp && o.OffsetMax != zp {
		x := float64(o.OffsetRatio.X) / float64(o.OffsetMax.X)
		y := float64(o.OffsetRatio.Y) / float64(o.OffsetMax.Y)
		buf = append(buf, "xr="...)
		buf = strconv.AppendFloat(buf, x, 'f', -1, 64)
		buf = appendComma(buf, escapeComma)
		buf = append(buf, "yr="...)
		buf = strconv.AppendFloat(buf, y, 'f', -1, 64)
		buf = appendComma(buf, escapeComma)
	}
	if o.OverlayOrigin != OriginDefault {
		buf = append(buf, "lg="...)
		buf = strconv.AppendInt(buf, int64(o.OverlayOrigin), 10)
		buf = appendComma(buf, escapeComma)
	}

	// mask
	if o.MaskType != "" {
		buf = append(buf, "mask="...)
		buf = append(buf, o.MaskType...)
		if o.PaddingMode != 0 {
			buf = append(buf, ':')
			buf = strconv.AppendInt(buf, int64(o.PaddingMode), 10)
		}
		buf = appendComma(buf, escapeComma)
	}

	if len(buf) > l {
		buf = buf[:len(buf)-1] // remove trailing comma
	}
	buf = append(buf, "%2F"...)
	buf = append(buf, url.QueryEscape(o.URL)...)
	return buf
}

func ParseConfig(s string) (config *Config, rest string, err error) {
	state := parseState{
		s:      s,
		config: &Config{},
	}
	return state.parseConfig()
}

type parseState struct {
	s      string
	idx    int
	config *Config

	// the signature that the user provided.
	signature string
}

func (s *parseState) parseConfig() (*Config, string, error) {
	if !s.hasParameter() {
		return s.config, s.rest(), nil
	}

	for {
		key := s.getKey()
		if key == "" {
			break
		}
		if !s.skipEqual() {
			return nil, "", fmt.Errorf("imageflux: missing '=' after key %q", key)
		}
		value := s.getValue()
		s.skipComma()
		if err := s.setValue(key, value); err != nil {
			return nil, "", err
		}
	}

	return s.config, s.rest(), nil
}

func (s *parseState) parseConfigAndVerifySignature(secret []byte) (*Config, string, error) {
	if !s.hasParameter() {
		buf := []byte(s.s)
		if err := s.verifySignature(secret, buf); err != nil {
			return nil, "", err
		}
		return s.config, s.rest(), nil
	}

	buf := make([]byte, 0, len(s.s))
	if len(s.s) == 0 || s.s[0] != '/' {
		buf = append(buf, '/')
	}
	buf = append(buf, s.s[:s.idx]...)

	hasParam := false
	for {
		start := s.idx
		key := s.getKey()
		if key == "" {
			break
		}
		if !s.skipEqual() {
			return nil, "", fmt.Errorf("imageflux: missing '=' after key %q", key)
		}
		value := s.getValue()
		s.skipComma()
		if err := s.setValue(key, value); err != nil {
			return nil, "", err
		}
		end := s.idx

		if key != "sig" {
			hasParam = true
			buf = append(buf, s.s[start:end]...)
		}
	}

	if hasParam {
		if len(buf) >= 1 && buf[len(buf)-1] == ',' {
			buf = buf[:len(buf)-1]
		} else if len(buf) >= 3 && string(buf[len(buf)-3:]) == "%2C" {
			buf = buf[:len(buf)-3]
		} else if len(buf) >= 3 && string(buf[len(buf)-3:]) == "%2c" {
			buf = buf[:len(buf)-3]
		}
	} else {
		buf = buf[:0]
	}
	buf = append(buf, s.rest()...)

	if err := s.verifySignature(secret, buf); err != nil {
		return nil, "", err
	}

	return s.config, s.rest(), nil
}

func (s *parseState) verifySignature(secret, data []byte) error {
	if strings.HasPrefix(s.signature, "1.") {
		// signature version 1
		sig, err := base64.URLEncoding.DecodeString(s.signature[len("1."):])
		if err != nil {
			return ErrInvalidSignature
		}

		w := hmac.New(sha256.New, secret)
		w.Write(data) // hash.hash never returns an error, so no need to check errors.
		sum := w.Sum(nil)

		if !hmac.Equal(sig, sum) {
			return ErrInvalidSignature
		}
		return nil
	}
	return ErrInvalidSignature
}

func (s *parseState) hasParameter() bool {
	i := s.idx
	if i >= len(s.s) {
		return false
	}

	// skip leading slash
	if s.s[i] == '/' {
		i++
	}

	// parameters may start with 'c/' or 'c!/'.
	if strings.HasPrefix(s.s[i:], "c/") {
		s.idx = i + len("c/")
		return true
	}
	if strings.HasPrefix(s.s[i:], "c!/") {
		s.idx += i + len("c!/")
		return true
	}

	// guess whether the string has parameters.
	// parameters always have '=', so we search for it.
	for ; i < len(s.s); i++ {
		if s.s[i] == '/' {
			// we didn't find any parameter.
			return false
		}

		if s.s[i] == '=' {
			// we might find a parameter.
			if s.s[s.idx] == '/' {
				s.idx++
			}
			return true
		}
	}
	return false
}

func (s *parseState) setValue(key, value string) error {
	var zr image.Rectangle

	switch key {
	// Width
	case "w":
		w, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid width %q", value)
		}
		s.config.Width = w

	// Height
	case "h":
		h, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid height %q", value)
		}
		s.config.Height = h

	// DisableEnlarge
	case "u":
		switch value {
		case "0":
			s.config.DisableEnlarge = true
		case "1":
			s.config.DisableEnlarge = false
		default:
			return fmt.Errorf("imageflux: invalid disable enlarge %q", value)
		}

	// AspectMode
	case "a":
		a, err := strconv.Atoi(value)
		if err != nil || a < 0 || AspectMode(a+1) >= aspectModeMax {
			return fmt.Errorf("imageflux: invalid aspect mode %q", value)
		}
		s.config.AspectMode = AspectMode(a + 1)

	// DevicePixelRatio
	case "dpr":
		dpr, err := strconv.ParseFloat(value, 64)
		if err != nil || dpr <= 0 || math.IsNaN(dpr) {
			return fmt.Errorf("imageflux: invalid device pixel ratio %q", value)
		}
		s.config.DevicePixelRatio = dpr

	// InputClip
	case "ic":
		v0, v1, v2, v3, ok := split4(value)
		if !ok {
			return fmt.Errorf("imageflux: invalid input clip %q", value)
		}
		minX, err0 := strconv.Atoi(v0)
		minY, err1 := strconv.Atoi(v1)
		maxX, err2 := strconv.Atoi(v2)
		maxY, err3 := strconv.Atoi(v3)
		ic := image.Rect(minX, minY, maxX, maxY)
		if err0 != nil || err1 != nil || err2 != nil || err3 != nil || ic == zr {
			return fmt.Errorf("imageflux: invalid input clip %q", value)
		}
		s.config.InputClip = ic

	// InputClipRatio
	case "icr":
		v0, v1, v2, v3, ok := split4(value)
		if !ok {
			return fmt.Errorf("imageflux: invalid input clip ratio %q", value)
		}
		minX, err0 := strconv.ParseFloat(v0, 64)
		minY, err1 := strconv.ParseFloat(v1, 64)
		maxX, err2 := strconv.ParseFloat(v2, 64)
		maxY, err3 := strconv.ParseFloat(v3, 64)
		icr := image.Rect(
			int(math.Round(minX*rectangleScale)),
			int(math.Round(minY*rectangleScale)),
			int(math.Round(maxX*rectangleScale)),
			int(math.Round(maxY*rectangleScale)),
		)
		if err0 != nil || err1 != nil || err2 != nil || err3 != nil || icr == zr {
			return fmt.Errorf("imageflux: invalid input clip ratio %q", value)
		}
		s.config.InputClipRatio = icr
		s.config.ClipMax = image.Pt(rectangleScale, rectangleScale)

	// InputOrigin
	case "ig":
		ig, err := strconv.Atoi(value)
		if err != nil || ig < 0 || Origin(ig) >= originMax {
			return fmt.Errorf("imageflux: invalid input origin %q", value)
		}
		s.config.InputOrigin = Origin(ig)

	// OutputClip
	case "oc", "c":
		v0, v1, v2, v3, ok := split4(value)
		if !ok {
			return fmt.Errorf("imageflux: invalid output clip %q", value)
		}
		minX, err0 := strconv.Atoi(v0)
		minY, err1 := strconv.Atoi(v1)
		maxX, err2 := strconv.Atoi(v2)
		maxY, err3 := strconv.Atoi(v3)
		oc := image.Rect(minX, minY, maxX, maxY)
		if err0 != nil || err1 != nil || err2 != nil || err3 != nil || oc == zr {
			return fmt.Errorf("imageflux: invalid input clip %q", value)
		}
		s.config.OutputClip = oc

	// OutputClipRatio
	case "ocr", "cr":
		v0, v1, v2, v3, ok := split4(value)
		if !ok {
			return fmt.Errorf("imageflux: invalid output clip ratio %q", value)
		}
		minX, err0 := strconv.ParseFloat(v0, 64)
		minY, err1 := strconv.ParseFloat(v1, 64)
		maxX, err2 := strconv.ParseFloat(v2, 64)
		maxY, err3 := strconv.ParseFloat(v3, 64)
		ocr := image.Rect(
			int(math.Round(minX*rectangleScale)),
			int(math.Round(minY*rectangleScale)),
			int(math.Round(maxX*rectangleScale)),
			int(math.Round(maxY*rectangleScale)),
		)
		if err0 != nil || err1 != nil || err2 != nil || err3 != nil || ocr == zr {
			return fmt.Errorf("imageflux: invalid input clip ratio %q", value)
		}
		s.config.OutputClipRatio = ocr
		s.config.ClipMax = image.Pt(rectangleScale, rectangleScale)

	// OutputOrigin
	case "og":
		og, err := strconv.Atoi(value)
		if err != nil || og < 0 || Origin(og) >= originMax {
			return fmt.Errorf("imageflux: invalid output origin %q", value)
		}
		s.config.OutputOrigin = Origin(og)

	// Origin
	case "g":
		g, err := strconv.Atoi(value)
		if err != nil || g < 0 || Origin(g) >= originMax {
			return fmt.Errorf("imageflux: invalid output origin %q", value)
		}
		s.config.Origin = Origin(g)

	// Background
	case "b":
		if len(value) == 6 {
			rgb, err := strconv.ParseUint(value, 16, 32)
			if err != nil {
				return fmt.Errorf("imageflux: invalid background %q", value)
			}
			s.config.Background = color.NRGBA{
				R: uint8(rgb >> 16),
				G: uint8(rgb >> 8),
				B: uint8(rgb),
				A: 0xff,
			}
		} else if len(value) == 8 {
			rgba, err := strconv.ParseUint(value, 16, 32)
			if err != nil {
				return fmt.Errorf("imageflux: invalid background %q", value)
			}
			s.config.Background = color.NRGBA{
				R: uint8(rgba >> 24),
				G: uint8(rgba >> 16),
				B: uint8(rgba >> 8),
				A: uint8(rgba),
			}
		} else {
			return fmt.Errorf("imageflux: invalid background %q", value)
		}

	// InputRotate
	case "ir":
		if value == "auto" {
			s.config.InputRotate = RotateAuto
		} else {
			ir, err := strconv.Atoi(value)
			if err != nil || Rotate(ir) < rotateMin || Rotate(ir) >= rotateMax {
				return fmt.Errorf("imageflux: invalid input rotate %q", value)
			}
			s.config.InputRotate = Rotate(ir)
		}

	// OutputRotate
	case "or", "r":
		if value == "auto" {
			s.config.OutputRotate = RotateAuto
		} else {
			ir, err := strconv.Atoi(value)
			if err != nil || Rotate(ir) < rotateMin || Rotate(ir) >= rotateMax {
				return fmt.Errorf("imageflux: invalid output rotate %q", value)
			}
			s.config.OutputRotate = Rotate(ir)
		}

	// Through
	case "through":
		t, err := parseThrough(value)
		if err != nil {
			return err
		}
		s.config.Through = t

	// TODO: Overlays

	// Format
	case "f":
		s.config.Format = Format(value)

	// Quality
	case "q":
		q, err := strconv.Atoi(value)
		if err != nil || q < 0 || q > 100 {
			return fmt.Errorf("imageflux: invalid quality %q", value)
		}
		s.config.Quality = q

	// DisableOptimization
	case "o":
		switch value {
		case "0":
			s.config.DisableOptimization = true
		case "1":
			s.config.DisableOptimization = false
		default:
			return fmt.Errorf("imageflux: invalid optimization %q", value)
		}

	// Lossless
	case "lossless":
		switch value {
		case "0":
			s.config.Lossless = false
		case "1":
			s.config.Lossless = true
		default:
			return fmt.Errorf("imageflux: invalid lossless %q", value)
		}

	// ExifOption
	case "s":
		v, err := strconv.Atoi(value)
		if err != nil || ExifOption(v) < exifOptionMin || ExifOption(v) >= exifOptionMax {
			return fmt.Errorf("imageflux: invalid exif option %q", value)
		}
		s.config.ExifOption = ExifOption(v)

	// Unsharp
	case "unsharp":
		unsharp, err := parseUnsharp(value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid unsharp %q", value)
		}
		s.config.Unsharp = unsharp

	// Blur
	case "blur":
		blur, err := parseBlur(value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid blur %q", value)
		}
		s.config.Blur = blur

	// GrayScale
	case "grayscale":
		grayscale, err := strconv.Atoi(value)
		if err != nil || grayscale < 0 || grayscale > 100 {
			return fmt.Errorf("imageflux: invalid grayscale %q", value)
		}
		s.config.GrayScale = grayscale

	// Sepia
	case "sepia":
		sepia, err := strconv.Atoi(value)
		if err != nil || sepia < 0 || sepia > 100 {
			return fmt.Errorf("imageflux: invalid sepia %q", value)
		}
		s.config.Sepia = sepia

	// Brightness
	case "brightness":
		brightness, err := strconv.Atoi(value)
		if err != nil || brightness < 0 {
			return fmt.Errorf("imageflux: invalid brightness %q", value)
		}
		s.config.Brightness = brightness - 100

	// Contrast
	case "contrast":
		contrast, err := strconv.Atoi(value)
		if err != nil || contrast < 0 {
			return fmt.Errorf("imageflux: invalid contrast %q", value)
		}
		s.config.Contrast = contrast - 100

	// Invert
	case "invert":
		switch value {
		case "0":
			s.config.Invert = false
		case "1":
			s.config.Invert = true
		default:
			return fmt.Errorf("imageflux: invalid invert %q", value)
		}

	case "expires":
		expires, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid expires %q", value)
		}
		if !expires.Before(nowFunc()) {
			return ErrExpired
		}

	case "sig":
		// if signature is already set, ignore this
		if s.signature == "" {
			s.signature = value
		}
	}
	return nil
}

// getKey returns the key at the current index and advances the index.
func (s *parseState) getKey() string {
	i := s.idx
	for i < len(s.s) {
		if s.s[i] == '=' || s.s[i] == '/' {
			break
		}
		i++
	}
	key := s.s[s.idx:i]
	s.idx = i
	return key
}

// skipEqual skips the '=' at the current index and returns true if it was found.
func (s *parseState) skipEqual() (skipped bool) {
	if s.idx < len(s.s) && s.s[s.idx] == '=' {
		s.idx++
		return true
	}
	return false
}

// getValue returns the value at the current index and advances the index.
func (s *parseState) getValue() string {
	var nest int
	i := s.idx
LOOP:
	for ; i < len(s.s); i++ {
		switch s.s[i] {
		case '(':
			nest++
		case ')':
			nest--
		case ',', '/':
			if nest == 0 {
				break LOOP
			}
		case '%':
			if nest != 0 {
				break
			}
			if i+3 < len(s.s) && (s.s[i:i+3] == "%2c" || s.s[i:i+3] == "%2C") {
				// "%2C" is encoded comma ','.
				break LOOP
			}
		}
	}
	value := s.s[s.idx:i]
	s.idx = i
	return value
}

// skipComma skips the ',' at the current index and returns true if it was found.
func (s *parseState) skipComma() (skipped bool) {
	if s.idx < len(s.s) && s.s[s.idx] == ',' {
		s.idx++
		return true
	}
	if s.idx+3 < len(s.s) && (s.s[s.idx:s.idx+3] == "%2c" || s.s[s.idx:s.idx+3] == "%2C") {
		// "%2C" is encoded comma ','.
		s.idx += 3
		return true
	}
	return false
}

func (s *parseState) rest() string {
	return s.s[s.idx:]
}

func split4(s string) (a, b, c, d string, ok bool) {
	idx1 := strings.IndexByte(s, ':')
	if idx1 < 0 {
		return
	}
	idx2 := strings.IndexByte(s[idx1+1:], ':')
	if idx2 < 0 {
		return
	}
	idx3 := strings.IndexByte(s[idx1+idx2+2:], ':')
	if idx3 < 0 {
		return
	}
	a = s[:idx1]
	b = s[idx1+1 : idx1+idx2+1]
	c = s[idx1+idx2+2 : idx1+idx2+idx3+2]
	d = s[idx1+idx2+idx3+3:]
	ok = true
	return
}
