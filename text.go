package imageflux

import (
	"image"
	"image/color"
	"net/url"
	"slices"
	"strconv"
)

// Text is a text to be used for the image.
type Text struct {
	// Font is the font to be used for the text.
	Font *Font

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

	// Text is the text string.
	Text string
}

func (t *Text) append(buf []byte, escapeComma bool) []byte {
	var zp image.Point

	if t == nil || t.Text == "" {
		return buf
	}
	buf = append(buf, '(')

	if t.Font != nil && t.Font.Name != "" {
		buf = append(buf, "font="...)
		buf = t.Font.append(buf)
		buf = appendComma(buf, escapeComma)
	}

	if t.Size != 0 {
		buf = append(buf, "size="...)
		buf = strconv.AppendFloat(buf, t.Size, 'f', -1, 64)
		buf = appendComma(buf, escapeComma)
	}

	if t.Foreground != nil {
		f := color.NRGBAModel.Convert(t.Foreground).(color.NRGBA)
		buf = append(buf, "f="...)
		if f.A == 0xff {
			// opaque foreground color
			buf = appendByte(buf, f.R)
			buf = appendByte(buf, f.G)
			buf = appendByte(buf, f.B)
		} else {
			// transparent foreground color
			buf = appendByte(buf, f.R)
			buf = appendByte(buf, f.G)
			buf = appendByte(buf, f.B)
			buf = appendByte(buf, f.A)
		}
		buf = appendComma(buf, escapeComma)
	}

	if t.Background != nil {
		b := color.NRGBAModel.Convert(t.Background).(color.NRGBA)
		buf = append(buf, "b="...)
		if b.A == 0xff {
			// opaque background color
			buf = appendByte(buf, b.R)
			buf = appendByte(buf, b.G)
			buf = appendByte(buf, b.B)
		} else {
			// transparent background color
			buf = appendByte(buf, b.R)
			buf = appendByte(buf, b.G)
			buf = appendByte(buf, b.B)
			buf = appendByte(buf, b.A)
		}
		buf = appendComma(buf, escapeComma)
	}

	if t.Width != 0 {
		buf = append(buf, "w="...)
		buf = strconv.AppendInt(buf, int64(t.Width), 10)
		buf = appendComma(buf, escapeComma)
	}
	if t.Height != 0 {
		buf = append(buf, "h="...)
		buf = strconv.AppendInt(buf, int64(t.Height), 10)
		buf = appendComma(buf, escapeComma)
	}

	if t.LineSpacing != 0 {
		buf = append(buf, "linespacing="...)
		buf = strconv.AppendFloat(buf, t.LineSpacing, 'f', -1, 64)
		buf = appendComma(buf, escapeComma)
	}

	if t.Align != 0 {
		buf = append(buf, "align="...)
		buf = strconv.AppendInt(buf, int64(t.Align), 10)
		buf = appendComma(buf, escapeComma)
	}

	if t.Direction != 0 {
		buf = append(buf, "dir="...)
		buf = strconv.AppendInt(buf, int64(t.Direction), 10)
		buf = appendComma(buf, escapeComma)
	}

	if t.Wrap != 0 {
		buf = append(buf, "wrap="...)
		buf = strconv.AppendInt(buf, int64(t.Wrap), 10)
		buf = appendComma(buf, escapeComma)
	}

	if t.Ellipsize {
		buf = append(buf, "ellipsize=1"...)
		buf = appendComma(buf, escapeComma)
	}

	if t.Justify {
		buf = append(buf, "justify=1"...)
		buf = appendComma(buf, escapeComma)
	}

	if t.Strike {
		buf = append(buf, "strike=1"...)
		buf = appendComma(buf, escapeComma)
	}

	if t.Offset != zp {
		buf = append(buf, "x="...)
		buf = strconv.AppendInt(buf, int64(t.Offset.X), 10)
		buf = appendComma(buf, escapeComma)
		buf = append(buf, "y="...)
		buf = strconv.AppendInt(buf, int64(t.Offset.Y), 10)
		buf = appendComma(buf, escapeComma)
	}
	if t.OffsetRatio != zp && t.OffsetMax != zp {
		x := float64(t.OffsetRatio.X) / float64(t.OffsetMax.X)
		y := float64(t.OffsetRatio.Y) / float64(t.OffsetMax.Y)
		buf = append(buf, "xr="...)
		buf = strconv.AppendFloat(buf, x, 'f', -1, 64)
		buf = appendComma(buf, escapeComma)
		buf = append(buf, "yr="...)
		buf = strconv.AppendFloat(buf, y, 'f', -1, 64)
		buf = appendComma(buf, escapeComma)
	}

	if t.OverlayOrigin != OriginDefault {
		buf = append(buf, "lg="...)
		buf = strconv.AppendInt(buf, int64(t.OverlayOrigin), 10)
		buf = appendComma(buf, escapeComma)
	}

	// mask
	if t.MaskType != "" {
		buf = append(buf, "mask="...)
		buf = append(buf, t.MaskType...)
		if t.PaddingMode != 0 {
			buf = append(buf, ':')
			buf = strconv.AppendInt(buf, int64(t.PaddingMode), 10)
		}
		buf = appendComma(buf, escapeComma)
	}

	// text MUST be the last parameter because it can contain any character.
	buf = append(buf, "text="...)
	buf = append(buf, url.PathEscape(t.Text)...)
	buf = append(buf, ')')
	return buf
}

// Font specifies the font to be used for the text.
type Font struct {
	// Name is the name of the font.
	Name string

	// Instance specifies the named instance of the variable font.
	// When len(Variables) > 0, Instance is ignored.
	Instance string

	// Variables is the variable values of the variable font.
	// The key of the map is the tag name of the variable font axis.
	// The value of the map is the value of the variable font axis.
	Variables map[string]float64
}

func (f *Font) append(buf []byte) []byte {
	if f == nil || f.Name == "" {
		return buf
	}
	if f.Instance == "" && len(f.Variables) == 0 {
		// If the instance and variables are not specified, we can use the font name as it is.
		name := url.PathEscape(f.Name)
		buf = append(buf, name...)
		return buf
	}

	if len(f.Variables) == 0 {
		// (font-name,instance=instance-name)
		name := url.PathEscape(f.Name)
		instance := url.PathEscape(f.Instance)
		buf = append(buf, '(')
		buf = append(buf, name...)
		buf = append(buf, ",instance="...)
		buf = append(buf, instance...)
		buf = append(buf, ')')
		return buf
	}

	// (font-name,var=tag1:value1,var=tag2:value2,...)
	tags := make([]string, 0, len(f.Variables))
	for tag := range f.Variables {
		tags = append(tags, tag)
	}
	slices.Sort(tags)
	name := url.PathEscape(f.Name)
	buf = append(buf, '(')
	buf = append(buf, name...)
	for _, tag := range tags {
		buf = append(buf, ",var="...)
		buf = append(buf, url.PathEscape(tag)...)
		buf = append(buf, ':')
		buf = strconv.AppendFloat(buf, f.Variables[tag], 'f', -1, 64)
	}
	buf = append(buf, ')')
	return buf
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
