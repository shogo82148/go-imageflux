package imageflux

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"io"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return []byte{}
	},
}

// Image is an image hosted on ImageFlux.
type Image struct {
	Path   string
	Proxy  *Proxy
	Config *Config
}

// Config is configure of image.
type Config struct {
	// Scaling Parameters.
	Width          int
	Height         int
	DisableEnlarge bool
	AspectMode     AspectMode
	Clip           image.Rectangle
	ClipRatio      image.Rectangle
	ClipMax        image.Point
	Origin         Origin
	Background     color.Color
	Rotate         Rotate
	Through        Through

	// Overlay Parameters.
	Overlay Overlay

	// Output Parameters.
	Format              Format
	Quality             int
	DisableOptimization bool
}

// Overlay is the configure of an overlay image.
type Overlay struct {
	URL         string
	Offset      image.Point
	OffsetRatio image.Point
	OffsetMax   image.Point
	Origin      Origin
}

// AspectMode is aspect mode.
type AspectMode int

const (
	// AspectModeDefault is the default value of aspect mode.
	AspectModeDefault AspectMode = iota

	// AspectModeScale holds the the aspect ratio of the input image,
	// and scales to fit in the specified size.
	AspectModeScale

	// AspectModeForceScale ignores the aspect ratio of the input image.
	AspectModeForceScale

	// AspectModeCrop holds the the aspect ratio of the input image,
	// and crops the image.
	AspectModeCrop

	// AspectModePad holds the the aspect ratio of the input image,
	// and fills the unfilled portion with the specified background color.
	AspectModePad
)

// Origin is the origin.
type Origin int

const (
	OriginDefault Origin = iota
	OriginTopLeft
	OriginTopCenter
	OriginTopRight
	OriginMiddleLeft
	OriginMiddleCenter
	OriginMiddleRight
	OriginBottomLeft
	OriginBottomCenter
	OriginBottomRight
)

func (o Origin) String() string {
	switch o {
	case OriginDefault:
		return "default"
	case OriginTopLeft:
		return "top-left"
	case OriginTopCenter:
		return "top-center"
	case OriginTopRight:
		return "top-right"
	case OriginMiddleLeft:
		return "middle-left"
	case OriginMiddleCenter:
		return "middle-center"
	case OriginMiddleRight:
		return "middle-right"
	case OriginBottomLeft:
		return "bottom-left"
	case OriginBottomCenter:
		return "bottom-center"
	case OriginBottomRight:
		return "bottom-right"
	}
	return ""
}

// Format is the format of the output image.
type Format string

const (
	// FormatAuto encodes the image by the same format with the input image.
	FormatAuto Format = "auto"

	// FormatJPEG encodes the image as a JPEG.
	FormatJPEG Format = "jpg"

	// FormatPNG encodes the image as a PNG.
	FormatPNG Format = "png"

	// FormatGIF encodes the image as a GIF.
	FormatGIF Format = "gif"

	// FormatWebPFromJPEG encodes the image as a WebP.
	// The input image should be a JPEG.
	FormatWebPFromJPEG Format = "webp:jpeg"

	// FormatWebPFromPNG encodes the image as a WebP.
	// The input image should be a PNG.
	FormatWebPFromPNG Format = "webp:png"
)

func (f Format) String() string {
	return string(f)
}

// Rotate rotates the image.
type Rotate int

const (
	// RotateDefault is the default value of Rotate. It is same as RotateTopLeft.
	RotateDefault Rotate = iota

	// RotateTopLeft does not anything.
	RotateTopLeft

	// RotateTopRight flips the image left and right.
	RotateTopRight

	// RotateBottomRight rotates the image 180 degrees.
	RotateBottomRight

	// RotateBottomLeft flips the image upside down.
	RotateBottomLeft

	// RotateLeftTop mirrors the image around the diagonal axis.
	RotateLeftTop

	// RotateRightTop rotates the image left 90 degrees.
	RotateRightTop

	// RotateRightBottom rotates the image 180 degrees and mirrors the image around the diagonal axis.
	RotateRightBottom

	// RotateLeftBottom rotates the image right 90 degrees.
	RotateLeftBottom

	// RotateAuto parses the Orientation of the Exif information and rorates the image.
	RotateAuto Rotate = -1
)

