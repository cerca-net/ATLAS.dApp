// Common JS for wallet, API, and utility functions

window.wallet = null;
// Get the selected port from localStorage or default to 8080
window.currentApiPort = parseInt(localStorage.getItem('selectedNodePort')) || 8080;

// --- WALLET MANAGEMENT ---
async function generateWallet() {
    try {
        const keyPair = await window.crypto.subtle.generateKey({
            name: 'ECDSA',
            namedCurve: 'P-256',
        }, true, ['sign', 'verify']);
        window.wallet = keyPair;
        await saveWallet();
        const address = await exportPublicKey(keyPair.publicKey);
        updateWalletDisplay(address);
        showStatus('New wallet generated successfully!', 'success');
        return address;
    } catch (error) {
        showStatus('Failed to generate wallet: ' + error.message, 'error');
    }
}

async function saveWallet() {
    if (!window.wallet) return;
    try {
        const privJwk = await window.crypto.subtle.exportKey('jwk', window.wallet.privateKey);
        const pubJwk = await window.crypto.subtle.exportKey('jwk', window.wallet.publicKey);
        localStorage.setItem('walletPriv', JSON.stringify(privJwk));
        localStorage.setItem('walletPub', JSON.stringify(pubJwk));
    } catch (error) {
        console.error('Failed to save wallet:', error);
    }
}

async function loadWallet() {
    try {
        const privJwk = localStorage.getItem('walletPriv');
        const pubJwk = localStorage.getItem('walletPub');
        if (!privJwk || !pubJwk) return null;
        const privateKey = await window.crypto.subtle.importKey('jwk', JSON.parse(privJwk), { name: 'ECDSA', namedCurve: 'P-256' }, true, ['sign']);
        const publicKey = await window.crypto.subtle.importKey('jwk', JSON.parse(pubJwk), { name: 'ECDSA', namedCurve: 'P-256' }, true, ['verify']);
        window.wallet = { privateKey, publicKey };
        return window.wallet;
    } catch (error) {
        console.error('Failed to load wallet:', error);
        return null;
    }
}

function exportPublicKey(pubKey) {
    return window.crypto.subtle.exportKey('spki', pubKey).then(buf => {
        return arrayBufferToHex(buf);
    });
}

function arrayBufferToHex(buffer) {
    return Array.from(new Uint8Array(buffer)).map(b => b.toString(16).padStart(2, '0')).join('');
}

async function getWalletAddress() {
    if (!window.wallet || !window.wallet.publicKey) return '(not loaded)';
    // Get the public key as SPKI format
    const pubKey = await window.crypto.subtle.exportKey('spki', window.wallet.publicKey);
    // Create a SHA-256 hash of the public key
    const hashBuffer = await crypto.subtle.digest('SHA-256', pubKey);
    // Take the last 20 bytes of the hash
    const hashArray = new Uint8Array(hashBuffer);
    const addressBytes = hashArray.slice(-20);
    // Convert to hex and prefix with 0x
    const hexString = Array.from(addressBytes).map(b => b.toString(16).padStart(2, '0')).join('');
    return '0x' + hexString;
}

function derToRS(der) {
    const bytes = new Uint8Array(der);
    if (bytes[0] !== 0x30) throw new Error('Invalid DER sequence');
    let offset = 2;
    if (bytes[1] & 0x80) {
        const n = bytes[1] & 0x7f;
        offset = 2 + n;
    }
    if (bytes[offset] !== 0x02) throw new Error('Invalid DER integer for r');
    let rLen = bytes[offset + 1];
    let r = bytes.slice(offset + 2, offset + 2 + rLen);
    offset = offset + 2 + rLen;
    if (bytes[offset] !== 0x02) throw new Error('Invalid DER integer for s');
    let sLen = bytes[offset + 1];
    let s = bytes.slice(offset + 2, offset + 2 + sLen);
    while (r.length > 32) r = r.slice(1);
    while (s.length > 32) s = s.slice(1);
    let rPad = new Uint8Array(32); rPad.set(r, 32 - r.length);
    let sPad = new Uint8Array(32); sPad.set(s, 32 - s.length);
    let rs = new Uint8Array(64); rs.set(rPad, 0); rs.set(sPad, 32);
    return Array.from(rs).map(b => b.toString(16).padStart(2, '0')).join('');
}

async function signTransaction(tx) {
    if (!window.wallet || !window.wallet.privateKey) {
        throw new Error('Wallet or private key not loaded. Please generate or load your wallet.');
    }
    if (typeof window.wallet.privateKey !== 'object') {
        throw new Error('Invalid private key object.');
    }
    const enc = new TextEncoder();
    const data = enc.encode(tx.Sender + tx.Recipient + tx.Amount + (tx.Fee || 0) + tx.Timestamp + tx.Nonce + (tx.Data || ''));
    const sig = await window.crypto.subtle.sign({ name: 'ECDSA', hash: 'SHA-256' }, window.wallet.privateKey, data);
    const signature = derToRS(sig);
    console.log('Signature (hex):', signature, 'Length:', signature.length);
    return signature;
}

