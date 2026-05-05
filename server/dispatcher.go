package server

import (
	"fmt"
	"net"
)

func Start(port string, numWorkers int , router *Router) {
	// Creating the TCP Listener on the given port
	listener, err := net.Listen("tcp" , port);

	if err != nil {
		fmt.Println("Failed to start server: " , err)
		return
	}

	// Ensuring that the socket closes when the server shuts down
	defer listener.Close()

	fmt.Printf("Server successfully listening on port %s...\n", port)

	// Creating the Buffered Channel (The Job Queue)
	// This holds up to 100 pending connections at a time

	jobQueue := make(chan net.Conn , 100)

	// These goroutines will just wait in the background for jobs
	for i := 1; i <= numWorkers; i++ {
		go worker(i , jobQueue , router)
	}

	// The Infinite Accept Loop
	for {
		// Blocks until a client connects
		conn , err := listener.Accept()

		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Pushing the new connection into the channel queue
		jobQueue <- conn
	}
}

