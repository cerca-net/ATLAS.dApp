// Blockchain Dashboard - Enhanced Version
let refreshInterval;

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
        
        const privateKey = await window.crypto.subtle.importKey('jwk', JSON.parse(privJwk), 
            { name: 'ECDSA', namedCurve: 'P-256' }, true, ['sign']);
        const publicKey = await window.crypto.subtle.importKey('jwk', JSON.parse(pubJwk), 
            { name: 'ECDSA', namedCurve: 'P-256' }, true, ['verify']);
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
    return await exportPublicKey(window.wallet.publicKey);
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

// --- UI UPDATES ---
function updateWalletDisplay(address) {
    const addressEl = document.getElementById('wallet-address');
    if (addressEl) {
        addressEl.innerHTML = `<span id='wallet-address-text'>${address}</span> <button id='copy-address-btn' style='margin-left:8px;padding:2px 8px;font-size:0.8rem;'>Copy</button>`;
        const copyBtn = document.getElementById('copy-address-btn');
        if (copyBtn) {
            copyBtn.onclick = function() {
                navigator.clipboard.writeText(address).then(() => {
                    showStatus('Address copied to clipboard!', 'success');
                });
            };
        }
    }
    fetchWalletBalance(address);
}

function showStatus(message, type = 'success') {
    const statusDiv = document.getElementById('send-status');
    if (!statusDiv) {
        // Silently return if status element doesn't exist (not all pages have this)
        return;
    }
    statusDiv.textContent = message;
    statusDiv.className = `status-message status-${type}`;
    setTimeout(() => {
        if (statusDiv) {
            statusDiv.textContent = '';
            statusDiv.className = '';
        }
    }, 5000);
}

async function fetchWalletBalance(address) {
    try {
        const response = await fetch(`http://localhost:${window.currentApiPort}/balance?address=${address}`);
        const data = await response.json();
        const el = document.getElementById('wallet-balance');
        if (el) {
            el.textContent = data.balance !== undefined ? data.balance : '0';
        }
    } catch (error) {
        console.error('Failed to fetch wallet balance:', error);
        const el = document.getElementById('wallet-balance');
        if (el) {
            el.textContent = 'Error';
        }
    }
}

async function fetchBlockchainData() {
    try {
        // Fetch blocks
        const blocksResponse = await fetch(`http://localhost:${window.currentApiPort}/blocks?limit=10`);
        const blocks = await blocksResponse.json();
        
        const tbody = document.getElementById('blockchain-table');
        if (!tbody) {
            console.warn('Blockchain table not found');
            return;
        }
        
        tbody.innerHTML = '';
        
        if (blocks.length === 0) {
            tbody.innerHTML = '<tr><td colspan="4" class="loading">No blocks found</td></tr>';
        } else {
            blocks.forEach(block => {
                const row = document.createElement('tr');
                const timestamp = new Date(block.Timestamp * 1000).toLocaleString();
                row.innerHTML = `
                    <td><strong>${block.Index}</strong></td>
                    <td class="hash-cell" title="${block.Hash}">${block.Hash.substring(0, 16)}...</td>
                    <td>${block.Transactions ? block.Transactions.length : 0}</td>
                    <td>${timestamp}</td>
                `;
                tbody.appendChild(row);
            });
        }
        
        // Update stats
        const blockCountEl = document.getElementById('block-count');
        const txCountEl = document.getElementById('tx-count');
        
        if (blockCountEl) blockCountEl.textContent = blocks.length;
        
        // Count total transactions
        const totalTx = blocks.reduce((sum, block) => sum + (block.Transactions ? block.Transactions.length : 0), 0);
        if (txCountEl) txCountEl.textContent = totalTx;
        
    } catch (error) {
        console.error('Failed to fetch blockchain data:', error);
    }
}

