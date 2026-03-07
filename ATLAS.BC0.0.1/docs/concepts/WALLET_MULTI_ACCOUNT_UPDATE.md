# Wallet Multi-Account Support Update

## Overview
This update adds multi-account (multi-wallet) support to `wallet.html`, allowing users to manage multiple wallets in the browser, switch between them, and perform all wallet actions (send, receive, export, import, delete) for any selected account.

## Key Changes

### 1. Multi-Account Storage
- Wallets are now stored as an array in localStorage under the key `walletAccounts`.
- Each wallet/account is an object: `{ name, privJwk, pubJwk }`.
- The selected account index is stored as `selectedAccountIndex` in localStorage.

### 2. UI Enhancements
- Added a dropdown (`<select>`) for account selection.
- Added buttons for:
  - Creating a new account
  - Exporting the selected account
  - Importing a new account
  - Deleting the selected account
  - Showing the QR code for the selected account
- The wallet overview, balance, and address now reflect the selected account.

### 3. Logic Changes
- All wallet actions (signing, sending, faucet, etc.) use the selected account's keypair.
- Account creation generates a new ECDSA keypair and adds it to the accounts array.
- Account import requires both private and public key JWKs (due to WebCrypto limitations).
- Account deletion is supported (at least one account must remain).
- Account export downloads the private key JWK as a JSON file.
- Account selection updates the UI and active wallet.

### 4. Validation and Security
- Keypairs are validated on import.
- All signing and transaction actions use the correct keypair for the selected account.
- The UI prevents deletion if only one account remains.

### 5. Compatibility
- Old single-wallet storage (`walletPriv`, `walletPub`) is no longer used.
- Existing users will have their first account auto-created on first load if no accounts exist.

## Next Steps / Recommendations
- (Optional) Add account renaming and labeling features.
- (Optional) Fetch and display real transaction history for each account.
- (Optional) Add password or mnemonic protection for accounts.

---

**This update enables robust multi-account management for blockchain wallet users in the browser.** 