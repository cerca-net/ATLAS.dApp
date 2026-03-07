//go:build test
// +build test

package main

import (
	"testing"
)

// TestMain sets up the test environment
func TestMain(m *testing.M) {
	// Set test mode before any tests run
	SetTestMode(true)
	
	// Run the tests
	m.Run()
}

func TestBasic(t *testing.T) {
	t.Log("✅ Basic test works!")
} 
