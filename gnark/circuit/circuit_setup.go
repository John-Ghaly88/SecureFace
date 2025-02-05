package circuit

import (
	"fmt"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

func SetupCircuit(provingKeyPath, verifyingKeyPath string) error {
	var c Circuit
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &c)
	if err != nil {
		return fmt.Errorf("failed to compile circuit: %v", err)
	}

	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		return fmt.Errorf("failed to setup keys: %v", err)
	}

	if err := saveProvingKey(pk, provingKeyPath); err != nil {
		return fmt.Errorf("failed to save proving key: %v", err)
	}

	if err := saveVerifyingKey(vk, verifyingKeyPath); err != nil {
		return fmt.Errorf("failed to save verifying key: %v", err)
	}

	return nil
}

func saveProvingKey(pk groth16.ProvingKey, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Error creating proving key file: %v", err)
	}
	defer file.Close()

	_, err = pk.WriteTo(file)
	if err != nil {
		return fmt.Errorf("Error writing proving key: %v", err)
	}

	fmt.Println("Debug: Proving key saved to", filename)
	return nil
}

func saveVerifyingKey(vk groth16.VerifyingKey, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = vk.WriteTo(file)
	return err
}
