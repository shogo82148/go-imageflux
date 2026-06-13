package imageflux

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"net/url"
	"slices"
	"strconv"
	"strings"
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

func (t *Text) String() string {
	return string(t.append(nil))
}

func (t *Text) append(buf []byte) []byte {
	var zp image.Point

	if t == nil || t.Text == "" {
		return buf
	}

	if t.Font != nil && t.Font.Name != "" {
		buf = append(buf, "font="...)
		buf = t.Font.append(buf)
		buf = appendComma(buf)
	}

	if t.Size != 0 {
		buf = append(buf, "size="...)
		buf = strconv.AppendFloat(buf, t.Size, 'f', -1, 64)
		buf = appendComma(buf)
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
		buf = appendComma(buf)
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
		buf = appendComma(buf)
	}

	if t.Width != 0 {
		buf = append(buf, "w="...)
		buf = strconv.AppendInt(buf, int64(t.Width), 10)
		buf = appendComma(buf)
	}
	if t.Height != 0 {
		buf = append(buf, "h="...)
		buf = strconv.AppendInt(buf, int64(t.Height), 10)
		buf = appendComma(buf)
	}

	if t.LineSpacing != 0 {
		buf = append(buf, "linespacing="...)
		buf = strconv.AppendFloat(buf, t.LineSpacing, 'f', -1, 64)
		buf = appendComma(buf)
	}

	if t.Align != 0 {
		buf = append(buf, "align="...)
		buf = strconv.AppendInt(buf, int64(t.Align), 10)
		buf = appendComma(buf)
	}

	if t.Direction != 0 {
		buf = append(buf, "dir="...)
		buf = strconv.AppendInt(buf, int64(t.Direction), 10)
		buf = appendComma(buf)
	}

	if t.Wrap != 0 {
		buf = append(buf, "wrap="...)
		buf = strconv.AppendInt(buf, int64(t.Wrap), 10)
		buf = appendComma(buf)
	}

	if t.Ellipsize {
		buf = append(buf, "ellipsize=1"...)
		buf = appendComma(buf)
	}

	if t.Justify {
		buf = append(buf, "justify=1"...)
		buf = appendComma(buf)
	}

	if t.Strike {
		buf = append(buf, "strike=1"...)
		buf = appendComma(buf)
	}

	if t.Offset != zp {
		buf = append(buf, "x="...)
		buf = strconv.AppendInt(buf, int64(t.Offset.X), 10)
		buf = appendComma(buf)
		buf = append(buf, "y="...)
		buf = strconv.AppendInt(buf, int64(t.Offset.Y), 10)
		buf = appendComma(buf)
	}
	if t.OffsetRatio != zp && t.OffsetMax.X != 0 && t.OffsetMax.Y != 0 {
		x := float64(t.OffsetRatio.X) / float64(t.OffsetMax.X)
		y := float64(t.OffsetRatio.Y) / float64(t.OffsetMax.Y)
		buf = append(buf, "xr="...)
		buf = strconv.AppendFloat(buf, x, 'f', -1, 64)
		buf = appendComma(buf)
		buf = append(buf, "yr="...)
		buf = strconv.AppendFloat(buf, y, 'f', -1, 64)
		buf = appendComma(buf)
	}

	if t.OverlayOrigin != OriginDefault {
		buf = append(buf, "lg="...)
		buf = strconv.AppendInt(buf, int64(t.OverlayOrigin), 10)
		buf = appendComma(buf)
	}

	// mask
	if t.MaskType != "" {
		buf = append(buf, "mask="...)
		buf = append(buf, t.MaskType...)
		if t.PaddingMode != 0 {
			buf = append(buf, ':')
			buf = strconv.AppendInt(buf, int64(t.PaddingMode), 10)
		}
		buf = appendComma(buf)
	}

	// text MUST be the last parameter because it can contain any character.
	buf = append(buf, "text="...)
	buf = append(buf, url.PathEscape(t.Text)...)
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

func (f *Font) String() string {
	if f == nil {
		return ""
	}
	buf := f.append(nil)
	return string(buf)
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
		buf = appendComma(buf)
		buf = append(buf, "instance="...)
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
		buf = appendComma(buf)
		buf = append(buf, "var="...)
		buf = append(buf, url.PathEscape(tag)...)
		buf = append(buf, ':')
		buf = strconv.AppendFloat(buf, f.Variables[tag], 'f', -1, 64)
	}
	buf = append(buf, ')')
	return buf
}

type parseFontState struct {
	s    string
	idx  int
	font *Font
}

func ParseFont(s string) (*Font, error) {
	state := &parseFontState{
		s:    s,
		idx:  0,
		font: &Font{},
	}
	return state.parseFont()
}

