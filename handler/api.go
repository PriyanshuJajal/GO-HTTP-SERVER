package handlers

import (
	"fmt"
	"GO-HTTP-SERVER/protocol" 
)

// GetStatus handles the GET /api/status route
func GetStatus(req protocol.HTTPRequest) protocol.HTTPResponse {
	return protocol.HTTPResponse{
		StatusCode:  200,
		ContentType: "application/json",
		Body: `{"status":"Server is running smoothly!"}`,
	}
}

// PostData handles the POST /api/data route
func PostData(req protocol.HTTPRequest) protocol.HTTPResponse {
	responseBody := fmt.Sprintf(`{"message":"Data received!", "yourData": %s}`, req.Body)
	return protocol.HTTPResponse{
		StatusCode:  201,
		ContentType: "application/json",
		Body: responseBody,
	}
}

func GetHome(req protocol.HTTPRequest) protocol.HTTPResponse {
	return  protocol.HTTPResponse{
		StatusCode: 200,
		ContentType: "text/html", 
		Body: `
			<html>
				<body style="font-family: Arial; text-align: center; margin-top: 50px;">
					<h1>Go Micro Server is Live!</h1>
					<p>Successfully rendering HTML from a custom TCP socket.</p>
				</body>
			</html>
		`,
	}
}