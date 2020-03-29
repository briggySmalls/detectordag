/*
 * Detectordag
 *
 * API for detectordag JAMStack dashboard
 *
 * API version: 1.0.0
 * Contact: briggySmalls90@gmail.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/briggysmalls/detectordag/shared/shadow"
	"github.com/gorilla/mux"
)

const (
	uuidRegex = `[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter(config *Config, db database.Client, shadow shadow.Client) *mux.Router {
	// Create the router
	router := mux.NewRouter().StrictSlash(true)
	// Create a handlerer
	s := server{
		config: *config,
		db:     db,
		shadow: shadow,
	}
	// Prepare the routes
	var routes = Routes{
		Route{
			"Index",
			"GET",
			"/v1/",
			Index,
		},

		Route{
			"GetAccount",
			strings.ToUpper("Get"),
			fmt.Sprintf("/v1/accounts/{accountId:%s}", uuidRegex),
			s.GetAccount,
		},

		Route{
			"GetDevices",
			strings.ToUpper("Get"),
			fmt.Sprintf("/v1/accounts/{accountId:%s}/devices", uuidRegex),
			s.GetDevices,
		},

		Route{
			"UpdateAccount",
			strings.ToUpper("Patch"),
			fmt.Sprintf("/v1/accounts/{accountId:%s}", uuidRegex),
			s.UpdateAccount,
		},

		Route{
			"Auth",
			strings.ToUpper("Post"),
			"/v1/auth",
			s.Auth,
		},

		Route{
			"UpdateDevice",
			strings.ToUpper("Patch"),
			"/v1/devices/{deviceId}",
			s.UpdateDevice,
		},
	}

	// Build the router
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	// Add OPTIONS routes for each one to allow CORS
	for _, route := range routes {
		var handler http.Handler
		handler = s.optionsHandler

		router.
			Methods("OPTIONS").
			Path(route.Pattern).
			Name(fmt.Sprintf("%sOptions", route.Name)).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
