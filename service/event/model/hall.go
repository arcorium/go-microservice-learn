package model

type Hall struct {
	Name     string `json:"name"`
	Location string `json:"location,omitempty" bson:"location,omitempty"`
	Capacity uint32 `json:"capacity"`
}
