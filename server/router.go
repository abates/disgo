package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc appHandler
}

type Routes []Route

var routes = Routes{
	Route{
		"Search",
		"POST",
		"/search",
		Search,
	},
	Route{
		"Add",
		"POST",
		"/add",
		Add,
	},
	Route{
		"Show",
		"GET",
		"/image/{hash}",
		Show,
	},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	return router
}
