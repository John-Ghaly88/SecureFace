package circuit

import (
    "bytes"
    "context"
    "encoding/hex"
    "fmt"
    "math/big"
    "os"

    "github.com/consensys/gnark-crypto/ecc"
    "github.com/consensys/gnark/backend/groth16"
    "github.com/joho/godotenv"
    "github.com/consensys/gnark/frontend"

    "ppba_project/gnark/db"
)

// VerifyProof fetches the stored proof from DB for username,
// then uses reproduced key (keyHex) to recompute proofHash and verify.
func VerifyProof(username string, keyHex string, db *database.DB) (bool, error) {
    // 1. Decode the key from hex
    keyBytes, err := hex.DecodeString(keyHex)
    if err != nil {
        return false, fmt.Errorf("VerifyProof: invalid key hex format: %w", err)
    }
    keyBig := new(big.Int).SetBytes(keyBytes)

    // 2. Retrieve proofBytes from DB
    ctx := context.Background()
    proofBytes, _, err := db.GetProofData(ctx, username)
    if err != nil {
        return false, fmt.Errorf("VerifyProof: DB retrieval error for user '%s': %w", username, err)
    }

    // 3. Load verifying key
    if err := godotenv.Load(); err != nil {
        // Not fatal, but log a warning
        fmt.Printf("VerifyProof: warning: .env not found: %v\n", err)
    }
    verifyingKeyPath := os.Getenv("VERIFYING_KEY_PATH")
    if verifyingKeyPath == "" {
        return false, fmt.Errorf("VerifyProof: VERIFYING_KEY_PATH not set in .env")
    }

    vkFile, err := os.Open(verifyingKeyPath)
    if err != nil {
        return false, fmt.Errorf("VerifyProof: error opening verifying key file: %w", err)
    }
    defer vkFile.Close()

    vk := groth16.NewVerifyingKey(ecc.BN254)
    if _, err := vk.ReadFrom(vkFile); err != nil {
        return false, fmt.Errorf("VerifyProof: error reading verifying key: %w", err)
    }

    // 4. Read proof from bytes
    proof := groth16.NewProof(ecc.BN254)
    buf := bytes.NewReader(proofBytes)
    if _, err := proof.ReadFrom(buf); err != nil {
        return false, fmt.Errorf("VerifyProof: error reading proof bytes: %w", err)
    }

    // 5. Compute the public input (mimcHash)
    proofHash := mimcHash(keyBig)

    // 6. Create public-only witness
    assignment := &Circuit{
        Proof: frontend.Variable(proofHash),
    }
    publicWitness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField(), frontend.PublicOnly())
    if err != nil {
        return false, fmt.Errorf("VerifyProof: error creating public witness: %w", err)
    }

    // 7. Perform the proof verification
    err = groth16.Verify(proof, vk, publicWitness)
    if err != nil {
        // Verification mismatch -> no internal error, but the proof is invalid.
        // Return (false, nil) to indicate "verification failed" without a system error.
        return false, nil
    }

    // Verification successful
    return true, nil
}

