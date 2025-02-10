package models

import (
	"errors"

	"github.com/google/uuid"
)

type Engine struct {
	EngineId      uuid.UUID `json:"engine_id"`
	Displacement  int64     `json:"displacement"`
	NoOfCylinders int64     `json:"no_of_cylinders"`
	CarRange      int64     `json:"car_range"`
}

type EngineRequest struct {
	Displacement  int64 `json:"displacement"`
	NoOfCylinders int64 `json:"no_of_cylinders"`
	CarRange      int64 `json:"car_range"`
}


func ValidateEngineRequest(engineRequest EngineRequest) error {
	if engineRequest.Displacement <= 0 {
		return errors.New("invalid displacement")
	}
	if engineRequest.NoOfCylinders <= 0 {
		return errors.New("invalid number of cylinders")
	}
	if engineRequest.CarRange <= 0 {
		return errors.New("invalid number of cylinders")
	}
	return nil
}
