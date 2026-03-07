# Cerca Network - Materialization Script
# This script starts the local devnet environment

Write-Host "Starting Cerca Network Devnet..." -ForegroundColor Cyan
Write-Host "Building containers (Blockchain and Frontend)..." -ForegroundColor Gray

docker compose -f docker-compose.devnet.yml up --build -d

if ($LASTEXITCODE -eq 0) {
    Write-Host "Devnet started successfully!" -ForegroundColor Green
    Write-Host "Blockchain API: http://localhost:8080/status" -ForegroundColor Yellow
    Write-Host "Frontend UI:     http://localhost:3000" -ForegroundColor Yellow
    Write-Host "`nTo see logs, run: docker compose -f docker-compose.devnet.yml logs -f" -ForegroundColor Gray
} else {
    Write-Host "Failed to start Devnet. Ensure Docker is running." -ForegroundColor Red
}
