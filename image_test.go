package imageflux

import "testing"

func TestImage(t *testing.T) {
	cases := []struct {
		image  *Image
		output string
	}{
		{
			&Image{
				Host: "p1-47e91401.imageflux.jp",
				Path: "/images/1.jpg",
			},
			"https://p1-47e91401.imageflux.jp/images/1.jpg",
		},
		{
			&Image{
				Host: "p1-47e91401.imageflux.jp",
				Path: "/images/1.jpg",
				Config: &Config{
					Width: 200,
				},
			},
			"https://p1-47e91401.imageflux.jp/c/w=200/images/1.jpg",
		},
		{
			&Image{
				Host: "p1-47e91401.imageflux.jp",
				Path: "/images/1.jpg",
				Config: &Config{
					Secret: "testsigningsecret",
				},
			},
			"https://p1-47e91401.imageflux.jp/c/sig=1.-Yd8m-5pXPihiZdlDATcwkkgjzPIC9gFHmmZ3JMxwS0=/images/1.jpg",
		},
		{
			&Image{
				Host: "p1-47e91401.imageflux.jp",
				Path: "/images/1.jpg",
				Config: &Config{
					Width:  200,
					Secret: "testsigningsecret",
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
	}

	for _, c := range cases {
		got := c.config.String()
		if got != c.output {
			t.Errorf("want %s, got %s", c.output, got)
		}
	}
}
