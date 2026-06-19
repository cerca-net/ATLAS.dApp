@echo off
cd /d "%~dp0"
echo Killing existing app node processes...
taskkill /IM main.exe /F >nul 2>&1
taskkill /IM atlas_node.exe /F >nul 2>&1

echo Starting App Node...
start "ATLAS App Node" powershell -NoExit -Command "Set-Location '%~dp0ATLAS.BC0.0.1'; $addr = try { (Get-Content -Path './.data_admin/multiaddr.txt' -Raw -ErrorAction Stop).Trim() } catch { '' }; go run cmd/main.go --port 8001 --api 8081 --datadir .data_app --bootstrap $addr"

echo Waiting for node to initialize...
timeout /t 5

echo Starting Flutter App (Chrome)...
start "ATLAS App" powershell -NoExit -Command "Set-Location '%~dp0cercaend'; flutter run -d chrome"

echo App system launched!
pause
