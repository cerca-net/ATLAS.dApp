package main

import (
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"

	"atlas-blockchain/pkg/wallet"
)

func main() {
	// 1. Treasury
	treasuryMnemonic := "canyon vision beer orange notice wrong savage coin fashion roam ranch weasel"
	treasuryWallet, err := wallet.NewWalletFromMnemonic(treasuryMnemonic)
	if err != nil {
		log.Fatalf("Treasury error: %v", err)
	}
	treasuryAddress := wallet.PublicKeyToAddress(treasuryWallet.PublicKey)

	// 2. Initial Validator
	// Let's use a fixed mnemonic for the initial bootstrap node as well so it's deterministic for testing/devnet
	validatorMnemonic := "apple banana cherry date elderberry fig grape honeydew kiwi lemon mango nectarine orange papaya quince raspberry strawberry tangerine umbrella violet watermelon xylophone yam zebra"
	validatorWallet, err := wallet.NewWalletFromMnemonic(validatorMnemonic)
	if err != nil {
		// Fallback to random if 24 words is needed
		validatorWallet, _, err = wallet.NewWalletWithMnemonic()
		if err != nil {
			log.Fatalf("Validator error: %v", err)
		}
	}
	validatorAddress := wallet.PublicKeyToAddress(validatorWallet.PublicKey)
	validatorPubKey := validatorWallet.PublicKeyStr()

	privKeyBytes, err := x509.MarshalECPrivateKey(validatorWallet.PrivateKey)
	if err != nil {
		log.Fatalf("Failed to marshal generic private key: %v", err)
	}
	hexKey := hex.EncodeToString(privKeyBytes)
	err = ioutil.WriteFile("genesis_validator.key", []byte(hexKey), 0600)
	if err != nil {
		log.Fatalf("Failed to write genesis.key: %v", err)
	}

	fmt.Printf("TREASURY_ADDRESS=%s\n", treasuryAddress)
	fmt.Printf("VALIDATOR_ADDRESS=%s\n", validatorAddress)
	fmt.Printf("VALIDATOR_PUBKEY=%s\n", validatorPubKey)
	fmt.Println("Initial validator hex key saved to genesis_validator.key")
}
