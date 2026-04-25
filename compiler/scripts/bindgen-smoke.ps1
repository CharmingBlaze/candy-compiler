Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$Root = (Resolve-Path "$PSScriptRoot\..\..").Path
Set-Location "$Root\compiler"

go run ./cmd/candywrap wrap --name mylib --output ../examples/bindgen/out ../examples/bindgen/mylib.h
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
go run ./cmd/candy -ast ../examples/bindgen/main.candy | Out-Null
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Write-Host "bindgen smoke: OK"
