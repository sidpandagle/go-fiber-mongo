package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	Name        string             `json:"name,omitempty" validate:"required"`
	Description string             `json:"description,omitempty" validate:"required"`
	Icon        string             `json:"icon,omitempty" validate:"required"`
}
