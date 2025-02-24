package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var bigData []map[string]interface{}

func loadBigData() {
	file, err := os.Open("./tess/big.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&bigData); err != nil {
		log.Fatal(err)
	}

	log.Println("Loaded big.json")
}

func handlerGoJSONGzip(w http.ResponseWriter, r *http.Request) {
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	if err := json.NewEncoder(bw).Encode(bigData); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "JSON encoding error: %v"}`, err), http.StatusInternalServerError)
		return
	}
}

func main() {
	loadBigData()

	http.HandleFunc("/go-json-gzip", handlerGoJSONGzip)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