func (s *parseFontState) parseFont() (*Font, error) {
	if s.idx >= len(s.s) || s.s[s.idx] != '(' {
		name, err := url.PathUnescape(s.s)
		if err != nil {
			return nil, fmt.Errorf("imageflux: invalid font name %q: %w", s.s, err)
		}
		s.font.Name = name
		return s.font, nil
	}
	s.idx++

	// parse font name
	name, err := url.PathUnescape(s.getValue())
	if err != nil {
		return nil, fmt.Errorf("imageflux: invalid font name %q: %w", s.s, err)
	}
	s.font.Name = name

	// parse parameters
	for s.idx < len(s.s) && s.s[s.idx] != ')' {
		if !s.skipComma() {
			return nil, fmt.Errorf("imageflux: unexpected character %q in font specification", s.s[s.idx])
		}

		key, foundEqual := s.getKey()
		if !foundEqual {
			return nil, fmt.Errorf("imageflux: missing '=' after key %q in font specification", key)
		}
		value := s.getValue()
		switch key {
		case "instance":
			instance, err := url.PathUnescape(value)
			if err != nil {
				return nil, fmt.Errorf("imageflux: invalid instance value %q: %w", value, err)
			}
			s.font.Instance = instance

		case "var":
			value, err := url.PathUnescape(value)
			if err != nil {
				return nil, fmt.Errorf("imageflux: invalid variable font specification %q: %w", value, err)
			}
			before, after, ok := strings.Cut(value, ":")
			if !ok {
				return nil, fmt.Errorf("imageflux: invalid variable font specification %q: missing ':'", value)
			}
			tag := before
			v, err := strconv.ParseFloat(after, 64)
			if err != nil {
				return nil, fmt.Errorf("imageflux: invalid variable font value %q: %w", after, err)
			}
			if math.IsNaN(v) || math.IsInf(v, 0) {
				return nil, fmt.Errorf("imageflux: invalid variable font value %q", after)
			}
			if s.font.Variables == nil {
				s.font.Variables = make(map[string]float64)
			}
			s.font.Variables[tag] = v

		default:
			return nil, fmt.Errorf("imageflux: unknown key %q in font specification", key)
		}
	}
	if s.idx >= len(s.s) || s.s[s.idx] != ')' {
		return nil, errors.New("imageflux: unexpected end of font specification")
	}
	s.idx++
	if s.idx < len(s.s) {
		return nil, fmt.Errorf("imageflux: extra characters after closing parenthesis in font specification: %q", s.s[s.idx:])
	}

	return s.font, nil
}

func (s *parseFontState) getKey() (key string, foundEqual bool) {
	i := s.idx
	for ; i < len(s.s); i++ {
		switch s.s[i] {
		case '=':
			key = s.s[s.idx:i]
			s.idx = i + 1
			return key, true
		case ')':
			key = s.s[s.idx:i]
			s.idx = i + 1
			return key, false
		}
	}
	return s.s[s.idx:i], false
}

func (s *parseFontState) getValue() string {
	i := s.idx
LOOP:
	for ; i < len(s.s); i++ {
		switch s.s[i] {
		case ',', ')':
			break LOOP
		case '%':
			if i+3 < len(s.s) && (s.s[i:i+3] == "%2c" || s.s[i:i+3] == "%2C") {
				break LOOP
			}
		}
	}
	value := s.s[s.idx:i]
	s.idx = i
	return value
}

