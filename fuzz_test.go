//go:build go1.18
// +build go1.18

package imageflux

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func FuzzParseConfig(f *testing.F) {
	for _, c := range parseConfigCases {
		f.Add(c.input)
	}

	f.Fuzz(func(t *testing.T, s string) {
		fixTime(t, time.Date(2023, 6, 24, 9, 23, 0, 0, time.UTC))

		c0, rest, err := ParseConfig(s)
		if err != nil {
			return
		}
		_ = rest
		s1 := c0.String()
		c1, _, err := ParseConfig(s1)
		if err != nil {
			t.Error(err)
			return
		}

		// The zero value of Format has same meaning as FormatAuto.
		if c0.Format == "" {
			c0.Format = FormatAuto
		}
		if c1.Format == "" {
			c1.Format = FormatAuto
		}

		if !reflect.DeepEqual(c0, c1) {
			t.Errorf("%q: c0 != c1: %v != %v", s, c0, c1)
		}
	})
}

func FuzzProxy_Parse(f *testing.F) {
	f.Add("/c!/w=100/images/1.jpg", "", "")
	f.Add("/c/sig=1.tiKX5u2kw6wp9zDgl1tLiOIi8IsoRIBw8fVgVc0yrNg=,w=200/images/1.jpg", "", "testsigningsecret")
	f.Add("/c/w=200,sig=1.tiKX5u2kw6wp9zDgl1tLiOIi8IsoRIBw8fVgVc0yrNg=/images/1.jpg", "", "testsigningsecret")
	f.Add("/c/w=200%2csig=1.tiKX5u2kw6wp9zDgl1tLiOIi8IsoRIBw8fVgVc0yrNg=/images/1.jpg", "", "testsigningsecret")
	f.Add("/c/w=200%2Csig=1.tiKX5u2kw6wp9zDgl1tLiOIi8IsoRIBw8fVgVc0yrNg=/images/1.jpg", "", "testsigningsecret")
	f.Add("/c/w=200/images/1.jpg", "1.tiKX5u2kw6wp9zDgl1tLiOIi8IsoRIBw8fVgVc0yrNg=", "testsigningsecret")
	f.Add("/c/sig=1.-Yd8m-5pXPihiZdlDATcwkkgjzPIC9gFHmmZ3JMxwS0=/images/1.jpg", "", "testsigningsecret")
	f.Add("/images/1.jpg", "1.-Yd8m-5pXPihiZdlDATcwkkgjzPIC9gFHmmZ3JMxwS0=", "testsigningsecret")

	f.Fuzz(func(t *testing.T, path, sig, secret string) {
		fixTime(t, time.Date(2023, 6, 24, 9, 23, 0, 0, time.UTC))

		p := &Proxy{
			Host:   "example.com",
			Secret: secret,
		}
		img0, err := p.Parse(path, sig)
		if err != nil {
			return
		}

		u := img0.SignedURL()
		_, err = p.Parse(strings.TrimPrefix(u, "https://example.com"), "")
		if err != nil {
			t.Error(err)
		}
	})
}
