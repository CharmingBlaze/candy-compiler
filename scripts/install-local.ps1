Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$root = (Resolve-Path "$PSScriptRoot\..").Path
$compilerDir = Join-Path $root "compiler"
$outDir = Join-Path $root "release\bin"

New-Item -ItemType Directory -Path $outDir -Force | Out-Null
Set-Location $compilerDir

Write-Host "Building Candy tools..."
go build -tags raylib -o (Join-Path $outDir "candy.exe") ./cmd/candy
go build -o (Join-Path $outDir "candywrap.exe") ./cmd/candywrap
go build -o (Join-Path $outDir "sweet.exe") ./cmd/sweet

Write-Host ""
Write-Host "Install complete."
Write-Host "Binaries:"
Write-Host "  $(Join-Path $outDir "candy.exe")"
Write-Host "  $(Join-Path $outDir "candywrap.exe")"
Write-Host "  $(Join-Path $outDir "sweet.exe")"
Write-Host ""
Write-Host "Optional: add to PATH for current session"
Write-Host "  `$env:Path = `"$outDir;`$env:Path`""
