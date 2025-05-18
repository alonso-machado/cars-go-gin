package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/username/go-car-service/internal/model"
	"github.com/username/go-car-service/pkg/logger"
)

// CarRepository defines the interface for car data operations
type CarRepository interface {
	Create(ctx context.Context, car *model.Car) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.Car, error)
	GetByName(ctx context.Context, name string) (*model.Car, error)
	GetByBrand(ctx context.Context, brand string) ([]*model.Car, error)
	GetByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]*model.Car, error)
	GetAll(ctx context.Context, page, pageSize int) ([]*model.Car, error)
	Update(ctx context.Context, car *model.Car) error
	Delete(ctx context.Context, id int64) error
}

type carRepository struct {
	db *sql.DB
}

// NewCarRepository creates a new instance of CarRepository
func NewCarRepository(db *sql.DB) CarRepository {
	return &carRepository{db: db}
}

// Create creates a new car in the database
func (r *carRepository) Create(ctx context.Context, car *model.Car) (int64, error) {
	query := `
		INSERT INTO cars (name, brand, manufacturing_value, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	now := time.Now()
	car.CreatedAt = now
	car.UpdatedAt = now

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		car.Name,
		car.Brand,
		car.ManufacturingValue,
		car.Description,
		car.CreatedAt,
		car.UpdatedAt,
	).Scan(&id)

	if err != nil {
		logger.LogSQLError(err, query, car.Name, car.Brand, car.ManufacturingValue, car.Description, now, now)
		return 0, fmt.Errorf("failed to create car: %v", err)
	}

	return id, nil
}

// GetByID retrieves a car by its ID
func (r *carRepository) GetByID(ctx context.Context, id int64) (*model.Car, error) {
	query := `
		SELECT id, name, brand, manufacturing_value, description, created_at, updated_at
		FROM cars
		WHERE id = $1 AND deleted_at IS NULL
	`

	var car model.Car
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&car.ID,
		&car.Name,
		&car.Brand,
		&car.ManufacturingValue,
		&car.Description,
		&car.CreatedAt,
		&car.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("car with ID %d not found", id)
		}
		logger.LogSQLError(err, query, id)
		return nil, fmt.Errorf("failed to get car: %v", err)
	}

	return &car, nil
}

// GetByName retrieves a car by its name
func (r *carRepository) GetByName(ctx context.Context, name string) (*model.Car, error) {
	query := `
		SELECT id, name, brand, manufacturing_value, description, created_at, updated_at
		FROM cars
		WHERE name = $1 AND deleted_at IS NULL
	`

	var car model.Car
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&car.ID,
		&car.Name,
		&car.Brand,
		&car.ManufacturingValue,
		&car.Description,
		&car.CreatedAt,
		&car.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("car with name %s not found", name)
		}
		logger.LogSQLError(err, query, name)
		return nil, fmt.Errorf("failed to get car by name: %v", err)
	}

	return &car, nil
}

// GetByBrand retrieves all cars by brand
func (r *carRepository) GetByBrand(ctx context.Context, brand string) ([]*model.Car, error) {
	query := `
		SELECT id, name, brand, manufacturing_value, description, created_at, updated_at
		FROM cars
		WHERE brand = $1 AND deleted_at IS NULL
	`

	rows, err := r.db.QueryContext(ctx, query, brand)
	if err != nil {
		logger.LogSQLError(err, query, brand)
		return nil, fmt.Errorf("failed to get cars by brand: %v", err)
	}
	defer rows.Close()

	var cars []*model.Car
	for rows.Next() {
		var car model.Car
		if err := rows.Scan(
			&car.ID,
			&car.Name,
			&car.Brand,
			&car.ManufacturingValue,
			&car.Description,
			&car.CreatedAt,
			&car.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan car row: %v", err)
		}
		cars = append(cars, &car)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating car rows: %v", err)
	}

	return cars, nil
}

// GetByPriceRange retrieves all cars within a price range
func (r *carRepository) GetByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]*model.Car, error) {
	query := `
		SELECT id, name, brand, manufacturing_value, description, created_at, updated_at
		FROM cars
		WHERE manufacturing_value BETWEEN $1 AND $2 AND deleted_at IS NULL
	`

	rows, err := r.db.QueryContext(ctx, query, minPrice, maxPrice)
	if err != nil {
		logger.LogSQLError(err, query, minPrice, maxPrice)
		return nil, fmt.Errorf("failed to get cars by price range: %v", err)
	}
	defer rows.Close()

	var cars []*model.Car
	for rows.Next() {
		var car model.Car
		if err := rows.Scan(
			&car.ID,
			&car.Name,
			&car.Brand,
			&car.ManufacturingValue,
			&car.Description,
			&car.CreatedAt,
			&car.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan car row: %v", err)
		}
		cars = append(cars, &car)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating car rows: %v", err)
	}

	return cars, nil
}

// GetAll retrieves all cars with pagination
func (r *carRepository) GetAll(ctx context.Context, page, pageSize int) ([]*model.Car, error) {
	offset := (page - 1) * pageSize

	query := `
		SELECT id, name, brand, manufacturing_value, description, created_at, updated_at
		FROM cars
		WHERE deleted_at IS NULL
		ORDER BY id
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		logger.LogSQLError(err, query, pageSize, offset)
		return nil, fmt.Errorf("failed to get all cars: %v", err)
	}
	defer rows.Close()

	var cars []*model.Car
	for rows.Next() {
		var car model.Car
		if err := rows.Scan(
			&car.ID,
			&car.Name,
			&car.Brand,
			&car.ManufacturingValue,
			&car.Description,
			&car.CreatedAt,
			&car.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan car row: %v", err)
		}
		cars = append(cars, &car)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating car rows: %v", err)
	}

	return cars, nil
}

// Update updates an existing car
func (r *carRepository) Update(ctx context.Context, car *model.Car) error {
	query := `
		UPDATE cars
		SET name = $1, brand = $2, manufacturing_value = $3, description = $4, updated_at = $5
		WHERE id = $6 AND deleted_at IS NULL
	`

	car.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(
		ctx,
		query,
		car.Name,
		car.Brand,
		car.ManufacturingValue,
		car.Description,
		car.UpdatedAt,
		car.ID,
	)

	if err != nil {
		logger.LogSQLError(err, query, car.Name, car.Brand, car.ManufacturingValue, car.Description, car.UpdatedAt, car.ID)
		return fmt.Errorf("failed to update car: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("car with ID %d not found", car.ID)
	}

	return nil
}

// Delete soft deletes a car by ID
func (r *carRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE cars
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		logger.LogSQLError(err, query, id)
		return fmt.Errorf("failed to delete car: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("car with ID %d not found", id)
	}

	return nil
}
