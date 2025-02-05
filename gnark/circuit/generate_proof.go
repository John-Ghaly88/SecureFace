package circuit

import (
	"bytes"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	bn254 "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/joho/godotenv"
)

// Hash a single big.Int using MiMC
func mimcHash(val *big.Int) *big.Int {
	f := bn254.NewMiMC()
	f.Write(val.Bytes())
	hash := f.Sum(nil)
	hashInt := new(big.Int).SetBytes(hash)
	return hashInt
}

// GenerateProof takes a username, a key (the biometric extract), and helper data.
// - username: returned as-is, no circuit involvement.
// - key: used in the circuit to generate a proof: Proof = MiMC(key).
// - helper: returned as-is, not used in the circuit.
func GenerateProof(username string, key *big.Int, helper_data []byte) ([]byte, string, []byte, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Get proving key path from .env
	provingKeyPath := os.Getenv("PROVING_KEY_PATH")
	if provingKeyPath == "" {
		return nil, username, helper_data, fmt.Errorf("PROVING_KEY_PATH not set in .env file")
	}

	// Load the proving key
	file, err := os.Open(provingKeyPath)
	if err != nil {
		return nil, username, helper_data, fmt.Errorf("Error opening proving key file: %v", err)
	}
	defer file.Close()

	pk := groth16.NewProvingKey(ecc.BN254)
	_, err = pk.ReadFrom(file)
	if err != nil {
		return nil, username, helper_data, fmt.Errorf("Error reading proving key: %v", err)
	}

	// Compute hash for proof = MiMC(key)
	proofHash := mimcHash(key)

	// Define the circuit assignment
	assignment := &Circuit{
		Key:   frontend.Variable(key),
		Proof: frontend.Variable(proofHash),
	}

	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Key: %v\n", key)
	fmt.Printf("Helper (JSON): %s\n", helper_data)
	fmt.Printf("Proof (Hash): %v\n", proofHash)

	// Compile the circuit
	var c Circuit
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &c)
	if err != nil {
		return nil, username, helper_data, fmt.Errorf("Error compiling circuit: %v", err)
	}

	// Create a full witness
	fullWitness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return nil, username, helper_data, fmt.Errorf("Error creating full witness: %v", err)
	}

	// Generate the proof
	proof, err := groth16.Prove(r1cs, pk, fullWitness)
	if err != nil {
		return nil, username, helper_data, fmt.Errorf("Error generating proof: %v", err)
	}

	// Serialize the proof
	var proofBytes bytes.Buffer
	_, err = proof.WriteTo(&proofBytes)
	if err != nil {
		return nil, username, helper_data, fmt.Errorf("Error serializing proof: %v", err)
	}

	// Return proof bytes, username, and helper
	return proofBytes.Bytes(), username, helper_data, nil
}
