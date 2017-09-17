package imageflux

import (
	"image"
	"image/color"
	"testing"
)

func TestImage(t *testing.T) {
	cases := []struct {
		image  *Image
		output string
	}{
		{
			&Image{
				Proxy: &Proxy{
					Host: "p1-47e91401.imageflux.jp",
				},
				Path: "/images/1.jpg",
			},
			"https://p1-47e91401.imageflux.jp/images/1.jpg",
		},
		{
			&Image{
				Proxy: &Proxy{
					Host: "p1-47e91401.imageflux.jp",
				},
				Path: "/images/1.jpg",
				Config: &Config{
					Width: 200,
				},
			},
			"https://p1-47e91401.imageflux.jp/c/w=200/images/1.jpg",
		},
		{
			&Image{
				Proxy: &Proxy{
					Host:   "p1-47e91401.imageflux.jp",
					Secret: "testsigningsecret",
				},
				Path: "/images/1.jpg",
			},
			"https://p1-47e91401.imageflux.jp/c/sig=1.-Yd8m-5pXPihiZdlDATcwkkgjzPIC9gFHmmZ3JMxwS0=/images/1.jpg",
		},
		{
			&Image{
				Proxy: &Proxy{
					Host:   "p1-47e91401.imageflux.jp",
					Secret: "testsigningsecret",
				},
				Path: "/images/1.jpg",
				Config: &Config{
					Width: 200,
				},
			},
			"https://p1-47e91401.imageflux.jp/c/sig=1.tiKX5u2kw6wp9zDgl1tLiOIi8IsoRIBw8fVgVc0yrNg=,w=200/images/1.jpg",
		},
	}

	for _, c := range cases {
		got := c.image.SignedURL().String()
		if got != c.output {
			t.Errorf("want %s, got %s", c.output, got)
		}
	}
}

func TestConfig(t *testing.T) {
	cases := []struct {
		config *Config
		output string
	}{
		{
			config: nil,
			output: "",
		},
		{
			config: &Config{},
			output: "",
		},
		{
			config: &Config{
				Width: 200,
			},
			output: "w=200",
		},
		{
			config: &Config{
				Height: 200,
			},
			output: "h=200",
		},
		{
			config: &Config{
				Width:  200,
				Height: 200,
			},
			output: "w=200,h=200",
		},
		{
			config: &Config{
				DisableEnlarge: true,
			},
			output: "u=0",
		},
		{
			config: &Config{
				AspectMode: AspectModeScale,
			},
			output: "a=0",
		},
		{
			config: &Config{
				Clip: image.Rect(100, 150, 200, 250),
			},
			output: "c=100:150:200:250",
		},
		{
			config: &Config{
				ClipRatio: image.Rect(25, 25, 75, 75),
				ClipMax:   image.Pt(100, 100),
			},
			output: "cr=0.25:0.25:0.75:0.75",
		},
		{
			config: &Config{
				Origin: OriginTopLeft,
			},
			output: "g=1",
		},
		{
			config: &Config{
				Background: color.Black,
			},
			output: "b=000000",
		},
		{
			config: &Config{
				Background: color.White,
			},
			output: "b=ffffff",
		},
		{
			config: &Config{
				Background: color.NRGBA{R: 255, G: 0, B: 0, A: 255},
			},
			output: "b=ff0000",
		},
		{
			config: &Config{
				Background: color.NRGBA{R: 0, G: 255, B: 0, A: 255},
			},
			output: "b=00ff00",
		},
		{
			config: &Config{
				Background: color.NRGBA{R: 0, G: 0, B: 255, A: 255},
			},
			output: "b=0000ff",
		},
		{
			config: &Config{
				Format: FormatWebPFromPNG,
			},
			output: "f=webp:png",
		},
		{
			config: &Config{
				Quality: 75,
			},
			output: "q=75",
		},
		{
			config: &Config{
				DisableOptimization: true,
			},
			output: "o=0",
		},
	}

	for _, c := range cases {
		got := c.config.String()
		if got != c.output {
			t.Errorf("want %s, got %s", c.output, got)
		}
	}
}
