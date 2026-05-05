package server

import (
	"fmt"
	"net"
	"GO-HTTP-SERVER/protocol"
)

func worker(id int, jobs <-chan net.Conn, router *Router) {
	for conn := range jobs {
		handleConnection(id, conn, router)
	}
}

func handleConnection(id int, conn net.Conn, router *Router) {
	defer conn.Close()

	// Parse the incoming bytes
	req , err := protocol.ParseRequest(conn)

	if err != nil {
		fmt.Printf("[Worker %d] Error parsing request: %v\n", id, err)
		return
	}

	fmt.Printf("[Worker %d] Processing: %s %s\n", id, req.Method, req.URI)

	// Route the request
	res := router.Route(req)

	// Send the bytes back
	responseBytes := protocol.BuildResponse(res)
	conn.Write(responseBytes)
}