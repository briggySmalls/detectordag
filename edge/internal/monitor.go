package internal

import (
  "github.com/stianeikeland/go-rpio/v4"
  "log"
)

const POWER_DETECT_PIN = 4

func monitor() {
  // Initialise pins
  err := rpio.Open()
  defer rpio.Close()
  // Configure pin as detecting input changes
  pin := rpio.Pin(POWER_DETECT_PIN)
  pin.Input()
  pin.Detect(rpio.AnyEdge)
  // Listen for changes
  for {
    if pin.EdgeDetected() {
      // Read the new state
      state := pin.Read()
      if state == gpio.Low {
        log.Print("Power back!")
      } else if state == gpio.High {
        log.Print("Power gone!")
      } else {
        log.Errorf("Unexpected state %i", state)
      }
    }
  }
}
