package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Event struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name"`
	StartDate uint64             `json:"startDate"`
	EndDate   uint64             `json:"endDate"`
	Duration  uint16             `json:"duration"`
	Location  Location           `json:"location"`
}
