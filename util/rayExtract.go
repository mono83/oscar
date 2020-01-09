package util

import "net/http"

// RayExtract method reads tracing identifiers from http headers response.
func RayExtract(h http.Header) (string, bool) {
	// Zipkin tracing ID
	if v := h.Get("X-B3-TraceId"); len(v) > 0 {
		return v, true
	}

	// Jaeger tracing ID
	if v := h.Get("uber-trace-id"); len(v) > 0 {
		return v, true
	}

	// Cloudflare
	if v := h.Get("Cf-Ray"); len(v) > 0 {
		return v, true
	}

	// Custom internal communication ray
	if v := h.Get("X-ITC-RayId"); len(v) > 0 {
		return v, true
	}

	return "", false
}

// RayExtractOrEmpty method reads tracing identifiers from http headers response.
func RayExtractOrEmpty(h http.Header) string {
	v, _ := RayExtract(h)
	return v
}
