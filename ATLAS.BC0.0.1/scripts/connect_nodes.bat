@echo off
setlocal

if "%1"=="" (
    echo Usage: connect_nodes.bat ^<Node1_MultiAddress^>
    echo.
    echo Example:
    echo   connect_nodes.bat /ip4/127.0.0.1/tcp/8000/p2p/12D3KooWLYq8QetDBFhFppY9A47V81ZH5B6cwLmcJtG3NFdcmF1C
    echo.
    echo This will connect Node 2 to Node 1 using the provided multiaddress.
    exit /b 1
)

set "NODE1_ADDR=%1"

echo Connecting Node 2 to Node 1...
echo Node 1 multiaddress: %NODE1_ADDR%
echo.

curl.exe -X POST http://localhost:8081/connect-peer -H "Content-Type: application/json" -d "{\"peer_address\": \"%NODE1_ADDR%\"}"

echo.
echo.
echo Connection request sent!
echo Check the node consoles for connection status.
