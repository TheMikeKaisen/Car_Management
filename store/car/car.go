package car

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/TheMikeKaisen/CarManagement/models"
	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

func new(db *sql.DB) Store {
	return Store{db: db}
}

func (s Store) GetCarById(ctx context.Context, id string) (models.Car, error) {

	var car models.Car

	query := `SELECT 
				c.id, c.name, c.year, c.brand, c.fuel_type, c.price, c.created_at, c.updated_at,
				e.id AS engine_id, e.displacement, e.no_of_cylinders, e.car_range
			FROM 
				car c 
			LEFT JOIN 
				engine e ON c.engine_id = e.id
			WHERE 
				c.id = $1;`

	row := s.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&car.ID, &car.Name, &car.Year, &car.Brand, &car.FuelType, &car.Price, &car.CreatedAt, &car.UpdatedAt,
		&car.Engine.EngineId, &car.Engine.Displacement, &car.Engine.NoOfCylinders, &car.Engine.CarRange,
	)

	if err != nil {
		return models.Car{}, err
	}

	return car, nil
}

func (s Store) GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error) {

	var cars []models.Car
	var query string

	if isEngine {
		query = `SELECT 
				c.id, c.name, c.year, c.brand, c.fuel_type, c.price, c.created_at, c.updated_at,
				e.id AS engine_id, e.displacement, e.no_of_cylinders, e.car_range
			FROM 
				car c 
			LEFT JOIN 
				engine e ON c.engine_id = e.id
			WHERE 
				c.brand = $1;`
	} else {
		query = `SELECT 
				id, name,year, brand, fuel_type, engine_id, price, created_at, updated_at
			FROM 
				car
			WHERE 
				brand = $1;`
	}

	// get the list of rows that matches the brand
	rows, queryErr := s.db.QueryContext(ctx, query, brand)
	if queryErr != nil {
		return nil, queryErr
	}
	defer rows.Close()

	// for each row, create a car and append it to the cars list
	for rows.Next() {
		var car models.Car
		if isEngine {
			err := rows.Scan(
				&car.ID, &car.Name, &car.Year, &car.Brand, &car.FuelType, &car.Price, &car.CreatedAt, &car.UpdatedAt,
				&car.Engine.EngineId, &car.Engine.Displacement, &car.Engine.NoOfCylinders, &car.Engine.CarRange,
			)
			if err != nil {
				return nil, err
			}
		} else {
			err := rows.Scan(
				&car.ID, &car.Name, &car.Year, &car.Brand, &car.FuelType, &car.Engine.EngineId, &car.Price, &car.CreatedAt, &car.UpdatedAt,
			)
			if err != nil {
				return nil, err
			}
		}

		// append the car to the cars list
		cars = append(cars, car)
	}

	// when query multiple rows at a time, there is change that it might lead to some error
	if err := rows.Err(); err != nil {
		fmt.Println("Error while querying rows.")
		return nil, err
	}

	return cars, nil

}

func (s Store) CreateCar(ctx context.Context, carReq models.CarRequest) (models.Car, error) {

	// check whether the engineId exists in the database or not
	var engineId uuid.UUID
	err := s.db.QueryRowContext(ctx, `SELECT id from engine WHERE id=$1`, carReq.Engine.EngineId).Scan(&engineId)

	if err != nil {
		// check if the err is no rows found err
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("No Engine with that Id present")
			return models.Car{}, errors.New("engine_id does not exists in the engine table")
		}
		fmt.Println("Error getting engine id")
	}

	// create a new car id
	carId := uuid.New()

	created_at := time.Now()
	updated_at := created_at

	// create a new car
	newCar := models.Car{
		ID:        carId,
		Name:      carReq.Name,
		Year:      carReq.Year,
		Brand:     carReq.Brand,
		FuelType:  carReq.FuelType,
		Engine:    carReq.Engine,
		Price:     carReq.Price,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
	}

	// to achieve atomicity, we can use the transactions function that postgres provides
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println("Transaction Error")
		return models.Car{}, nil
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `INSERT INTO car 
				(id, name, year, brand, fuel_type, engine, price, create_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7. $8. $9)
			RETURNING id, name, year, brand, fuel_type, engine, price, created_at, updated_at`

	var createdCar models.Car
	scanErr := tx.QueryRowContext(
		ctx, query,

		&newCar.ID,
		&newCar.Name,
		&newCar.Year,
		&newCar.Brand,
		&newCar.FuelType,
		&newCar.Price,
		&newCar.CreatedAt,
		&newCar.UpdatedAt,
	).Scan(
		&createdCar.ID,
		&createdCar.Name,
		&createdCar.Year,
		&createdCar.Brand,
		&createdCar.FuelType,
		&createdCar.Price,
		&createdCar.CreatedAt,
		&createdCar.UpdatedAt,
	)

	if scanErr != nil {
		fmt.Println("Error scanning the car")
		return models.Car{}, scanErr
	}

	return createdCar, nil

}

func (s Store) UpdateCar(ctx context.Context, id string, carReq *models.CarRequest) (models.Car, error) {
	var updateCar models.Car

	// use transaction -> Either everything will complete or none will!
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Car{}, nil
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `
		UPDATE car
		SET name = $2, year=$3, brand=$4, fuel_type=$5, engine_id=$6, price=$7, updated_at=$8
		WHERE id=$1
		RETURNING id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at 
	`
	err = tx.QueryRowContext(ctx, query,
		id,
		&carReq.Name,
		&carReq.Year,
		&carReq.Brand,
		&carReq.FuelType,
		&carReq.Engine.EngineId,
		&carReq.Price,
		time.Now(),
	).Scan(
		&updateCar.ID,
		&updateCar.Name,
		&updateCar.Year,
		&updateCar.Brand,
		&updateCar.FuelType,
		&updateCar.Engine.EngineId,
		&updateCar.Price,
		&updateCar.CreatedAt,
		&updateCar.UpdatedAt,
	)

	if err != nil {
		fmt.Println("Error updating car")
		return models.Car{}, err
	}

	return updateCar, nil

}

func (s Store) DeleteCar(ctx context.Context, id string) (models.Car, error) {
	// start transaction -> either all or none
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println("Error while starting transaction!")
		return models.Car{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	var deletedCar models.Car

	returnQuery := `
		SELECT * from car WHERE id=$1;
	`
	err = tx.QueryRowContext(ctx, returnQuery,
		id,
	).Scan(
		&deletedCar.ID,
		&deletedCar.Name,
		&deletedCar.Year,
		&deletedCar.Brand,
		&deletedCar.FuelType,
		&deletedCar.Engine.EngineId,
		&deletedCar.Price,
		&deletedCar.CreatedAt,
		&deletedCar.UpdatedAt,
	)
	if err != nil {
		fmt.Println("Error while returning car values")
		return models.Car{}, nil
	}

	deleteQuery := `
		DELETE FROM car
		WHERE id=$1
	`

	result, err := tx.ExecContext(ctx, deleteQuery,
		id,
	)

	if err != nil {
		fmt.Println("Error while deleting car")
		return models.Car{}, nil
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Row affect error!")
		return models.Car{}, nil
	}

	if rowsAffected == 0 {
		println("Id do not exist!")
		return models.Car{}, errors.New("id do not exist")
	}

	return deletedCar, nil

}
