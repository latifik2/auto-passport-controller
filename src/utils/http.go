package utils

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

func MakeServiceApiCall(host string, endpoint string, b64Creds string, scheme string) ([]byte, error) {

	_ = b64Creds
	// log.Printf("Airflow Collector initialized with credentials: %s", b64Creds)

	url := scheme + "://" + host + endpoint

	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating HTTP request: %v", err))
		return []byte{}, err
	}

	// req.Header.Add("Authorization", "Basic "+b64Creds)

	slog.Info(fmt.Sprintf("Making API request to host: %s", host))

	resp, err := httpClient.Do(req)
	if err != nil {
		slog.Error(fmt.Sprintf("Error making HTTP request: %v", err))
		return []byte{}, err
	}
	defer resp.Body.Close()

	slog.Info(fmt.Sprintf("Response from host %s: %s", host, resp.Status))

	body, _ := io.ReadAll(resp.Body)

	slog.Debug(fmt.Sprintf("Airflow reponse Body: %s", body))

	return body, nil
}

func GetCredentials(serviceType string) (string, string) {
	return "auto-passport", "7lXrvaRCwJgvEVq3Gzpu0jqHUqCbUF6B"
}
