param(
    [switch]$KeepData = $false
)

Write-Host "==========================="
Write-Host " ATLAS CercaChain Devnet "
Write-Host "==========================="

if (-not $KeepData) {
    Write-Host "[*] Cleaning up old data directories..."
    if (Test-Path "data_node1") { Remove-Item -Recurse -Force "data_node1" }
    if (Test-Path "data_node2") { Remove-Item -Recurse -Force "data_node2" }
    if (Test-Path "node1.log") { Remove-Item -Force "node1.log" }
}

Write-Host "[*] Ensuring data directories exist..."
New-Item -ItemType Directory -Force -Path "data_node1" | Out-Null
New-Item -ItemType Directory -Force -Path "data_node2" | Out-Null

# We use absolute paths or relative to execution
$exePath = ".\build\atlas-node.exe"

Write-Host "[*] Starting Node 1 (Genesis Validator) on Port 8000 / API 8080..."
Start-Process -FilePath $exePath -ArgumentList "-port 8000 -api 8080 -datadir data_node1 -validator-key genesis_validator.key -genesis genesis.json" -RedirectStandardOutput "node1.log" -RedirectStandardError "node1_err.log" -WindowStyle Hidden

Write-Host "[*] Waiting for Node 1 to bindings and generate PeerID..."
Start-Sleep -Seconds 6

$multiaddr = $null
if (Test-Path "node1.log") {
    $nodeIdLine = Get-Content "node1.log" | Where-Object { $_ -match '\[libp2p\] Node ID: (.*)' } | Select-Object -First 1
    if ($nodeIdLine -match '\[libp2p\] Node ID: ([A-Za-z0-9]+)') {
        $nodeId = $matches[1]
        $multiaddr = "/ip4/127.0.0.1/tcp/8000/p2p/$nodeId"
    }
}

if ($multiaddr) {
    Write-Host "[+] Node 1 Multiaddress identified: $multiaddr"
    Write-Host "[*] Starting Node 2 (Observer Node) on Port 8001 / API 8081..."
    
    $env:NODE1_MULTIADDR = $multiaddr
    Start-Process -FilePath $exePath -ArgumentList "-port 8001 -api 8081 -datadir data_node2 -genesis genesis.json -validator=false" -RedirectStandardOutput "node2.log" -RedirectStandardError "node2_err.log" -WindowStyle Hidden
    
    Write-Host "[+] Devnet running in background!"
    Write-Host "[+] Check node1.log and node2.log for output. Run 'Stop-Process -Name atlas-node' to terminate."
} else {
    Write-Host "[!] Failed to detect Node 1 multiaddress. Check node1.log for errors."
}
