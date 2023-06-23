//go:build go1.18
// +build go1.18

package imageflux

import (
	"reflect"
	"testing"
)

func FuzzParseConfig(f *testing.F) {
	f.Add("w=100")
	f.Fuzz(func(t *testing.T, s string) {
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
		if !reflect.DeepEqual(c0, c1) {
			t.Errorf("c0 != c1: %v != %v", c0, c1)
		}
	})
}
