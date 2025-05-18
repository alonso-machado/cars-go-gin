package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/username/go-car-service/internal/model"
	"github.com/username/go-car-service/internal/repository"
	"github.com/username/go-car-service/pkg/logger"
)

// CarService defines the interface for car business logic
type CarService interface {
	CreateCar(ctx context.Context, req *model.CarRequest) (*model.CarResponse, error)
	GetCarByID(ctx context.Context, id int64) (*model.CarResponse, error)
	GetCarByName(ctx context.Context, name string) (*model.CarResponse, error)
	GetCarsByBrand(ctx context.Context, brand string) ([]*model.CarResponse, error)
	GetCarsByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]*model.CarResponse, error)
	GetAllCars(ctx context.Context, page, pageSize int) ([]*model.CarResponse, error)
	UpdateCar(ctx context.Context, id int64, req *model.CarRequest) (*model.CarResponse, error)
	DeleteCar(ctx context.Context, id int64) error
}

type carService struct {
	repo repository.CarRepository
}

// NewCarService creates a new instance of CarService
func NewCarService(repo repository.CarRepository) CarService {
	return &carService{repo: repo}
}

// CreateCar creates a new car
func (s *carService) CreateCar(ctx context.Context, req *model.CarRequest) (*model.CarResponse, error) {
	// Validate request
	if err := validateCarRequest(req); err != nil {
		return nil, err
	}

	// Convert request to model
	car := req.ToModel()

	// Check if car with the same name already exists
	existingCar, err := s.repo.GetByName(ctx, car.Name)
	if err == nil && existingCar != nil {
		return nil, fmt.Errorf("car with name %s already exists", car.Name)
	}

	// Create car in repository
	id, err := s.repo.Create(ctx, car)
	if err != nil {
		logger.Errorf("Failed to create car: %v", err)
		return nil, fmt.Errorf("failed to create car: %v", err)
	}

	// Get the created car
	createdCar, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Errorf("Failed to fetch created car: %v", err)
		return nil, fmt.Errorf("failed to fetch created car: %v", err)
	}

	return createdCar.ToResponse(), nil
}

// GetCarByID retrieves a car by its ID
func (s *carService) GetCarByID(ctx context.Context, id int64) (*model.CarResponse, error) {
	if id <= 0 {
		return nil, errors.New("invalid car ID")
	}

	car, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Errorf("Failed to get car by ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to get car: %v", err)
	}

	return car.ToResponse(), nil
}

// GetCarByName retrieves a car by its name
func (s *carService) GetCarByName(ctx context.Context, name string) (*model.CarResponse, error) {
	if name == "" {
		return nil, errors.New("car name cannot be empty")
	}

	car, err := s.repo.GetByName(ctx, name)
	if err != nil {
		logger.Errorf("Failed to get car by name %s: %v", name, err)
		return nil, fmt.Errorf("failed to get car: %v", err)
	}

	return car.ToResponse(), nil
}

// GetCarsByBrand retrieves all cars by brand
func (s *carService) GetCarsByBrand(ctx context.Context, brand string) ([]*model.CarResponse, error) {
	if brand == "" {
		return nil, errors.New("brand name cannot be empty")
	}

	cars, err := s.repo.GetByBrand(ctx, brand)
	if err != nil {
		logger.Errorf("Failed to get cars by brand %s: %v", brand, err)
		return nil, fmt.Errorf("failed to get cars by brand: %v", err)
	}

	return toCarResponses(cars), nil
}

// GetCarsByPriceRange retrieves all cars within a price range
func (s *carService) GetCarsByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]*model.CarResponse, error) {
	if minPrice < 0 || maxPrice < 0 || minPrice > maxPrice {
		return nil, errors.New("invalid price range")
	}

	cars, err := s.repo.GetByPriceRange(ctx, minPrice, maxPrice)
	if err != nil {
		logger.Errorf("Failed to get cars by price range %.2f-%.2f: %v", minPrice, maxPrice, err)
		return nil, fmt.Errorf("failed to get cars by price range: %v", err)
	}

	return toCarResponses(cars), nil
}

// GetAllCars retrieves all cars with pagination
func (s *carService) GetAllCars(ctx context.Context, page, pageSize int) ([]*model.CarResponse, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 || pageSize > 100 {
		pageSize = 10 // Default page size
	}

	cars, err := s.repo.GetAll(ctx, page, pageSize)
	if err != nil {
		logger.Errorf("Failed to get all cars (page %d, size %d): %v", page, pageSize, err)
		return nil, fmt.Errorf("failed to get all cars: %v", err)
	}

	return toCarResponses(cars), nil
}

// UpdateCar updates an existing car
func (s *carService) UpdateCar(ctx context.Context, id int64, req *model.CarRequest) (*model.CarResponse, error) {
	if id <= 0 {
		return nil, errors.New("invalid car ID")
	}

	// Validate request
	if err := validateCarRequest(req); err != nil {
		return nil, err
	}

	// Check if car exists
	existingCar, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Errorf("Failed to find car with ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to find car: %v", err)
	}

	// Update car fields
	existingCar.UpdateFromRequest(req)

	// Update car in repository
	if err := s.repo.Update(ctx, existingCar); err != nil {
		logger.Errorf("Failed to update car with ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to update car: %v", err)
	}

	// Get the updated car
	updatedCar, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Errorf("Failed to fetch updated car with ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to fetch updated car: %v", err)
	}

	return updatedCar.ToResponse(), nil
}

// DeleteCar deletes a car by ID
func (s *carService) DeleteCar(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid car ID")
	}

	// Check if car exists
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		logger.Errorf("Failed to find car with ID %d: %v", id, err)
		return fmt.Errorf("failed to find car: %v", err)
	}

	// Delete car from repository
	if err := s.repo.Delete(ctx, id); err != nil {
		logger.Errorf("Failed to delete car with ID %d: %v", id, err)
		return fmt.Errorf("failed to delete car: %v", err)
	}

	return nil
}

// validateCarRequest validates the car request
func validateCarRequest(req *model.CarRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	if req.Name == "" {
		return errors.New("car name is required")
	}

	if req.Brand == "" {
		return errors.New("car brand is required")
	}

	if req.ManufacturingValue <= 0 {
		return errors.New("manufacturing value must be greater than 0")
	}

	if req.ManufacturingValue >= 15000000 {
		return errors.New("manufacturing value must be less than 15,000,000")
	}

	return nil
}

// toCarResponses converts a slice of Car to a slice of CarResponse
func toCarResponses(cars []*model.Car) []*model.CarResponse {
	responses := make([]*model.CarResponse, 0, len(cars))
	for _, car := range cars {
		responses = append(responses, car.ToResponse())
	}
	return responses
}
