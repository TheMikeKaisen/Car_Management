package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Car struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Year      string    `json:"year"`
	Brand     string    `json:"brand"`
	FuelType  string    `json:"fuel_type"`
	Engine    Engine    `json:"engine"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CarRequest struct {
	Name     string  `json:"name"`
	Year     string  `json:"year"`
	Brand    string  `json:"brand"`
	FuelType string  `json:"fuel_type"`
	Engine   Engine  `json:"engine"`
	Price    float64 `json:"price"`
}

// Call all other validate functions
func ValidateRequest(carReq CarRequest) error {
	if err := ValidateNameBrandPrice(carReq.Name, carReq.Brand, carReq.Price); err != nil {
		return err
	}
	if err := ValidateYear(carReq.Year); err != nil {
		return err
	}
	if err := ValidateEngine(carReq.Engine); err != nil {
		return err
	}
	if err := ValidateFuelType(carReq.FuelType); err != nil {
		return err
	}
	return nil
}


func ValidateNameBrandPrice(name string, brand string, price float64) error {
	// validate name
	if name == "" {
		return errors.New("name is required")
	}

	// validate brand
	if brand == "" {
		return errors.New("brand is required")
	}

	// validate price
	if price <= 0 {
		return errors.New("enter valid price")
	}

	return nil
}

func ValidateYear(year string) error {
	if year == "" {
		return errors.New("year is required")
	}

	yearInt, convErr := strconv.Atoi(year)
	if convErr != nil {
		return errors.New("year must be a valid number")
	}

	currentYear := time.Now()
	if yearInt < 1950 || yearInt > currentYear.Year() {
		return errors.New("enter valid year")
	}
	return nil
}

func ValidateFuelType(fuelType string) error {
	// validate fuelType
	validateFuelTypes := []string{"Petrol", "Electric", "Diesel", "Hybrid"}
	for _, validType := range validateFuelTypes {
		if fuelType == validType {
			return nil
		}
	}
	return errors.New("enter a valid fuel type")
}

func ValidateEngine(engine Engine) error {
	if engine.EngineId == uuid.Nil {
		return errors.New("Engine Id is required")
	}
	if engine.Displacement <= 0 {
		return errors.New("displacement must be greater than zero")
	}
	if engine.NoOfCylinders <= 0 {
		return errors.New("number of cylinders must be greater than zero")
	}
	if engine.CarRange <= 0 {
		return errors.New("car range must be greater than zero")
	}
	return nil
}
