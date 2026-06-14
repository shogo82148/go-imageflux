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
	"(DriveFlux)Extra",                // extra characters after closing parenthesis
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

		// Basic case
		{
			text: &Text{
				Height: 100,
				Width:  400,
				Size:   12,
				Text:   "Hello, world!",
			},
			expected: "font=%2Csize=12%2Cw=400%2Ch=100%2Ctext=Hello%2C%20world%21",
		},
	}

	for _, c := range cases {
		if got := c.text.String(); got != c.expected {
			t.Errorf("%#v: Text.String() = %q, want %q", c.text, got, c.expected)
		}
	}
}

var parseTextCases = []struct {
	input    string
	expected *Text
}{
	// Basic case
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

	// use %2C instead of comma
	{
		input: "font=%E6%96%B0%E3%82%B4%20R%2Csize=12%2Cw=400%2Ch=100%2Ctext=Hello%2C%20world%21",
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

	// empty font name
	{
		input: "font=,size=12,w=400,h=100,text=Hello%2C%20world%21",
		expected: &Text{
			Font:   &Font{},
			Height: 100,
			Width:  400,
			Size:   12,
			Text:   "Hello, world!",
		},
	},

	// foreground and background colors
	{
		input: "font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,f=000000,b=ffffff,text=Hello%2C%20world%21",
		expected: &Text{
			Font: &Font{
				Name: "新ゴ R",
			},
			Height:     100,
			Width:      400,
			Size:       12,
			Foreground: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
			Background: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
			Text:       "Hello, world!",
		},
	},

	// line spacing
	{
		input: "font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,linespacing=1.5,text=Hello%2C%20world%21",
		expected: &Text{
			Font: &Font{
				Name: "新ゴ R",
			},
			Height:      100,
			Width:       400,
			Size:        12,
			LineSpacing: 1.5,
			Text:        "Hello, world!",
		},
	},

	// align
	{
		input: "font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,align=1,text=Hello%2C%20world%21",
		expected: &Text{
			Font: &Font{
				Name: "新ゴ R",
			},
			Height: 100,
			Width:  400,
			Size:   12,
			Align:  TextAlignCenter,
			Text:   "Hello, world!",
		},
	},

	// direction
	{
		input: "font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,dir=1,text=Hello%2C%20world%21",
		expected: &Text{
			Font: &Font{
				Name: "新ゴ R",
			},
			Height:    100,
			Width:     400,
			Size:      12,
			Direction: TextDirectionLTR,
			Text:      "Hello, world!",
		},
	},

	// wrap
	{
		input: "font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,wrap=1,text=Hello%2C%20world%21",
		expected: &Text{
			Font: &Font{
				Name: "新ゴ R",
			},
			Height: 100,
			Width:  400,
			Size:   12,
			Wrap:   TextWrapChar,
			Text:   "Hello, world!",
		},
	},

	// ellipsize
	{
		input: "font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,ellipsize=0,text=Hello%2C%20world%21",
		expected: &Text{
			Font: &Font{
				Name: "新ゴ R",
			},
			Height:    100,
			Width:     400,
			Size:      12,
			Ellipsize: false,
			Text:      "Hello, world!",
		},
	},
	{
		input: "font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,ellipsize=1,text=Hello%2C%20world%21",
		expected: &Text{
			Font: &Font{
				Name: "新ゴ R",
			},
			Height:    100,
			Width:     400,
			Size:      12,
			Ellipsize: true,
			Text:      "Hello, world!",
		},
	},

	// justify
	{
		input: "font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,justify=1,text=Hello%2C%20world%21",
		expected: &Text{
			Font: &Font{
				Name: "新ゴ R",
			},
			Height:  100,
			Width:   400,
			Size:    12,
			Justify: true,
			Text:    "Hello, world!",
		},
	},

	// strike
	{
		input: "font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,strike=1,text=Hello%2C%20world%21",
		expected: &Text{
			Font: &Font{
				Name: "新ゴ R",
			},
			Height: 100,
			Width:  400,
			Size:   12,
			Strike: true,
			Text:   "Hello, world!",
		},
	},
}

func TestParseText(t *testing.T) {
	for _, c := range parseTextCases {
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

var parseTextErrorCases = []string{
	// invalid font name
	"font=%XX,size=12,w=400,h=100,text=Hello%2C%20world%21",

	// invalid percent encoding in text
	"font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,text=Hello%2C%20world%XX",

	// missing text parameter
	"font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100",

	// invalid color
	"font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,f=ZZZZZZ,text=Hello%2C%20world%21",
	"font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,b=ZZZZZZ,text=Hello%2C%20world%21",

	// size: syntax error
	"font=%E6%96%B0%E3%82%B4%20R,size=invalid,w=400,h=100,text=Hello%2C%20world%21",

	// width and height: syntax error
	"font=%E6%96%B0%E3%82%B4%20R,size=12,w=invalid,h=100,text=Hello%2C%20world%21",
	"font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=invalid,text=Hello%2C%20world%21",

	// line spacing: syntax error
	"font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,linespacing=invalid,text=Hello%2C%20world%21",

	// align: syntax error
	"font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,align=invalid,text=Hello%2C%20world%21",

	// direction: syntax error
	"font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,dir=invalid,text=Hello%2C%20world%21",

	// wrap: syntax error
	"font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,wrap=invalid,text=Hello%2C%20world%21",

	// out of range
	"font=%E6%96%B0%E3%82%B4%20R,size=-1,w=-1,h=-1,text=Hello%2C%20world%21",
	"font=%E6%96%B0%E3%82%B4%20R,size=NaN,w=400,h=100,text=Hello%2C%20world%21",
	"font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,align=-1,dir=-1,wrap=-1,text=Hello%2C%20world%21",
	"font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,linespacing=NaN,text=Hello%2C%20world%21",

	// ellipsize: syntax error
	"font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,ellipsize=invalid,text=Hello%2C%20world%21",

	// justify: syntax error
	"font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,justify=invalid,text=Hello%2C%20world%21",

	// strike: syntax error
	"font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,strike=invalid,text=Hello%2C%20world%21",

	// unknown key
	"font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,unknown=value,text=Hello%2C%20world%21",
}

func TestParseText_Error(t *testing.T) {
	for _, c := range parseTextErrorCases {
		if _, err := ParseText(c); err == nil {
			t.Errorf("ParseText(%q) did not return error", c)
		}
	}
}

func FuzzParseText(f *testing.F) {
	for _, c := range parseTextCases {
		f.Add(c.input)
	}
	for _, c := range parseTextErrorCases {
		f.Add(c)
	}

	f.Fuzz(func(t *testing.T, s string) {
		text, err := ParseText(s)
		if err != nil {
			return
		}
		s1 := text.String()
		text2, err := ParseText(s1)
		if err != nil {
			t.Errorf("ParseText(%q) returned error: %v", s1, err)
			return
		}
		if diff := cmp.Diff(text, text2); diff != "" {
			t.Errorf("ParseText(%q) = (-text / +text2) %s", s, diff)
		}
	})
}
