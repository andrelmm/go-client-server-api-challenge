package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	serverAddress = "http://localhost:8080"
	clientTimeout = 300 * time.Millisecond
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), clientTimeout)
	defer cancel()

	// Fetching quotation
	quotation, err := fetchQuotation(ctx)
	if err != nil {
		log.Fatal("Error fetching quotation: ", err)
	}

	// Saving quotation to file
	if err := saveQuotationToFile(quotation); err != nil {
		log.Fatal("Error saving quotation: ", err)
	}
}

func fetchQuotation(ctx context.Context) (string, error) {
	client := &http.Client{}

	// Creating request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/cotacao", serverAddress), nil)
	if err != nil {
		return "", err
	}

	// Sending request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Reading response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parsing response
	var quotation map[string]string
	if err := json.Unmarshal(body, &quotation); err != nil {
		return "", err
	}

	bid, ok := quotation["bid"]
	if !ok {
		return "", fmt.Errorf("bid not found in quotation")
	}

	return bid, nil
}

func saveQuotationToFile(quotation string) error {
	file, err := os.Create("quotation.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "Dólar: %s\n", quotation)

	if err != nil {
		return err
	}
	return nil
}
