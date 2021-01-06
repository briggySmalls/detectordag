package shadow

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInvalid(t *testing.T) {
	testStrings := []string{
		`{"metadata":{"reported":{"connection":{"timestamp":1584803417},"status":{"timestamp":1584803414}}},"state":{"reported":{"connection":"dummy","status":"off"}},"timestamp":1584810789,"version":50}`,
		`{"metadata":{"reported":{"connection":{"timestamp":1584803417},"status":{"timestamp":1584803414}}},"state":{"reported":{"connection":"connected","status":"dummy"}},"timestamp":1584810789,"version":50}`,
	}
	for _, str := range testStrings {
		// Unpack the payload
		var shadowSchema DeviceShadowSchema
		_, err := shadowSchema.Extract([]byte(str))
		// Expect an error
		assert.NotNil(t, err)
	}
}

func TestSuccess(t *testing.T) {
	payload := `
	{
	  "metadata": {
	    "reported": {
	      "connection": {
	        "timestamp": 1584803417
	      },
	      "status": {
	        "timestamp": 1584803414
	      }
	    }
	  },
	  "state": {
	    "reported": {
	      "connection": {
			  "status": "connected",
			  "transientId": "f5dc1874-5ba1-4727-8366-35d8278ea3e4",
			  "updated": 1584803417,
	      "status": "off"
	    }
	  },
	  "timestamp": 1584810789,
	  "version": 50
	}`
	// Unpack the payload
	var shadowSchema DeviceShadowSchema
	shadow, err := shadowSchema.Extract([]byte(payload))
	assert.Nil(t, err)
	// Assert the general values
	assert.Equal(t, time.Unix(1584810789, 0), shadow.Time)
	assert.Equal(t, 50, shadow.Version)
	// Assert the connection values
	assert.Equal(t, CONNECTION_STATUS_CONNECTED, shadow.Connection.Status)
	assert.Equal(t, time.Unix(1584803417, 0), shadow.Connection.Updated)
	assert.Equal(t, "f5dc1874-5ba1-4727-8366-35d8278ea3e4", shadow.Connection.TransientID)
	// Assert the power values
	assert.Equal(t, POWER_STATUS_OFF, shadow.Power.Value)
	assert.Equal(t, time.Unix(1584803414, 0), shadow.Power.Updated)
}
