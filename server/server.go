package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"time"
)

const (
	serverPort      = ":8080"
	serverTimeOut   = 200 * time.Millisecond
	databaseTimeout = 10 * time.Millisecond
)

type Quotation struct {
	Bid string `json:"bid"`
}

func main() {

	// Opening connection with database
	db, err := sql.Open("sqlite3", "./quotation.db")
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	defer db.Close()

	// Create table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS quotation (id INTEGER PRIMARY KEY AUTOINCREMENT, bid TEXT, created_at DATETIME)`)
	if err != nil {
		log.Fatal("Error creating table: ", err)
	}

	// Defining routes
	http.HandleFunc("/cotacao", quotationHandler(db))

	// Starting server
	log.Println("Starting server on port", serverPort)
	http.ListenAndServe(serverPort, nil)
}

func quotationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		ctx, cancel := context.WithTimeout(ctx, serverTimeOut)
		defer cancel()

		// Fetching quotation
		quotation, err := fetchQuotation(ctx)
		if err != nil {
			http.Error(w, "Error fetching quotation", http.StatusInternalServerError)
			return
		}

		// Saving quotation into database
		if err := createQuotation(db, quotation); err != nil {
			http.Error(w, "Error saving quotation", http.StatusInternalServerError)
			return
		}

		sendJSONResponse(w, quotation)
	}
}

func createQuotation(db *sql.DB, quotation *Quotation) error {
	dbCtx, dbCancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer dbCancel()

	_, err := db.ExecContext(dbCtx, "INSERT INTO quotation (bid, created_at) VALUES (?, ?)", quotation.Bid, time.Now())
	return err
}

func fetchQuotation(ctx context.Context) (*Quotation, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	var data map[string]Quotation
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	quotation, ok := data["USDBRL"]
	if !ok {
		return nil, fmt.Errorf("exchange rate not found in response")
	}

	return &quotation, nil
}

func sendJSONResponse(w http.ResponseWriter, quotation *Quotation) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quotation)
}
