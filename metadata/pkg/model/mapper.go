package model

import (
	"movieapp.com/gen"
)

// MetadataToProto converts a Metadata struct into a
// generated proto counterpart.
func MetadataToProto(m *Metadata) *gen.Metadata {
	return &gen.Metadata{
		Id:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Director:    m.Director,
	}
}
