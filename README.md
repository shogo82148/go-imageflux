# go-imageflux

[![Build Status](https://github.com/shogo82148/go-imageflux/workflows/test/badge.svg?branch=main)](https://github.com/shogo82148/go-imageflux/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/shogo82148/go-imageflux.svg)](https://pkg.go.dev/github.com/shogo82148/go-imageflux)

URL builder and parser for [ImageFlux](https://imageflux.sakura.ad.jp/).

## Usage

[ImageFlux](https://imageflux.sakura.ad.jp/) is Image Conversion & Distribution Engine.
This allows you to easily generate images optimized for each device based on a single source image,
and delivers them quickly and with high quality.

The imageflux package builds and parse URLs for ImageFlux.

### Build URL

In ImageFlux, parameters for image transformation are embedded in the URL.

```go
proxy := &imageflux.Proxy{
    Host: "demo.imageflux.jp",
}
cfg := &imageflux.Config{
    // resize the image to 200px width.
    Width: 200,

    // convert the image to WebP format.
    Format: imageflux.FormatWebPAuto,
}
u := proxy.Image("/images/1.jpg", cfg).SignedURL()
fmt.Println(u)

// Output:
// https://demo.imageflux.jp/c/w=200,f=webp:auto/images/1.jpg
```

### Build Signed URL

By attaching a signature to the transformation parameters,
it prevents third parties from rewriting the URL.

```go
proxy := &imageflux.Proxy{
    Host:   "demo.imageflux.jp",
    Secret: "testsigningsecret",
}
cfg := &imageflux.Config{
    // resize the image to 200px width.
    Width: 200,
}
u := proxy.Image("/images/1.jpg", cfg).SignedURL()
fmt.Println(u)

// Output:
// https://demo.imageflux.jp/c/sig=1.tiKX5u2kw6wp9zDgl1tLiOIi8IsoRIBw8fVgVc0yrNg=,w=200/images/1.jpg
```

### Parse URL

```go
proxy := &imageflux.Proxy{}
image, err := proxy.Parse("/c/w=200/images/1.jpg", "")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("path = %s\n", image.Path)
fmt.Printf("width = %d\n", image.Config.Width)

// Output:
// path = /images/1.jpg
// width = 200
```

## References

- [ImageFlux](https://imageflux.sakura.ad.jp/) (written in Japanese)
- [The document of ImageFlux](https://console.imageflux.jp/docs/) (written in Japanese)
