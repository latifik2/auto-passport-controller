package utils

import (
	"io"
	"log"
	"net/http"
)

func MakeApiCall(url string, b64Creds string) ([]byte, error) {

	_ = b64Creds
	// log.Printf("Airflow Collector initialized with credentials: %s", b64Creds)

	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating HTTP request: %v", err)
		return []byte{}, err
	}

	// req.Header.Add("Authorization", "Basic "+b64Creds)

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error making HTTP request: %v", err)
		return []byte{}, err
	}
	defer resp.Body.Close()

	log.Printf("Airflow response status: %s", resp.Status)
	body, _ := io.ReadAll(resp.Body)
	log.Printf("Airflow reponse Body: %s", body)

	return body, nil
}

func GetCredentials(serviceType string) (string, string) {
	return "auto-passport", "7lXrvaRCwJgvEVq3Gzpu0jqHUqCbUF6B"
}
