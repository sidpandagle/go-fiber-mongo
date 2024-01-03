package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Report struct {
	Id           primitive.ObjectID `json:"id,omitempty"`
	Category     string             `json:"category,omitempty" validate:"required"`
	Name         string             `json:"name,omitempty" validate:"required"`
	Url          string             `json:"url,omitempty" validate:"required"`
	Description  string             `json:"description,omitempty" validate:"required"`
	Summary      string             `json:"summary,omitempty" validate:"required"`
	Segmentation string             `json:"segmentation,omitempty" validate:"required"`
	TOC          string             `json:"toc,omitempty" validate:"required"`
	Methodology  string             `json:"methodology,omitempty" validate:"required"`
	Price        string             `json:"price,omitempty" validate:"required"`
	CreatedAt    string             `json:"createdAt,omitempty" validate:"required"`
	UpdatedAt    string             `json:"updatedAt,omitempty" validate:"required"`
	Status       string             `json:"status,omitempty" validate:"required"`
}
