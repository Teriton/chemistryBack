package handler

import (
	"encoding/json"
	"log"
)

func errorMessage(err error) []byte {
	jsonString, err := json.Marshal(map[string]string{"error": err.Error()})
	if err != nil {
		log.Fatal(err)
	}
	return jsonString
}
