package http_helper

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Wrapper to parse JSON from HTTP requests
func GetJSONFromRequest(r *http.Request, obj interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithField("err", err).Error("Failed to read from body.")
		return err
	}
	if err := json.Unmarshal(body, obj); err != nil {
		log.WithField("err", err).Error("Parse JSON failed.")
		return err
	}
	return nil
}

// Wrapper to set JSON to HTTP response
func SetJSONToResponse(w http.ResponseWriter, obj interface{}) error {
	jsonStruct, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonStruct)
	return nil
}

func ParseJSONBody(obj interface{}) ([]byte, error) {
	jsonStruct, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return jsonStruct, nil
}
