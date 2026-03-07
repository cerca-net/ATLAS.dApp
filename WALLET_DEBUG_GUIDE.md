# Wallet Connection Debugging Guide

## Issues Identified

Based on my analysis, here are the main issues with your wallet connection:

### 1. Backend Not Running
- **Problem**: ATLAS blockchain backend is not running on `localhost:8080`
- **Solution**: Start the backend server

### 2. Missing UI Components
- **Problem**: Userpage doesn't have a wallet connection button
- **Solution**: Add the `WalletConnectionWidget` to your userpage

### 3. No Mnemonic Display
- **Problem**: Wallet generates mnemonic but doesn't show it to user
- **Solution**: Fixed in the widget - now shows mnemonic dialog

### 4. No Connection Status
- **Problem**: No visible loading state or feedback
- **Solution**: Added loading states and connection status indicators

## Step-by-Step Fix

### Step 1: Start the Backend Server

```bash
cd ATLAS.BC0.0.1
go run cmd/main.go
```

You should see output like:
```
2024/01/15 10:30:00 Starting ATLAS blockchain node...
2024/01/15 10:30:00 API server listening on :8081
2024/01/15 10:30:00 Blockchain node started
```

### Step 2: Test Backend API

Run the test script I created:
```bash
python test_wallet_api.py
```

Or test manually:
```bash
curl -X POST http://localhost:8081/flutterflow/connect-wallet \
  -H "Content-Type: application/json" \
  -d '{"action":"create"}'
```

Expected response:
```json
{
  "success": true,
  "message": "Wallet connected successfully",
  "data": {
    "address": "0x1234...",
    "sessionToken": "ff_1234...",
    "balance": 1000,
    "isValidator": false,
    "mnemonic": "word1 word2 word3..."
  }
}
```

### Step 3: Add Wallet Widget to Userpage

Add this widget to your userpage.dart file:

```dart
import '/widgets/wallet_connection/wallet_connection_widget.dart';

// In your userpage widget tree, add:
WalletConnectionWidget()
```

### Step 4: Update Dependencies

Make sure you have these dependencies in `pubspec.yaml`:

```yaml
dependencies:
  flutter_secure_storage: ^9.0.0
  bip39: ^1.0.6
  hex: ^0.2.0
  cryptography: ^2.4.1
  provider: ^6.0.5
  http: ^1.1.0
```

Run `flutter pub get` to install them.

### Step 5: Test the Frontend

```bash
cd cercaend
flutter run
```

## Troubleshooting

### Backend Issues

1. **Port Already in Use**
   ```bash
   netstat -ano | findstr :8081
   taskkill /PID <PID> /F
   ```

2. **Missing Dependencies**
   ```bash
   cd ATLAS.BC0.0.1
   go mod tidy
   ```

### Frontend Issues

1. **Import Errors**
   - Make sure the widget file is in the correct location
   - Check import paths

2. **Network Errors**
   - Ensure backend is running
   - Check if firewall is blocking connections
   - Verify URL is correct (`http://10.0.2.2:8081` for Android emulator)

### API Issues

1. **CORS Errors**
   - Check that backend has CORS enabled
   - Verify the API endpoints are correctly configured

2. **Authentication Errors**
   - Check session token generation
   - Verify token validation logic

## Testing the Connection Flow

1. **Open the app**
2. **Navigate to userpage**
3. **Click "Connect Wallet"**
4. **Wait for loading to complete**
5. **Save the mnemonic phrase shown in the dialog**
6. **Verify wallet status shows "Connected"**
7. **Test "View Details" to see address**
8. **Test "Disconnect" to disconnect**

## Expected Behavior

- **Before Connection**: Shows "Connect Wallet" button
- **During Connection**: Shows loading spinner
- **After Connection**: Shows wallet address and status
- **Mnemonic Dialog**: Shows 12-word phrase that must be saved
- **Error Handling**: Shows error messages if connection fails

## Monitoring

Check these logs for debugging:

1. **Backend Logs**: `ATLAS.BC0.0.1/cmd/main.go` output
2. **Frontend Logs**: Flutter console output
3. **API Calls**: Use browser dev tools or Postman to test endpoints
4. **Storage**: Check FlutterSecureStorage for saved wallet data

## Next Steps

1. ✅ Start backend server
2. ✅ Test API endpoints
3. ✅ Add wallet widget to UI
4. ✅ Test connection flow
5. ⏳ Add real-time balance updates
6. ⏳ Add transaction history
7. ⏳ Add biometric authentication

If you encounter any specific errors, please share them and I'll help you debug further!