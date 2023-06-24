//go:build go1.18
// +build go1.18

package imageflux

import (
	"reflect"
	"testing"
)

func FuzzParseConfig(f *testing.F) {
	f.Add("w=100")
	f.Add("w=100,h=200")
	f.Add("a=3")
	f.Add("dpr=5")
	f.Add("ic=100:150:200:250")
	f.Add("icr=0.25:0.25:0.75:0.75")
	f.Add("ic=100:150:200:250,ig=5")
	f.Add("oc=100:150:200:250")
	f.Add("c=100:150:200:250")
	f.Add("ocr=0.25:0.25:0.75:0.75")
	f.Add("cr=0.25:0.25:0.75:0.75")
	f.Add("oc=100:150:200:250,og=5")
	f.Add("g=1")
	f.Add("b=000000")
	f.Add("b=ffffff")
	f.Add("b=FFFFFF")
	f.Add("b=ff0000")
	f.Add("b=00ff00")
	f.Add("b=0000ff")
	f.Add("b=ffffff00")
	f.Add("b=ffffff80")
	f.Add("ir=8")
	f.Add("ir=auto")
	f.Add("or=8")
	f.Add("or=auto")
	f.Add("r=8")
	f.Add("r=auto")
	f.Add("through=jpg")
	f.Add("through=webp:gif:png:jpg")
	f.Add("f=webp:png")
	f.Add("q=75")
	f.Add("o=0")
	f.Add("lossless=1")
	f.Add("s=2")
	f.Add("unsharp=10x1")
	f.Add("unsharp=10x1+1+0.5")
	f.Add("grayscale=0")
	f.Add("grayscale=100")
	f.Add("sepia=0")
	f.Add("sepia=100")
	f.Add("brightness=0")
	f.Add("brightness=200")
	f.Add("contrast=0")
	f.Add("contrast=200")
	f.Add("invert=1")

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
