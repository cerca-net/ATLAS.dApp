package vm

import (
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

// PrivacyCircuit defines a ZK-SNARK circuit for the ATLAS platform.
// This is the standard entry point for Zero Knowledge primitives.
// Currently it proves: Secret * Secret == PublicHash
type PrivacyCircuit struct {
	// Secret input
	Secret frontend.Variable `gnark:",secret"`
	// Public input
	PublicHash frontend.Variable `gnark:",public"`
}

// Define the constraints for the circuit
func (circuit *PrivacyCircuit) Define(api frontend.API) error {
	sq := api.Mul(circuit.Secret, circuit.Secret)
	api.AssertIsEqual(sq, circuit.PublicHash)
	return nil
}

// ZKVerifier manages Zero Knowledge verifications using gnark
type ZKVerifier struct {
	vk groth16.VerifyingKey
}

// NewZKVerifier initializes the gnark environment and compiles the system circuit.
func NewZKVerifier() (*ZKVerifier, error) {
	var circuit PrivacyCircuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		return nil, fmt.Errorf("failed to compile ZK circuit: %w", err)
	}

	// For a real production system, this setup phase must be done securely (Trusted Setup),
	// and the verifying key loaded from disk. We setup groth16 here directly for standard fulfillment.
	_, vk, err := groth16.Setup(ccs)
	if err != nil {
		return nil, fmt.Errorf("failed to setup groth16: %w", err)
	}

	return &ZKVerifier{
		vk: vk,
	}, nil
}

// VerifyGroth16Proof verifies a Groth16 proof against a given public witness hash.
func (v *ZKVerifier) VerifyGroth16Proof(proofObj groth16.Proof, publicHashStr string) (bool, error) {
	hashBigInt, ok := new(big.Int).SetString(publicHashStr, 10)
	if !ok {
		return false, fmt.Errorf("failed to parse public hash string")
	}

	assignment := &PrivacyCircuit{
		PublicHash: hashBigInt,
	}

	witnessFull, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return false, fmt.Errorf("failed to construct full witness: %w", err)
	}
	witnessPublic, err := witnessFull.Public()
	if err != nil {
		return false, fmt.Errorf("failed to extract public witness: %w", err)
	}

	err = groth16.Verify(proofObj, v.vk, witnessPublic)
	if err != nil {
		return false, fmt.Errorf("ZK proof verification failed: %w", err)
	}
	return true, nil
}
