package server

import (
	"github.com/briggysmalls/detectordag/api/swagger/tokens"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/briggysmalls/detectordag/shared/email"
	"github.com/briggysmalls/detectordag/shared/iot"
	"github.com/briggysmalls/detectordag/shared/shadow"
)

type Params struct {
	Db     database.Client
	Email  email.Client
	Shadow shadow.Client
	IoT    iot.Client
	Tokens tokens.Tokens
}
