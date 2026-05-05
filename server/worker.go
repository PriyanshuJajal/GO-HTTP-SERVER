package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

// worker listens to the job queue and processes incoming connections.
func worker(id int , jobs <- chan net.Conn , router *Router) {

	// It blocks until a connection is pushed into the channel
	for conn := range jobs {
		handleConnection(id , conn , router)
	}
}

// handleConnection does the actual HTTP parsing and responding.
func handleConnection(id int , conn net.Conn , router *Router) {
	defer conn.Close()

	// Wrap the raw socket in a Buffered Reader
	reader := bufio.NewReader(conn)

	// Parse the Request Line (e.g., "GET /api/data HTTP/1.1\r\n")
	reqLine , err := reader.ReadString('\n') 

	if err != nil {
		fmt.Printf("[Worker %d] Error reading socket: %v\n", id, err)
		return
	}

	// Clean up the string and split it to get the Method and URI
	reqLine = strings.TrimSpace(reqLine)
	parts := strings.Split(reqLine , " ")

	if len(parts) < 3 {
		return
	}

	method := parts[0]
	uri := parts[1]

	fmt.Printf("[Worker %d] Processing: %s %s\n", id, method, uri)

	// Parsing Headers to find the \r\n\r\n delimiter and Content-Length
	contentLen := 0
	for {
		line , err := reader.ReadString('\n')

		if err != nil {
			return
		}

		// Trim the \r\n off the end of the line
		line = strings.TrimSpace(line)

		// An empty line means we hit the \r\n\r\n delimiter! Headers are done.
		if line == "" {
			break
		}

		// If it's a POST request, we MUST find the Content-Length
		if strings.HasPrefix(strings.ToLower(line) , "content-length:") {
			headerParts := strings.Split(line , ":")

			if len(headerParts) == 2 {
				lenStr := strings.TrimSpace(headerParts[1])
				contentLen , _ = strconv.Atoi(lenStr)
			}
		}
	}

	// Read the POST Body (if it exists)
	body := ""
	if method == "POST" && contentLen > 0 {
		bodyBytes := make([]byte , contentLen)

		// Read exactly that many bytes from the raw socket
		_ , err := io.ReadFull(reader , bodyBytes)

		if err == nil {
			body = string(bodyBytes)
			fmt.Printf("[Worker %d] POST Body: %s\n", id, body)
		}
	}

	req := HTTPRequest {
		Method: method,
		URI: uri,
		Body: body,
	}

	res := router.Route(req)

	// Map status codes to exact HTTP status phrases
	statusPhrase := "200 OK"

	switch res.StatusCode {
	case 201:
		statusPhrase = "201 created"
	case 404:
		statusPhrase = "404 Not Found"
	case 500:
		statusPhrase = "500 Internal Server Error"
	}

	httpResponse := fmt.Sprintf(
		"HTTP/1.1 %s\r\nContent-Length: %d\r\nContent-Type: %s\r\nConnection: close\r\n\r\n%s",
		statusPhrase, len(res.Body), res.ContentType, res.Body,
	)

	// Write the bytes back over the network to the client
	conn.Write([]byte(httpResponse))
}

