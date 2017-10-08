package main

import "google.golang.org/grpc/metadata"

// Metadata represents the metadata that is pass through grpc, but with additional methods
// to extract the values
type Metadata struct {
	val metadata.MD
}

// Get a value based on the key name
func (m Metadata) Get(key string) string {
	v, ok := m.val[key]
	if ok && len(v) > 0 {
		return v[0]
	}
	return ""
}
