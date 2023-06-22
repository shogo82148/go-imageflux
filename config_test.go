package imageflux

import (
	"image"
	"image/color"
	"reflect"
	"testing"
)

func BenchmarkConfig(b *testing.B) {
	config := &Config{
		Width:          100,
		Height:         100,
		DisableEnlarge: true,
		AspectMode:     AspectModePad,
		Clip:           image.Rect(0, 0, 100, 100),
		ClipRatio:      image.Rect(0, 0, 100, 100),
		ClipMax:        image.Pt(100, 100),
		Origin:         OriginBottomRight,
		Background:     color.Black,
		Rotate:         RotateLeftBottom,
		Through:        ThroughJPEG | ThroughPNG | ThroughGIF,
		Overlays: []Overlay{
			{
				URL:            "http://example.com/",
				Width:          100,
				Height:         100,
				DisableEnlarge: true,
				AspectMode:     AspectModePad,
				Clip:           image.Rect(0, 0, 100, 100),
				ClipRatio:      image.Rect(0, 0, 100, 100),
				ClipMax:        image.Pt(100, 100),
				Origin:         OriginBottomRight,
				Background:     color.Black,
				Rotate:         RotateLeftBottom,
				Offset:         image.Pt(100, 100),
				OffsetRatio:    image.Pt(100, 100),
				OffsetMax:      image.Pt(100, 100),
				OverlayOrigin:  OriginBottomRight,
			},
		},
		Format:              FormatWebPFromPNG,
		Quality:             75,
		DisableOptimization: true,
		Unsharp: Unsharp{
			Radius:    10,
			Sigma:     1.0,
			Gain:      1.0,
			Threshold: 0.5,
		},
		Blur: Blur{
			Radius: 10,
			Sigma:  1.0,
		},
	}
	for i := 0; i < b.N; i++ {
		_ = config.String()
	}
}

