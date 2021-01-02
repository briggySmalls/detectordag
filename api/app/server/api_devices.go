package server

import (
	"net/http"
)

func (s *server) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	var err error
	// Ensure the auth middleware provided us with the account ID
	accountID, err := getAccountId(r.Context())
	if err != nil {
		SetError(w, ErrAccountIDMissing, http.StatusInternalServerError)
		return
	}
	// Get the device
	id := mux.Vars(r)["deviceId"]
	device, err := s.iot.GetThing(id)
	if err != nil {
		SetError(w, err, http.StatusInternalServerError)
		return
	}
	// Ensure the user is allowed to modify this device
	if device.AccountID != accountID {
		SetError(w, err, http.StatusForbidden)
		return
	}
	// Try to parse the body
	var updates models.MutableDevice
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		SetError(w, err, http.StatusBadRequest)
		return
	}
	// Update the name
	device, err := s.iot.UpdateThing()
	if err != nil {
		SetError(w, err, http.StatusInternalServerError)
		return
	}
	// Request the shadow
	shdw, err := s.shadow.Get(device.DeviceId)
	if err != nil {
		SetError(w, err, http.StatusInternalServerError)
		return
	}
	// Build the payload
	payload := models.Device{
		Name:     device.Name,
		DeviceId: device.DeviceId,
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
