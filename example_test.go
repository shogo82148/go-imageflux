package imageflux_test

import (
	"fmt"
	"log"

	"github.com/shogo82148/go-imageflux"
)

func ExampleImage_SignedURL() {
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
}

func ExampleImage_SignedURL_signed() {
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
}

func ExampleImage_SignedURLWithoutComma() {
	proxy := &imageflux.Proxy{
		Host: "demo.imageflux.jp",
	}
	cfg := &imageflux.Config{
		// resize the image to 200px width.
		Width: 200,

		// convert the image to WebP format.
		Format: imageflux.FormatWebPAuto,
	}
	u := proxy.Image("/images/1.jpg", cfg).SignedURLWithoutComma()
	fmt.Println(u)

	// Output:
	// https://demo.imageflux.jp/c/w=200%2Cf=webp:auto/images/1.jpg
}

func ExampleProxy_Parse() {
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
}

func ExampleProxy_Parse_signed() {
	proxy := &imageflux.Proxy{
		Secret: "testsigningsecret",
	}
	image, err := proxy.Parse("/c/w=200/images/1.jpg", "1.tiKX5u2kw6wp9zDgl1tLiOIi8IsoRIBw8fVgVc0yrNg=")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("path = %s\n", image.Path)
	fmt.Printf("width = %d\n", image.Config.Width)

	// Output:
	// path = /images/1.jpg
	// width = 200
}
