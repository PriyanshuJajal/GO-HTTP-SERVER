package server

import (
	"GO-HTTP-SERVER/protocol"
)

// A custom type for the Handler Functions
type RouteHandler func(protocol.HTTPRequest) protocol.HTTPResponse

// The Router Struct holds a map of routes
type Router struct {
	// The key will be "METHOD URI" (e.g., "GET /users")
	routes map[string]RouteHandler
}

// NewRouter initializes the map in memory
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]RouteHandler),
	}
}


// Handle registers a new route in the map
func (r *Router) Handle(method, uri string, handler RouteHandler) {
	key := method + " " + uri
	r.routes[key] = handler
}


// Route matches an incoming request to a registered handler
func (r *Router) Route(req protocol.HTTPRequest) protocol.HTTPResponse {
	key := req.Method + " " + req.URI

	if handler , exists := r.routes[key]; exists {
		return handler(req)
	}

	return protocol.HTTPResponse {
		StatusCode: 404,
		ContentType: "text/plain",
		Body: "404 Not Found\n",
	}
}