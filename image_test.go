package imageflux

import (
	"testing"
)

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

func TestImage(t *testing.T) {
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
			"https://demo.imageflux.jp/images/1.jpg",
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
			"https://demo.imageflux.jp/c/sig=1.-Yd8m-5pXPihiZdlDATcwkkgjzPIC9gFHmmZ3JMxwS0=/images/1.jpg",
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
			"https://demo.imageflux.jp/c/sig=1.-Yd8m-5pXPihiZdlDATcwkkgjzPIC9gFHmmZ3JMxwS0=/images/1.jpg",
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
	}

	for _, c := range cases {
		if got := c.image.SignedURL(); got != c.output {
			t.Errorf("want %s, got %s", c.output, got)
		}
	}
}
