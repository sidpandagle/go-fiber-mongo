package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Swatch struct {
	Id        primitive.ObjectID `json:"id,omitempty"`
	Name      string             `json:"name,omitempty" validate:"required"`
	Tags      string             `json:"tags,omitempty" validate:"required"`
	CreatedAt string             `json:"createdAt,omitempty" validate:"required"`
	Likes     int64              `json:"likes,omitempty" validate:"required"`
}
