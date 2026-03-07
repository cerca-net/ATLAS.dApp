package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	// Test the balance endpoint directly
	resp, err := http.Get("http://localhost:8081/balance?address=0x1234567890abcdef1234567890abcdef12345678")
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Headers: %v\n", resp.Header)
	fmt.Printf("Body: %s\n", string(body))

	// Also test with a POST request to see if that's the issue
	resp2, err := http.Post("http://localhost:8081/balance?address=0x1234567890abcdef1234567890abcdef12345678", "application/json", strings.NewReader(""))
	if err != nil {
		fmt.Printf("Error making POST request: %v\n", err)
		return
	}
	defer resp2.Body.Close()

	body2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		fmt.Printf("Error reading POST response: %v\n", err)
		return
	}

	fmt.Printf("\nPOST Status: %d\n", resp2.StatusCode)
	fmt.Printf("POST Headers: %v\n", resp2.Header)
	fmt.Printf("POST Body: %s\n", string(body2))
}
