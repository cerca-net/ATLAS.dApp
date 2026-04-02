#!/bin/bash
# =============================================================================
# ATLAS CercaChain - Frontend Production Build Script
# =============================================================================
# This script builds the Flutter Web App for production deployment.

echo "=================================================="
echo "    Building CercaEnd Frontend (Web) for PROD     "
echo "=================================================="

# Check if flutter is installed
if ! command -v flutter &> /dev/null
then
    echo "⚠️  Flutter SDK could not be found. Please install Flutter and ensure it is in your PATH."
    exit 1
fi

# Define API URL (defaults to production or local Node 1)
API_URL=${1:-"http://localhost:8080"}
echo "[*] Using Blockchain API URL: $API_URL"

# Navigate to the frontend directory
cd cercaend || { echo "Directory 'cercaend' not found. Run from repo root."; exit 1; }

echo "[*] Fetching pub dependencies..."
flutter pub get

echo "[*] Compiling Flutter Web Release..."
flutter build web --release --dart-define=BLOCKCHAIN_API_URL=$API_URL

echo "=================================================="
echo "✅ Build Complete!"
echo "The production-ready assets are located in: cercaend/build/web"
echo "You can host this folder via any static web server (NGINX, Apache, Vercel, S3)."
echo "=================================================="
