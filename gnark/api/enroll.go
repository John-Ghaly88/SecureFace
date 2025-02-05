package api

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"net/http"

	"ppba_project/gnark/circuit"
	"ppba_project/gnark/db"
)

// HelperEntry represents a single helper array with its data, shape, and dtype.
type HelperEntry struct {
	Data  string `json:"data"`
	Shape []int  `json:"shape"`
	Dtype string `json:"dtype"`
}

type requestPayload struct {
	Username string        `json:"username"`
	Key      string        `json:"key"`
	Helper   []HelperEntry `json:"helper"` // Changed from string to a structured type
}

type responsePayload struct {
	Username string `json:"username"`
	Status   string `json:"status"`
	Message  string `json:"message,omitempty"`
}

// handleEnroll is the HTTP handler that receives username, key, and helper data,
// generates a proof, and stores username, proof, and helper in the database.
func HandleEnroll(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req requestPayload
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Convert the key from string to big.Int
		keyBig := new(big.Int)
		if _, ok := keyBig.SetString(req.Key, 16); !ok {
			http.Error(w, "Invalid key format; must be a hex string", http.StatusBadRequest)
			return
		}

		// Marshal the helper entries back into JSON bytes so we can store them easily.
		// This preserves the structure so it can be reconstructed later.
		helperBytes, err := json.Marshal(req.Helper)
		if err != nil {
			http.Error(w, "Invalid helper format", http.StatusBadRequest)
			return
		}

		// Generate the proof (assuming your GenerateProof can handle the raw bytes and store them).
		proof, username, helperData, err := circuit.GenerateProof(req.Username, keyBig, helperBytes)
		if err != nil {
			log.Printf("Error generating proof: %v", err)
			http.Error(w, "Failed to generate proof", http.StatusInternalServerError)
			return
		}

		// Save to database
		ctx := context.Background()
		if err = db.SaveProofData(ctx, username, proof, helperData); err != nil {
			log.Printf("Error saving to database: %v", err)
			http.Error(w, "Failed to save data", http.StatusInternalServerError)
			return
		}

		// Respond
		resp := responsePayload{
			Username: username,
			Status:   "success",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
