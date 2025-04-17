package b64_header_parser

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

type Config struct {
	HeaderName      string   `json:"headerName,omitempty"`
	HeaderSeparator string   `json:"headerSeparator,omitempty"`
	AllowedValues   []string `json:"allowedValues,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		HeaderName:      "X-Custom-Header",
		HeaderSeparator: ":",
		AllowedValues:   []string{},
	}
}

type HeaderDecode struct {
	next            http.Handler
	name            string
	headerName      string
	headerSeparator string
	allowedValues   []string
}

func New(_ context.Context, config *Config, _ string) (http.Handler, error) {
	if config.HeaderName == "" {
		return nil, fmt.Errorf("headerName is required")
	}
	if config.HeaderSeparator == "" {
		config.HeaderSeparator = ":"
	}

	return &HeaderDecode{
		headerName:      config.HeaderName,
		headerSeparator: config.HeaderSeparator,
		allowedValues:   config.AllowedValues,
	}, nil
}

func (h *HeaderDecode) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	encodedHeader := req.Header.Get(h.headerName)

	if encodedHeader == "" {
		fmt.Println("Encoded header not found")
		h.next.ServeHTTP(rw, req)
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(encodedHeader)
	if err != nil {
		fmt.Printf("Failed to decode header: %v", err)
		http.Error(rw, "Invalid base64 header", http.StatusBadRequest)
		return
	}
	parts := strings.SplitN(string(decoded), h.headerSeparator, 2)
	if len(parts) != 2 {
		fmt.Println("Decoded header format incorrect")
		http.Error(rw, "Invalid header format", http.StatusBadRequest)
		return
	}

	value := strings.TrimSpace(parts[1])

	if len(h.allowedValues) > 0 && !slices.Contains(h.allowedValues, value) {
		fmt.Printf("Value %s not in allowed list", value)
		http.Error(rw, "Forbidden", http.StatusForbidden)
		return
	}
	fmt.Printf("Value %s present in allowed list", value)

	h.next.ServeHTTP(rw, req)
}
