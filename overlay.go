package imageflux

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"net/url"
	"strconv"
	"strings"
)

// Overlay is the configure of an overlay image.
type Overlay struct {
	// Path is a path for overlay image.
	Path string

	// URL is an url for overlay image.
	//
	// Deprecated: Use Path instead.
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
		b := color.NRGBAModel.Convert(o.Background).(color.NRGBA)
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

	// remove trailing comma
	if len(buf) != l {
		if escapeComma {
			buf = buf[:len(buf)-3]
		} else {
			buf = buf[:len(buf)-1]
		}
	}

	path := o.Path
	if path == "" {
		path = o.URL
	}
	if path == "" {
		return buf
	}
	if path[0] != '/' {
		buf = append(buf, "%2F"...)
	}
	return append(buf, url.PathEscape(path)...)
}

type overlayParseState struct {
	s       string
	idx     int
	overlay *Overlay
}

// ParseOverlay parses an overlay image.
func ParseOverlay(s string) (*Overlay, error) {
	ss, err := url.PathUnescape(s)
	if err != nil {
		return nil, err
	}

	state := overlayParseState{
		s:       ss,
		overlay: &Overlay{},
	}
	return state.parseOverlay()
}

// getKey returns the key at the current index and advances the index.
func (s *overlayParseState) getKey() (key string, foundEqual bool) {
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
func (s *overlayParseState) getValue() string {
	i := s.idx
	for ; i < len(s.s); i++ {
		if s.s[i] == ',' || s.s[i] == '/' {
			break
		}
	}
	value := s.s[s.idx:i]
	s.idx = i
	return value
}

// skipComma skips the ',' at the current index and returns true if it was found.
func (s *overlayParseState) skipComma() (skipped bool) {
	if s.idx < len(s.s) && s.s[s.idx] == ',' {
		s.idx++
		return true
	}
	return false
}

func (s *overlayParseState) parseOverlay() (*Overlay, error) {
	for {
		key, foundEqual := s.getKey()
		if !foundEqual {
			if key != "" {
				return nil, fmt.Errorf("imageflux: missing '=' after key %q", key)
			}
			break
		}
		value := s.getValue()
		s.skipComma()
		if err := s.setValue(key, value); err != nil {
			return nil, err
		}
	}
	s.overlay.Path = s.s[s.idx:]
	if !strings.HasPrefix(s.overlay.Path, "/") {
		s.overlay.Path = "/" + s.overlay.Path
	}
	return s.overlay, nil
}

func (s *overlayParseState) setValue(key, value string) error {
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
		s.overlay.Width = w

	// Height
	case "h":
		h, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid height %q: %w", value, err)
		}
		if h <= 0 {
			return fmt.Errorf("imageflux: invalid height %q: validation error", value)
		}
		s.overlay.Height = h

	// DisableEnlarge
	case "u":
		switch value {
		case "0":
			s.overlay.DisableEnlarge = true
		case "1":
			s.overlay.DisableEnlarge = false
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
		s.overlay.AspectMode = AspectMode(a + 1)

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
		s.overlay.InputClip = ic

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
		s.overlay.InputClipRatio = icr
		s.overlay.ClipMax = image.Pt(rectangleScale, rectangleScale)

	// InputOrigin
	case "ig":
		ig, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid input origin %q: %w", value, err)
		}
		if ig < 0 || Origin(ig) >= originMax {
			return fmt.Errorf("imageflux: invalid input origin %q: validation error", value)
		}
		s.overlay.InputOrigin = Origin(ig)

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
		s.overlay.OutputClip = oc

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
		s.overlay.OutputClipRatio = ocr
		s.overlay.ClipMax = image.Pt(rectangleScale, rectangleScale)

	// OutputOrigin
	case "og":
		og, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid output origin %q: %w", value, err)
		}
		if og < 0 || Origin(og) >= originMax {
			return fmt.Errorf("imageflux: invalid output origin %q: validation error", value)
		}
		s.overlay.OutputOrigin = Origin(og)

	// Origin
	case "g":
		g, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid origin %q: %w", value, err)
		}
		if g < 0 || Origin(g) >= originMax {
			return fmt.Errorf("imageflux: invalid origin %q: validation error", value)
		}
		s.overlay.Origin = Origin(g)

	// Background
	case "b":
		if len(value) == 6 {
			rgb, err := strconv.ParseUint(value, 16, 32)
			if err != nil {
				return fmt.Errorf("imageflux: invalid background %q", value)
			}
			s.overlay.Background = color.NRGBA{
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
			s.overlay.Background = color.NRGBA{
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
			s.overlay.InputRotate = RotateAuto
		} else {
			ir, err := strconv.Atoi(value)
			if err != nil || Rotate(ir) < rotateMin || Rotate(ir) >= rotateMax {
				return fmt.Errorf("imageflux: invalid input rotate %q", value)
			}
			s.overlay.InputRotate = Rotate(ir)
		}

	// OutputRotate
	case "or", "r":
		if value == "auto" {
			s.overlay.OutputRotate = RotateAuto
		} else {
			ir, err := strconv.Atoi(value)
			if err != nil || Rotate(ir) < rotateMin || Rotate(ir) >= rotateMax {
				return fmt.Errorf("imageflux: invalid output rotate %q", value)
			}
			s.overlay.OutputRotate = Rotate(ir)
		}
	}
	return nil
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
