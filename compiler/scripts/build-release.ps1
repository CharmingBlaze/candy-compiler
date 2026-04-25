param(
  [Parameter(Mandatory = $true)][string]$LlvmRoot,
  [Parameter(Mandatory = $true)][string]$OutDir,
  [string]$Version = "",
  [string]$RaylibRuntimeDir = ""
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$projectRoot = Resolve-Path (Join-Path $repoRoot "..")
New-Item -ItemType Directory -Path $OutDir -Force | Out-Null

if ([string]::IsNullOrWhiteSpace($Version)) {
  try {
    $Version = (git -C $projectRoot describe --tags --always).Trim()
  } catch {
    $Version = "dev"
  }
}

$stdlibFiles = Get-ChildItem -Path (Join-Path $repoRoot "candy_stdlib") -Filter "*.go" | Sort-Object FullName
$combined = New-Object System.Text.StringBuilder
foreach ($f in $stdlibFiles) {
  $h = (Get-FileHash -Path $f.FullName -Algorithm SHA256).Hash
  [void]$combined.Append($h)
}
$bytes = [System.Text.Encoding]::UTF8.GetBytes($combined.ToString())
$sha = [System.Security.Cryptography.SHA256]::Create()
$hashBytes = $sha.ComputeHash($bytes)
$stdlibHash = -join ($hashBytes | ForEach-Object { $_.ToString("x2") })

$stageRoot = Join-Path $OutDir ".release-stage-$([System.Guid]::NewGuid().ToString('N'))"
$bundleDir = Join-Path $stageRoot "candy-release"
New-Item -ItemType Directory -Path (Join-Path $bundleDir "bin") -Force | Out-Null
New-Item -ItemType Directory -Path (Join-Path $bundleDir "lib") -Force | Out-Null
New-Item -ItemType Directory -Path (Join-Path $bundleDir "docs") -Force | Out-Null

Set-Location $repoRoot
$ldflags = "-X main.BuildVersion=$Version -X main.BuildStdlibHash=$stdlibHash"
go build -tags raylib -ldflags "$ldflags" -o (Join-Path $bundleDir "bin\candy.exe") ./cmd/candy
go build -ldflags "$ldflags" -o (Join-Path $bundleDir "bin\candywrap.exe") ./cmd/candywrap
go build -ldflags "$ldflags" -o (Join-Path $bundleDir "bin\sweet.exe") ./cmd/sweet

& (Join-Path $PSScriptRoot "bundle-llvm.ps1") `
  -CandyBinary (Join-Path $bundleDir "bin\candy.exe") `
  -CandywrapBinary (Join-Path $bundleDir "bin\candywrap.exe") `
  -LlvmRoot $LlvmRoot `
  -OutDir $bundleDir `
  -RaylibRuntimeDir $RaylibRuntimeDir

"stdlib_hash=$stdlibHash`nversion=$Version`n" | Set-Content -Path (Join-Path $bundleDir "lib\STDLIB_MANIFEST.txt") -Encoding ASCII
$manifestText = Get-Content -Path (Join-Path $bundleDir "lib\STDLIB_MANIFEST.txt") -Raw
if ($manifestText -notmatch [regex]::Escape($stdlibHash)) {
  throw "stdlib manifest integrity check failed"
}

if (Test-Path (Join-Path $projectRoot "docs")) {
  Copy-Item (Join-Path $projectRoot "docs\*") (Join-Path $bundleDir "docs") -Recurse -Force
}
if (Test-Path (Join-Path $projectRoot "templates")) {
  Copy-Item (Join-Path $projectRoot "templates") (Join-Path $bundleDir "templates") -Recurse -Force
} else {
  New-Item -ItemType Directory -Path (Join-Path $bundleDir "templates") -Force | Out-Null
}
Copy-Item (Join-Path $projectRoot "examples") (Join-Path $bundleDir "examples") -Recurse -Force

$readme = @"
# Candy Portable Release

Version: $Version  
Stdlib hash: $stdlibHash

This release is self-contained:

- binaries in `bin/`
- native backend toolchain in `toolchain/`
- compatibility copy in `llvm/`
- docs in `docs/`
- templates in `templates/`

Quick checks:

```powershell
.\bin\candy.exe doctor
.\bin\candy.exe --help
.\bin\candywrap.exe wrap --help
.\bin\sweet.exe convert --help
```
"@
$readme | Set-Content -Path (Join-Path $bundleDir "README.md") -Encoding UTF8

$archive = Join-Path $OutDir ("candy-{0}-windows-x64.zip" -f $Version)
Compress-Archive -Path (Join-Path $stageRoot "candy-release") -DestinationPath $archive -Force

Write-Host "Release bundle ready:"
Write-Host "  folder: $bundleDir"
Write-Host "  archive: $archive"
