package imageflux

import (
	"image"
	"image/color"
	"net/url"
	"strconv"
	"testing"
)

// Overlay is the configure of an overlay image.
type Overlay struct {
	// URL is an url for overlay image.
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
	buf = append(buf, "%2F"...)
	buf = append(buf, url.QueryEscape(o.URL)...)
	return buf
}

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
