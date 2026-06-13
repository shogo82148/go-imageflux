package imageflux

import (
	"image"
	"image/color"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseFont(t *testing.T) {
	cases := []struct {
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
	}

	for _, c := range cases {
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

func TestText(t *testing.T) {
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
			expected: "font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,text=Hello%2C%20world%21",
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
			expected: "font=(DriveFlux,instance=B%20Italic),size=12,w=400,h=100,text=Hello%2C%20world%21",
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
			expected: "font=(DriveFlux,var=CNTR:0,var=SMTH:0,var=slnt:-16,var=wght:700),size=12,w=400,h=100,text=Hello%2C%20world%21",
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
			expected: "font=%E6%96%B0%E3%82%B4%20R,size=12,f=112233,b=445566,w=400,h=100,text=Hello%2C%20world%21",
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
			expected: "font=%E6%96%B0%E3%82%B4%20R,size=12,f=1122337f,b=4455667f,w=400,h=100,text=Hello%2C%20world%21",
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
			expected: "font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,linespacing=1.5,align=1,dir=1,wrap=2,ellipsize=1,justify=1,strike=1,text=Hello%2C%20world%21",
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
			expected: "font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,x=10,y=20,text=Hello%2C%20world%21",
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
			expected: "font=%E6%96%B0%E3%82%B4%20R,size=12,w=400,h=100,xr=0.5,yr=0.5,text=Hello%2C%20world%21",
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
			expected: "font=%E6%96%B0%E3%82%B4%20R,size=30,f=000000,b=ffffff,w=400,h=100,mask=black:1,text=Hello%2C%20world%21",
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
