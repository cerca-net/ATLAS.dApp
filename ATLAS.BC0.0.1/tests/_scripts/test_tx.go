package main

import (
	"atlas-blockchain/pkg/transaction"
	"atlas-blockchain/pkg/wallet"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	// Create a new wallet to send from
	w, err := wallet.NewWallet()
	if err != nil {
		panic(err)
	}

	senderAddr := wallet.PublicKeyToAddress(w.PublicKey)
	recipientAddr := "0x1234567890123456789012345678901234567890"

	fmt.Printf("Sender Initialized: %s\n", senderAddr)

	// Get some funds from faucet first
	faucetData := map[string]string{"address": senderAddr}
	faucetJson, _ := json.Marshal(faucetData)
	resp, err := http.Post("http://localhost:8080/faucet", "application/json", bytes.NewBuffer(faucetJson))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("Faucet response status:", resp.Status)

	// Wait for state to update
	time.Sleep(2 * time.Second)

	// Create and sign transaction
	tx := &transaction.Transaction{
		Type:            transaction.TxTypeRegular,
		Sender:          senderAddr,
		SenderPublicKey: w.PublicKeyStr(),
		Recipient:       recipientAddr,
		Amount:          100,
		Fee:             10,
		Timestamp:       time.Now().Unix(),
		Nonce:           0,
		Data:            "Live test transaction",
	}

	if err := w.SignTransaction(tx); err != nil {
		panic(err)
	}

	// Submit transaction
	txJson, _ := json.Marshal(tx)
	resp, err = http.Post("http://localhost:8080/submit-transaction", "application/json", bytes.NewBuffer(txJson))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Submit response status: %s, body: %s\n", resp.Status, string(body))
}
