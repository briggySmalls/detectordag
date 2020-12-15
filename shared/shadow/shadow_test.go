package shadow

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUnmarshal(t *testing.T) {
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
	      "connection": "connected",
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
	assert.Equal(t, CONNECTION_STATUS_CONNECTED, shadow.Connection.Value)
	assert.Equal(t, time.Unix(1584803417, 0), shadow.Connection.Updated)
	// Assert the power values
	assert.Equal(t, POWER_STATUS_OFF, shadow.Power.Value)
	assert.Equal(t, time.Unix(1584803414, 0), shadow.Power.Updated)
}