func (r Rotate) String() string {
	switch r {
	case RotateDefault:
		return "default"
	case RotateTopLeft:
		return "top-left"
	case RotateTopRight:
		return "top-right"
	case RotateBottomRight:
		return "bottom-right"
	case RotateBottomLeft:
		return "bottom-left"
	case RotateLeftTop:
		return "left-top"
	case RotateRightTop:
		return "right-top"
	case RotateRightBottom:
		return "right-bottom"
	case RotateLeftBottom:
		return "left-bottom"
	case RotateAuto:
		return "auto"
	}
	return ""
}

// Through is an image format list for skipping converting.
type Through int

const (
	// ThroughJPEG skips converting JPEG images.
	ThroughJPEG Through = 1 << iota

	// ThroughPNG skips converting PNG images.
	ThroughPNG

	// ThroughGIF skips converting GIF images.
	ThroughGIF
)

func (t Through) String() string {
	var buf [12]byte
	return string(t.append(buf[:]))
}

func (t Through) append(buf []byte) []byte {
	if t == 0 {
		return buf
	}
	if (t & ThroughJPEG) != 0 {
		buf = append(buf, "jpg:"...)
	}
	if (t & ThroughPNG) != 0 {
		buf = append(buf, "png:"...)
	}
	if (t & ThroughGIF) != 0 {
		buf = append(buf, "gif:"...)
	}
	return buf[:len(buf)-1]
}

func (c *Config) String() string {
	if c == nil {
		return ""
	}

	buf := bufPool.Get().([]byte)[:0]

	if c.Width != 0 {
		buf = append(buf, 'w', '=')
		buf = strconv.AppendInt(buf, int64(c.Width), 10)
		buf = append(buf, ',')
	}
	if c.Height != 0 {
		buf = append(buf, 'h', '=')
		buf = strconv.AppendInt(buf, int64(c.Height), 10)
		buf = append(buf, ',')
	}
	if c.DisableEnlarge {
		buf = append(buf, 'u', '=', '0', ',')
	}
	if c.AspectMode != AspectModeDefault {
		buf = append(buf, 'a', '=')
		buf = strconv.AppendInt(buf, int64(c.AspectMode-1), 10)
		buf = append(buf, ',')
	}
	if c.Clip != image.ZR {
		buf = append(buf, 'c', '=')
		buf = strconv.AppendInt(buf, int64(c.Clip.Min.X), 10)
		buf = append(buf, ':')
		buf = strconv.AppendInt(buf, int64(c.Clip.Min.Y), 10)
		buf = append(buf, ':')
		buf = strconv.AppendInt(buf, int64(c.Clip.Max.X), 10)
		buf = append(buf, ':')
		buf = strconv.AppendInt(buf, int64(c.Clip.Max.Y), 10)
		buf = append(buf, ',')
	}
	if c.ClipRatio != image.ZR && c.ClipMax != image.ZP {
		x1 := float64(c.ClipRatio.Min.X) / float64(c.ClipMax.X)
		y1 := float64(c.ClipRatio.Min.Y) / float64(c.ClipMax.Y)
		x2 := float64(c.ClipRatio.Max.X) / float64(c.ClipMax.X)
		y2 := float64(c.ClipRatio.Max.Y) / float64(c.ClipMax.Y)
		buf = append(buf, 'c', 'r', '=')
		buf = strconv.AppendFloat(buf, x1, 'f', -1, 64)
		buf = append(buf, ':')
		buf = strconv.AppendFloat(buf, y1, 'f', -1, 64)
		buf = append(buf, ':')
		buf = strconv.AppendFloat(buf, x2, 'f', -1, 64)
		buf = append(buf, ':')
		buf = strconv.AppendFloat(buf, y2, 'f', -1, 64)
		buf = append(buf, ',')
	}
	if c.Origin != OriginDefault {
		buf = append(buf, 'g', '=')
		buf = strconv.AppendInt(buf, int64(c.Origin), 10)
		buf = append(buf, ',')
	}
	if c.Background != nil {
		r, g, b, a := c.Background.RGBA()
		if a == 0xffff {
			buf = append(buf, 'b', '=')
			buf = appendByte(buf, byte(r>>8))
			buf = appendByte(buf, byte(g>>8))
			buf = appendByte(buf, byte(b>>8))
			buf = append(buf, ',')
		} else if a == 0 {
			buf = append(buf, "b=000000,"...)
		} else {
			c := fmt.Sprintf("b=%02x%02x%02x%02x,", r>>8, g>>8, b>>8, a>>8)
			buf = append(buf, c...)
		}
	}
	if c.Rotate != RotateDefault {
		if c.Rotate == RotateAuto {
			buf = append(buf, "r=auto,"...)
		} else {
			buf = append(buf, "r="...)
			buf = strconv.AppendInt(buf, int64(c.Rotate), 10)
			buf = append(buf, ',')
		}
	}
	if c.Through != 0 {
		buf = append(buf, "through="...)
		buf = c.Through.append(buf)
		buf = append(buf, ',')
	}

	if overlay := c.Overlay.String(); overlay != "" {
		buf = append(buf, overlay...)
		buf = append(buf, ',')
	}

	if c.Format != "" {
		buf = append(buf, 'f', '=')
		buf = append(buf, c.Format...)
		buf = append(buf, ',')
	}
	if c.Quality != 0 {
		buf = append(buf, 'q', '=')
		buf = strconv.AppendInt(buf, int64(c.Quality), 10)
		buf = append(buf, ',')
	}
	if c.DisableOptimization {
		buf = append(buf, 'o', '=', '0', ',')
	}

	var ret string
	if len(buf) != 0 {
		ret = string(buf[:len(buf)-1])
	}
	bufPool.Put(buf)
	return ret
}

