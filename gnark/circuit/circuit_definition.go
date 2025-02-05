package circuit

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

type Circuit struct {
	Key   frontend.Variable `gnark:",private"`
	Proof frontend.Variable `gnark:",public"`
}

func (circuit *Circuit) Define(api frontend.API) error {
	h, _ := mimc.NewMiMC(api)
	h.Write(circuit.Key)
	computedHash := h.Sum()

	// Proof must equal the hash of Key
	api.AssertIsEqual(circuit.Proof, computedHash)
	return nil
}
