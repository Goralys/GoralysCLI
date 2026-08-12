# This script was generated using AI
# The model used was: Claude Sonnet 5 [reasoning: low]
# Might need a human rewrite one day
# Do not hesitate if you have the knowledge and faith to do it :)

$ErrorActionPreference = "Stop"

$RepoBase = "https://cli.goralys.fr/release"
$BinaryName = "goralys-cli"
$Version = if ($args.Count -gt 0) { $args[0] } else { "latest" }
$InstallDir = if ($env:GORALYS_INSTALL_DIR) { $env:GORALYS_INSTALL_DIR } else { "$env:USERPROFILE\.goralys\bin" }

# --- detect arch ---
$arch = if ([Environment]::Is64BitOperatingSystem) {
    if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }
} else {
    Write-Error "Unsupported architecture: 32-bit Windows is not supported."
    exit 1
}

$asset = "$BinaryName-windows-$arch.exe"
$url = "$RepoBase/$Version/$asset"

# --- download ---
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
$dest = Join-Path $InstallDir "$BinaryName.exe"

Write-Host "Downloading $asset ($Version)..."

try {
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    Invoke-WebRequest -Uri $url -OutFile $dest -UseBasicParsing
} catch {
    Write-Error "Failed to download from $url"
    Write-Error "Check that this platform/architecture is supported and the version exists."
    exit 1
}

Write-Host "Installed to $dest"

# --- add to PATH (user scope, no admin required) ---
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")

if ($currentPath -notlike "*$InstallDir*") {
    $newPath = if ($currentPath) { "$currentPath;$InstallDir" } else { $InstallDir }
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Host "Added $InstallDir to your user PATH."
    Write-Host "Restart your terminal for this to take effect."
} else {
    Write-Host "$InstallDir is already in your PATH."
}

Write-Host ""
Write-Host "Done! Restart your terminal, then verify with: $BinaryName --version"