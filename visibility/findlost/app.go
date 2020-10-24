package main

import (
	"context"
	"fmt"
	"github.com/briggysmalls/detectordag/shared"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"github.com/briggysmalls/detectordag/visibility"
	"log"
	"time"
)

const lastSeenDurationHours = 24

var lastSeenDuration time.Duration

type app struct {
	iot              iot.Client
	email            visibility.EmailClient
	shadow           shadow.Client
	lastSeenDuration time.Duration
}

// handleRequest handles a lambda call
func (a *app) runJob(ctx context.Context) error {
	// Print out handler parameters
	log.Print("Context: ", ctx)
	// Request all devices that are considered 'visible'
	devices, err := a.iot.GetThingsByVisibility(true)
	if err != nil {
		return shared.LogErrorAndReturn(err)
	}
	log.Printf("Checking %d devices for updated visibility status", len(devices))
	// Iterate through devices
	for _, device := range devices {
		// Fetch the shadow
		shdw, err := a.shadow.Get(device.DeviceId)
		if err != nil {
			return shared.LogErrorAndReturn(err)
		}
		// Check we have a reported status
		_, ok := shdw.State.Reported["status"].(bool)
		if !ok {
			return fmt.Errorf("%s doesn't have status", visibility.DeviceString(device))
		}
		// Check when the device was last seen
		lastSeen := shdw.Metadata.Reported["status"].Timestamp.Time
		if time.Now().Before(lastSeen.Add(lastSeenDuration)) {
			// This device was seen recently enough
			continue
		}
		if !device.Visibility {
			// We searched for visible devices, something weird has happened
			log.Printf("%s already marked lost despite searching for visible devices", visibility.DeviceString(device))
			continue
		}
		// Mark the device as lost
		err = a.iot.SetVisibiltyState(device.DeviceId, false)
		if err != nil {
			return shared.LogErrorAndReturn(err)
		}
		// Email to say so
		err = a.email.SendVisibilityStatus(device, lastSeen, false)
		if err != nil {
			return shared.LogErrorAndReturn(err)
		}
	}
	// Return the response
	return err
}
