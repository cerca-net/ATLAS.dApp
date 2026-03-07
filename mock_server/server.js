const express = require('express');
const app = express();
const port = 8080;

app.use(express.json());

// Mock database
let walletData = {
  'test_address': {
    address: 'test_address',
    balance: 1000,
    nonce: 0,
    transactions: [],
  },
};

// HELPER: Generate mock session token
const generateSessionToken = (address) => `ff_${address.substring(0, 8)}_${Date.now()}`;

// ----- Standard Endpoints -----

app.get('/status', (req, res) => {
  res.json({
    blockHeight: 123,
    txPoolSize: 5,
    isValidator: false,
    mode: 'observer'
  });
});

app.get('/balance', (req, res) => {
  const address = req.query.address;
  const wallet = walletData[address] || { address, balance: 0 };
  res.json({ address, balance: wallet.balance });
});

app.get('/nonce', (req, res) => {
  const address = req.query.address;
  const wallet = walletData[address] || { address, nonce: 0 };
  res.json({ address, nonce: wallet.nonce });
});

app.post('/faucet', (req, res) => {
  const address = req.body.address;
  if (!walletData[address]) {
    walletData[address] = { address, balance: 0, nonce: 0, transactions: [] };
  }
  walletData[address].balance += 1000;
  res.json({ success: true, message: '1000 tokens credited', new_balance: walletData[address].balance });
});

// ----- FlutterFlow Integration Endpoints -----

app.post('/flutterflow/connect-wallet', (req, res) => {
  const { action, address, privateKey } = req.body;
  let targetAddress = address || '0x' + Math.random().toString(16).slice(2, 42);
  
  if (!walletData[targetAddress]) {
    walletData[targetAddress] = { address: targetAddress, balance: 1000, nonce: 0, transactions: [] };
  }

  res.json({
    success: true,
    data: {
      address: targetAddress,
      sessionToken: generateSessionToken(targetAddress),
      balance: walletData[targetAddress].balance,
      isValidator: false
    }
  });
});

app.get('/flutterflow/wallet-info', (req, res) => {
  const address = req.query.address;
  const wallet = walletData[address];
  if (wallet) {
    res.json({
      success: true,
      data: {
        address: wallet.address,
        balance: wallet.balance,
        nonce: wallet.nonce,
        recentTransactions: wallet.transactions.slice(-5)
      }
    });
  } else {
    res.status(404).json({ success: false, message: 'Wallet not found' });
  }
});

app.post('/flutterflow/send-transaction', (req, res) => {
  const { from, to, amount, signature } = req.body;
  
  if (!walletData[from] || walletData[from].balance < amount) {
    return res.status(400).json({ success: false, message: 'Insufficient funds or unknown sender' });
  }

  if (!walletData[to]) {
    walletData[to] = { address: to, balance: 0, nonce: 0, transactions: [] };
  }

  // Process transaction
  walletData[from].balance -= amount;
  walletData[to].balance += amount;
  walletData[from].nonce += 1;

  const tx = {
    hash: '0x' + Math.random().toString(16).slice(2),
    sender: from,
    recipient: to,
    amount,
    timestamp: Math.floor(Date.now() / 1000)
  };

  walletData[from].transactions.push(tx);
  walletData[to].transactions.push(tx);

  res.json({
    success: true,
    message: 'Transaction submitted successfully',
    data: { transactionHash: tx.hash }
  });
});

app.get('/flutterflow/transaction-history', (req, res) => {
  const address = req.query.address;
  const wallet = walletData[address];
  res.json({
    success: true,
    data: {
      address,
      transactions: wallet ? wallet.transactions : []
    }
  });
});

app.listen(port, () => {
  console.log(`Mock server listening at http://localhost:${port}`);
});
