package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `json:"id,omitempty"`
	Name      string             `json:"name,omitempty" validate:"required"`
	Location  string             `json:"location,omitempty" validate:"required"`
	Title     string             `json:"title,omitempty" validate:"required"`
	Email     string             `json:"email,omitempty" validate:"required"`
	Password  string             `json:"password,omitempty" validate:"required"`
	CreatedAt string             `json:"createdAt,omitempty" validate:"required"`
}
