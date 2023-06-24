package imageflux

// Proxy is a proxy of ImageFlux.
type Proxy struct {
	Host string

	// Secret is signing secret.
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

// Parse parses the path and returns the image.
func (p *Proxy) Parse(path string, signature string) (*Image, error) {
	state := parseState{
		s:         path,
		config:    &Config{},
		signature: signature,
	}

	if p.Secret == "" {
		c, rest, err := state.parseConfig()
		if err != nil {
			return nil, err
		}
		return &Image{
			Proxy:  p,
			Path:   rest,
			Config: c,
		}, nil
	}

	c, rest, err := state.parseConfigAndVerifySignature([]byte(p.Secret))
	if err != nil {
		return nil, err
	}
	return &Image{
		Proxy:  p,
		Path:   rest,
		Config: c,
	}, nil
}
