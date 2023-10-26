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

type RatingEvent struct {
	UserID     UserId          `json:"userId"`
	RecordID   RecordId        `json:"recordId"`
	RecordType RecordType      `json:"recordType"`
	Value      RatingValue     `json:"value"`
	EventType  RatingEventType `json:"eventType"`
}

// RatingEventType defines the type of a rating event.
type RatingEventType string

// Rating event types.
const (
	RatingEventTypePut    = "put"
	RatingEventTypeDelete = "delete"
)
