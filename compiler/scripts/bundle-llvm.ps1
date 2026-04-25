param(
  [Parameter(Mandatory = $true)][string]$CandyBinary,
  [Parameter(Mandatory = $true)][string]$CandywrapBinary,
  [Parameter(Mandatory = $true)][string]$LlvmRoot,
  [Parameter(Mandatory = $true)][string]$OutDir,
  [string]$RaylibRuntimeDir = ""
)

$ErrorActionPreference = "Stop"

if (!(Test-Path $CandyBinary)) { throw "Candy binary not found: $CandyBinary" }
if (!(Test-Path $CandywrapBinary)) { throw "Candywrap binary not found: $CandywrapBinary" }
if (!(Test-Path $LlvmRoot)) { throw "LLVM root not found: $LlvmRoot" }
$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$manifest = $env:CANDY_LLVM_MANIFEST
if ([string]::IsNullOrWhiteSpace($manifest)) {
  $manifest = Join-Path $PSScriptRoot "llvm-manifest.windows-x64.txt"
}
if (Test-Path $manifest) {
  $patterns = Get-Content $manifest | Where-Object { $_ -and $_.Trim() -ne "" }
  foreach ($pat in $patterns) {
    $trim = $pat.Trim()
    if ($trim.Contains("*")) {
      $foundMatches = Get-ChildItem -Path $LlvmRoot -Recurse -File -ErrorAction SilentlyContinue |
        Where-Object { $_.FullName -like (Join-Path $LlvmRoot $trim) }
      if (!$foundMatches -or $foundMatches.Count -eq 0) {
        throw "missing required LLVM artifact for pattern: $pat"
      }
      continue
    }
    $exact = Join-Path $LlvmRoot $trim
    if (!(Test-Path $exact)) {
      throw "missing required LLVM artifact: $pat"
    }
  }
}

New-Item -ItemType Directory -Path $OutDir -Force | Out-Null
$binOut = Join-Path $OutDir "bin"
New-Item -ItemType Directory -Path $binOut -Force | Out-Null
Copy-Item $CandyBinary (Join-Path $binOut (Split-Path $CandyBinary -Leaf)) -Force
Copy-Item $CandywrapBinary (Join-Path $binOut (Split-Path $CandywrapBinary -Leaf)) -Force

$toolchainOut = Join-Path $OutDir "toolchain"
New-Item -ItemType Directory -Path $toolchainOut -Force | Out-Null
Copy-Item (Join-Path $LlvmRoot "bin") (Join-Path $toolchainOut "bin") -Recurse -Force
if (Test-Path (Join-Path $LlvmRoot "lib")) {
  Copy-Item (Join-Path $LlvmRoot "lib") (Join-Path $toolchainOut "lib") -Recurse -Force
}

# Backward-compatible copy for older bundles/scripts expecting .\llvm.
$llvmOut = Join-Path $OutDir "llvm"
if (!(Test-Path $llvmOut)) {
  Copy-Item $toolchainOut $llvmOut -Recurse -Force
}

if (Test-Path (Join-Path $repoRoot "licenses\LLVM-LICENSE.txt")) {
  New-Item -ItemType Directory -Path (Join-Path $OutDir "licenses") -Force | Out-Null
  Copy-Item (Join-Path $repoRoot "licenses\LLVM-LICENSE.txt") (Join-Path $OutDir "licenses\LLVM-LICENSE.txt") -Force
} else {
  throw "missing required license file: licenses\\LLVM-LICENSE.txt"
}

if (-not [string]::IsNullOrWhiteSpace($RaylibRuntimeDir) -and (Test-Path $RaylibRuntimeDir)) {
  $raylibOut = Join-Path $OutDir "raylib-runtime"
  New-Item -ItemType Directory -Path $raylibOut -Force | Out-Null
  Copy-Item (Join-Path $RaylibRuntimeDir "*") $raylibOut -Recurse -Force
}

$readme = @"
Candy Portable Bundle (Windows)
===============================

This folder is self-contained. No global installs are required.

Included:
- bin\candy.exe
- bin\candywrap.exe
- sweet.exe
- toolchain/ (clang + toolchain used by candy -build)
- llvm/ (compatibility copy for older bundles)
- licenses/
- optional raylib-runtime/ (if provided at packaging time)

Usage:
- .\bin\candy.exe script.candy
- .\bin\candywrap.exe wrap --name mylib --output .\bindings mylib.h
- .\sweet.exe convert --name mylib --output .\bindings mylib.h
"@
$readme | Set-Content -Path (Join-Path $OutDir "README_PORTABLE.txt") -Encoding ASCII

Write-Host "Bundled package created at $OutDir"
