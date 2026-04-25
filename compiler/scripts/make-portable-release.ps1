param(
  [Parameter(Mandatory = $true)][string]$LlvmRoot,
  [Parameter(Mandatory = $true)][string]$OutDir,
  [string]$RaylibRuntimeDir = ""
)

$ErrorActionPreference = "Stop"
$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$releaseRoot = Resolve-Path (Join-Path $repoRoot "..")

Set-Location $repoRoot

$binDir = Join-Path $OutDir "bin"
New-Item -ItemType Directory -Path $binDir -Force | Out-Null

go build -tags raylib -o (Join-Path $binDir "candy.exe") ./cmd/candy
go build -o (Join-Path $binDir "candywrap.exe") ./cmd/candywrap
go build -o (Join-Path $binDir "sweet.exe") ./cmd/sweet

$bundleDir = Join-Path $OutDir "portable-windows-x64"
& (Join-Path $PSScriptRoot "bundle-llvm.ps1") `
  -CandyBinary (Join-Path $binDir "candy.exe") `
  -CandywrapBinary (Join-Path $binDir "candywrap.exe") `
  -LlvmRoot $LlvmRoot `
  -OutDir $bundleDir `
  -RaylibRuntimeDir $RaylibRuntimeDir

Copy-Item (Join-Path $binDir "sweet.exe") (Join-Path $bundleDir "sweet.exe") -Force

Copy-Item (Join-Path $releaseRoot "examples") (Join-Path $bundleDir "examples") -Recurse -Force
Copy-Item (Join-Path $releaseRoot "docs") (Join-Path $bundleDir "docs") -Recurse -Force

Write-Host "Portable release ready at $bundleDir"
