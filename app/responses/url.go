package responses

import (
	"crypto/tls"
	"fmt"
)

type UrlConfig struct {
	Tls  *tls.ConnectionState
	Host string
	Path string
}

// FullUrl is a function to get the full url from UrlConfig
func (c *UrlConfig) FullUrl() string {
	scheme := "http"
	if c.Tls != nil {
		scheme = "https"
	}
	hostname := c.Host + c.Path

	return fmt.Sprintf("%s://%s", scheme, hostname)
}
