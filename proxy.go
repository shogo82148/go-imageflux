package imageflux

// Proxy is a proxy of ImageFlux.
type Proxy struct {
	// Host is the host of the proxy server.
	Host string

	// SecretBytes is signing secret.
	SecretBytes []byte

	// Secret is signing secret.
	//
	// Deprecated: Use SecretBytes instead.
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

	secret := p.SecretBytes
	if len(secret) == 0 && p.Secret != "" {
		secret = []byte(p.Secret)
	}

	if len(secret) == 0 {
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

	c, rest, err := state.parseConfigAndVerifySignature(secret)
	if err != nil {
		return nil, err
	}
	return &Image{
		Proxy:  p,
		Path:   rest,
		Config: c,
	}, nil
}
