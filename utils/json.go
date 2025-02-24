package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")

	// pr, pw := io.Pipe()

	// go func() {
	// 	defer pw.Close()
	// 	if err := json.NewEncoder(pw).Encode(v); err != nil {
	// 		pw.CloseWithError(err)
	// 		return
	// 	}
	// }()

	// if _, err := io.Copy(w, pr); err != nil {
	// 	http.Error(w, fmt.Sprintf(`{"error": "data transfer error: %v"}`, err), http.StatusInternalServerError)
	// 	return
	// }
	// if err := json.NewEncoder(w).Encode(v); err != nil {
	// 	http.Error(w, fmt.Sprintf(`{"error": "JSON encoding error: %v"}`, err), http.StatusInternalServerError)
	// }
	// bw := bufio.NewWriter(w)
	// defer bw.Flush()

	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "JSON encoding error: %v"}`, err), http.StatusInternalServerError)
	}
}
