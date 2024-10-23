package common

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

// provide a json response to handlers
func ResponseJson(w http.ResponseWriter, jsonStruct any, httpStatus int) {
	// respond json
	jsonOutput, err := json.Marshal(jsonStruct)
	if err != nil {
		logrus.WithError(err).Error("parsing json")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write(jsonOutput)
}
