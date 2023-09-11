package model

//rating value indicates rated value for a movie
type RatingValue int

type RecordType string

type RecordId string

const (
	RecordTypeMovie = RecordType("movie")
)

type UserId string

type Rating struct {
	RecordId   string      `json:"recordId"`
	RecordType string      `json:"recordType`
	Value      RatingValue `json:"value"`
	UserId     UserId      `json:"userId"`
}
