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
	    "desired": {
	      "status": {
	        "timestamp": 1584003580
	      }
	    },
	    "reported": {
	      "connection": {
	        "timestamp": 1584803417
	      },
	      "status": {
	        "timestamp": 1584803417
	      }
	    }
	  },
	  "state": {
	    "desired": {
	      "status": true
	    },
	    "reported": {
	      "connection": true,
	      "status": false
	    }
	  },
	  "timestamp": 1584810789,
	  "version": 50
	}`
	// Unpack the payload
	var connState ConnectionStateSchema
	err := connState.Load([]byte(payload))
	assert.Nil(t, err)
	// Flatten
	flat := connState.Flatten()
	assert.Equal(t, true, flat.State)
	assert.Equal(t, time.Unix(1584803417, 0), flat.Updated)
}
