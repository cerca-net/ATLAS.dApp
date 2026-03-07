@echo off
setlocal

echo Creating data directories...
if not exist "data\node1" mkdir "data\node1"
if not exist "data\node2" mkdir "data\node2"

echo Starting Node 1 (Genesis / Activator)...
start "ATLAS Node 1 (Genesis)" cmd /k "go run cmd/main.go -port 8000 -api 8080 -datadir data/node1 -key nodekey.priv -validator-key validator.hex"

echo Waiting for Node 1 to initialize...
timeout /t 5

echo.
echo Starting Node 2 (User / Peer)...
echo To connect Node 2 to Node 1, you can set NODE1_MULTIADDR variable with Node 1's address.
echo Example multiaddress format: /ip4/127.0.0.1/tcp/8000/p2p/12D3Koo...
echo.
echo For now, starting Node 2 without pre-configured peer connection.
echo You can manually connect them via the API after startup.
start "ATLAS Node 2 (User)" cmd /k "go run cmd/main.go -port 8001 -api 8081 -datadir data/node2 -key nodekey.priv -validator-key validator.hex -peers 10"

echo.
echo ========================================
echo Network started!
echo ========================================
echo Node 1 API: http://localhost:8080
echo Node 2 API: http://localhost:8081
echo.
echo NEXT STEPS:
echo 1. Check Node 1's console for its multiaddress (looks like /ip4/127.0.0.1/tcp/8000/p2p/12D3Koo...)
echo 2. Use the API or manual connection to connect the nodes
echo.
