package visibility

import (
	"fmt"
	"github.com/briggysmalls/detectordag/shared/iot"
)

func DeviceString(device *iot.Device) string {
	return fmt.Sprintf("Device '%s' ('%s')", device.DeviceId, device.Name)
}
