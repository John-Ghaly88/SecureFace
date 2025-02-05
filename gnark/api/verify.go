package api

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"

    "ppba_project/gnark/circuit"
    "ppba_project/gnark/db"
)

func HandleVerify(db *database.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            msg := "Method not allowed"
            log.Println("handleVerify:", msg)
            http.Error(w, msg, http.StatusMethodNotAllowed)
            return
        }

        var req struct {
            Username string `json:"username"`
            Key      string `json:"key"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            msg := fmt.Sprintf("JSON decode error: %v", err)
            log.Println("handleVerify:", msg)
            http.Error(w, msg, http.StatusBadRequest)
            return
        }

        log.Printf("handleVerify: verifying user '%s' with key '%s'", req.Username, req.Key)

        // Call the circuit.VerifyProof
        verified, err := circuit.VerifyProof(req.Username, req.Key, db)
        if err != nil {
            // This means an internal error occurred, not just a failed verification
            msg := fmt.Sprintf("internal error in VerifyProof: %v", err)
            log.Println("handleVerify:", msg)
            http.Error(w, msg, http.StatusInternalServerError)
            return
        }

        if !verified {
            // Means the proof check failed (e.g. invalid proof)
            msg := "Verification failed for user " + req.Username
            log.Println("handleVerify:", msg)
            http.Error(w, msg, http.StatusUnauthorized)
            return
        }

        // If we reach here, verification succeeded
        log.Printf("handleVerify: verification succeeded for user '%s'", req.Username)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{"status": "verified"})
    }
}

