package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Location struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name"`
	Address   string             `json:"address"`
	Country   string             `json:"country"`
	OpenTime  uint8              `json:"openTime"`
	CloseTime uint8              `json:"closeTime"`
	Halls     []Hall             `json:"halls,omitempty" bson:"halls,omitempty"`
}
