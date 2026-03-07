#!/usr/bin/env python3
import requests
import json
import time

# Test script for wallet connection API
BASE_URL = "http://localhost:8080"


def test_backend_connection():
    """Test if the backend server is running"""
    try:
        response = requests.get(f"{BASE_URL}/status", timeout=5)
        if response.status_code == 200:
            print("✅ Backend server is running")
            print(f"   Status: {response.json()}")
            return True
        else:
            print(f"❌ Backend returned status: {response.status_code}")
            return False
    except requests.exceptions.ConnectionError:
        print("❌ Backend server is not running")
        print(
            "   Please start the backend with: cd ATLAS.BC0.0.1 && go run cmd/main.go"
        )
        return False
    except Exception as e:
        print(f"❌ Error connecting to backend: {e}")
        return False


def test_wallet_creation():
    """Test wallet creation endpoint"""
    print("\n🔧 Testing wallet creation...")

    payload = {"action": "create"}

    try:
        response = requests.post(
            f"{BASE_URL}/flutterflow/connect-wallet", json=payload, timeout=10
        )

        print(f"   Status Code: {response.status_code}")
        print(f"   Response: {response.text}")

        if response.status_code == 200:
            result = response.json()
            if result.get("success"):
                print("✅ Wallet creation successful")
                return result.get("data", {})
            else:
                print(f"❌ Wallet creation failed: {result.get('message')}")
                return None
        else:
            print(f"❌ API request failed with status {response.status_code}")
            return None

    except Exception as e:
        print(f"❌ Error testing wallet creation: {e}")
        return None


def test_wallet_info(address):
    """Test wallet info endpoint"""
    print(f"\n🔧 Testing wallet info for address: {address}")

    try:
        response = requests.get(
            f"{BASE_URL}/flutterflow/wallet-info?address={address}", timeout=10
        )

        print(f"   Status Code: {response.status_code}")
        print(f"   Response: {response.text}")

        if response.status_code == 200:
            result = response.json()
            if result.get("success"):
                print("✅ Wallet info retrieval successful")
                return result.get("data", {})
            else:
                print(f"❌ Wallet info retrieval failed: {result.get('message')}")
                return None
        else:
            print(f"❌ API request failed with status {response.status_code}")
            return None

    except Exception as e:
        print(f"❌ Error testing wallet info: {e}")
        return None


def main():
    print("🚀 Starting Wallet Connection Tests")
    print("=" * 50)

    # Test 1: Backend connection
    if not test_backend_connection():
        return

    # Test 2: Wallet creation
    wallet_data = test_wallet_creation()
    if wallet_data:
        address = wallet_data.get("address")
        print(f"   Generated Address: {address}")

        # Test 3: Wallet info
        info_data = test_wallet_info(address)
        if info_data:
            balance = info_data.get("balance", 0)
            print(f"   Wallet Balance: {balance}")

    print("\n" + "=" * 50)
    print("🎯 Test Summary")
    print("If all tests passed, your backend API is working correctly.")
    print("If tests failed, check the backend logs and ensure it's running.")


if __name__ == "__main__":
    main()
