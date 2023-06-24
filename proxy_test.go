package imageflux

import (
	"reflect"
	"testing"
	"time"
)

func TestProxy_Parse(t *testing.T) {
	fixTime(t, time.Date(2023, 6, 24, 9, 23, 0, 0, time.UTC))

	cases := []struct {
		input     string
		signature string
		proxy     *Proxy
		want      *Config
		path      string
	}{
		{
			input: "/images/1.jpg",
			proxy: &Proxy{},
			want:  &Config{},
			path:  "/images/1.jpg",
		},
		{
			input: "/c/sig=1.tiKX5u2kw6wp9zDgl1tLiOIi8IsoRIBw8fVgVc0yrNg=,w=200/images/1.jpg",
			proxy: &Proxy{
				Secret: "testsigningsecret",
			},
			want: &Config{
				Width: 200,
			},
			path: "/images/1.jpg",
		},
		{
			input: "/c/w=200,sig=1.tiKX5u2kw6wp9zDgl1tLiOIi8IsoRIBw8fVgVc0yrNg=/images/1.jpg",
			proxy: &Proxy{
				Secret: "testsigningsecret",
			},
			want: &Config{
				Width: 200,
			},
			path: "/images/1.jpg",
		},
		{
			input:     "/c/w=200/images/1.jpg",
			signature: "1.tiKX5u2kw6wp9zDgl1tLiOIi8IsoRIBw8fVgVc0yrNg=",
			proxy: &Proxy{
				Secret: "testsigningsecret",
			},
			want: &Config{
				Width: 200,
			},
			path: "/images/1.jpg",
		},
		{
			input:     "/c/sig=1.invalid,w=200/images/1.jpg",
			signature: "1.tiKX5u2kw6wp9zDgl1tLiOIi8IsoRIBw8fVgVc0yrNg=",
			proxy: &Proxy{
				Secret: "testsigningsecret",
			},
			want: &Config{
				Width: 200,
			},
			path: "/images/1.jpg",
		},
		{
			input: "/c/sig=1.-Yd8m-5pXPihiZdlDATcwkkgjzPIC9gFHmmZ3JMxwS0=/images/1.jpg",
			proxy: &Proxy{
				Secret: "testsigningsecret",
			},
			want: &Config{},
			path: "/images/1.jpg",
		},
		{
			input:     "/images/1.jpg",
			signature: "1.-Yd8m-5pXPihiZdlDATcwkkgjzPIC9gFHmmZ3JMxwS0=",
			proxy: &Proxy{
				Secret: "testsigningsecret",
			},
			want: &Config{},
			path: "/images/1.jpg",
		},
	}

	for _, c := range cases {
		got, err := c.proxy.Parse(c.input, c.signature)
		if err != nil {
			t.Errorf("%q: unexpected error: %v", c.input, err)
			continue
		}
		if !reflect.DeepEqual(got.Config, c.want) {
			t.Errorf("%q: unexpected config: want %#v, got %#v", c.input, c.want, got)
		}
		if got.Path != c.path {
			t.Errorf("%q: want %s, got %s", c.input, c.path, got.Path)
		}
	}
}

// test signature validation errors
func TestProxy_Parse_sig_error(t *testing.T) {
	fixTime(t, time.Date(2023, 6, 24, 9, 23, 0, 0, time.UTC))

	cases := []struct {
		input     string
		signature string
		proxy     *Proxy
	}{
		{
			// signature mismatch
			input: "/c/sig=1.tiKX5u2kw6wp9zDgl1tLiOIi8IsoRIBw8fVgVc0yrNg=/images/1.jpg",
			proxy: &Proxy{
				Secret: "testsigningsecret",
			},
		},
		{
			// invalid signature version
			input: "/c/sig=Z.tiKX5u2kw6wp9zDgl1tLiOIi8IsoRIBw8fVgVc0yrNg=,w=200/images/1.jpg",
			proxy: &Proxy{
				Secret: "testsigningsecret",
			},
		},
		{
			// base64 decode error
			input: "/c/sig=1.A/images/1.jpg",
			proxy: &Proxy{
				Secret: "testsigningsecret",
			},
		},
		{
			input:     "/c/sig=1.tiKX5u2kw6wp9zDgl1tLiOIi8IsoRIBw8fVgVc0yrNg=,w=200/images/1.jpg",
			signature: "1.-Yd8m-5pXPihiZdlDATcwkkgjzPIC9gFHmmZ3JMxwS0=",
			proxy: &Proxy{
				Secret: "testsigningsecret",
			},
		},
	}

	for _, c := range cases {
		_, err := c.proxy.Parse(c.input, c.signature)
		if err != ErrInvalidSignature {
			t.Errorf("%q: want ErrInvalidSignature, got %v", c.input, err)
		}
	}
}
