package imageflux

import (
	"image"
	"image/color"
	"reflect"
	"testing"
)

func TestOverlay(t *testing.T) {
	cases := []struct {
		overlay *Overlay
		output  string
	}{
		{
			overlay: &Overlay{
				Width:  100,
				Height: 200,
				Path:   "images/1.png",
			},
			output: "w=100,h=200%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				Width:  100,
				Height: 200,
				Path:   "/images/1.png",
			},
			output: "w=100,h=200%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				DisableEnlarge: true,
				Path:           "images/1.png",
			},
			output: "u=0%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				AspectMode: AspectModeScale,
				Path:       "images/1.png",
			},
			output: "a=0%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				InputClip: image.Rect(100, 200, 300, 400),
				Path:      "images/1.png",
			},
			output: "ic=100:200:300:400%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				InputClipRatio: image.Rect(25, 25, 75, 75),
				ClipMax:        image.Pt(100, 100),
				Path:           "images/1.png",
			},
			output: "icr=0.25:0.25:0.75:0.75%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				InputClip:   image.Rect(100, 200, 300, 400),
				InputOrigin: OriginMiddleCenter,
				Path:        "images/1.png",
			},
			output: "ic=100:200:300:400,ig=5%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				OutputClip: image.Rect(100, 150, 200, 250),
				Path:       "images/1.png",
			},
			output: "oc=100:150:200:250%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				// for backward compatibility,
				// you can use Clip instead of OutputClip.
				Clip: image.Rect(100, 150, 200, 250),
				Path: "images/1.png",
			},
			output: "oc=100:150:200:250%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				// If you specify both Clip and OutputClip,
				// OutputClip is used.
				OutputClip: image.Rect(100, 150, 200, 250),
				Clip:       image.Rect(200, 250, 300, 350),
				Path:       "images/1.png",
			},
			output: "oc=100:150:200:250%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				OutputClipRatio: image.Rect(25, 25, 75, 75),
				ClipMax:         image.Pt(100, 100),
				Path:            "images/1.png",
			},
			output: "ocr=0.25:0.25:0.75:0.75%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				// for backward compatibility,
				// you can use ClipRatio instead of OutputClipRatio.
				ClipRatio: image.Rect(25, 25, 75, 75),
				ClipMax:   image.Pt(100, 100),
				Path:      "images/1.png",
			},
			output: "ocr=0.25:0.25:0.75:0.75%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				// If you specify both ClipRatio and OutputClipRatio,
				// OutputClipRatio is used.
				OutputClipRatio: image.Rect(25, 25, 75, 75),
				ClipRatio:       image.Rect(35, 35, 85, 85),
				ClipMax:         image.Pt(100, 100),
				Path:            "images/1.png",
			},
			output: "ocr=0.25:0.25:0.75:0.75%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				OutputClip:   image.Rect(100, 150, 200, 250),
				OutputOrigin: OriginMiddleCenter,
				Path:         "images/1.png",
			},
			output: "oc=100:150:200:250,og=5%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				Origin: OriginTopLeft,
				Path:   "images/1.png",
			},
			output: "g=1%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				Background: color.Black,
				Path:       "images/1.png",
			},
			output: "b=000000%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				Background: color.NRGBA{R: 255, G: 255, B: 255, A: 128},
				Path:       "images/1.png",
			},
			output: "b=ffffff80%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				InputRotate: RotateLeftBottom,
				Path:        "images/1.png",
			},
			output: "ir=8%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				InputRotate: RotateAuto,
				Path:        "images/1.png",
			},
			output: "ir=auto%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				OutputRotate: RotateLeftBottom,
				Path:         "images/1.png",
			},
			output: "or=8%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				OutputRotate: RotateAuto,
				Path:         "images/1.png",
			},
			output: "or=auto%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				// for backward compatibility,
				// you can use Rotate instead of OutputRotate.
				Rotate: RotateLeftBottom,
				Path:   "images/1.png",
			},
			output: "or=8%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				// If you specify both Rotate and OutputRotate,
				// OutputRotate is used.
				OutputRotate: RotateAuto,
				Rotate:       RotateLeftBottom,
				Path:         "images/1.png",
			},
			output: "or=auto%2Fimages%2F1.png",
		},

		// offset
		{
			overlay: &Overlay{
				Offset: image.Pt(100, 200),
				Path:   "images/1.png",
			},
			output: "x=100,y=200%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				OffsetRatio: image.Pt(25, 75),
				OffsetMax:   image.Pt(100, 100),
				Path:        "images/1.png",
			},
			output: "xr=0.25,yr=0.75%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				OverlayOrigin: OriginBottomCenter,
				Path:          "images/1.png",
			},
			output: "lg=8%2Fimages%2F1.png",
		},
	}
	for _, c := range cases {
		if got := c.overlay.String(); got != c.output {
			t.Errorf("%#v: want %q, got %q", c.overlay, c.output, got)
		}
	}
}

