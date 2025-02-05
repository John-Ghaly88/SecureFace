package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ppba_project/gnark/db"
)

func HandleRetrieve(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		username := r.URL.Query().Get("username")
		if username == "" {
			http.Error(w, "username query param required", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		proofBytes, helperBytes, err := db.GetProofData(ctx, username)
		if err != nil {
			http.Error(w, "No data found for user", http.StatusNotFound)
			return
		}

		// We store helper as raw JSON in DB, so just pass it along
		// or re-encode if necessary
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]string{
			"proof":  fmt.Sprintf("%x", proofBytes),
			"helper": string(helperBytes), // Already JSON if stored that way
		}
		json.NewEncoder(w).Encode(resp)
	}
}