async function fetchNetworkData() {
    try {
        const response = await fetch(`http://localhost:${window.currentApiPort}/peers`);
        const peers = await response.json();
        
        const peerCountEl = document.getElementById('peer-count');
        const peerList = document.getElementById('peer-list');
        
        if (peerCountEl) peerCountEl.textContent = peers.length;
        
        if (!peerList) {
            console.warn('Peer list not found');
            return;
        }
        
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
        const peerCountEl = document.getElementById('peer-count');
        if (peerCountEl) peerCountEl.textContent = 'Error';
    }
}

async function fetchRecentTransactions() {
    const container = document.getElementById('recent-transactions');
    if (!container) {
        // Silently return if container doesn't exist (not all pages have this)
        return;
    }
    try {
        const response = await fetch(`http://localhost:${window.currentApiPort}/mempool?limit=5`);
        const transactions = await response.json();
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
        if (container) {
            container.innerHTML = '<div class="loading">Error loading transactions</div>';
        }
    }
}

async function fetchValidatorInfo() {
    try {
        const address = await getWalletAddress();
        const response = await fetch(`http://localhost:${window.currentApiPort}/validator?address=${address}`);
        const data = await response.json();
        
        const stakeEl = document.getElementById('stake-amount');
        const rankEl = document.getElementById('validator-rank');
        
        if (stakeEl) stakeEl.textContent = data.stake || '0';
        if (rankEl) rankEl.textContent = data.rank || 'Not a validator';
    } catch (error) {
        console.error('Failed to fetch validator info:', error);
        const stakeEl = document.getElementById('stake-amount');
        const rankEl = document.getElementById('validator-rank');
        
        if (stakeEl) stakeEl.textContent = 'Error';
        if (rankEl) rankEl.textContent = 'Error';
    }
}

// --- EVENT HANDLERS ---
document.addEventListener('DOMContentLoaded', async () => {
    // Initialize wallet
    let priv = localStorage.getItem('walletPriv');
    let pub = localStorage.getItem('walletPub');
    if (priv && pub) {
        await loadWallet();
    } else {
        await generateWallet();
    }
    const address = await getWalletAddress();
    updateWalletDisplay(address);
    // Initial data fetch
    await Promise.all([
        fetchWalletBalance(address),
        fetchRecentTransactions(),
        fetchNetworkArchitecture(),
        updateNodeStatus()
    ]);
    // Set up auto-refresh
    refreshInterval = setInterval(async () => {
        await Promise.all([
            fetchWalletBalance(address),
            fetchRecentTransactions(),
            fetchNetworkArchitecture()
        ]);
    }, 10000); // Refresh every 10 seconds
    // Event listeners - with null checks
    const faucetBtn = document.getElementById('faucet-btn');
    if (faucetBtn) {
        faucetBtn.addEventListener('click', async () => {
            const addr = await getWalletAddress();
            try {
                const response = await fetch(`http://localhost:${window.currentApiPort}/faucet`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ address: addr })
                });
                const resp = await response.json();
                showStatus(resp.status || 'Faucet tokens credited!', 'success');
                setTimeout(() => fetchWalletBalance(addr), 2000);
            } catch (error) {
                showStatus('Failed to get faucet tokens', 'error');
            }
        });
    }
    
    const generateWalletBtn = document.getElementById('generate-wallet-btn');
    if (generateWalletBtn) {
        generateWalletBtn.addEventListener('click', generateWallet);
    }
    
    const exportWalletBtn = document.getElementById('export-wallet-btn');
    if (exportWalletBtn) {
        exportWalletBtn.addEventListener('click', () => {
            if (window.wallet) {
                const dataStr = JSON.stringify({
                    private: localStorage.getItem('walletPriv'),
                    public: localStorage.getItem('walletPub')
                });
                const dataBlob = new Blob([dataStr], {type: 'application/json'});
                const url = URL.createObjectURL(dataBlob);
                const link = document.createElement('a');
                link.href = url;
                link.download = 'wallet.json';
                link.click();
                URL.revokeObjectURL(url);
                showStatus('Wallet exported successfully!', 'success');
            }
        });
    }
    
    const importWalletBtn = document.getElementById('import-wallet-btn');
    const importWalletFile = document.getElementById('import-wallet-file');
    if (importWalletBtn && importWalletFile) {
        importWalletBtn.addEventListener('click', () => {
            importWalletFile.click();
        });
        importWalletFile.addEventListener('change', async (e) => {
            const file = e.target.files[0];
            if (file) {
                try {
                    const text = await file.text();
                    const data = JSON.parse(text);
                    localStorage.setItem('walletPriv', data.private);
                    localStorage.setItem('walletPub', data.public);
                    await loadWallet();
                    const address = await getWalletAddress();
                    updateWalletDisplay(address);
                    showStatus('Wallet imported successfully!', 'success');
                } catch (error) {
                    showStatus('Failed to import wallet', 'error');
                }
            }
        });
    }
    
    const sendForm = document.getElementById('send-form');
    if (sendForm) {
        sendForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const to = document.getElementById('to-address').value;
            const amount = parseInt(document.getElementById('amount').value, 10);
            const address = await getWalletAddress();
            const tx = {
                Sender: address,
                Recipient: to,
                Amount: amount,
                Timestamp: Math.floor(Date.now() / 1000),
            };
            try {
                tx.Signature = await signTransaction(tx);
                const response = await fetch(`http://localhost:${window.currentApiPort}/submit-transaction`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(tx)
                });
                const resp = await response.json();
                showStatus(resp.status || 'Transaction submitted successfully!', 'success');
                e.target.reset();
            } catch (error) {
                showStatus('Transaction failed: ' + error.message, 'error');
            }
        });
    }
});