function showStatus(message, type = 'success', target = 'send-status') {
    const statusDiv = document.getElementById(target);
    if (!statusDiv) return;
    statusDiv.textContent = message;
    statusDiv.className = `status-message status-${type}`;
    setTimeout(() => {
        statusDiv.textContent = '';
        statusDiv.className = '';
    }, 5000);
}

async function fetchWalletBalance(address) {
    try {
        const response = await fetch(`http://localhost:${window.currentApiPort}/balance?address=${address}`);
        const data = await response.json();
        const el = document.getElementById('wallet-balance');
        if (el) el.textContent = data.balance !== undefined ? data.balance : '0';
    } catch (error) {
        const el = document.getElementById('wallet-balance');
        if (el) el.textContent = 'Error';
    }
}

async function fetchRecentTransactions() {
    try {
        const response = await fetch(`http://localhost:${window.currentApiPort}/mempool?limit=5`);
        const transactions = await response.json();
        const container = document.getElementById('recent-transactions');
        if (transactions.length === 0) {
            container.innerHTML = '<div class="loading">No recent transactions</div>';
            return;
        }
        container.innerHTML = '';
        transactions.forEach(tx => {
            const txDiv = document.createElement('div');
            txDiv.style.cssText = `background: #f7fafc; padding: 15px; margin-bottom: 10px; border-radius: 8px; border-left: 4px solid #4299e1;`;
            txDiv.innerHTML = `<div style="font-weight: 600; margin-bottom: 5px;">${tx.Amount} tokens</div><div style="font-size: 0.9rem; color: #718096;">From: ${tx.Sender.substring(0, 16)}...<br>To: ${tx.Recipient.substring(0, 16)}...</div>`;
            container.appendChild(txDiv);
        });
    } catch (error) {
        console.error('Failed to fetch transactions:', error);
    }
}

async function fetchNetworkData() {
    try {
        const response = await fetch(`http://localhost:${window.currentApiPort}/peers`);
        const peers = await response.json();
        document.getElementById('peer-count').textContent = peers.length;
        const peerList = document.getElementById('peer-list');
        peerList.innerHTML = '';
        if (peers.length === 0) {
            peerList.innerHTML = '<li>No peers connected</li>';
        } else {
            peers.forEach(peer => {
                const li = document.createElement('li');
                li.textContent = peer;
                peerList.appendChild(li);
            });
        }
    } catch (error) {
        console.error('Failed to fetch network data:', error);
        document.getElementById('peer-count').textContent = 'Error';
    }
}

async function fetchBlockchainData() {
    try {
        const blocksResponse = await fetch(`http://localhost:${window.currentApiPort}/blocks?limit=10`);
        const blocks = await blocksResponse.json();
        const tbody = document.getElementById('blockchain-table');
        if (tbody) {
            tbody.innerHTML = '';
            if (blocks.length === 0) {
                tbody.innerHTML = '<tr><td colspan="4" class="loading">No blocks found</td></tr>';
            } else {
                blocks.forEach(block => {
                    const row = document.createElement('tr');
                    const timestamp = new Date(block.Timestamp * 1000).toLocaleString();
                    row.innerHTML = `<td><strong>${block.Index}</strong></td><td class="hash-cell" title="${block.Hash}">${block.Hash.substring(0, 16)}...</td><td>${block.Transactions ? block.Transactions.length : 0}</td><td>${timestamp}</td>`;
                    tbody.appendChild(row);
                });
            }
            document.getElementById('block-count').textContent = blocks.length;
            const totalTx = blocks.reduce((sum, block) => sum + (block.Transactions ? block.Transactions.length : 0), 0);
            document.getElementById('tx-count').textContent = totalTx;
        }
    } catch (error) {
        console.error('Failed to fetch blockchain data:', error);
    }
}

async function fetchValidatorInfo() {
    try {
        const address = await getWalletAddress();
        const response = await fetch(`http://localhost:${window.currentApiPort}/validator?address=${address}`);
        const data = await response.json();
        document.getElementById('stake-amount').textContent = data.stake || '0';
        document.getElementById('validator-rank').textContent = data.rank || 'Not a validator';
    } catch (error) {
        document.getElementById('stake-amount').textContent = 'Error';
        document.getElementById('validator-rank').textContent = 'Error';
    }
}

// Export functions for use in other scripts
window.generateWallet = generateWallet;
window.saveWallet = saveWallet;
window.loadWallet = loadWallet;
window.exportPublicKey = exportPublicKey;
window.arrayBufferToHex = arrayBufferToHex;
window.getWalletAddress = getWalletAddress;
window.derToRS = derToRS;
window.signTransaction = signTransaction;
window.showStatus = showStatus;
window.fetchWalletBalance = fetchWalletBalance;
window.fetchRecentTransactions = fetchRecentTransactions;
window.fetchNetworkData = fetchNetworkData;
window.fetchBlockchainData = fetchBlockchainData;
window.fetchValidatorInfo = fetchValidatorInfo; 