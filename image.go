package imageflux

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/url"
	"path"
	"strconv"
	"strings"
)

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
