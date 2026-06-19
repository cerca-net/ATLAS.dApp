@echo off
cd /d "%~dp0"
echo Killing existing processes...
taskkill /IM main.exe /F >nul 2>&1
taskkill /IM atlas_node.exe /F >nul 2>&1

echo Starting Node 1 (Port: 8000, API: 8080)...
start "ATLAS Node 1" powershell -NoExit -Command "Set-Location '%~dp0ATLAS.BC0.0.1'; go run cmd/main.go -port 8000 -api 8080 -datadir ./data_node1 -validator-key node1.key"

echo Waiting 5 seconds before starting Node 2...
timeout /t 5

echo Starting Node 2 (Port: 8001, API: 8081)...
start "ATLAS Node 2" powershell -NoExit -Command "Set-Location '%~dp0ATLAS.BC0.0.1'; $addr = try { (Get-Content -Path './data_node1/multiaddr.txt' -Raw -ErrorAction Stop).Trim() } catch { '' }; go run cmd/main.go -port 8001 -api 8081 -datadir ./data_node2 -validator-key node2.key -bootstrap $addr"

echo Two local nodes launched!
echo They will use MDNS to discover each other automatically and synchronize the blockchain.
pause
