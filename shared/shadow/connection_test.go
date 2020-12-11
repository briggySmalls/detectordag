package shadow

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	payload := `
	{
		"state": {
			"desired":{
				"status":true,
				"connection":true
			},
			"reported":{
				"status":false
			},
			"delta":{
				"status":true
			}
		},
		"metadata":{
			"desired":{
				"status":{
					"timestamp":1584003580
				}
			},
			"reported":{
				"status":{
					"timestamp":1584803417
				},
				"connection":{
					"timestamp":1584803417
				}
			}
		},
		"version":50,
		"timestamp":1584810789
	}`
	// Unpack the payload
	var connState ConnectionStateSchema
	err := connState.Load([]byte(payload))
	assert.Nil(t, err)
	// Flatten
	flat := connState.Flatten()
	assert.Equal(t, flat.State, true)
	assert.Equal(t, flat.Timestamp, 1584803417)
}
