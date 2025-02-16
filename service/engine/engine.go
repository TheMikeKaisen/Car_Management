package engine

import (
	"context"
	"errors"

	"github.com/TheMikeKaisen/CarManagement/models"
	"github.com/TheMikeKaisen/CarManagement/store"
)

type EngineService struct {
	store store.EngineStoreInterface
}

func NewEngineStore(store store.EngineStoreInterface) *EngineService {
	return &EngineService{store: store}
}

func (e *EngineService) CreateEngine(ctx context.Context, engineReq *models.EngineRequest) (models.Engine, error) {

	// validate the incoming engine
	validateErr := models.ValidateEngineRequest(*engineReq)
	if validateErr != nil {
		return models.Engine{}, validateErr
	}

	// call create engine function
	newEngine, createErr := e.store.CreateEngine(ctx, engineReq)
	if createErr != nil {
		return models.Engine{}, createErr
	}

	return newEngine, nil
}

func (e *EngineService) GetEngineById(ctx context.Context, engineId string) (models.Engine, error) {

	// validate engineId
	if engineId == "" {
		return models.Engine{}, errors.New("id cannot be empty")
	}

	engine, err := e.GetEngineById(ctx, engineId)
	if err != nil {
		return models.Engine{}, err
	}
	return engine, nil
}

func (e *EngineService) UpdateEngine(ctx context.Context, engineId string, engineReq *models.EngineRequest) (models.Engine, error) {

	// validate the incoming engine
	validateErr := models.ValidateEngineRequest(*engineReq)
	if validateErr != nil {
		return models.Engine{}, validateErr
	}

	updatedEngine, updateErr := e.store.UpdateEngine(ctx, engineId, engineReq)

	if updateErr != nil {
		return models.Engine{}, updateErr
	}

	return updatedEngine, nil
}

func (e *EngineService) DeleteEngine(ctx context.Context, engineId string) (*models.Engine, error) {
	// check if id is empty
	if engineId == "" {
		return nil, errors.New("engine id cannot be empty")
	}

	deletedEngine, deleteErr := e.store.DeleteEngine(ctx, engineId)
	if deleteErr != nil {
		return nil, deleteErr
	}

	return &deletedEngine, nil
}
