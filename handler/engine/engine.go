package engine

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/TheMikeKaisen/CarManagement/models"
	"github.com/TheMikeKaisen/CarManagement/service"
	"github.com/gorilla/mux"
)

type EngineHandler struct {
	service service.EngineServiceInterface
}

func NewCarHandler(service service.EngineServiceInterface) *EngineHandler {
	return &EngineHandler{service: service}
}

func (e *EngineHandler) CreateEngine(w http.ResponseWriter, r *http.Request) {
	// create context
	ctx := r.Context()

	// read the request body
	engineBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		log.Print("Error while reading the body: ", err)
		return
	}

	// unmarshal the data: byte[] -> json
	var body models.EngineRequest
	err = json.Unmarshal(engineBody, &body)
	if err != nil {
		w.WriteHeader(500)
		log.Print("Error unmarshaling: ", err)
		return
	}

	response, err := e.service.CreateEngine(ctx, &body)
	if err != nil {
		w.WriteHeader(500)
		log.Print("Error creating engine: ", err)
		return
	}

	// marshal the data
	responseBody, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(500)
		log.Print("Error while Marshaling: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	_, err = w.Write(responseBody)
	if err != nil {
		w.WriteHeader(500)
		log.Print("Error while writing the response: ", err)
		return
	}

}

func (e *EngineHandler) GetEngineById(w http.ResponseWriter, r *http.Request) {
	// create context
	ctx := r.Context()

	// extract id
	id := mux.Vars(r)["id"]

	resp, err := e.service.GetEngineById(ctx, id)
	if err != nil {
		w.WriteHeader(500)
		log.Print("Error while getting the engine: ", err)
		return
	}

	// marshal the data
	engineBody, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(500)
		log.Print("Error while marshaling : ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	_, err = w.Write(engineBody)
	if err != nil {
		w.WriteHeader(500)
		log.Print("Error while writing the repsonse: ", err)
		return
	}
}

func (e *EngineHandler) UpdateEngine(w http.ResponseWriter, r *http.Request) {
	// create the context
	ctx := r.Context()

	// extract id
	id := mux.Vars(r)["id"]

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		log.Print("Error while reading the request body: ", err)
		return
	}

	// unmasharshal the data
	var engineReqBody models.EngineRequest
	err = json.Unmarshal(reqBody, &engineReqBody)
	if err != nil {
		w.WriteHeader(500)
		log.Print("Error while marshaling: ", err)
		return
	}

	respBody, err := e.service.UpdateEngine(ctx, id, &engineReqBody)
	if err != nil {
		w.WriteHeader(500)
		log.Print("Error while updating the engine: ", err)
		return
	}

	// marshal the data
	engineBody, err := json.Marshal(respBody)
	if err != nil {
		w.WriteHeader(500)
		log.Print("Error while marshaling: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	_, err = w.Write(engineBody)
	if err != nil {
		w.WriteHeader(500)
		log.Print("Error while writing the response: ", err)
		return
	}

}

func (e *EngineHandler) DeleteEngine(w http.ResponseWriter, r *http.Request){
	// create context
	ctx := r.Context()

	// extract id
	id := mux.Vars(r)["id"]

	deletedEngine, err := e.service.DeleteEngine(ctx, id)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Error deleting the engine: ", err)
		return 
	}

	// marshal
	responseBody, err := json.Marshal(deletedEngine)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Error marshaling body: ", err)
		return 
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(responseBody)
}
