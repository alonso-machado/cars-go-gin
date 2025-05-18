package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/username/go-car-service/internal/model"
	"github.com/username/go-car-service/internal/service"
	"github.com/username/go-car-service/pkg/logger"
)

// CarHandler handles HTTP requests related to cars
type CarHandler struct {
	carService service.CarService
}

// NewCarHandler creates a new instance of CarHandler
func NewCarHandler(carService service.CarService) *CarHandler {
	return &CarHandler{carService: carService}
}

// RegisterRoutes registers car routes
func (h *CarHandler) RegisterRoutes(router *gin.RouterGroup) {
	carsGroup := router.Group("/cars")
	{
		carsGroup.GET("", h.GetAllCars)
		carsGroup.GET("/:id", h.GetCarByID)
		carsGroup.GET("/name/:name", h.GetCarByName)
		carsGroup.GET("/brand/:brand", h.GetCarsByBrand)
		carsGroup.GET("/price-range", h.GetCarsByPriceRange)
		carsGroup.POST("", h.CreateCar)
		carsGroup.PUT("/:id", h.UpdateCar)
		carsGroup.DELETE("/:id", h.DeleteCar)
	}
}

// CreateCar handles POST /api/v1/cars
// @Summary Create a new car
// @Description Create a new car with the input payload
// @Tags cars
// @Accept  json
// @Produce  json
// @Param car body model.CarRequest true "Car object that needs to be added"
// @Success 201 {object} model.CarResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /cars [post]
func (h *CarHandler) CreateCar(c *gin.Context) {
	var req model.CarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	car, err := h.carService.CreateCar(c.Request.Context(), &req)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to create car", err)
		return
	}

	c.JSON(http.StatusCreated, car)
}

// GetCarByID handles GET /api/v1/cars/:id
// @Summary Get a car by ID
// @Description Get a car by its ID
// @Tags cars
// @Accept  json
// @Produce  json
// @Param id path int true "Car ID"
// @Success 200 {object} model.CarResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /cars/{id} [get]
func (h *CarHandler) GetCarByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		handleError(c, http.StatusBadRequest, "Invalid car ID", err)
		return
	}

	car, err := h.carService.GetCarByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			handleError(c, http.StatusNotFound, "Car not found", err)
		} else {
			handleError(c, http.StatusInternalServerError, "Failed to get car", err)
		}
		return
	}

	c.JSON(http.StatusOK, car)
}

// GetCarByName handles GET /api/v1/cars/name/:name
// @Summary Get a car by name
// @Description Get a car by its name
// @Tags cars
// @Accept  json
// @Produce  json
// @Param name path string true "Car Name"
// @Success 200 {object} model.CarResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /cars/name/{name} [get]
func (h *CarHandler) GetCarByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		handleError(c, http.StatusBadRequest, "Car name is required", nil)
		return
	}

	car, err := h.carService.GetCarByName(c.Request.Context(), name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			handleError(c, http.StatusNotFound, "Car not found", err)
		} else {
			handleError(c, http.StatusInternalServerError, "Failed to get car", err)
		}
		return
	}

	c.JSON(http.StatusOK, car)
}

// GetCarsByBrand handles GET /api/v1/cars/brand/:brand
// @Summary Get cars by brand
// @Description Get all cars for a specific brand
// @Tags cars
// @Accept  json
// @Produce  json
// @Param brand path string true "Brand Name"
// @Success 200 {array} model.CarResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /cars/brand/{brand} [get]
func (h *CarHandler) GetCarsByBrand(c *gin.Context) {
	brand := c.Param("brand")
	if brand == "" {
		handleError(c, http.StatusBadRequest, "Brand name is required", nil)
		return
	}

	cars, err := h.carService.GetCarsByBrand(c.Request.Context(), brand)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to get cars by brand", err)
		return
	}

	c.JSON(http.StatusOK, cars)
}

// GetCarsByPriceRange handles GET /api/v1/cars/price-range
// @Summary Get cars by price range
// @Description Get all cars within a specified price range
// @Tags cars
// @Accept  json
// @Produce  json
// @Param startPrice query number true "Minimum price"
// @Param finalPrice query number true "Maximum price"
// @Success 200 {array} model.CarResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /cars/price-range [get]
func (h *CarHandler) GetCarsByPriceRange(c *gin.Context) {
	startPrice, err := strconv.ParseFloat(c.Query("startPrice"), 64)
	if err != nil || startPrice < 0 {
		handleError(c, http.StatusBadRequest, "Invalid start price", err)
		return
	}

	finalPrice, err := strconv.ParseFloat(c.Query("finalPrice"), 64)
	if err != nil || finalPrice < 0 || finalPrice < startPrice {
		handleError(c, http.StatusBadRequest, "Invalid final price", err)
		return
	}

	cars, err := h.carService.GetCarsByPriceRange(c.Request.Context(), startPrice, finalPrice)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to get cars by price range", err)
		return
	}

	c.JSON(http.StatusOK, cars)
}

// GetAllCars handles GET /api/v1/cars
// @Summary Get all cars
// @Description Get a list of all cars with pagination
// @Tags cars
// @Accept  json
// @Produce  json
// @Param page query int false "Page number (default 1)"
// @Param pageSize query int false "Number of items per page (default 10, max 100)"
// @Success 200 {array} model.CarResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /cars [get]
func (h *CarHandler) GetAllCars(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	cars, err := h.carService.GetAllCars(c.Request.Context(), page, pageSize)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to get cars", err)
		return
	}

	c.JSON(http.StatusOK, cars)
}

// UpdateCar handles PUT /api/v1/cars/:id
// @Summary Update an existing car
// @Description Update an existing car with the input payload
// @Tags cars
// @Accept  json
// @Produce  json
// @Param id path int true "Car ID"
// @Param car body model.CarRequest true "Car object that needs to be updated"
// @Success 200 {object} model.CarResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /cars/{id} [put]
func (h *CarHandler) UpdateCar(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		handleError(c, http.StatusBadRequest, "Invalid car ID", err)
		return
	}

	var req model.CarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	car, err := h.carService.UpdateCar(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			handleError(c, http.StatusNotFound, "Car not found", err)
		} else {
			handleError(c, http.StatusInternalServerError, "Failed to update car", err)
		}
		return
	}

	c.JSON(http.StatusOK, car)
}

// DeleteCar handles DELETE /api/v1/cars/:id
// @Summary Delete a car
// @Description Delete a car by its ID
// @Tags cars
// @Accept  json
// @Produce  json
// @Param id path int true "Car ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /cars/{id} [delete]
func (h *CarHandler) DeleteCar(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		handleError(c, http.StatusBadRequest, "Invalid car ID", err)
		return
	}

	err = h.carService.DeleteCar(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			handleError(c, http.StatusNotFound, "Car not found", err)
		} else {
			handleError(c, http.StatusInternalServerError, "Failed to delete car", err)
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// ErrorResponse represents an error response
// @Description Error response with message and optional error details
type ErrorResponse struct {
	Success bool   `json:"success" example:false`
	Message string `json:"message" example:"An error occurred"`
	Error   string `json:"error,omitempty" example:"error details"`
}

// handleError is a helper function to handle errors consistently
func handleError(c *gin.Context, statusCode int, message string, err error) {
	logger.Errorf("Error: %v, Details: %v", message, err)

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Message: message,
		Error:   errMsg,
	})
}
