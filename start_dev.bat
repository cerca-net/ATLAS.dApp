@echo off
cd /d "%~dp0"
echo Killing existing processes...
taskkill /IM main.exe /F >nul 2>&1
taskkill /IM atlas_node.exe /F >nul 2>&1

echo Starting Blockchain Node (Universe of One)...
start "ATLAS Node" powershell -NoExit -Command "Set-Location '%~dp0ATLAS.BC0.0.1'; go run cmd/main.go"

echo Waiting for node to initialize...
timeout /t 5

echo Opening Faucet Monitor...
start "" "%~dp0ATLAS.BC0.0.1\tools\faucet_monitor.html"

echo Starting Flutter App (Chrome)...
start "ATLAS App" powershell -NoExit -Command "Set-Location '%~dp0cercaend'; flutter run -d chrome"

echo All systems launched. Check the separate windows.
pause
