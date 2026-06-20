@echo off
cd /d "%~dp0"
echo Killing existing processes...
taskkill /IM atlas_dev.exe /F >nul 2>&1

echo Building Dev Node binary...
cd ATLAS.BC0.0.1
go build -o atlas_dev.exe cmd/main.go
cd ..

echo Starting Blockchain Node (Universe of One)...
start "ATLAS Node" powershell -NoExit -Command "Set-Location '%~dp0ATLAS.BC0.0.1'; .\atlas_dev.exe"

echo Waiting for node to initialize...
timeout /t 5

echo Starting React Admin Panel...
start "Cerca Admin Panel" powershell -NoExit -Command "Set-Location '%~dp0cerca-admin-panel'; npm run dev"

echo Starting Flutter App (Chrome)...
start "ATLAS App" powershell -NoExit -Command "Set-Location '%~dp0cercaend'; flutter run -d chrome"

echo All systems launched. Check the separate windows.
pause
