param(
  [Parameter(Mandatory = $true)][string]$Message
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$root = (Resolve-Path "$PSScriptRoot\..").Path
Set-Location $root

git add .
git commit -m $Message
git push origin HEAD:main

Write-Host "Pushed to origin/main"
