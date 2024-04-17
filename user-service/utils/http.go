package utils

import (
	"net/http"

	"github.com/rs/zerolog"
)

func SendGetRequest(url string, logger zerolog.Logger) (*http.Response, error) {
	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error().Msg("Failed to create HTTP request")
		return nil, err
	}

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error().Msg("Failed to send HTTP request")
		return nil, err
	}

	return resp, nil
}
