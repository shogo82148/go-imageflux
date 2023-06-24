package imageflux

import (
	"testing"
	"time"
)

var jst *time.Location = time.FixedZone("Asia/Tokyo", 9*60*60)

func BenchmarkImage(b *testing.B) {
	img := &Image{
		Proxy: &Proxy{
			Host:   "demo.imageflux.jp",
			Secret: "testsigningsecret",
		},
		Path: "/images/1.jpg",
		Config: &Config{
			Width: 200,
		},
	}
	for i := 0; i < b.N; i++ {
		img.SignedURL()
	}
}

func TestImage_SignedURL(t *testing.T) {
	cases := []struct {
		image  *Image
		output string
	}{
		{
			&Image{
				Proxy: &Proxy{
					Host: "demo.imageflux.jp",
				},
				Path: "/images/1.jpg",
			},
			"https://demo.imageflux.jp/c/f=auto/images/1.jpg",
		},
		{
			&Image{
				Proxy: &Proxy{
					Host: "demo.imageflux.jp",
				},
				Path: "/images/1.jpg",
				Config: &Config{
					Width: 200,
				},
			},
			"https://demo.imageflux.jp/c/w=200/images/1.jpg",
		},
		{
			&Image{
				Proxy: &Proxy{
					Host:   "demo.imageflux.jp",
					Secret: "testsigningsecret",
				},
				Path: "/images/1.jpg",
			},
			"https://demo.imageflux.jp/c/sig=1.tbCHoq4CHTiwxkfATFMnYqrJ7jcjG4D34B_oPQkzf-k=,f=auto/images/1.jpg",
		},
		{
			&Image{
				Proxy: &Proxy{
					Host:   "demo.imageflux.jp",
					Secret: "testsigningsecret",
				},
				Path:   "/images/1.jpg",
				Config: &Config{},
			},
			"https://demo.imageflux.jp/c/sig=1.tbCHoq4CHTiwxkfATFMnYqrJ7jcjG4D34B_oPQkzf-k=,f=auto/images/1.jpg",
		},
		{
			&Image{
				Proxy: &Proxy{
					Host:   "demo.imageflux.jp",
					Secret: "testsigningsecret",
				},
				Path: "/images/1.jpg",
				Config: &Config{
					Width: 200,
				},
			},
			"https://demo.imageflux.jp/c/sig=1.tiKX5u2kw6wp9zDgl1tLiOIi8IsoRIBw8fVgVc0yrNg=,w=200/images/1.jpg",
		},
		{
			&Image{
				Proxy: &Proxy{
					Host: "demo.imageflux.jp",
				},
				Path:    "/images/1.jpg",
				Expires: time.Date(2023, 6, 24, 18, 23, 0, 123456789, jst),
			},
			"https://demo.imageflux.jp/c/f=auto,expires=2023-06-24T09:23:00Z/images/1.jpg",
		},
		{
			&Image{
				Proxy: &Proxy{
					Host:   "demo.imageflux.jp",
					Secret: "testsigningsecret",
				},
				Path: "/images/1.jpg",
				Config: &Config{
					Width: 200,
				},
				Expires: time.Date(2023, 6, 24, 18, 23, 0, 123456789, jst),
			},
			"https://demo.imageflux.jp/c/sig=1.dFGx33tPqUTZLhzxcbOY5_f-afI9EBDga8rwbmMsW2o=,w=200,expires=2023-06-24T09:23:00Z/images/1.jpg",
		},
		{
			&Image{
				Proxy: &Proxy{
					Host: "demo.imageflux.jp",
				},
				Path: "/bridge.jpg",
				Config: &Config{
					Width: 400,
					Overlays: []Overlay{
						{
							Width: 300,
							URL:   "images/1.png",
						},
					},
					Format: FormatWebPAuto,
				},
			},
			"https://demo.imageflux.jp/c/w=400,l=(w=300%2Fimages%2F1.png),f=webp:auto/bridge.jpg",
		},
	}

	for _, c := range cases {
		if got := c.image.SignedURL(); got != c.output {
			t.Errorf("want %s, got %s", c.output, got)
		}
	}
}

