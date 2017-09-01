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
	}

	for _, c := range cases {
		got := c.image.String()
		if got != c.output {
			t.Errorf("want %s, got %s", c.output, got)
		}
	}
}
