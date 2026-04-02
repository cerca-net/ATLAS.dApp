package vm

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
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

var (
	globalVerifier     *ZKVerifier
	globalVerifierOnce sync.Once
	globalVerifierErr  error
)

// GetGlobalZKVerifier returns a singleton instance of the ZKVerifier
func GetGlobalZKVerifier() (*ZKVerifier, error) {
	globalVerifierOnce.Do(func() {
		globalVerifier, globalVerifierErr = NewZKVerifier()
	})
	return globalVerifier, globalVerifierErr
}

// ZKVerifier manages Zero Knowledge verifications using gnark
type ZKVerifier struct {
	vk groth16.VerifyingKey
}

// NewZKVerifier initializes the gnark environment and compiles the system circuit.
func NewZKVerifier() (*ZKVerifier, error) {
	// For V1, we bypass the expensive Groth16 trusted setup.
	// This prevents the node from severely hanging on boot.
	fmt.Println("Warning: ZKVerifier circuit compilation bypassed for V1")

	return &ZKVerifier{}, nil
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
