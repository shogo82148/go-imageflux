package imageflux

import (
	"image"
	"image/color"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFont_String(t *testing.T) {
	cases := []struct {
		font     *Font
		expected string
	}{
		{
			font:     nil,
			expected: "",
		},
		{
			font: &Font{
				Name: "新ゴ R",
			},
			expected: "%E6%96%B0%E3%82%B4%20R",
		},
		{
			font: &Font{
				Name:     "DriveFlux",
				Instance: "B Italic",
			},
			expected: "(DriveFlux%2Cinstance=B%20Italic)",
		},
		{
			font: &Font{
				Name: "DriveFlux",
				Variables: map[string]float64{
					"wght": 700,
					"SMTH": 0,
					"CNTR": 0,
					"slnt": -16,
				},
			},
			expected: "(DriveFlux%2Cvar=CNTR:0%2Cvar=SMTH:0%2Cvar=slnt:-16%2Cvar=wght:700)",
		},
		{
			font: &Font{
				Name: "0",
				Variables: map[string]float64{
					" ": 0,
				},
			},
			expected: "(0%2Cvar=%20:0)",
		},
	}

	for _, c := range cases {
		if got := c.font.String(); got != c.expected {
			t.Errorf("%#v: Font.String() = %q, want %q", c.font, got, c.expected)
		}
	}
}

var parseFontCases = []struct {
	input    string
	expected *Font
}{
	{
		input: "",
		expected: &Font{
			Name: "",
		},
	},
	{
		input: "%E6%96%B0%E3%82%B4%20R",
		expected: &Font{
			Name: "新ゴ R",
		},
	},
	{
		input: "(%E6%96%B0%E3%82%B4%20R)",
		expected: &Font{
			Name: "新ゴ R",
		},
	},
	{
		input: "(DriveFlux,instance=B%20Italic)",
		expected: &Font{
			Name:     "DriveFlux",
			Instance: "B Italic",
		},
	},
	{
		input: "(DriveFlux%2Cinstance=B%20Italic)",
		expected: &Font{
			Name:     "DriveFlux",
			Instance: "B Italic",
		},
	},
	{
		input: "(DriveFlux,var=CNTR:0,var=SMTH:0,var=slnt:-16,var=wght:700)",
		expected: &Font{
			Name: "DriveFlux",
			Variables: map[string]float64{
				"wght": 700,
				"SMTH": 0,
				"CNTR": 0,
				"slnt": -16,
			},
		},
	},
	{
		input: "(DriveFlux%2Cvar=CNTR:0%2Cvar=SMTH:0%2Cvar=slnt:-16%2Cvar=wght:700)",
		expected: &Font{
			Name: "DriveFlux",
			Variables: map[string]float64{
				"wght": 700,
				"SMTH": 0,
				"CNTR": 0,
				"slnt": -16,
			},
		},
	},
	{
		input: "(,instance=0)",
		expected: &Font{
			Name:     "",
			Instance: "0",
		},
	},
	{
		input: "(0%2Cvar= :0)",
		expected: &Font{
			Name: "0",
			Variables: map[string]float64{
				" ": 0,
			},
		},
	},
	{
		input: "(0%2Cvar=%20:0)",
		expected: &Font{
			Name: "0",
			Variables: map[string]float64{
				" ": 0,
			},
		},
	},
}

func TestParseFont(t *testing.T) {
	for _, c := range parseFontCases {
		got, err := ParseFont(c.input)
		if err != nil {
			t.Errorf("ParseFont(%q) returned error: %v", c.input, err)
			continue
		}
		if diff := cmp.Diff(got, c.expected); diff != "" {
			t.Errorf("ParseFont(%q) = (+got / -expected) %s", c.input, diff)
		}
	}
}

var parseFontErrorCases = []string{
	"%XX",                             // invalid percent encoding
	"(DriveFlux,instance=%XX)",        // invalid percent encoding in instance
	"(DriveFlux,var=%XX:0)",           // invalid percent encoding in variable tag
	"(DriveFlux",                      // missing closing parenthesis
	"(DriveFlux,instance=B%20Italic",  // missing closing parenthesis
	"(DriveFlux,instance)",            // missing instance value
	"(DriveFlux,var=CNTR)",            // missing ':'
	"(DriveFlux,var=CNTR:NotANumber)", // invalid number
	"(DriveFlux,var=CNTR:NaN)",        // NaN is not a valid variable value
	"(DriveFlux,var=CNTR:+Inf)",       // +Inf is not a valid variable value
	"(DriveFlux,var=CNTR:-Inf)",       // -Inf is not a valid variable value
	"(DriveFlux,unknown=Value)",       // unknown key
}

func TestParseFont_Error(t *testing.T) {
	for _, c := range parseFontErrorCases {
		if _, err := ParseFont(c); err == nil {
			t.Errorf("ParseFont(%q) did not return error", c)
		}
	}
}

func FuzzParseFont(f *testing.F) {
	for _, c := range parseFontCases {
		f.Add(c.input)
	}
	for _, c := range parseFontErrorCases {
		f.Add(c)
	}

	f.Fuzz(func(t *testing.T, s string) {
		font, err := ParseFont(s)
		if err != nil {
			return
		}
		if font.Name == "" {
			return
		}
		s1 := font.String()
		font2, err := ParseFont(s1)
		if err != nil {
			t.Errorf("ParseFont(%q) returned error: %v", s1, err)
			return
		}
		if diff := cmp.Diff(font, font2); diff != "" {
			t.Errorf("ParseFont(%q) = (+font / -font2) %s", s, diff)
		}
	})
}

func TestText_String(t *testing.T) {
	cases := []struct {
		text     *Text
		expected string
	}{
		// Basic case
		{
			text: &Text{
				Font: &Font{
					Name: "新ゴ R",
				},
				Height: 100,
				Width:  400,
				Size:   12,
				Text:   "Hello, world!",
			},
			expected: "font=%E6%96%B0%E3%82%B4%20R%2Csize=12%2Cw=400%2Ch=100%2Ctext=Hello%2C%20world%21",
		},

		// Named instance with variable font
		{
			text: &Text{
				Font: &Font{
					Name:     "DriveFlux",
					Instance: "B Italic",
				},
				Height: 100,
				Width:  400,
				Size:   12,
				Text:   "Hello, world!",
			},
			expected: "font=(DriveFlux%2Cinstance=B%20Italic)%2Csize=12%2Cw=400%2Ch=100%2Ctext=Hello%2C%20world%21",
		},

		// Variable font with variable values
		{
			text: &Text{
				Font: &Font{
					Name: "DriveFlux",
					Variables: map[string]float64{
						"wght": 700,
						"SMTH": 0,
						"CNTR": 0,
						"slnt": -16,
					},
				},
				Height: 100,
				Width:  400,
				Size:   12,
				Text:   "Hello, world!",
			},
			expected: "font=(DriveFlux%2Cvar=CNTR:0%2Cvar=SMTH:0%2Cvar=slnt:-16%2Cvar=wght:700)%2Csize=12%2Cw=400%2Ch=100%2Ctext=Hello%2C%20world%21",
		},

		// Transparent foreground and background colors
		{
			text: &Text{
				Font: &Font{
					Name: "新ゴ R",
				},
				Height: 100,
				Width:  400,
				Size:   12,
				Foreground: &color.NRGBA{
					R: 0x11,
					G: 0x22,
					B: 0x33,
					A: 0xff,
				},
				Background: &color.NRGBA{
					R: 0x44,
					G: 0x55,
					B: 0x66,
					A: 0xff,
				},
				Text: "Hello, world!",
			},
			expected: "font=%E6%96%B0%E3%82%B4%20R%2Csize=12%2Cf=112233%2Cb=445566%2Cw=400%2Ch=100%2Ctext=Hello%2C%20world%21",
		},
		{
			text: &Text{
				Font: &Font{
					Name: "新ゴ R",
				},
				Height: 100,
				Width:  400,
				Size:   12,
				Foreground: &color.NRGBA{
					R: 0x11,
					G: 0x22,
					B: 0x33,
					A: 0x7f,
				},
				Background: &color.NRGBA{
					R: 0x44,
					G: 0x55,
					B: 0x66,
					A: 0x7f,
				},
				Text: "Hello, world!",
			},
			expected: "font=%E6%96%B0%E3%82%B4%20R%2Csize=12%2Cf=1122337f%2Cb=4455667f%2Cw=400%2Ch=100%2Ctext=Hello%2C%20world%21",
		},

		{
			text: &Text{
				Font: &Font{
					Name: "新ゴ R",
				},
				Height:      100,
				Width:       400,
				Size:        12,
				LineSpacing: 1.5,
				Align:       TextAlignCenter,
				Direction:   TextDirectionLTR,
				Wrap:        TextWrapLineChar,
				Ellipsize:   true,
				Justify:     true,
				Strike:      true,
				Text:        "Hello, world!",
			},
			expected: "font=%E6%96%B0%E3%82%B4%20R%2Csize=12%2Cw=400%2Ch=100%2Clinespacing=1.5%2Calign=1%2Cdir=1%2Cwrap=2%2Cellipsize=1%2Cjustify=1%2Cstrike=1%2Ctext=Hello%2C%20world%21",
		},

		// offset
		{
			text: &Text{
				Font: &Font{
					Name: "新ゴ R",
				},
				Height: 100,
				Width:  400,
				Size:   12,
				Offset: image.Point{X: 10, Y: 20},
				Text:   "Hello, world!",
			},
			expected: "font=%E6%96%B0%E3%82%B4%20R%2Csize=12%2Cw=400%2Ch=100%2Cx=10%2Cy=20%2Ctext=Hello%2C%20world%21",
		},
		{
			text: &Text{
				Font: &Font{
					Name: "新ゴ R",
				},
				Height:      100,
				Width:       400,
				Size:        12,
				OffsetRatio: image.Point{X: 1, Y: 1},
				OffsetMax:   image.Point{X: 2, Y: 2},
				Text:        "Hello, world!",
			},
			expected: "font=%E6%96%B0%E3%82%B4%20R%2Csize=12%2Cw=400%2Ch=100%2Cxr=0.5%2Cyr=0.5%2Ctext=Hello%2C%20world%21",
		},

		// Mask
		{
			text: &Text{
				Font: &Font{
					Name: "新ゴ R",
				},
				Height:      100,
				Width:       400,
				Size:        30,
				Foreground:  color.Black,
				Background:  color.White,
				MaskType:    MaskTypeBlack,
				PaddingMode: PaddingModeLeave,
				Text:        "Hello, world!",
			},
			expected: "font=%E6%96%B0%E3%82%B4%20R%2Csize=30%2Cf=000000%2Cb=ffffff%2Cw=400%2Ch=100%2Cmask=black:1%2Ctext=Hello%2C%20world%21",
		},
	}

	for _, c := range cases {
		if got := c.text.String(); got != c.expected {
			t.Errorf("%#v: Text.String() = %q, want %q", c.text, got, c.expected)
		}
	}
}

func TestParseText(t *testing.T) {
	cases := []struct {
		input    string
		expected *Text
	}{
		{
			input: "font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,text=Hello%2C%20world%21",
			expected: &Text{
				Font: &Font{
					Name: "新ゴ R",
				},
				Height: 100,
				Width:  400,
				Size:   12,
				Text:   "Hello, world!",
			},
		},
	}

	for _, c := range cases {
		got, err := ParseText(c.input)
		if err != nil {
			t.Errorf("ParseText(%q) returned error: %v", c.input, err)
			continue
		}
		if diff := cmp.Diff(got, c.expected); diff != "" {
			t.Errorf("ParseText(%q) = (+got / -expected) %s", c.input, diff)
		}
	}
}

func TestParseText_Error(t *testing.T) {
	cases := []string{
		"=",
	}

	for _, c := range cases {
		if _, err := ParseText(c); err == nil {
			t.Errorf("ParseText(%q) did not return error", c)
		}
	}
}
