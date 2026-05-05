package main

import (
	"GO-HTTP-SERVER/server"
	"fmt"
)

func main() {
	port := ":8080"
	numWorkers := 10 // Spawn 10 concurrent threads (goroutines)

	router := server.NewRouter()

	// Defining a GET Route
	router.Handle("GET" , "/api/status" , func(req server.HTTPRequest) server.HTTPResponse {
		return server.HTTPResponse{
			StatusCode: 200,
			ContentType: "application/json",
			Body: `{"status":"Server is running smoothly!"}`,
		}
	})

	// Defining a POST Route
	router.Handle("POST" , "/api/data" , func(req server.HTTPRequest) server.HTTPResponse {
		responseBody := fmt.Sprintf(`{"message":"Data received!", "yourData": %s}`, req.Body)
		
		return server.HTTPResponse{
			StatusCode: 201,
			ContentType: "application/json",
			Body: responseBody,
		}
	})

	fmt.Println("=====================================")
	fmt.Println("   Starting Go Micro HTTP Server     ")
	fmt.Println("=====================================")
	fmt.Printf("\nThread Pool Size: %d\n", numWorkers)

	// Call the Start function from dispatcher.go
	server.Start(port , numWorkers , router)
}