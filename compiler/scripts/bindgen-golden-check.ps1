Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$Root = (Resolve-Path "$PSScriptRoot\..\..").Path
Set-Location "$Root\compiler"

Write-Host "Running candy_bindgen golden drift check..."
go test ./candy_bindgen -run TestGolden_
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
Write-Host "bindgen golden check: OK"

