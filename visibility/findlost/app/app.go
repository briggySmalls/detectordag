package app

import (
	"context"
	"fmt"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"github.com/briggysmalls/detectordag/visibility"
	"log"
	"time"
)

type app struct {
	iot              iot.Client
	email            visibility.EmailClient
	shadow           shadow.Client
	lastSeenDuration time.Duration
}

type App interface {
	RunJob(ctx context.Context) error
}

func New(
	iot iot.Client,
	email visibility.EmailClient,
	shadow shadow.Client,
	lastSeenDuration time.Duration,
) App {
	return &app{
		iot:              iot,
		email:            email,
		shadow:           shadow,
		lastSeenDuration: lastSeenDuration,
	}
}

// handleRequest handles a lambda call
func (a *app) RunJob(ctx context.Context) error {
	// Request all devices that are considered 'visible'
	devices, err := a.iot.GetThingsByVisibility(true)
	if err != nil {
		return err
	}
	log.Printf("Checking %d devices for updated visibility status", len(devices))
	// Iterate through devices
	for _, device := range devices {
		// Process all the devices, logging any errors
		if err := a.processVisibleDevice(device); err != nil {
			log.Print(err)
		}
	}
	return nil
}

func (a *app) processVisibleDevice(device *iot.Device) error {
	// We shouldn't ever process a lost device
	if !device.Visibility {
		log.Printf("%s already marked lost despite searching for visible devices", visibility.DeviceString(device))
		return nil
	}
	// Fetch the shadow
	shdw, err := a.shadow.Get(device.DeviceId)
	if err != nil {
		return err
	}
	// Check we have a reported status
	_, ok := shdw.State.Reported["status"].(bool)
	if !ok {
		return fmt.Errorf("%s doesn't have status", visibility.DeviceString(device))
	}
	// Check when the device was last seen
	lastSeen := shdw.Metadata.Reported["status"].Timestamp.Time
	if time.Now().Before(lastSeen.Add(a.lastSeenDuration)) {
		// This device was seen recently enough
		return nil
	}
	// Mark the device as lost
	log.Print("Device '%s' identifed as lost", device.DeviceId)
	lostStatus := false
	err = a.iot.SetVisibiltyState(device.DeviceId, lostStatus)
	if err != nil {
		return err
	}
	// Email to say so
	err = a.email.SendVisibilityStatus(device, lastSeen, lostStatus)
	if err != nil {
		return err
	}
	return nil
}
