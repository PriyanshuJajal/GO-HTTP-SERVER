package protocol

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type HTTPRequest struct {
	Method string
	URI    string
	Body   string
}

type HTTPResponse struct {
	StatusCode  int
	ContentType string
	Body        string
}

// ParseRequest handles all the raw byte reading and HTTP formatting
func ParseRequest(conn net.Conn) (HTTPRequest, error) {
	// Wrap the raw socket in a Buffered Reader
	reader := bufio.NewReader(conn)

	// Parsing the Request Line (e.g., "GET /api/data HTTP/1.1\r\n")
	reqLine, err := reader.ReadString('\n')
	if err != nil {
		return HTTPRequest{}, err
	}

	// Clean up the string and split it to get the Method and URI
	parts := strings.Split(strings.TrimSpace(reqLine), " ")

	if len(parts) < 3 {
		return HTTPRequest{}, fmt.Errorf("invalid request line")
	}
	
	req := HTTPRequest{
		Method: parts[0],
		URI:    parts[1],
	}

	// Parsing Headers to find the \r\n\r\n delimiter and Content-Length
	contentLength := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return HTTPRequest{}, err
		}

		// Trim the \r\n off the end of the line
		line = strings.TrimSpace(line)

		// An empty line means we hit the \r\n\r\n delimiter! Headers are done.
		if line == "" {
			break 
		}

		// If it's a POST request, we MUST find the Content-Length
		if strings.HasPrefix(strings.ToLower(line), "content-length:") {
			headerParts := strings.Split(line, ":")
			if len(headerParts) == 2 {
				contentLength, _ = strconv.Atoi(strings.TrimSpace(headerParts[1]))
			}
		}
	}

	if req.Method == "POST" && contentLength > 0 {
		bodyBytes := make([]byte, contentLength)
		_, err := io.ReadFull(reader, bodyBytes)
		if err == nil {
			req.Body = string(bodyBytes)
		}
	}

	return req , nil
}

// BuildResponse turns our struct back into raw network bytes
func BuildResponse(res HTTPResponse) []byte {
	// Map status codes to exact HTTP status phrases
	statusPhrase := "200 OK"

	switch res.StatusCode {
	case 201:
		statusPhrase = "201 Created"
	case 404:
		statusPhrase = "404 Not Found"
	case 500:
		statusPhrase = "500 Internal Server Error"
	}

	httpResponse := fmt.Sprintf(
		"HTTP/1.1 %s\r\nContent-Length: %d\r\nContent-Type: %s\r\nConnection: close\r\n\r\n%s",
		statusPhrase, len(res.Body), res.ContentType, res.Body,
	)

	return []byte(httpResponse)
}