package model

import (
	"database/sql"
	"time"
)

// Car represents a car in the system
type Car struct {
	ID                int64          `json:"id" db:"id"`
	Name              string         `json:"name" db:"name"`
	Brand             string         `json:"brand" db:"brand"`
	ManufacturingValue float64        `json:"manufacturing_value" db:"manufacturing_value"`
	Description       sql.NullString `json:"description,omitempty" db:"description"`
	CreatedAt         time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at" db:"updated_at"`
}

// CarRequest represents the request payload for creating/updating a car
type CarRequest struct {
	Name              string  `json:"name" binding:"required"`
	Brand             string  `json:"brand" binding:"required"`
	ManufacturingValue float64 `json:"manufacturing_value" binding:"required,gt=0,lt=15000000"`
	Description       *string `json:"description,omitempty"`
}

// CarResponse represents the response payload for a car
type CarResponse struct {
	ID                int64   `json:"id"`
	Name              string  `json:"name"`
	Brand             string  `json:"brand"`
	ManufacturingValue float64 `json:"manufacturing_value"`
	Description       *string `json:"description,omitempty"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

// ToResponse converts a Car model to a CarResponse
toResponse(car *Car) *CarResponse {
	var desc *string
	if car.Description.Valid {
		desc = &car.Description.String
	}

	return &CarResponse{
		ID:                car.ID,
		Name:              car.Name,
		Brand:             car.Brand,
		ManufacturingValue: car.ManufacturingValue,
		Description:       desc,
		CreatedAt:         car.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         car.UpdatedAt.Format(time.RFC3339),
	}
}

// ToModel converts a CarRequest to a Car model
func (cr *CarRequest) ToModel() *Car {
	var desc sql.NullString
	if cr.Description != nil {
		desc = sql.NullString{String: *cr.Description, Valid: true}
	}

	return &Car{
		Name:              cr.Name,
		Brand:             cr.Brand,
		ManufacturingValue: cr.ManufacturingValue,
		Description:       desc,
	}
}

// UpdateFromRequest updates a Car model from a CarRequest
func (c *Car) UpdateFromRequest(req *CarRequest) {
	c.Name = req.Name
	c.Brand = req.Brand
	c.ManufacturingValue = req.ManufacturingValue
	if req.Description != nil {
		c.Description = sql.NullString{String: *req.Description, Valid: true}
	} else {
		c.Description = sql.NullString{Valid: false}
	}
}
