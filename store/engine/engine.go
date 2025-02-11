package engine

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/TheMikeKaisen/CarManagement/models"
	"github.com/google/uuid"
)

type Engine struct {
	db *sql.DB
}

func new(db *sql.DB) Engine {
	return Engine{db: db}
}

func (e Engine) CreateEngine(ctx context.Context, engineReq models.EngineRequest) (models.Engine, error) {

	// start transaction -> either all or none!
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println("Error while starting transaction")
		return models.Engine{}, nil
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// to store engine
	var createdEngine models.Engine

	// create an engine id
	engineId := uuid.New()

	query := `
		INSERT INTO engine(id, displacement, no_of_cylinders, car_range)
		VALUES($1, $2, $3, $4)
		RETURNING id, displacement, no_of_cylinders, car_range
	`

	err = tx.QueryRowContext(ctx, query,
		engineId,
		engineReq.Displacement,
		engineReq.NoOfCylinders,
		engineReq.CarRange,
	).Scan(
		&createdEngine.EngineId,
		&createdEngine.Displacement,
		&createdEngine.NoOfCylinders,
		&createdEngine.CarRange,
	)

	if err != nil {
		fmt.Println("Error while creating an engine")
		return models.Engine{}, nil
	}

	return createdEngine, nil

}

func (e Engine) GetEngineById(ctx context.Context, engineId string) (models.Engine, error) {

	
	// parse string id into uuid.UUID
	id, err := uuid.Parse(engineId)
	if err != nil {
		fmt.Println("error while parsing id")
		return models.Engine{}, err
	}

	// to store engine
	var getEngine models.Engine

	// query
	getEngineQuery := `
		SELECT id, displacement, no_of_cylinders, car_range
		from engine
		WHERE id=$1
	`

	err = e.db.QueryRowContext(ctx, getEngineQuery, id).Scan(
		&getEngine.EngineId,
		&getEngine.Displacement,
		&getEngine.NoOfCylinders,
		&getEngine.CarRange,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No engine with the given id")
			return models.Engine{}, errors.New("no engine with the given id")
		}
		fmt.Println("Error while getting engine")
		return models.Engine{}, err
	}
	return getEngine, nil
}

func (e Engine) UpdateEngine(ctx context.Context, engineId string, engineReq models.EngineRequest) (models.Engine, error) {

	// parse string id into uuid.UUID
	id, err := uuid.Parse(engineId)
	if err != nil {
		fmt.Println("error while parsing id")
		return models.Engine{}, err
	}

	// start transaction
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println("Error while starting transaction")
		return models.Engine{}, nil
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// store updated engine
	var updatedEngine models.Engine

	updateEngineQuery := `
		UPDATE engine
		SET displacement=$2, no_of_cylinders=$3, car_range=$4
		WHERE id=$1
		RETURNING id, displacement, no_of_cylinders, car_range
	`
	err = tx.QueryRowContext(ctx, updateEngineQuery, 
		id, 
		engineReq.Displacement, 
		engineReq.NoOfCylinders, 
		engineReq.CarRange,
	).Scan(
		&updatedEngine.EngineId,
		&updatedEngine.Displacement,
		&updatedEngine.NoOfCylinders,
		&updatedEngine.CarRange,
	)

	if err!= nil {
		fmt.Println("Error while updating engine")
		return models.Engine{}, err
	}

	return updatedEngine, nil


}

func (e Engine) DeleteEngine(ctx context.Context, engineId string)(models.Engine, error) {

	// parse string id into uuid.UUID
	id, err := uuid.Parse(engineId)
	if err != nil {
		fmt.Println("error while parsing id")
		return models.Engine{}, err
	}

	// start the transaction
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println("Error starting a transactions")
		return models.Engine{}, nil
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// store deleted engine
	var deletedEngine models.Engine

	// get the engine
	getEngineQuery := `
		SELECT id, displacement, no_of_cylinders, car_range
		FROM engine 
		WHERE id=$1
	`
	err = tx.QueryRowContext(ctx, getEngineQuery, 
		id,
	).Scan(
		&deletedEngine.EngineId, 
		&deletedEngine.Displacement, 
		&deletedEngine.NoOfCylinders, 
		&deletedEngine.CarRange,
	)

	if err != nil {
		fmt.Println("Error while storing.")
		return models.Engine{}, err
	}

	// query
	deleteEngineQuery := `
		DELETE FROM engine
		WHERE id=$1
	`

	result, err :=tx.ExecContext(ctx, deleteEngineQuery, id)
	if err != nil {
		fmt.Println("Error while deleting engine")
		return models.Engine{}, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected==0 {
		fmt.Println("engine with the given id does not exist")
		return models.Engine{}, errors.New("engine with the given id does not exist")
	}

	return deletedEngine, nil


	
}
