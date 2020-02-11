package internal

import (
  rpio "github.com/stianeikeland/go-rpio/v4"
  "log"
  "time"
)

const POWER_DETECT_PIN = 4

func Monitor() error {
  // Initialise pins
  log.Print("Initialising GPIO")
  err := rpio.Open()
  if err != nil {
    return err
  }
  defer rpio.Close()
  // Configure pin as detecting input changes
  log.Print("Configuring pin")
  pin := rpio.Pin(POWER_DETECT_PIN)
  pin.Input()
  pin.Detect(rpio.AnyEdge)
  // Listen for changes
  log.Print("Listening...")
  for {
    if pin.EdgeDetected() {
      // Read the new state
      state := pin.Read()
      if state == rpio.Low {
        log.Print("Power gone!")
      } else if state == rpio.High {
        log.Print("Power back!")
      } else {
        log.Print("Unexpected state %i", state)
      }
      time.Sleep(time.Second / 10)
    }
  }
}
