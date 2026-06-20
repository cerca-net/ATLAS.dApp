@echo off
cd /d "%~dp0"
echo Killing existing admin/treasury node processes...
taskkill /IM atlas_treasury.exe /F >nul 2>&1

echo Building Treasury Node binary...
cd ATLAS.BC0.0.1
go build -o atlas_treasury.exe cmd/main.go
cd ..

echo Starting Treasury Node...
start "ATLAS Treasury Node" powershell -NoExit -Command "Set-Location '%~dp0ATLAS.BC0.0.1'; .\atlas_treasury.exe --port 8000 --api 8080 --datadir .data_admin"

echo Waiting for node to initialize...
timeout /t 5

echo Starting React Admin Panel...
start "Cerca Admin Panel" powershell -NoExit -Command "Set-Location '%~dp0cerca-admin-panel'; npm run dev"

echo Admin system launched!
pause
