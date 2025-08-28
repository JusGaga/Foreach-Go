# Script de test pour l'API WordMon Go
Write-Host "=== Test de l'API WordMon Go ===" -ForegroundColor Green

# Test 1: Status
Write-Host "`n1. Test du status..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/status" -Method Get
    Write-Host "Status: $($response | ConvertTo-Json -Depth 3)" -ForegroundColor Green
} catch {
    Write-Host "Erreur status: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 2: Créer un joueur
Write-Host "`n2. Test création joueur..." -ForegroundColor Yellow
try {
    $playerData = @{
        name = "Sacha"
    } | ConvertTo-Json
    
    $response = Invoke-RestMethod -Uri "http://localhost:8080/players" -Method Post -Body $playerData -ContentType "application/json"
    Write-Host "Joueur créé: $($response | ConvertTo-Json -Depth 3)" -ForegroundColor Green
    
    $playerId = $response.id
} catch {
    Write-Host "Erreur création joueur: $($_.Exception.Message)" -ForegroundColor Red
    $playerId = "p1" # Fallback pour les tests suivants
}

# Test 3: Récupérer un joueur
Write-Host "`n3. Test récupération joueur..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/players/$playerId" -Method Get
    Write-Host "Joueur récupéré: $($response | ConvertTo-Json -Depth 3)" -ForegroundColor Green
} catch {
    Write-Host "Erreur récupération joueur: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 4: Vérifier le spawn courant
Write-Host "`n4. Test spawn courant..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/spawn/current" -Method Get
    Write-Host "Spawn actuel: $($response | ConvertTo-Json -Depth 3)" -ForegroundColor Green
} catch {
    Write-Host "Erreur spawn courant: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 5: Leaderboard
Write-Host "`n5. Test leaderboard..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/leaderboard?limit=5" -Method Get
    Write-Host "Leaderboard: $($response | ConvertTo-Json -Depth 3)" -ForegroundColor Green
} catch {
    Write-Host "Erreur leaderboard: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== Tests terminés ===" -ForegroundColor Green
