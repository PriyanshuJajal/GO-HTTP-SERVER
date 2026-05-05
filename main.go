package main

import (
	"fmt"
	"GO-HTTP-SERVER/handler"
	"GO-HTTP-SERVER/server"
)

func main() {
	port := ":8080"
	numWorkers := 10

	router := server.NewRouter()

	// Routes for GET & POST request
	router.Handle("GET", "/api/status", handlers.GetStatus)
	router.Handle("POST", "/api/data", handlers.PostData)

	fmt.Println("=====================================")
	fmt.Println("   Starting Go Micro HTTP Server     ")
	fmt.Println("=====================================")
	fmt.Printf("\n Thread Pool Size: %d\n", numWorkers)

	server.Start(port , numWorkers , router)
}