func appendByte(buf []byte, b byte) []byte {
	const digits = "0123456789abcdef"
	return append(buf, digits[b>>4], digits[b&0x0F])
}

func (a AspectMode) String() string {
	switch a {
	case AspectModeDefault:
		return "default"
	case AspectModeScale:
		return "scale"
	case AspectModeForceScale:
		return "force-scale"
	case AspectModePad:
		return "pad"
	}
	return ""
}

// URL returns the URL of the image.
func (img *Image) URL() *url.URL {
	p := img.Path
	if c := img.Config.String(); c != "" {
		p = path.Join("c", c, p)
	}

	return &url.URL{
		Scheme: "https",
		Host:   img.Proxy.Host,
		Path:   p,
	}
}

// SignedURL returns the URL of the image with the signature.
func (img *Image) SignedURL() *url.URL {
	u, s := img.urlAndSign()
	if s == "" {
		return u
	}

	if strings.HasPrefix(u.Path, "/c/") {
		u.Path = "/c/sig=" + s + "," + u.Path[len("/c/"):]
		return u
	}

	if strings.HasPrefix(u.Path, "/c!/") {
		u.Path = "/c!/sig=" + s + "," + u.Path[len("/c!/"):]
		return u
	}

	u.Path = "/c/sig=" + s + u.Path
	return u
}

// Sign returns the signature.
func (img *Image) Sign() string {
	_, s := img.urlAndSign()
	return s
}

func (img *Image) urlAndSign() (*url.URL, string) {
	u := img.URL()
	if img.Proxy == nil || img.Proxy.Secret == "" {
		return u, ""
	}

	p := u.Path
	if len(p) < 1 || p[0] != '/' {
		p = "/" + p
		u.Path = p
	}
	mac := hmac.New(sha256.New, []byte(img.Proxy.Secret))
	io.WriteString(mac, p)

	return u, "1." + base64.URLEncoding.EncodeToString(mac.Sum(nil))
}

func (img *Image) String() string {
	return img.URL().String()
}

func (o Overlay) String() string {
	var buf []byte
	if o.URL != "" {
		buf = append(buf, 'l', '=')
		buf = append(buf, url.QueryEscape(o.URL)...)
		buf = append(buf, ',')
	}
	if o.Offset != image.ZP {
		buf = append(buf, 'l', 'x', '=')
		buf = strconv.AppendInt(buf, int64(o.Offset.X), 10)
		buf = append(buf, ',', 'l', 'y', '=')
		buf = strconv.AppendInt(buf, int64(o.Offset.Y), 10)
		buf = append(buf, ',')
	}
	if o.OffsetRatio != image.ZP && o.OffsetMax != image.ZP {
		x := float64(o.OffsetRatio.X) / float64(o.OffsetMax.X)
		y := float64(o.OffsetRatio.Y) / float64(o.OffsetMax.Y)
		buf = append(buf, 'l', 'x', 'r', '=')
		buf = strconv.AppendFloat(buf, x, 'f', -1, 64)
		buf = append(buf, ',', 'l', 'y', 'r', '=')
		buf = strconv.AppendFloat(buf, y, 'f', -1, 64)
		buf = append(buf, ',')
	}
	if o.Origin != OriginDefault {
		buf = append(buf, 'l', 'g', '=')
		buf = strconv.AppendInt(buf, int64(o.Origin), 10)
		buf = append(buf, ',')
	}
	if len(buf) == 0 {
		return ""
	}
	return string(buf[:len(buf)-1])
}
