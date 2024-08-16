package main

import (
	"encoding/json"
	"net/http"
)

func decodePayload[V any](r *http.Request, payLoad V) (V, error) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&payLoad)
	if err != nil {
		var zero V
		return zero, err
	}
	return payLoad, nil
}
