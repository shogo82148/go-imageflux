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
	"strconv"
	"strings"
	"time"
)

const rectangleScale = 65536

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

	// Expires is the time when the image expires.
	Expires time.Time

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
	Overlays []*Overlay

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

	// Text is the text to be used for the image.
	Text []*Text
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
		if sigma <= 0 || math.IsNaN(sigma) || math.IsInf(sigma, 0) {
			return u, errors.New("imageflux: invalid unsharp format")
		}
		u.Sigma = sigma
		return u, nil
	}
	sigma, err := strconv.ParseFloat(s[:idx], 64)
	if err != nil {
		return Unsharp{}, fmt.Errorf("imageflux: invalid unsharp format: %w", err)
	}
	if sigma <= 0 || math.IsNaN(sigma) || math.IsInf(sigma, 0) {
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
	if math.IsNaN(gain) || math.IsInf(gain, 0) {
		return Unsharp{}, errors.New("imageflux: invalid unsharp format")
	}
	u.Gain = gain
	s = s[idx+1:]

	// threshold
	threshold, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return Unsharp{}, fmt.Errorf("imageflux: invalid unsharp format: %w", err)
	}
	if threshold <= 0 || threshold >= 1 || math.IsNaN(threshold) {
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
	if sigma <= 0 || math.IsNaN(sigma) || math.IsInf(sigma, 0) {
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

func newFormat(s string) (Format, error) {
	// validate the input with the regexp /[a-z]+(:[a-z]+)*/.
	colon := true
	for _, ch := range []byte(s) {
		if ch == ':' {
			if colon {
				// double colons are detected. it's an error.
				return "", fmt.Errorf("imageflux: invalid format %q", s)
			}
			colon = true
			continue
		}
		if ch < 'a' || ch > 'z' {
			return "", fmt.Errorf("imageflux: invalid format %q", s)
		}
		colon = false
	}
	if colon {
		return "", fmt.Errorf("imageflux: invalid format %q", s)
	}
	return Format(s), nil
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
		default:
			return 0, fmt.Errorf("imageflux: unknown through format: %s", v)
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

// Text is a text to be used for the image.
type Text struct {
	// Font is the font to be used for the text.
	Font string

	// Size is the size of the text.
	Size float64

	// Foreground is the foreground color of the text.
	Foreground color.Color

	// Background is the background color of the text.
	Background color.Color

	// Width is the width of the text.
	Width int

	// Height is the height of the text.
	Height int

	// LineSpacing is the line spacing of the text.
	LineSpacing float64

	// Align is the alignment of the text.
	Align TextAlign

	// Direction is the direction of the text.
	Direction TextDirection

	// Wrap is the wrap mode of the text.
	Wrap TextWrap

	// Ellipsize is true if the text should be ellipsized.
	Ellipsize bool

	// Justify is true if the text should be justified.
	Justify bool

	// Strike is true if the text should be struck through.
	Strike bool

	// Text is the text string.
	Text string
}

// TextAlign specifies the alignment of the text.
type TextAlign int

const (
	// TextAlignLeft aligns the text to the left.
	TextAlignLeft TextAlign = 0

	// TextAlignCenter aligns the text to the center.
	TextAlignCenter TextAlign = 1

	// TextAlignRight aligns the text to the right.
	TextAlignRight TextAlign = 2
)

// TextDirection specifies the direction of the text.
type TextDirection int

const (
	// TextDirectionAuto is the default value of TextDirection.
	TextDirectionAuto TextDirection = 0

	// TextDirectionLTR is left to right.
	TextDirectionLTR TextDirection = 1

	// TextDirectionRTL is right to left.
	TextDirectionRTL TextDirection = 2
)

// TextWrap specifies the wrap mode of the text.
type TextWrap int

const (
	// TextWrapLine is the default value of TextWrap.
	TextWrapLine TextWrap = 0

	// TextWrapChar is character wrap.
	TextWrapChar TextWrap = 1

	// TextWrapLineChar is line and character wrap.
	TextWrapLineChar TextWrap = 2
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
	if !c.Expires.IsZero() {
		buf = append(buf, "expires="...)
		buf = c.Expires.In(time.UTC).AppendFormat(buf, time.RFC3339)
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
		key, foundEqual := s.getKey()
		if !foundEqual {
			if key != "" {
				return nil, "", fmt.Errorf("imageflux: missing '=' after key %q", key)
			}
			break
		}
		value, err := s.getValue()
		if err != nil {
			return nil, "", err
		}
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
		key, foundEqual := s.getKey()
		if !foundEqual {
			if key != "" {
				return nil, "", fmt.Errorf("imageflux: missing '=' after key %q", key)
			}
			break
		}
		value, err := s.getValue()
		if err != nil {
			return nil, "", err
		}
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
			return fmt.Errorf("imageflux: invalid width %q: %w", value, err)
		}
		if w <= 0 {
			return fmt.Errorf("imageflux: invalid width %q: validation error", value)
		}
		s.config.Width = w

	// Height
	case "h":
		h, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid height %q: %w", value, err)
		}
		if h <= 0 {
			return fmt.Errorf("imageflux: invalid height %q: validation error", value)
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
		if err != nil {
			return fmt.Errorf("imageflux: invalid aspect mode %q: %w", value, err)
		}
		if a < 0 || AspectMode(a+1) >= aspectModeMax {
			return fmt.Errorf("imageflux: invalid aspect mode %q: validation error", value)
		}
		s.config.AspectMode = AspectMode(a + 1)

	// DevicePixelRatio
	case "dpr":
		dpr, err := strconv.ParseFloat(value, 64)
		if err != nil || dpr <= 0 || math.IsNaN(dpr) || math.IsInf(dpr, 0) {
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
		ok = ok && err0 == nil && err1 == nil && err2 == nil && err3 == nil && icr != zr
		ok = ok && minX >= 0 && minX <= 1 && minY >= 0 && minY <= 1 && maxX >= 0 && maxX <= 1 && maxY >= 0 && maxY <= 1
		if !ok {
			return fmt.Errorf("imageflux: invalid input clip ratio %q", value)
		}
		s.config.InputClipRatio = icr
		s.config.ClipMax = image.Pt(rectangleScale, rectangleScale)

	// InputOrigin
	case "ig":
		ig, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid input origin %q: %w", value, err)
		}
		if ig < 0 || Origin(ig) >= originMax {
			return fmt.Errorf("imageflux: invalid input origin %q: validation error", value)
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
		ok = ok && err0 == nil && err1 == nil && err2 == nil && err3 == nil && ocr != zr
		ok = ok && minX >= 0 && minX <= 1 && minY >= 0 && minY <= 1 && maxX >= 0 && maxX <= 1 && maxY >= 0 && maxY <= 1
		if !ok {
			return fmt.Errorf("imageflux: invalid input clip ratio %q", value)
		}

		s.config.OutputClipRatio = ocr
		s.config.ClipMax = image.Pt(rectangleScale, rectangleScale)

	// OutputOrigin
	case "og":
		og, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid output origin %q: %w", value, err)
		}
		if og < 0 || Origin(og) >= originMax {
			return fmt.Errorf("imageflux: invalid output origin %q: validation error", value)
		}
		s.config.OutputOrigin = Origin(og)

	// Origin
	case "g":
		g, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid origin %q: %w", value, err)
		}
		if g < 0 || Origin(g) >= originMax {
			return fmt.Errorf("imageflux: invalid origin %q: validation error", value)
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

	// Overlays
	case "l":
		if len(value) < 2 || value[0] != '(' || value[len(value)-1] != ')' {
			return fmt.Errorf("imageflux: invalid overlays %q", value)
		}
		value = value[1 : len(value)-1]
		overlay, err := ParseOverlay(value)
		if err != nil {
			return err
		}
		s.config.Overlays = append(s.config.Overlays, overlay)

	// Format
	case "f":
		f, err := newFormat(value)
		if err != nil {
			return err
		}
		s.config.Format = f

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

	// Expires
	case "expires":
		expires, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid expires %q", value)
		}
		if !expires.After(nowFunc()) {
			return ErrExpired
		}
		s.config.Expires = expires

	case "sig":
		// if signature is already set, ignore this
		if s.signature == "" {
			s.signature = value
		}
	}
	return nil
}

// getKey returns the key at the current index and advances the index.
func (s *parseState) getKey() (key string, foundEqual bool) {
	i := s.idx
	for ; i < len(s.s); i++ {
		switch s.s[i] {
		case '=':
			key = s.s[s.idx:i]
			s.idx = i + 1
			foundEqual = true
			return
		case '/', ',':
			key = s.s[s.idx:i]
			s.idx = i
			foundEqual = false
			return
		}
	}
	return s.s[s.idx:i], false
}

// getValue returns the value at the current index and advances the index.
func (s *parseState) getValue() (string, error) {
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
	if nest != 0 {
		return "", errors.New("imageflux: invalid value: parenthesis is not closed")
	}
	value := s.s[s.idx:i]
	s.idx = i
	return value, nil
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
