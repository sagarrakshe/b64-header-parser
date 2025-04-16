package plugin

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
)

type Config struct {
	HeaderName string `json:"headerName,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		HeaderName: "X-Custom-Header",
	}
}

type HeaderDecode struct {
	next       http.Handler
	name       string
	headerName string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &HeaderDecode{
		next:       next,
		name:       name,
		headerName: config.HeaderName,
	}, nil

}

func (h *HeaderDecode) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	encoded := req.Header.Get(h.headerName)
	if encoded != "" {
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			fmt.Printf("[plugin-%s] Failed to decode header %s: %v\n", h.name, h.headerName, err)
		} else {
			fmt.Printf("[plugin-%s] Decoded header %s: %s\n", h.name, h.headerName, string(decoded))
		}
	} else {
		fmt.Printf("[plugin-%s] Header %s not found\n", h.name, h.headerName)
	}

	h.next.ServeHTTP(rw, req)
}
