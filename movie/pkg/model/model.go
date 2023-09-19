package model

import (
	"movieapp.com/metadata/pkg/model"
)

type MovieDetails struct {
	Rating   float64        `json:"rating"`
	MetaData model.Metadata `json:"meta"`
}
