package main

import (
	"GO-HTTP-SERVER/server"
	"fmt"
)

func main() {
	port := "8080"
	numWorkers := 10 // Spawn 10 concurrent threads (goroutines)

	fmt.Println("=====================================")
	fmt.Println("   Starting Go Micro HTTP Server     ")
	fmt.Println("=====================================")
	fmt.Printf("Thread Pool Size: %d\n", numWorkers)

	// Call the Start function from dispatcher.go
	server.Start(port , numWorkers)
}