func TestImage_SignedURLWithoutComma(t *testing.T) {
	cases := []struct {
		image  *Image
		output string
	}{
		{
			&Image{
				Proxy: &Proxy{
					Host: "demo.imageflux.jp",
				},
				Path: "/images/1.jpg",
			},
			"https://demo.imageflux.jp/c/f=auto/images/1.jpg",
		},
		{
			&Image{
				Proxy: &Proxy{
					Host: "demo.imageflux.jp",
				},
				Path: "/images/1.jpg",
				Config: &Config{
					Width: 200,
				},
			},
			"https://demo.imageflux.jp/c/w=200/images/1.jpg",
		},
		{
			&Image{
				Proxy: &Proxy{
					Host: "demo.imageflux.jp",
				},
				Path: "/images/1.jpg",
				Config: &Config{
					Width:  200,
					Height: 200,
				},
			},
			"https://demo.imageflux.jp/c/w=200%2Ch=200/images/1.jpg",
		},
		{
			&Image{
				Proxy: &Proxy{
					Host:   "demo.imageflux.jp",
					Secret: "testsigningsecret",
				},
				Path: "/images/1.jpg",
			},
			"https://demo.imageflux.jp/c/sig=1.tbCHoq4CHTiwxkfATFMnYqrJ7jcjG4D34B_oPQkzf-k=%2Cf=auto/images/1.jpg",
		},
		{
			&Image{
				Proxy: &Proxy{
					Host:   "demo.imageflux.jp",
					Secret: "testsigningsecret",
				},
				Path:   "/images/1.jpg",
				Config: &Config{},
			},
			"https://demo.imageflux.jp/c/sig=1.tbCHoq4CHTiwxkfATFMnYqrJ7jcjG4D34B_oPQkzf-k=%2Cf=auto/images/1.jpg",
		},
		{
			&Image{
				Proxy: &Proxy{
					Host:   "demo.imageflux.jp",
					Secret: "testsigningsecret",
				},
				Path: "/images/1.jpg",
				Config: &Config{
					Width: 200,
				},
			},
			"https://demo.imageflux.jp/c/sig=1.tiKX5u2kw6wp9zDgl1tLiOIi8IsoRIBw8fVgVc0yrNg=%2Cw=200/images/1.jpg",
		},
		{
			&Image{
				Proxy: &Proxy{
					Host: "demo.imageflux.jp",
				},
				Path:    "/images/1.jpg",
				Expires: time.Date(2023, 6, 24, 18, 23, 0, 123456789, jst),
			},
			"https://demo.imageflux.jp/c/f=auto%2Cexpires=2023-06-24T09:23:00Z/images/1.jpg",
		},
		{
			&Image{
				Proxy: &Proxy{
					Host:   "demo.imageflux.jp",
					Secret: "testsigningsecret",
				},
				Path: "/images/1.jpg",
				Config: &Config{
					Width: 200,
				},
				Expires: time.Date(2023, 6, 24, 18, 23, 0, 123456789, jst),
			},
			"https://demo.imageflux.jp/c/sig=1.Aa05y5VnlhocCF-RABA2--P7-4kc8E9LqJ86BqGosqw=%2Cw=200%2Cexpires=2023-06-24T09:23:00Z/images/1.jpg",
		},
		{
			&Image{
				Proxy: &Proxy{
					Host: "demo.imageflux.jp",
				},
				Path: "/bridge.jpg",
				Config: &Config{
					Width: 400,
					Overlays: []Overlay{
						{
							Width: 300,
							URL:   "images/1.png",
						},
					},
					Format: FormatWebPAuto,
				},
			},
			"https://demo.imageflux.jp/c/w=400%2Cl=(w=300%2Fimages%2F1.png)%2Cf=webp:auto/bridge.jpg",
		},
	}

	for _, c := range cases {
		if got := c.image.SignedURLWithoutComma(); got != c.output {
			t.Errorf("want %s, got %s", c.output, got)
		}
	}
}

func TestImage_String(t *testing.T) {
	cases := []struct {
		image  *Image
		output string
	}{
		{
			&Image{
				Proxy: &Proxy{
					Host: "demo.imageflux.jp",
				},
				Path: "/images/1.jpg",
			},
			"https://demo.imageflux.jp/c/f=auto/images/1.jpg",
		},
		{
			&Image{
				Proxy: &Proxy{
					Host: "demo.imageflux.jp",
				},
				Path: "/images/1.jpg",
				Config: &Config{
					Width: 200,
				},
			},
			"https://demo.imageflux.jp/c/w=200/images/1.jpg",
		},
		{
			&Image{
				Proxy: &Proxy{
					Host: "demo.imageflux.jp",
				},
				Path:    "/images/1.jpg",
				Expires: time.Date(2023, 6, 24, 18, 23, 0, 123456789, jst),
			},
			"https://demo.imageflux.jp/c/f=auto,expires=2023-06-24T09:23:00Z/images/1.jpg",
		},
	}

	for _, c := range cases {
		if got := c.image.String(); got != c.output {
			t.Errorf("want %s, got %s", c.output, got)
		}
	}
}