func TestParseOverlay(t *testing.T) {
	cases := []struct {
		input string
		want  *Overlay
	}{
		{
			input: "",
			want: &Overlay{
				Path: "/",
			},
		},
		{
			input: "w=100%2Fimages%2F1.png",
			want: &Overlay{
				Width: 100,
				Path:  "/images/1.png",
			},
		},
		{
			input: "w=100,h=200%2Fimages%2F1.png",
			want: &Overlay{
				Width:  100,
				Height: 200,
				Path:   "/images/1.png",
			},
		},
		{
			input: "a=3%2Fimages%2F1.png",
			want: &Overlay{
				AspectMode: AspectModePad,
				Path:       "/images/1.png",
			},
		},

		// clipping parameters
		{
			input: "ic=100:150:200:250%2Fimages%2F1.png",
			want: &Overlay{
				InputClip: image.Rect(100, 150, 200, 250),
				Path:      "/images/1.png",
			},
		},
		{
			input: "icr=0.25:0.25:0.75:0.75%2Fimages%2F1.png",
			want: &Overlay{
				InputClipRatio: image.Rect(16384, 16384, 49152, 49152),
				ClipMax:        image.Pt(65536, 65536),
				Path:           "/images/1.png",
			},
		},
		{
			input: "ic=100:150:200:250,ig=5%2Fimages%2F1.png",
			want: &Overlay{
				InputClip:   image.Rect(100, 150, 200, 250),
				InputOrigin: OriginMiddleCenter,
				Path:        "/images/1.png",
			},
		},
		{
			input: "oc=100:150:200:250%2Fimages%2F1.png",
			want: &Overlay{
				OutputClip: image.Rect(100, 150, 200, 250),
				Path:       "/images/1.png",
			},
		},
		{
			// for backward compatibility, you can use "c" instead of "oc".
			input: "c=100:150:200:250%2Fimages%2F1.png",
			want: &Overlay{
				OutputClip: image.Rect(100, 150, 200, 250),
				Path:       "/images/1.png",
			},
		},
		{
			input: "ocr=0.25:0.25:0.75:0.75%2Fimages%2F1.png",
			want: &Overlay{
				OutputClipRatio: image.Rect(16384, 16384, 49152, 49152),
				ClipMax:         image.Pt(65536, 65536),
				Path:            "/images/1.png",
			},
		},
		{
			// for backward compatibility, you can use "cr" instead of "ocr".
			input: "cr=0.25:0.25:0.75:0.75%2Fimages%2F1.png",
			want: &Overlay{
				OutputClipRatio: image.Rect(16384, 16384, 49152, 49152),
				ClipMax:         image.Pt(65536, 65536),
				Path:            "/images/1.png",
			},
		},
		{
			input: "oc=100:150:200:250,og=5%2Fimages%2F1.png",
			want: &Overlay{
				OutputClip:   image.Rect(100, 150, 200, 250),
				OutputOrigin: OriginMiddleCenter,
				Path:         "/images/1.png",
			},
		},

		{
			input: "g=1%2Fimages%2F1.png",
			want: &Overlay{
				Origin: OriginTopLeft,
				Path:   "/images/1.png",
			},
		},
		{
			input: "b=000000%2Fimages%2F1.png",
			want: &Overlay{
				Background: color.NRGBA{R: 0, G: 0, B: 0, A: 0xff},
				Path:       "/images/1.png",
			},
		},
		{
			input: "b=ffffff%2Fimages%2F1.png",
			want: &Overlay{
				Background: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
				Path:       "/images/1.png",
			},
		},
		{
			input: "b=FFFFFF%2Fimages%2F1.png",
			want: &Overlay{
				Background: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
				Path:       "/images/1.png",
			},
		},
		{
			input: "b=ff0000%2Fimages%2F1.png",
			want: &Overlay{
				Background: color.NRGBA{R: 0xff, G: 0, B: 0, A: 0xff},
				Path:       "/images/1.png",
			},
		},
		{
			input: "b=00ff00%2Fimages%2F1.png",
			want: &Overlay{
				Background: color.NRGBA{R: 0, G: 0xff, B: 0, A: 0xff},
				Path:       "/images/1.png",
			},
		},
		{
			input: "b=0000ff%2Fimages%2F1.png",
			want: &Overlay{
				Background: color.NRGBA{R: 0, G: 0, B: 0xff, A: 0xff},
				Path:       "/images/1.png",
			},
		},
		{
			input: "b=ffffff00%2Fimages%2F1.png",
			want: &Overlay{
				Background: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x00},
				Path:       "/images/1.png",
			},
		},
		{
			input: "b=ffffff80%2Fimages%2F1.png",
			want: &Overlay{
				Background: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x80},
				Path:       "/images/1.png",
			},
		},

		// rotation
		{
			input: "ir=8%2Fimages%2F1.png",
			want: &Overlay{
				InputRotate: RotateLeftBottom,
				Path:        "/images/1.png",
			},
		},
		{
			input: "ir=auto%2Fimages%2F1.png",
			want: &Overlay{
				InputRotate: RotateAuto,
				Path:        "/images/1.png",
			},
		},
		{
			input: "or=8%2Fimages%2F1.png",
			want: &Overlay{
				OutputRotate: RotateLeftBottom,
				Path:         "/images/1.png",
			},
		},
		{
			input: "or=auto%2Fimages%2F1.png",
			want: &Overlay{
				OutputRotate: RotateAuto,
				Path:         "/images/1.png",
			},
		},
		{
			// for backward compatibility,
			// you can use "r" instead of "or".
			input: "r=8%2Fimages%2F1.png",
			want: &Overlay{
				OutputRotate: RotateLeftBottom,
				Path:         "/images/1.png",
			},
		},
		{
			// for backward compatibility,
			// you can use "r" instead of "or".
			input: "r=auto%2Fimages%2F1.png",
			want: &Overlay{
				OutputRotate: RotateAuto,
				Path:         "/images/1.png",
			},
		},
	}

	for _, c := range cases {
		got, err := ParseOverlay(c.input)
		if err != nil {
			t.Errorf("%q: unexpected %v", c.input, err)
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%q: want %#v, got %#v", c.input, c.want, got)
		}
	}
}
