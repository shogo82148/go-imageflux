package imageflux

import (
	"net/url"
	"path"
	"strconv"
)

// Image is an image hosted on ImageFlux.
type Image struct {
	Host   string
	Path   string
	Config *Config
}

// Config is configure of image.
type Config struct {
	// Scaling Parameters.
	Width          int
	Height         int
	DisableEnlarge bool
	AspectMode     AspectMode

	// Signed URL Parameters.
	Secret string
}

type AspectMode int

const (
	AspectModeDefault AspectMode = iota
	AspectModeScale
	AspectModeForceScale
	AspectModeCrop
	AspectModePad
)

func (c *Config) String() string {
	if c == nil {
		return ""
	}

	var buf []byte
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

	if len(buf) == 0 {
		return ""
	}
	return string(buf[:len(buf)-1])
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

func (img *Image) URL() *url.URL {
	p := img.Path
	if c := img.Config.String(); c != "" {
		p = path.Join("c", c, p)
	}

	return &url.URL{
		Scheme: "https",
		Host:   img.Host,
		Path:   p,
	}
}

func (img *Image) SignedURL() *url.URL {
	return nil
}

func (img *Image) Sign() string {
	return ""
}

func (img *Image) String() string {
	return img.URL().String()
}
