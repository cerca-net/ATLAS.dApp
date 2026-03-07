# Test Backup System
Write-Host "🧪 Testing Backup System..." -ForegroundColor Green

# Start the blockchain node
Write-Host "🚀 Starting blockchain node..." -ForegroundColor Yellow
Start-Process -FilePath ".\blockchain.exe" -ArgumentList "-port", "8000", "-api", "8080", "-validator" -WindowStyle Hidden

# Wait for node to start
Write-Host "⏳ Waiting for node to start..." -ForegroundColor Yellow
Start-Sleep -Seconds 15

# Test backup status
Write-Host "📊 Testing backup status..." -ForegroundColor Yellow
try {
    $status = Invoke-RestMethod -Uri "http://localhost:8080/backup/status" -Method GET
    Write-Host "✅ Backup Status:" -ForegroundColor Green
    $status | ConvertTo-Json -Depth 3
} catch {
    Write-Host "❌ Failed to get backup status: $($_.Exception.Message)" -ForegroundColor Red
}

# Test backup list
Write-Host "📋 Testing backup list..." -ForegroundColor Yellow
try {
    $backups = Invoke-RestMethod -Uri "http://localhost:8080/backup/list" -Method GET
    Write-Host "✅ Backup List:" -ForegroundColor Green
    $backups | ConvertTo-Json -Depth 3
} catch {
    Write-Host "❌ Failed to get backup list: $($_.Exception.Message)" -ForegroundColor Red
}

# Test manual backup creation
Write-Host "📦 Testing manual backup creation..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/backup/create" -Method POST
    Write-Host "✅ Manual Backup Response:" -ForegroundColor Green
    $response | ConvertTo-Json -Depth 3
} catch {
    Write-Host "❌ Failed to create manual backup: $($_.Exception.Message)" -ForegroundColor Red
}

# Wait and check status again
Write-Host "⏳ Waiting for backup to complete..." -ForegroundColor Yellow
Start-Sleep -Seconds 10

try {
    $status2 = Invoke-RestMethod -Uri "http://localhost:8080/backup/status" -Method GET
    Write-Host "✅ Updated Backup Status:" -ForegroundColor Green
    $status2 | ConvertTo-Json -Depth 3
} catch {
    Write-Host "❌ Failed to get updated backup status: $($_.Exception.Message)" -ForegroundColor Red
}

# Check backup files
Write-Host "📁 Checking backup files..." -ForegroundColor Yellow
if (Test-Path "backups") {
    $backupFiles = Get-ChildItem "backups" -Filter "*.gz"
    Write-Host "✅ Found $($backupFiles.Count) backup files:" -ForegroundColor Green
    foreach ($file in $backupFiles) {
        Write-Host "  - $($file.Name) ($([math]::Round($file.Length/1KB, 2)) KB)" -ForegroundColor Cyan
    }
} else {
    Write-Host "⚠️ No backups directory found" -ForegroundColor Yellow
}

Write-Host "🎉 Backup system test completed!" -ForegroundColor Green 