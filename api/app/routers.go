// Detectordag
//
// API for detectordag backend
//
//    Schemes: https
//    Host: detectorkag.tk
//    BasePath: /v1
//    Version: 0.0.1
//    Contact: Sam Briggs<detectordag@sambriggs.dev>
//
//    Consumes:
//    - application/json
//
//    Produces:
//    - application/json
//
// swagger:meta
package app

import (
	"fmt"
	"net/http"

	"github.com/briggysmalls/detectordag/api/app/server"
	"github.com/briggysmalls/detectordag/api/app/tokens"
	"github.com/briggysmalls/detectordag/shared/iot"
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

func NewRouter(iot iot.Client, server server.Server, tokens tokens.Tokens) *mux.Router {
	// Create the router
	router := mux.NewRouter().StrictSlash(true)
	// Create subrouter for 'v1'
	api := router.PathPrefix("/v1").Subrouter()
	// Add the non-auth routes
	nonAuthRoutes := Routes{
		// swagger:route POST /auth authentication auth
		//
		// Obtain token for the site
		//
		//     Responses:
		//       200: tokenResponse
		//       403: authFailedResponse
		Route{
		   "Auth",
			http.MethodPost,
			"/auth",
			server.Auth,
		},
	}
	addRoutes(api, nonAuthRoutes)

	// Create subrouter for accounts
	accounts := api.PathPrefix("/accounts").Subrouter()
	addRoutes(accounts, Routes{
		// swagger:route GET /accounts/{accountId} accounts getAccount
		//
		// Get account details
		//
		//     Responses:
		//       200: getAccountResponse
		//       400: accountNotFoundResponse
		//       401: unauthenticatedResponse
		//       403: unauthorizedResponse
		Route{
			"GetAccount",
			http.MethodGet,
			fmt.Sprintf("/{accountId:%s}", uuidRegex),
			server.GetAccount,
		},
		// swagger:route GET /accounts/{accountId}/devices accounts getDevices
		//
		// Get all devices associated with the user's account
		//
		//     Responses:
		//       200: getDevicesResponse
		//       400: accountNotFoundResponse
		//       401: unauthenticatedResponse
		//       403: unauthorizedResponse
		Route{
			"GetDevices",
			http.MethodGet,
			fmt.Sprintf("/{accountId:%s}/devices", uuidRegex),
			server.GetDevices,
		},
		// swagger:route PATCH /accounts/{accountId} accounts updateAccount
		//
		// Update an account
		//
		// Update account configuration
		//
		//     Responses:
		//       200: getAccountResponse
		//       400: accountNotFoundResponse
		//       401: unauthenticatedResponse
		//       403: unauthorizedResponse
		Route{
			"UpdateAccount",
			http.MethodPatch,
			fmt.Sprintf("/{accountId:%s}", uuidRegex),
			server.UpdateAccount,
		},
		// swagger:route PUT /accounts/{accountId}/devices/{deviceId} accounts registerDevice
		//
		// Register a new device
		//
		// Register a new device to the account
		//
		//     Responses:
		//       200: getDeviceResponse
		//       400: accountNotFoundResponse
		//       401: unauthenticatedResponse
		//       403: unauthorizedResponse
		Route{
			"RegisterDevice",
			http.MethodPut,
			fmt.Sprintf("/{accountId:%s}/devices/{deviceId:%s}", uuidRegex, uuidRegex),
			server.RegisterDevice,
		},
	})

	// Create subrouter for devices
	devices := api.PathPrefix("/devices").Subrouter()
	addRoutes(devices, Routes{
		// swagger:route PATCH /devices/{deviceId} devices updateDevice
		//
		// Update a device
		//
		// Update device configuration
		//
		//     Responses:
		//       200: getDeviceResponse
		//       400: deviceNotFoundResponse
		//       401: unauthenticatedResponse
		//       403: unauthorizedResponse
		Route{
			"UpdateDevice",
			http.MethodPatch,
			fmt.Sprintf("/{deviceId:%s}", uuidRegex),
			server.UpdateDevice,
		},
	})

	// Add CORS header on all responses
	api.Use(mux.CORSMethodMiddleware(api))
	api.Use(corsMiddleware)

	// Add authentication middleware
	a := auth{
		tokens: tokens,
		iot:    iot,
	}
	accounts.Use(a.middleware)
	devices.Use(a.middleware)

	// Return the router
	return router
}

func addRoutes(router *mux.Router, routes []Route) {
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		router.
			Methods(route.Method, http.MethodOptions).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
}
