package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func infoMessage(w http.ResponseWriter, message string, statusCode int) {
	jsonString, err := json.Marshal(map[string]string{"info": message})
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(statusCode)
	w.Write(jsonString)
}

func checkError(w http.ResponseWriter, err error, statusCode int) bool {
	if err != nil {
		errorMessageJSON := errorMessage(err)
		w.WriteHeader(statusCode)
		w.Write(errorMessageJSON)
		return true
	}
	return false
}

func errorMessage(err error) []byte {
	jsonString, err := json.Marshal(map[string]string{"error": err.Error()})
	if err != nil {
		log.Fatal(err)
	}
	return jsonString
}
