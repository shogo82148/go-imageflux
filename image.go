package imageflux

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"sync"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 0, 32)
		return &buf
	},
}

// Image is an image hosted on ImageFlux.
type Image struct {
	Path   string
	Proxy  *Proxy
	Config *Config
}

// SignedURL returns the URL of the image with the signature.
func (img *Image) SignedURL() string {
	path, s := img.pathAndSign()
	if s == "" {
		return "https://" + img.Proxy.Host + path
	}
	if strings.HasPrefix(path, "/c/") {
		return "https://" + img.Proxy.Host + "/c/sig=" + s + "," + strings.TrimPrefix(path, "/c/")
	}
	return "https://" + img.Proxy.Host + "/c/sig=" + s + path
}

// Sign returns the signature.
func (img *Image) Sign() string {
	_, s := img.pathAndSign()
	return s
}

func (img *Image) pathAndSign() (string, string) {
	pbuf := bufPool.Get().(*[]byte)
	buf := (*pbuf)[:0]
	buf = append(buf, "/c/"...)
	buf = img.Config.append(buf)
	if len(buf) == len("/c/") {
		buf = buf[:0]
	}
	if len(img.Path) == 0 || img.Path[0] != '/' {
		buf = append(buf, '/')
	}
	buf = append(buf, img.Path...)
	path := string(buf)

	if img.Proxy.Secret == "" {
		*pbuf = buf
		bufPool.Put(pbuf)
		return path, ""
	}

	mac := hmac.New(sha256.New, []byte(img.Proxy.Secret))
	mac.Write(buf)
	buf = mac.Sum(buf[:0])
	buf2 := make([]byte, len("1.")+base64.URLEncoding.EncodedLen(len(buf)))
	buf2[0] = '1'
	buf2[1] = '.'
	base64.URLEncoding.Encode(buf2[2:], buf)

	*pbuf = buf
	bufPool.Put(pbuf)
	return path, string(buf2[:])
}

func (img *Image) String() string {
	pbuf := bufPool.Get().(*[]byte)
	buf := (*pbuf)[:0]

	buf = append(buf, "https://"...)
	buf = append(buf, img.Proxy.Host...)
	buf = append(buf, "/c/"...)
	buf = img.Config.append(buf)
	if len(img.Path) == 0 || img.Path[0] != '/' {
		buf = append(buf, '/')
	}
	buf = append(buf, img.Path...)
	str := string(buf)
	*pbuf = buf
	bufPool.Put(pbuf)
	return str
}
