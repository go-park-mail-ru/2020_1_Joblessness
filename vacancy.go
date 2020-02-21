package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type HelloResponse struct {
	Message string `json:"message"`
}

func HelloName(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]
	response := HelloResponse{
		Message: fmt.Sprintf("Hello %s!", name),
	}
	encoder := json.NewEncoder(w)
	err := encoder.Encode(response)

	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		w.Write([]byte("{}"))
		return
	}

	//jsonResponse(w, response, http.StatusOK)
}
