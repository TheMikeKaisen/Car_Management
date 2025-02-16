package store

import (
	"context"

	"github.com/TheMikeKaisen/CarManagement/models"
)

type CarStoreInterface interface {
	GetCarById(ctx context.Context, id string) (models.Car, error)
	GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error)
	CreateCar(ctx context.Context, carReq models.CarRequest) (models.Car, error)
	UpdateCar(ctx context.Context, id string, carReq *models.CarRequest) (models.Car, error)
	DeleteCar(ctx context.Context, id string) (models.Car, error)
}

type EngineStoreInterface interface{
	CreateEngine(ctx context.Context, engineReq *models.EngineRequest) (models.Engine, error) 

	GetEngineById(ctx context.Context, engineId string) (models.Engine, error)

	UpdateEngine(ctx context.Context, engineId string, engineReq *models.EngineRequest) (models.Engine, error)

	DeleteEngine(ctx context.Context, engineId string)(models.Engine, error)
}