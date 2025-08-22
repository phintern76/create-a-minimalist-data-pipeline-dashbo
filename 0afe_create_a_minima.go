package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Dataset struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Data        []byte `json:"data"`
}

type Pipeline struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Datasets     []Dataset  `json:"datasets"`
	Transformers []string   `json:"transformers"`
	Model        string     `json:"model"`
	Results      []byte     `json:"results"`
	CreatedAt    int64      `json:"created_at"`
	UpdatedAt    int64      `json:"updated_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var pipelines []Pipeline

func getPipeline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	for _, pipeline := range pipelines {
		if pipeline.ID == id {
			json.NewEncoder(w).Encode(pipeline)
			return
		}
	}
	json.NewEncoder(w).Encode(ErrorResponse{Error: "Pipeline not found"})
}

func getPipelines(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(pipelines)
}

func createPipeline(w http.ResponseWriter, r *http.Request) {
	var pipeline Pipeline
	err := json.NewDecoder(r.Body).Decode(&pipeline)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pipelines = append(pipelines, pipeline)
	json.NewEncoder(w).Encode(pipeline)
}

func updatePipeline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var pipeline Pipeline
	err := json.NewDecoder(r.Body).Decode(&pipeline)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for i, p := range pipelines {
		if p.ID == id {
			pipelines[i] = pipeline
			json.NewEncoder(w).Encode(pipeline)
			return
		}
	}
	json.NewEncoder(w).Encode(ErrorResponse{Error: "Pipeline not found"})
}

func deletePipeline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	for i, p := range pipelines {
		if p.ID == id {
			pipelines = append(pipelines[:i], pipelines[i+1:]...)
			return
		}
	}
	json.NewEncoder(w).Encode(ErrorResponse{Error: "Pipeline not found"})
}

func main() {
	router := mux.NewRouter()

	pipelines = []Pipeline{
		{ID: "1", Name: "Pipeline 1", Description: "Description 1", Datasets: []Dataset{{Name: "Dataset 1", Description: "Description 1"}}},
		{ID: "2", Name: "Pipeline 2", Description: "Description 2", Datasets: []Dataset{{Name: "Dataset 2", Description: "Description 2"}}},
	}

	router.HandleFunc("/pipelines", getPipelines).Methods("GET")
	router.HandleFunc("/pipelines/{id}", getPipeline).Methods("GET")
	router.HandleFunc("/pipelines", createPipeline).Methods("POST")
	router.HandleFunc("/pipelines/{id}", updatePipeline).Methods("PUT")
	router.HandleFunc("/pipelines/{id}", deletePipeline).Methods("DELETE")

	fmt.Println("Listening on port 8000")
	http.ListenAndServe(":8000", router)
}