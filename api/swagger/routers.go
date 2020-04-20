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

	"github.com/briggysmalls/detectordag/api/swagger/server"
	"github.com/briggysmalls/detectordag/shared/database"
	"github.com/briggysmalls/detectordag/shared/email"
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

type RouterConfig struct {
	db     database.Client
	shadow shadow.Client
	email  email.Client
}

type Routes []Route

func NewRouter(s server.Server) *mux.Router {
	// Create the router
	router := mux.NewRouter().StrictSlash(true)
	// Create subrouter for 'v1'
	api := router.PathPrefix("/v1").Subrouter()
	// Prepare the routes
	var routes = Routes{
		Route{
			"GetAccount",
			http.MethodGet,
			fmt.Sprintf("/accounts/{accountId:%s}", uuidRegex),
			s.GetAccount,
		},

		Route{
			"GetDevices",
			http.MethodGet,
			fmt.Sprintf("/accounts/{accountId:%s}/devices", uuidRegex),
			s.GetDevices,
		},

		Route{
			"UpdateAccount",
			http.MethodPatch,
			fmt.Sprintf("/accounts/{accountId:%s}", uuidRegex),
			s.UpdateAccount,
		},

		Route{
			"Auth",
			http.MethodPost,
			"/auth",
			s.Auth,
		},

		Route{
			"UpdateDevice",
			http.MethodPatch,
			"/devices/{deviceId}",
			s.UpdateDevice,
		},
	}

	// Build the router
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc

		api.
			Methods(route.Method, http.MethodOptions).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	// Add CORS header on all responses
	api.Use(mux.CORSMethodMiddleware(api))
	api.Use(corsMiddleware)

	return router
}
