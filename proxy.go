package imageflux

// Proxy is a proxy of ImageFlux.
type Proxy struct {
	Host   string
	Secret string
}

// Image returns an image served via the proxy.
func (p *Proxy) Image(path string, config *Config) *Image {
	return &Image{
		Proxy:  p,
		Path:   path,
		Config: config,
	}
}