func (s *parseFontState) skipComma() (skipped bool) {
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

// TextAlign specifies the alignment of the text.
type TextAlign int

const (
	textAlignMin TextAlign = 0

	// TextAlignLeft aligns the text to the left.
	TextAlignLeft TextAlign = 0

	// TextAlignCenter aligns the text to the center.
	TextAlignCenter TextAlign = 1

	// TextAlignRight aligns the text to the right.
	TextAlignRight TextAlign = 2

	textAlignMax TextAlign = 3
)

// TextDirection specifies the direction of the text.
type TextDirection int

const (
	textDirectionMin TextDirection = 0

	// TextDirectionAuto is the default value of TextDirection.
	TextDirectionAuto TextDirection = 0

	// TextDirectionLTR is left to right.
	TextDirectionLTR TextDirection = 1

	// TextDirectionRTL is right to left.
	TextDirectionRTL TextDirection = 2

	textDirectionMax TextDirection = 3
)

// TextWrap specifies the wrap mode of the text.
type TextWrap int

const (
	textWrapMin TextWrap = 0

	// TextWrapLine is the default value of TextWrap.
	// Break lines at allowed breakpoints.
	// Allowed breakpoints are based on [UAX #14].
	//
	// [UAX #14]: https://unicode.org/reports/tr14/
	TextWrapLine TextWrap = 0

	// TextWrapChar allows line breaks between arbitrary characters.
	// Here, characters refer to grapheme clusters as defined in [UAX #29].
	//
	// [UAX #29]: https://unicode.org/reports/tr29/
	TextWrapChar TextWrap = 1

	// TextWrapLineChar attempts to wrap lines at breakable positions,
	// but if that is not possible, it wraps lines between arbitrary characters.
	TextWrapLineChar TextWrap = 2

	textWrapMax TextWrap = 3
)

type textParseState struct {
	s    string
	idx  int
	text *Text
}

func ParseText(s string) (*Text, error) {
	state := &textParseState{
		s:    s,
		idx:  0,
		text: &Text{},
	}
	return state.parseText()
}

func (s *textParseState) getKey() (key string, foundEqual bool) {
	i := s.idx
	for ; i < len(s.s); i++ {
		switch s.s[i] {
		case '=':
			key = s.s[s.idx:i]
			s.idx = i + 1
			return key, true
		}
	}
	return s.s[s.idx:i], false
}

func (s *textParseState) getValue() string {
	i := s.idx
LOOP:
	for ; i < len(s.s); i++ {
		switch s.s[i] {
		case ',':
			break LOOP
		case '%':
			if i+3 < len(s.s) && (s.s[i:i+3] == "%2c" || s.s[i:i+3] == "%2C") {
				break LOOP
			}
		}
	}
	value := s.s[s.idx:i]
	s.idx = i
	return value
}

func (s *textParseState) skipComma() (skipped bool) {
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

func (s *textParseState) parseText() (*Text, error) {
	foundText := false
	for s.idx < len(s.s) {
		key, foundEqual := s.getKey()
		if !foundEqual {
			if key != "" {
				return nil, fmt.Errorf("imageflux: missing '=' after key %q", key)
			}
			return nil, errors.New("imageflux: unexpected ','")
		}
		if key == "text" {
			foundText = true
			break
		}
		value := s.getValue()
		if err := s.setValue(key, value); err != nil {
			return nil, err
		}
		if !s.skipComma() {
			return nil, fmt.Errorf("imageflux: unexpected character after key %q", key)
		}
	}
	if !foundText {
		return nil, errors.New("imageflux: missing text parameter")
	}
	text := s.s[s.idx:]
	text, err := url.PathUnescape(text)
	if err != nil {
		return nil, fmt.Errorf("imageflux: invalid text value %q: %w", text, err)
	}
	s.text.Text = text

	// validate the parameters.
	var errs []error
	if s.text.Width <= 0 {
		errs = append(errs, fmt.Errorf("imageflux: width must be positive, but got %d", s.text.Width))
	}
	if s.text.Height <= 0 {
		errs = append(errs, fmt.Errorf("imageflux: height must be positive, but got %d", s.text.Height))
	}
	if s.text.Size <= 0 {
		errs = append(errs, fmt.Errorf("imageflux: size must be positive, but got %f", s.text.Size))
	}
	if math.IsNaN(s.text.Size) || math.IsInf(s.text.Size, 0) {
		errs = append(errs, fmt.Errorf("imageflux: invalid size value %f", s.text.Size))
	}
	if math.IsNaN(s.text.LineSpacing) || math.IsInf(s.text.LineSpacing, 0) {
		errs = append(errs, fmt.Errorf("imageflux: invalid line spacing value %f", s.text.LineSpacing))
	}
	if s.text.Align < textAlignMin || s.text.Align >= textAlignMax {
		errs = append(errs, fmt.Errorf("imageflux: align value must be between %d and %d, but got %d", textAlignMin, textAlignMax, s.text.Align))
	}
	if s.text.Direction < textDirectionMin || s.text.Direction >= textDirectionMax {
		errs = append(errs, fmt.Errorf("imageflux: direction value must be between %d and %d, but got %d", textDirectionMin, textDirectionMax, s.text.Direction))
	}
	if s.text.Wrap < textWrapMin || s.text.Wrap >= textWrapMax {
		errs = append(errs, fmt.Errorf("imageflux: wrap value must be between %d and %d, but got %d", textWrapMin, textWrapMax, s.text.Wrap))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return s.text, nil
}

func (s *textParseState) setValue(key, value string) error {
	switch key {
	case "font":
		font, err := ParseFont(value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid font specification %q: %w", value, err)
		}
		s.text.Font = font

	case "size":
		size, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("imageflux: invalid size value %q: %w", value, err)
		}
		s.text.Size = size

	case "f":
		c, err := parseColor(value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid foreground %q: %w", value, err)
		}
		s.text.Foreground = c

	case "b":
		c, err := parseColor(value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid background %q: %w", value, err)
		}
		s.text.Background = c

	case "w":
		w, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid width value %q: %w", value, err)
		}
		s.text.Width = w

	case "h":
		h, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid height value %q: %w", value, err)
		}
		s.text.Height = h

	case "linespacing":
		lineSpacing, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("imageflux: invalid line spacing value %q: %w", value, err)
		}
		s.text.LineSpacing = lineSpacing

	case "align":
		align, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid align value %q: %w", value, err)
		}
		s.text.Align = TextAlign(align)

	case "dir":
		dir, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid direction value %q: %w", value, err)
		}
		s.text.Direction = TextDirection(dir)

	case "wrap":
		wrap, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid wrap value %q: %w", value, err)
		}
		s.text.Wrap = TextWrap(wrap)

	case "ellipsize":
		v, err := parseBoolean(value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid ellipsize value %q: %w", value, err)
		}
		s.text.Ellipsize = v

	case "justify":
		v, err := parseBoolean(value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid justify value %q: %w", value, err)
		}
		s.text.Justify = v

	case "strike":
		v, err := parseBoolean(value)
		if err != nil {
			return fmt.Errorf("imageflux: invalid strike value %q: %w", value, err)
		}
		s.text.Strike = v

	default:
		return fmt.Errorf("imageflux: unknown key %q in text specification", key)
	}
	return nil
}
