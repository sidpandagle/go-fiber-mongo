package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Habit struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	Activity string             `json:"activity,omitempty" validate:"required"`
	Status   string             `json:"status,omitempty" validate:"required"`
	UserId   string             `json:"userId,omitempty" validate:"required"`
	Date     string             `json:"date,omitempty" validate:"required"`
}
