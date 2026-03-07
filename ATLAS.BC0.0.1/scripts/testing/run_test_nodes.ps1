# Blockchain Test Environment - PowerShell Script
# This script starts multiple blockchain nodes for testing

Write-Host "Starting Blockchain Test Environment..." -ForegroundColor Green

# Kill any existing processes
Write-Host "Cleaning up existing processes..." -ForegroundColor Yellow
Get-Process -Name "blockchain" -ErrorAction SilentlyContinue | Stop-Process -Force
Get-Process -Name "go" -ErrorAction SilentlyContinue | Stop-Process -Force

# Build the blockchain
Write-Host "Building blockchain..." -ForegroundColor Yellow
go build -o blockchain.exe .

# Create test directories
if (!(Test-Path "test_logs")) {
    New-Item -ItemType Directory -Path "test_logs"
}

Write-Host "Starting blockchain nodes..." -ForegroundColor Green

# Start Node 1
Write-Host "Starting Node 1 (Peer: 8001, API: 8081)" -ForegroundColor Green
$args1 = @("-port", "8001", "-api", "8081", "-peers", "5")
Start-Process -FilePath ".\blockchain.exe" -ArgumentList $args1 -RedirectStandardOutput "test_logs\node1.log" -RedirectStandardError "test_logs\node1_error.log" -WindowStyle Hidden

# Start Node 2
Write-Host "Starting Node 2 (Peer: 8002, API: 8082)" -ForegroundColor Green
$args2 = @("-port", "8002", "-api", "8082", "-peers", "5")
Start-Process -FilePath ".\blockchain.exe" -ArgumentList $args2 -RedirectStandardOutput "test_logs\node2.log" -RedirectStandardError "test_logs\node2_error.log" -WindowStyle Hidden

# Start Node 3
Write-Host "Starting Node 3 (Peer: 8003, API: 8083)" -ForegroundColor Green
$args3 = @("-port", "8003", "-api", "8083", "-peers", "5")
Start-Process -FilePath ".\blockchain.exe" -ArgumentList $args3 -RedirectStandardOutput "test_logs\node3.log" -RedirectStandardError "test_logs\node3_error.log" -WindowStyle Hidden

Write-Host "Waiting for nodes to initialize..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

Write-Host "Creating test wallets..." -ForegroundColor Green

# Create test wallets using faucet
try {
    $body1 = '{"address": "test_wallet_1"}'
    $body2 = '{"address": "test_wallet_2"}'
    $body3 = '{"address": "test_wallet_3"}'
    
    Invoke-RestMethod -Uri "http://localhost:8081/faucet" -Method POST -ContentType "application/json" -Body $body1 | Out-Null
    Invoke-RestMethod -Uri "http://localhost:8082/faucet" -Method POST -ContentType "application/json" -Body $body2 | Out-Null
    Invoke-RestMethod -Uri "http://localhost:8083/faucet" -Method POST -ContentType "application/json" -Body $body3 | Out-Null
    Write-Host "Test wallets created successfully" -ForegroundColor Green
} catch {
    Write-Host "Warning: Could not create test wallets. Nodes may still be starting up." -ForegroundColor Yellow
}

Write-Host "Waiting for initial block production..." -ForegroundColor Yellow
Start-Sleep -Seconds 15

Write-Host "Running test transactions..." -ForegroundColor Green

# Send test transactions
try {
    $tx1 = @{
        from = "test_wallet_1"
        to = "test_wallet_2"
        amount = 50
        timestamp = [DateTimeOffset]::Now.ToUnixTimeSeconds()
        fee = 1
    } | ConvertTo-Json

    $tx2 = @{
        from = "test_wallet_2"
        to = "test_wallet_3"
        amount = 30
        timestamp = [DateTimeOffset]::Now.ToUnixTimeSeconds()
        fee = 1
    } | ConvertTo-Json

    $tx3 = @{
        from = "test_wallet_3"
        to = "test_wallet_1"
        amount = 20
        timestamp = [DateTimeOffset]::Now.ToUnixTimeSeconds()
        fee = 1
    } | ConvertTo-Json

    Invoke-RestMethod -Uri "http://localhost:8081/submit-transaction" -Method POST -ContentType "application/json" -Body $tx1 | Out-Null
    Start-Sleep -Seconds 5
    Invoke-RestMethod -Uri "http://localhost:8082/submit-transaction" -Method POST -ContentType "application/json" -Body $tx2 | Out-Null
    Start-Sleep -Seconds 5
    Invoke-RestMethod -Uri "http://localhost:8083/submit-transaction" -Method POST -ContentType "application/json" -Body $tx3 | Out-Null
    
    Write-Host "Test transactions submitted successfully" -ForegroundColor Green
} catch {
    Write-Host "Warning: Could not submit test transactions. Nodes may still be starting up." -ForegroundColor Yellow
}

Write-Host "Waiting for transactions to be processed..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

Write-Host ""
Write-Host "Test environment is ready!" -ForegroundColor Green
Write-Host ""
Write-Host "Blockchain nodes are running on:" -ForegroundColor Cyan
Write-Host "   Node 1: API http://localhost:8081" -ForegroundColor White
Write-Host "   Node 2: API http://localhost:8082" -ForegroundColor White
Write-Host "   Node 3: API http://localhost:8083" -ForegroundColor White
Write-Host ""
Write-Host "You can now open the UI files in your browser:" -ForegroundColor Cyan
Write-Host "   Main Hub: file://$((Get-Location).Path)/frontend/index.html" -ForegroundColor White
Write-Host "   Node Dashboard: file://$((Get-Location).Path)/frontend/node-dashboard.html" -ForegroundColor White
Write-Host "   Blockchain Explorer: file://$((Get-Location).Path)/frontend/explorer.html" -ForegroundColor White
Write-Host "   Wallet Interface: file://$((Get-Location).Path)/frontend/wallet.html" -ForegroundColor White
Write-Host ""
Write-Host "To view logs:" -ForegroundColor Cyan
Write-Host "   Get-Content test_logs\node1.log -Wait" -ForegroundColor White
Write-Host "   Get-Content test_logs\node2.log -Wait" -ForegroundColor White
Write-Host "   Get-Content test_logs\node3.log -Wait" -ForegroundColor White
Write-Host ""
Write-Host "To stop all nodes, press Ctrl+C" -ForegroundColor Yellow
Write-Host ""

# Function to cleanup on exit
function Cleanup {
    Write-Host ""
    Write-Host "Stopping blockchain nodes..." -ForegroundColor Yellow
    Get-Process -Name "blockchain" -ErrorAction SilentlyContinue | Stop-Process -Force
    Write-Host "Test environment stopped!" -ForegroundColor Green
    exit 0
}

# Set up signal handler
Register-EngineEvent PowerShell.Exiting -Action { Cleanup }

# Keep the script running
Write-Host "Test environment is running. Press Ctrl+C to stop all nodes." -ForegroundColor Green
Write-Host ""

try {
    while ($true) {
        Start-Sleep -Seconds 1
    }
} catch {
    Cleanup
} 