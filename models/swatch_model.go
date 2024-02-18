package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Swatch struct {
	Id             primitive.ObjectID `json:"id,omitempty"`
	Name           string             `json:"name,omitempty" validate:"required"`
	Tags           []string           `json:"tags,omitempty"`
	Likes          int64              `json:"likes,omitempty"`
	CreatedAt      time.Time          `json:"createdAt,omitempty"`
	TimeDifference time.Duration      `json:"timeDifference,omitempty"`
}
