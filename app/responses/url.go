package responses

import (
	"crypto/tls"
	"net/url"
)

type UrlConfig struct {
	Tls   *tls.ConnectionState
	Host  string
	Path  string
	Query url.Values
}
