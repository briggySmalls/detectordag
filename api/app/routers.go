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
		Route{
			"GetAccount",
			http.MethodGet,
			fmt.Sprintf("/{accountId:%s}", uuidRegex),
			server.GetAccount,
		},
		Route{
			"GetDevices",
			http.MethodGet,
			fmt.Sprintf("/{accountId:%s}/devices", uuidRegex),
			server.GetDevices,
		},
		Route{
			"UpdateAccount",
			http.MethodPatch,
			fmt.Sprintf("/{accountId:%s}", uuidRegex),
			server.UpdateAccount,
		},
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