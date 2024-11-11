package service

import (
	"net/http"
	"time"
)

// HttpRequest defines the interface for making HTTP requests.
// This interface includes a method for performing GET requests.
type HttpRequest interface {
	Get(url string) (http.Response, error) // Method to perform a GET request
}

// httpRequest is a concrete implementation of HttpRequest.
// It holds an HTTP client configured with a timeout.
type httpRequest struct {
	client http.Client // HTTP client for making requests
}

// NewHttpRequestService creates a new instance of httpRequest.
// It initializes the HTTP client with a specified request timeout.
//
// Parameters:
//   - requestTimeout: Duration to set the timeout for HTTP requests.
//
// Returns:
//   - An instance of HttpRequest.
func NewHttpRequestService(requestTimeout time.Duration) HttpRequest {
	client := http.Client{
		Timeout: requestTimeout, // Set the timeout for the HTTP client
	}

	return &httpRequest{
		client: client, // Return an instance of httpRequest with the configured client
	}
}

// Get performs a GET request to the specified URL.
//
// Parameters:
//   - url: The URL to send the GET request to.
//
// Returns:
//   - The HTTP response and any error encountered during the request.
func (hr *httpRequest) Get(url string) (http.Response, error) {
	req, err := http.NewRequest("GET", url, nil) // Create a new GET request
	if err != nil {
		return http.Response{}, err // Return an empty response and the error
	}

	resp, err := hr.client.Do(req) // Execute the GET request using the HTTP client
	if err != nil {
		return http.Response{}, err // Return an empty response and the error
	}

	return *resp, nil // Return the response from the GET request
}