window.addEventListener('beforeunload', () => {
    if (refreshInterval) {
        clearInterval(refreshInterval);
    }
}); 

async function fetchNetworkArchitecture() {
    try {
        console.log(`Fetching network architecture from http://localhost:${window.currentApiPort}/network/architecture`);
        const response = await fetch(`http://localhost:${window.currentApiPort}/network/architecture`);
        console.log('Response status:', response.status);
        if (!response.ok) throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        const architecture = await response.json();
        console.log('Network architecture data:', architecture);
        
        // Update Node Types card
        const nodeTypesCard = document.querySelector('[data-architecture="nodeTypes"]');
        if (nodeTypesCard) {
            const validators = architecture.nodeTypes.validators;
            const observers = architecture.nodeTypes.observers;
            const fullNodes = architecture.nodeTypes.fullNodes;
            
            nodeTypesCard.innerHTML = `
                <h4 style="color:#1e3c72; margin-bottom:8px;">Node Types</h4>
                <div style="font-size:0.97rem; color:#444;">
                    <div style="margin-bottom:8px;"><strong>Validators:</strong> ${validators.count} active (${validators.active} online)</div>
                    <div style="margin-bottom:8px;"><strong>Observers:</strong> ${observers.count} connected</div>
                    <div style="margin-bottom:8px;"><strong>Full Nodes:</strong> ${fullNodes.count} total</div>
                    <div style="margin-top:12px; font-size:0.9rem; color:#666;">${validators.description}</div>
                </div>
            `;
        }
        
        // Update P2P Protocol card
        const p2pCard = document.querySelector('[data-architecture="p2pProtocol"]');
        if (p2pCard) {
            const p2p = architecture.p2pProtocol;
            p2pCard.innerHTML = `
                <h4 style="color:#1e3c72; margin-bottom:8px;">P2P Protocol</h4>
                <div style="font-size:0.97rem; color:#444;">
                    <div style="margin-bottom:8px;"><strong>Type:</strong> ${p2p.type} v${p2p.version}</div>
                    <div style="margin-bottom:8px;"><strong>Discovery:</strong> ${p2p.discovery}</div>
                    <div style="margin-bottom:8px;"><strong>Transport:</strong> ${p2p.transport}</div>
                    <div style="margin-top:12px; font-size:0.9rem; color:#666;">${p2p.description}</div>
                </div>
            `;
        }
        
        // Update Consensus Mechanism card
        const consensusCard = document.querySelector('[data-architecture="consensusMechanism"]');
        if (consensusCard) {
            const consensus = architecture.consensusMechanism;
            consensusCard.innerHTML = `
                <h4 style="color:#1e3c72; margin-bottom:8px;">Consensus Mechanism</h4>
                <div style="font-size:0.97rem; color:#444;">
                    <div style="margin-bottom:8px;"><strong>Type:</strong> ${consensus.type}</div>
                    <div style="margin-bottom:8px;"><strong>Block Time:</strong> ${consensus.blockTime}</div>
                    <div style="margin-bottom:8px;"><strong>Finality:</strong> ${consensus.finality}</div>
                    <div style="margin-top:12px; font-size:0.9rem; color:#666;">${consensus.description}</div>
                </div>
            `;
        }
        
        // Update Network Topology card
        const topologyCard = document.querySelector('[data-architecture="networkTopology"]');
        if (topologyCard) {
            const topology = architecture.networkTopology;
            topologyCard.innerHTML = `
                <h4 style="color:#1e3c72; margin-bottom:8px;">Network Topology</h4>
                <div style="font-size:0.97rem; color:#444;">
                    <div style="margin-bottom:8px;"><strong>Type:</strong> ${topology.type}</div>
                    <div style="margin-bottom:8px;"><strong>Max Peers:</strong> ${topology.maxPeers}</div>
                    <div style="margin-bottom:8px;"><strong>Connections:</strong> ${topology.connections}</div>
                    <div style="margin-top:12px; font-size:0.9rem; color:#666;">${topology.description}</div>
                </div>
            `;
        }
        
        // Update Security Features card
        const securityCard = document.querySelector('[data-architecture="securityFeatures"]');
        if (securityCard) {
            const security = architecture.securityFeatures;
            securityCard.innerHTML = `
                <h4 style="color:#1e3c72; margin-bottom:8px;">Security Features</h4>
                <div style="font-size:0.97rem; color:#444;">
                    <div style="margin-bottom:8px;"><strong>Encryption:</strong> ${security.encryption}</div>
                    <div style="margin-bottom:8px;"><strong>Authentication:</strong> ${security.authentication}</div>
                    <div style="margin-bottom:8px;"><strong>Rate Limiting:</strong> ${security.rateLimiting}</div>
                    <div style="margin-bottom:8px;"><strong>Slashing:</strong> ${security.slashing}</div>
                    <div style="margin-top:12px; font-size:0.9rem; color:#666;">${security.description}</div>
                </div>
            `;
        }
        
    } catch (error) {
        console.error('Error fetching network architecture:', error);
        // Fallback to static content if API fails
        showStatus(`Network architecture error: ${error.message}`, 'error');
        
        // Show error in the cards
        const cards = document.querySelectorAll('[data-architecture]');
        cards.forEach(card => {
            card.innerHTML = `
                <h4 style="color:#1e3c72; margin-bottom:8px;">${card.querySelector('h4').textContent}</h4>
                <div style="font-size:0.97rem; color:#e53e3e;">
                    <strong>Connection Error:</strong> ${error.message}<br>
                    <small>Make sure the blockchain is running on port ${window.currentApiPort}</small>
                </div>
            `;
        });
    }
}

async function updateNodeStatus() {
    try {
        const response = await fetch(`http://localhost:${window.currentApiPort}/status`);
        if (!response.ok) throw new Error('Failed to fetch node status');
        const status = await response.json();
        
        const statusElement = document.getElementById('node-status');
        if (statusElement) {
            statusElement.textContent = `Connected to Node ${window.currentApiPort} (${status.mode})`;
        }
    } catch (error) {
        console.error('Error updating node status:', error);
        const statusElement = document.getElementById('node-status');
        if (statusElement) {
            statusElement.textContent = `Connected to Node ${window.currentApiPort}`;
        }
    }
}