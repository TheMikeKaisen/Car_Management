package car

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/TheMikeKaisen/CarManagement/models"
	"github.com/TheMikeKaisen/CarManagement/service"
	"github.com/gorilla/mux"
)

type CarHandler struct {
	service service.CarServiceInterface
}

func NewCarHandler(service service.CarServiceInterface) *CarHandler {
	return &CarHandler{service: service}
}

func (c *CarHandler) GetCarById(w http.ResponseWriter, r *http.Request) {

	// take out the id params from url
	vars := mux.Vars(r)
	id := vars["id"]

	// create a context
	ctx := r.Context()

	// call GetCarById service
	car, getErr := c.service.GetCarById(ctx, id)
	if getErr != nil {
		w.WriteHeader(500)
		log.Println("Server Error: ", getErr)
		return
	}

	body, err := json.Marshal(car)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Server Error: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	// write in the response body
	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Server Error: ", err)
		return
	}

}

func (c *CarHandler) GetCarByBrand(w http.ResponseWriter, r *http.Request) {
	// context
	ctx := r.Context()

	// get the brand and isEngine string from url
	brand := r.URL.Query().Get("brand")
	isEngine := r.URL.Query().Get("isEngine") == "true"

	resp, err := c.service.GetCarByBrand(ctx, brand, isEngine)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Error", err)
		return
	}

	body, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err = w.Write(body)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Error", err)
		return
	}

}

func (c *CarHandler) CreateCar(w http.ResponseWriter, r *http.Request) {

	var carBody models.CarRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Error", err)
		return
	}

	err = json.Unmarshal(body, &carBody)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Error", err)
		return
	}

	// create context
	ctx := r.Context()
	createdCar, err := c.service.CreateCar(ctx, carBody)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Unable to create Car", err)
		return
	}

	// marshall the data
	car, err := json.Marshal(createdCar)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Error marshalling the data", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err = w.Write(car)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Error writing to response: ", err)
		return
	}
}

func (c *CarHandler) UpdateCar(w http.ResponseWriter, r *http.Request) {

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Error Reading from request body: ", err)
		return
	}

	var carBody models.CarRequest
	err = json.Unmarshal(reqBody, &carBody)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Error while unmarshalling: ", err)
		return
	}

	// create context
	ctx := r.Context()

	// extract id
	id := mux.Vars(r)["id"]
	updatedCar, err := c.service.UpdateCar(ctx, id, &carBody)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Error updating the car: ", err)
		return
	}

	// marshal the data to send as response
	response, err := json.Marshal(updatedCar)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Error marshaling: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	w.Write(response)

}

func (c *CarHandler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	// create a context
	ctx := r.Context()

	// extract id
	id := mux.Vars(r)["id"]

	deletedCar, err := c.service.DeleteCar(ctx, id)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Error while Deleting the car: ", err)
		return
	}

	// marshal the response
	response, err := json.Marshal(deletedCar)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Error while marshaling: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err = w.Write(response)
	if err != nil {
		w.WriteHeader(500)
		log.Println("error while writing to response: ", err)
		return
	}

}
