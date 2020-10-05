package server

import (
	"encoding/json"
	"errors"
	"github.com/briggysmalls/detectordag/api/app/models"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/gorilla/mux"
	"net/http"
)

func (s *server) GetAccount(w http.ResponseWriter, r *http.Request) {
	// Ensure the auth middleware provided us with the account ID
	accountID, err := getAccountId(r.Context())
	if err != nil {
		SetError(w, ErrAccountIDMissing, http.StatusInternalServerError)
		return
	}
	// Request the account
	account, err := s.db.GetAccountById(string(accountID))
	if err != nil {
		SetError(w, err, http.StatusInternalServerError)
		return
	}
	// Build the response
	payload, err := s.createAccountPayload(account)
	if err != nil {
		SetError(w, err, http.StatusInternalServerError)
		return
	}
	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

func (s *server) GetDevices(w http.ResponseWriter, r *http.Request) {
	// Ensure the auth middleware provided us with the account ID
	accountID, err := getAccountId(r.Context())
	if err != nil {
		SetError(w, ErrAccountIDMissing, http.StatusInternalServerError)
		return
	}
	// Fetch the devices associated with the account
	devices, err := s.iot.GetThingsByAccount(accountID)
	if err != nil {
		SetError(w, err, http.StatusInternalServerError)
		return
	}
	// Request each device's shadow
	payload := make([]models.Device, len(devices))
	for i, device := range devices {
		// Request the shadow
		shdw, err := s.shadow.Get(device.DeviceId)
		if err != nil {
			SetError(w, err, http.StatusInternalServerError)
			return
		}
		// Coerce the data into the right form
		status, ok := shdw.State.Reported["status"].(bool)
		if !ok {
			SetError(w, err, http.StatusInternalServerError)
		}
		// Build the payload
		payload[i] = models.Device{
			Name:     device.Name,
			DeviceId: device.DeviceId,
			Updated:  shdw.Metadata.Reported["status"].Timestamp.Time,
			State:    &models.DeviceState{Power: status},
		}
	}
	// Prepare the JSON response
	body, err := json.Marshal(payload)
	if err != nil {
		SetError(w, err, http.StatusInternalServerError)
	}
	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (s *server) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	// Ensure the auth middleware provided us with the account ID
	accountID, err := getAccountId(r.Context())
	if err != nil {
		SetError(w, ErrAccountIDMissing, http.StatusInternalServerError)
		return
	}
	// Parse the emails from the request
	var emails models.Emails
	err = json.NewDecoder(r.Body).Decode(&emails)
	if err != nil {
		SetError(w, err, http.StatusBadRequest)
		return
	}
	// Request that emails are verified
	err = s.email.VerifyEmailsIfNecessary(emails.Emails)
	if err != nil {
		SetError(w, err, http.StatusInternalServerError)
		return
	}
	// Update the database
	account, err := s.db.UpdateAccountEmails(accountID, emails.Emails)
	if err != nil {
		SetError(w, err, http.StatusInternalServerError)
		return
	}
	// Build the response
	payload, err := s.createAccountPayload(account)
	if err != nil {
		SetError(w, err, http.StatusInternalServerError)
		return
	}
	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

func (s *server) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	// Ensure the auth middleware provided us with the account ID
	accountID, err := getAccountId(r.Context())
	if err != nil {
		SetError(w, ErrAccountIDMissing, http.StatusInternalServerError)
		return
	}
	// Parse the device name from the request
	var params models.MutableDevice
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		SetError(w, err, http.StatusBadRequest)
		return
	}
	// Pull out the device ID
	vars := mux.Vars(r)
	deviceID, ok := vars["deviceId"]
	if !ok {
		SetError(w, errors.New("Device ID not found in URL"), http.StatusBadRequest)
		return
	}
	// Make a request to register a device
	device, certs, err := s.iot.RegisterThing(accountID, deviceID, params.Name)
	if err != nil {
		SetError(w, err, http.StatusInternalServerError)
		return
	}
	// Construct the response
	payload := models.DeviceRegistered{
		Name:     device.Name,
		DeviceId: device.DeviceId,
		Certificate: &models.DeviceRegisteredCertificate{
			Certificate: certs.Certificate,
			PublicKey:   certs.Public,
			PrivateKey:  certs.Private,
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		SetError(w, err, http.StatusInternalServerError)
		return
	}
	// Send response
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// Create account payload from database response
func (s *server) createAccountPayload(account *database.Account) ([]byte, error) {
	// Build the response
	payload := models.Account{
		Username: account.Username,
		Emails:   account.Emails,
	}
	// Ensure empty slices appear as '[]' in JSON
	if payload.Emails == nil {
		payload.Emails = make([]string, 0)
	}
	// Prepare the JSON response
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return body, nil
}
