package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Sale struct {
	Product  string  `json:"product"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Date     string  `json:"date"`
	Error    string  `json:"error,omitempty"`
}

var products = []string{"Laptop", "Smartphone", "Tablet", "Headphones", "Monitor"}

func randomSale() (Sale, error) {
	rand.Seed(time.Now().UnixNano())
	if rand.Float32() < 0.1 {
		return Sale{}, fmt.Errorf("simulated error")
	}
	return Sale{
		Product:  products[rand.Intn(len(products))],
		Price:    rand.Float64() * 1000,
		Quantity: rand.Intn(10) + 1,
		Date:     time.Now().Format(time.RFC3339),
	}, nil
}

func logToFile(filename string, data []byte) {
	logDir := "../logs"
	os.MkdirAll(logDir, os.ModePerm)
	filePath := filepath.Join(logDir, filename)

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v\n", err)
		return
	}
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		log.Printf("Error writing to log file: %v\n", err)
	}
	if _, err := f.WriteString("\n"); err != nil {
		log.Printf("Error writing newline to log file: %v\n", err)
	}
}

func saleHandler(w http.ResponseWriter, r *http.Request) {
	sale, err := randomSale()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		errorLog := Sale{
			Product:  "N/A",
			Price:    0,
			Quantity: 0,
			Date:     time.Now().Format(time.RFC3339),
			Error:    err.Error(),
		}
		jsonData, _ := json.Marshal(errorLog)
		logToFile("errors.log", jsonData)
		return
	}
	jsonData, err := json.Marshal(sale)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		errorLog := Sale{
			Product:  "N/A",
			Price:    0,
			Quantity: 0,
			Date:     time.Now().Format(time.RFC3339),
			Error:    err.Error(),
		}
		jsonData, _ = json.Marshal(errorLog)
		logToFile("errors.log", jsonData)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
	logToFile("sales.log", jsonData)
	log.Printf("Generated sale: %s\n", jsonData)
}

func main() {
	http.HandleFunc("/sale", saleHandler)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
