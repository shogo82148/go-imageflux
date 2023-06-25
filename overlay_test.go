package imageflux

import (
	"image"
	"image/color"
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
				URL:    "images/1.png",
			},
			output: "w=100,h=200%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				DisableEnlarge: true,
				URL:            "images/1.png",
			},
			output: "u=0%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				AspectMode: AspectModeScale,
				URL:        "images/1.png",
			},
			output: "a=0%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				InputClip: image.Rect(100, 200, 300, 400),
				URL:       "images/1.png",
			},
			output: "ic=100:200:300:400%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				InputClipRatio: image.Rect(25, 25, 75, 75),
				ClipMax:        image.Pt(100, 100),
				URL:            "images/1.png",
			},
			output: "icr=0.25:0.25:0.75:0.75%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				InputClip:   image.Rect(100, 200, 300, 400),
				InputOrigin: OriginMiddleCenter,
				URL:         "images/1.png",
			},
			output: "ic=100:200:300:400,ig=5%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				OutputClip: image.Rect(100, 150, 200, 250),
				URL:        "images/1.png",
			},
			output: "oc=100:150:200:250%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				// for backward compatibility,
				// you can use Clip instead of OutputClip.
				Clip: image.Rect(100, 150, 200, 250),
				URL:  "images/1.png",
			},
			output: "oc=100:150:200:250%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				// If you specify both Clip and OutputClip,
				// OutputClip is used.
				OutputClip: image.Rect(100, 150, 200, 250),
				Clip:       image.Rect(200, 250, 300, 350),
				URL:        "images/1.png",
			},
			output: "oc=100:150:200:250%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				OutputClipRatio: image.Rect(25, 25, 75, 75),
				ClipMax:         image.Pt(100, 100),
				URL:             "images/1.png",
			},
			output: "ocr=0.25:0.25:0.75:0.75%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				// for backward compatibility,
				// you can use ClipRatio instead of OutputClipRatio.
				ClipRatio: image.Rect(25, 25, 75, 75),
				ClipMax:   image.Pt(100, 100),
				URL:       "images/1.png",
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
				URL:             "images/1.png",
			},
			output: "ocr=0.25:0.25:0.75:0.75%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				OutputClip:   image.Rect(100, 150, 200, 250),
				OutputOrigin: OriginMiddleCenter,
				URL:          "images/1.png",
			},
			output: "oc=100:150:200:250,og=5%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				Origin: OriginTopLeft,
				URL:    "images/1.png",
			},
			output: "g=1%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				Background: color.Black,
				URL:        "images/1.png",
			},
			output: "b=000000%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				Background: color.NRGBA{R: 255, G: 255, B: 255, A: 128},
				URL:        "images/1.png",
			},
			output: "b=ffffff80%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				InputRotate: RotateLeftBottom,
				URL:         "images/1.png",
			},
			output: "ir=8%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				InputRotate: RotateAuto,
				URL:         "images/1.png",
			},
			output: "ir=auto%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				OutputRotate: RotateLeftBottom,
				URL:          "images/1.png",
			},
			output: "or=8%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				OutputRotate: RotateAuto,
				URL:          "images/1.png",
			},
			output: "or=auto%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				// for backward compatibility,
				// you can use Rotate instead of OutputRotate.
				Rotate: RotateLeftBottom,
				URL:    "images/1.png",
			},
			output: "or=8%2Fimages%2F1.png",
		},
		{
			overlay: &Overlay{
				// If you specify both Rotate and OutputRotate,
				// OutputRotate is used.
				OutputRotate: RotateAuto,
				Rotate:       RotateLeftBottom,
				URL:          "images/1.png",
			},
			output: "or=auto%2Fimages%2F1.png",
		},
	}
	for _, c := range cases {
		if got := c.overlay.String(); got != c.output {
			t.Errorf("%#v: want %q, got %q", c.overlay, c.output, got)
		}
	}
}