func TestConfig(t *testing.T) {
	cases := []struct {
		config *Config
		output string
	}{
		{
			config: nil,
			output: "",
		},
		{
			config: &Config{},
			output: "",
		},
		{
			config: &Config{
				Width: 200,
			},
			output: "w=200",
		},
		{
			config: &Config{
				Height: 200,
			},
			output: "h=200",
		},
		{
			config: &Config{
				Width:  200,
				Height: 200,
			},
			output: "w=200,h=200",
		},
		{
			config: &Config{
				DisableEnlarge: true,
			},
			output: "u=0",
		},
		{
			config: &Config{
				AspectMode: AspectModeScale,
			},
			output: "a=0",
		},
		{
			config: &Config{
				Clip: image.Rect(100, 150, 200, 250),
			},
			output: "c=100:150:200:250",
		},
		{
			config: &Config{
				ClipRatio: image.Rect(25, 25, 75, 75),
				ClipMax:   image.Pt(100, 100),
			},
			output: "cr=0.25:0.25:0.75:0.75",
		},
		{
			config: &Config{
				Origin: OriginTopLeft,
			},
			output: "g=1",
		},
		{
			config: &Config{
				Background: color.Black,
			},
			output: "b=000000",
		},
		{
			config: &Config{
				Background: color.White,
			},
			output: "b=ffffff",
		},
		{
			config: &Config{
				Background: color.NRGBA{R: 255, G: 0, B: 0, A: 255},
			},
			output: "b=ff0000",
		},
		{
			config: &Config{
				Background: color.NRGBA{R: 0, G: 255, B: 0, A: 255},
			},
			output: "b=00ff00",
		},
		{
			config: &Config{
				Background: color.NRGBA{R: 0, G: 0, B: 255, A: 255},
			},
			output: "b=0000ff",
		},
		{
			config: &Config{
				Background: color.NRGBA{R: 255, G: 255, B: 255, A: 0},
			},
			output: "b=000000",
		},
		{
			config: &Config{
				Background: color.NRGBA{R: 0, G: 0, B: 0, A: 128},
			},
			output: "b=00000080",
		},
		{
			config: &Config{
				Rotate: RotateLeftBottom,
			},
			output: "r=8",
		},
		{
			config: &Config{
				Rotate: RotateAuto,
			},
			output: "r=auto",
		},
		{
			config: &Config{
				Through: ThroughJPEG,
			},
			output: "through=jpg",
		},
		{
			config: &Config{
				Through: ThroughJPEG | ThroughGIF,
			},
			output: "through=jpg:gif",
		},
		{
			config: &Config{
				Through: ThroughJPEG | ThroughPNG | ThroughGIF,
			},
			output: "through=jpg:png:gif",
		},
		{
			config: &Config{
				Overlays: []Overlay{{
					URL: "http://example.com/",
				}},
			},
			output: "l=(%2fhttp%3A%2F%2Fexample.com%2F)",
		},
		{
			config: &Config{
				Overlays: []Overlay{{
					URL:    "http://example.com/",
					Offset: image.Pt(100, 200),
				}},
			},
			output: "l=(x=100,y=200%2fhttp%3A%2F%2Fexample.com%2F)",
		},
		{
			config: &Config{
				Overlays: []Overlay{{
					URL:         "http://example.com/",
					OffsetRatio: image.Pt(25, 75),
					OffsetMax:   image.Pt(100, 100),
				}},
			},
			output: "l=(xr=0.25,yr=0.75%2fhttp%3A%2F%2Fexample.com%2F)",
		},
		{
			config: &Config{
				Overlays: []Overlay{{
					URL:           "http://example.com/",
					OverlayOrigin: OriginTopLeft,
				}},
			},
			output: "l=(lg=1%2fhttp%3A%2F%2Fexample.com%2F)",
		},
		{
			config: &Config{
				Overlays: []Overlay{
					{
						URL:    "http://example.com/1.png",
						Offset: image.Pt(100, 200),
					},
					{
						URL:    "http://example.com/2.png",
						Offset: image.Pt(200, 100),
					},
				},
			},
			output: "l=(x=100,y=200%2fhttp%3A%2F%2Fexample.com%2F1.png),l=(x=200,y=100%2fhttp%3A%2F%2Fexample.com%2F2.png)",
		},
		{
			config: &Config{
				Format: FormatWebPFromPNG,
			},
			output: "f=webp:png",
		},
		{
			config: &Config{
				Quality: 75,
			},
			output: "q=75",
		},
		{
			config: &Config{
				DisableOptimization: true,
			},
			output: "o=0",
		},
		{
			config: &Config{
				Unsharp: Unsharp{
					Radius: 10,
					Sigma:  1.0,
				},
			},
			output: "unsharp=10x1",
		},
		{
			config: &Config{
				Unsharp: Unsharp{
					Radius:    10,
					Sigma:     1.0,
					Gain:      1.0,
					Threshold: 0.5,
				},
			},
			output: "unsharp=10x1+1+0.5",
		},
		{
			config: &Config{
				Blur: Blur{
					Radius: 10,
					Sigma:  1.0,
				},
			},
			output: "blur=10x1",
		},
	}

	for _, c := range cases {
		if got := c.config.String(); got != c.output {
			t.Errorf("%#v: want %s, got %s", c.config, c.output, got)
		}
	}
}

func TestParseConfig(t *testing.T) {
	cases := []struct {
		input string
		want  *Config
		rest  string
	}{
		{
			input: "",
			want:  &Config{},
		},
		{
			input: "w=100",
			want: &Config{
				Width: 100,
			},
		},
		{
			input: "w=100,h=200",
			want: &Config{
				Width:  100,
				Height: 200,
			},
		},
		{
			input: "a=3",
			want: &Config{
				AspectMode: AspectModePad,
			},
		},
	}

	for _, c := range cases {
		got, rest, err := ParseConfig(c.input)
		if err != nil {
			t.Errorf("%q: unexpected error: %s", c.input, err)
			continue
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%q: want %#v, got %#v", c.input, c.want, got)
		}
		if rest != c.rest {
			t.Errorf("%q: want %q, got %q", c.input, c.rest, rest)
		}
	}
}
