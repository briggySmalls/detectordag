package server

import (
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/briggysmalls/detectordag/shared/email"
	"github.com/briggysmalls/detectordag/shared/shadow"
)

type Params struct {
	Config Config
	Db     database.Client
	Email  email.Client
	Shadow shadow.Client
}
