package imageflux

import (
	"image"
	"image/color"
	"reflect"
	"testing"
)

func BenchmarkConfig(b *testing.B) {
	config := &Config{
		Width:           100,
		Height:          100,
		DisableEnlarge:  true,
		AspectMode:      AspectModePad,
		InputClip:       image.Rect(0, 0, 100, 100),
		InputClipRatio:  image.Rect(0, 0, 100, 100),
		OutputClip:      image.Rect(0, 0, 100, 100),
		OutputClipRatio: image.Rect(0, 0, 100, 100),
		ClipMax:         image.Pt(100, 100),
		Origin:          OriginBottomRight,
		Background:      color.Black,
		InputRotate:     RotateLeftBottom,
		OutputRotate:    RotateLeftBottom,
		Through:         ThroughJPEG | ThroughPNG | ThroughGIF | ThroughWebP,
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
				DevicePixelRatio: 5,
			},
			output: "dpr=5",
		},

		// clipping parameters
		{
			config: &Config{
				InputClip: image.Rect(100, 150, 200, 250),
			},
			output: "ic=100:150:200:250",
		},
		{
			config: &Config{
				InputClipRatio: image.Rect(25, 25, 75, 75),
				ClipMax:        image.Pt(100, 100),
			},
			output: "icr=0.25:0.25:0.75:0.75",
		},
		{
			config: &Config{
				InputClip:   image.Rect(100, 150, 200, 250),
				InputOrigin: OriginMiddleCenter,
			},
			output: "ic=100:150:200:250,ig=5",
		},
		{
			config: &Config{
				OutputClip: image.Rect(100, 150, 200, 250),
			},
			output: "oc=100:150:200:250",
		},
		{
			config: &Config{
				// for backward compatibility,
				// you can use Clip instead of OutputClip.
				Clip: image.Rect(100, 150, 200, 250),
			},
			output: "oc=100:150:200:250",
		},
		{
			config: &Config{
				// If you specify both Clip and OutputClip,
				// OutputClip is used.
				OutputClip: image.Rect(100, 150, 200, 250),
				Clip:       image.Rect(200, 250, 300, 350),
			},
			output: "oc=100:150:200:250",
		},
		{
			config: &Config{
				OutputClipRatio: image.Rect(25, 25, 75, 75),
				ClipMax:         image.Pt(100, 100),
			},
			output: "ocr=0.25:0.25:0.75:0.75",
		},
		{
			config: &Config{
				// for backward compatibility,
				// you can use ClipRatio instead of OutputClipRatio.
				ClipRatio: image.Rect(25, 25, 75, 75),
				ClipMax:   image.Pt(100, 100),
			},
			output: "ocr=0.25:0.25:0.75:0.75",
		},
		{
			config: &Config{
				// If you specify both ClipRatio and OutputClipRatio,
				// OutputClipRatio is used.
				OutputClipRatio: image.Rect(25, 25, 75, 75),
				ClipRatio:       image.Rect(35, 35, 85, 85),
				ClipMax:         image.Pt(100, 100),
			},
			output: "ocr=0.25:0.25:0.75:0.75",
		},
		{
			config: &Config{
				OutputClip:   image.Rect(100, 150, 200, 250),
				OutputOrigin: OriginMiddleCenter,
			},
			output: "oc=100:150:200:250,og=5",
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
			output: "b=ffffff00",
		},
		{
			config: &Config{
				Background: color.NRGBA{R: 255, G: 255, B: 255, A: 128},
			},
			output: "b=ffffff80",
		},

		// rotation
		{
			config: &Config{
				InputRotate: RotateLeftBottom,
			},
			output: "ir=8",
		},
		{
			config: &Config{
				InputRotate: RotateAuto,
			},
			output: "ir=auto",
		},
		{
			config: &Config{
				OutputRotate: RotateLeftBottom,
			},
			output: "or=8",
		},
		{
			config: &Config{
				OutputRotate: RotateAuto,
			},
			output: "or=auto",
		},
		{
			config: &Config{
				// for backward compatibility,
				// you can use Rotate instead of OutputRotate.
				Rotate: RotateLeftBottom,
			},
			output: "or=8",
		},
		{
			config: &Config{
				// for backward compatibility,
				// you can use Rotate instead of OutputRotate.
				Rotate: RotateAuto,
			},
			output: "or=auto",
		},
		{
			config: &Config{
				// If you specify both Rotate and OutputRotate,
				// OutputRotate is used.
				OutputRotate: RotateAuto,
				Rotate:       RotateLeftBottom,
			},
			output: "or=auto",
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
				Through: ThroughJPEG | ThroughPNG | ThroughGIF | ThroughWebP,
			},
			output: "through=jpg:png:gif:webp",
		},

		// overlays
		{
			config: &Config{
				Overlays: []Overlay{{
					URL: "images/1.png",
				}},
			},
			output: "l=(%2Fimages%2F1.png)",
		},
		{
			config: &Config{
				Overlays: []Overlay{{
					URL:    "images/1.png",
					Offset: image.Pt(100, 200),
				}},
			},
			output: "l=(x=100,y=200%2Fimages%2F1.png)",
		},
		{
			config: &Config{
				Overlays: []Overlay{{
					URL:         "images/1.png",
					OffsetRatio: image.Pt(25, 75),
					OffsetMax:   image.Pt(100, 100),
				}},
			},
			output: "l=(xr=0.25,yr=0.75%2Fimages%2F1.png)",
		},
		{
			config: &Config{
				Overlays: []Overlay{{
					URL:           "images/1.png",
					OverlayOrigin: OriginTopLeft,
				}},
			},
			output: "l=(lg=1%2Fimages%2F1.png)",
		},
		{
			config: &Config{
				Overlays: []Overlay{
					{
						URL:    "images/1.png",
						Offset: image.Pt(100, 200),
					},
					{
						URL:    "images/2.png",
						Offset: image.Pt(200, 100),
					},
				},
			},
			output: "l=(x=100,y=200%2Fimages%2F1.png),l=(x=200,y=100%2Fimages%2F2.png)",
		},
		{
			config: &Config{
				Overlays: []Overlay{{
					URL:      "images/1.png",
					MaskType: MaskTypeWhite,
				}},
			},
			output: "l=(mask=white%2Fimages%2F1.png)",
		},
		{
			config: &Config{
				Overlays: []Overlay{{
					URL:         "images/1.png",
					MaskType:    MaskTypeAlpha,
					PaddingMode: PaddingModeLeave,
				}},
			},
			output: "l=(mask=alpha:1%2Fimages%2F1.png)",
		},

		// output format
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
				Lossless: true,
			},
			output: "lossless=1",
		},
		{
			config: &Config{
				ExifOption: ExifOptionKeepOrientation,
			},
			output: "s=2",
		},

		// image filters
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
		{
			config: &Config{
				GrayScale: 100,
			},
			output: "grayscale=100",
		},
		{
			config: &Config{
				Sepia: 100,
			},
			output: "sepia=100",
		},
		{
			config: &Config{
				Brightness: 100,
			},
			output: "brightness=200",
		},
		{
			config: &Config{
				Contrast: 100,
			},
			output: "contrast=200",
		},
		{
			config: &Config{
				Invert: true,
			},
			output: "invert=1",
		},
	}

	for _, c := range cases {
		if got := c.config.String(); got != c.output {
			t.Errorf("%#v: want %q, got %q", c.config, c.output, got)
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
		{
			input: "dpr=5",
			want: &Config{
				DevicePixelRatio: 5,
			},
		},

		// clipping parameters
		{
			input: "ic=100:150:200:250",
			want: &Config{
				InputClip: image.Rect(100, 150, 200, 250),
			},
		},
		{
			input: "icr=0.25:0.25:0.75:0.75",
			want: &Config{
				InputClipRatio: image.Rect(25, 25, 75, 75),
				ClipMax:        image.Pt(100, 100),
			},
		},
		{
			input: "ic=100:150:200:250,ig=5",
			want: &Config{
				InputClip:   image.Rect(100, 150, 200, 250),
				InputOrigin: OriginMiddleCenter,
			},
		},
		{
			input: "oc=100:150:200:250",
			want: &Config{
				OutputClip: image.Rect(100, 150, 200, 250),
			},
		},
		{
			// for backward compatibility, you can use "c" instead of "oc".
			input: "c=100:150:200:250",
			want: &Config{
				OutputClip: image.Rect(100, 150, 200, 250),
			},
		},
		{
			input: "ocr=0.25:0.25:0.75:0.75",
			want: &Config{
				OutputClipRatio: image.Rect(25, 25, 75, 75),
				ClipMax:         image.Pt(100, 100),
			},
		},
		{
			// for backward compatibility, you can use "cr" instead of "ocr".
			input: "cr=0.25:0.25:0.75:0.75",
			want: &Config{
				OutputClipRatio: image.Rect(25, 25, 75, 75),
				ClipMax:         image.Pt(100, 100),
			},
		},
		{
			input: "oc=100:150:200:250,og=5",
			want: &Config{
				OutputClip:   image.Rect(100, 150, 200, 250),
				OutputOrigin: OriginMiddleCenter,
			},
		},

		{
			input: "g=1",
			want: &Config{
				Origin: OriginTopLeft,
			},
		},
		{
			input: "b=000000",
			want: &Config{
				Background: color.NRGBA{
					R: 0,
					G: 0,
					B: 0,
					A: 0xff,
				},
			},
		},
		{
			input: "b=ffffff",
			want: &Config{
				Background: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
			},
		},
		{
			input: "b=FFFFFF",
			want: &Config{
				Background: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
			},
		},
		{
			input: "b=ff0000",
			want: &Config{
				Background: color.NRGBA{R: 0xff, G: 0, B: 0, A: 0xff},
			},
		},
		{
			input: "b=00ff00",
			want: &Config{
				Background: color.NRGBA{R: 0, G: 0xff, B: 0, A: 0xff},
			},
		},
		{
			input: "b=0000ff",
			want: &Config{
				Background: color.NRGBA{R: 0, G: 0, B: 0xff, A: 0xff},
			},
		},
		{
			input: "b=ffffff00",
			want: &Config{
				Background: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x00},
			},
		},
		{
			input: "b=ffffff80",
			want: &Config{
				Background: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x80},
			},
		},

		// rotation
		{
			input: "ir=8",
			want: &Config{
				InputRotate: RotateLeftBottom,
			},
		},
		{
			input: "ir=auto",
			want: &Config{
				InputRotate: RotateAuto,
			},
		},
		{
			input: "or=8",
			want: &Config{
				OutputRotate: RotateLeftBottom,
			},
		},
		{
			input: "or=auto",
			want: &Config{
				OutputRotate: RotateAuto,
			},
		},
		{
			// for backward compatibility,
			// you can use "r" instead of "or".
			input: "r=8",
			want: &Config{
				OutputRotate: RotateLeftBottom,
			},
		},
		{
			input: "r=auto",
			want: &Config{
				OutputRotate: RotateAuto,
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
