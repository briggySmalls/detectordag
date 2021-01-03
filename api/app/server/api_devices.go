package server

import (
	"encoding/json"
	"github.com/briggysmalls/detectordag/api/app/models"
	"github.com/gorilla/mux"
	"net/http"
)

func (s *server) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	var err error
	// Get the device ID
	id := mux.Vars(r)["deviceId"]
	// Try to parse the body
	var updates models.MutableDevice
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		SetError(w, err, http.StatusBadRequest)
		return
	}
	// Update the name
	shdw, err := s.shadow.UpdateName(id, updates.Name)
	if err != nil {
		SetError(w, err, http.StatusInternalServerError)
		return
	}
	// Build the payload
	payload := models.Device{
		Name:     shdw.Name,
		DeviceId: id,
		State: &models.DeviceState{
			Power:   shdw.Power.Value,
			Updated: shdw.Power.Updated,
		},
		Connection: &models.DeviceConnection{
			Status:  shdw.Connection.Value,
			Updated: shdw.Connection.Updated,
		},
	}
	// Build response content
	body, err := json.Marshal(payload)
	if err != nil {
		SetError(w, err, http.StatusInternalServerError)
	}
	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
