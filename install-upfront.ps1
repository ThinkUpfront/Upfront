# install-upfront.ps1 — Download and install upfront on Windows.
# Usage: irm https://raw.githubusercontent.com/brennhill/upfront/main/install-upfront.ps1 | iex

$ErrorActionPreference = "Stop"

$Repo = "brennhill/upfront"
$InstallDir = if ($env:INSTALL_DIR) { $env:INSTALL_DIR } else { "$env:LOCALAPPDATA\upfront\bin" }

function Write-Info($msg) { Write-Host "==> $msg" -ForegroundColor Blue }
function Write-Warn($msg) { Write-Host "warning: $msg" -ForegroundColor Yellow }

# Get latest version
Write-Info "Installing upfront - audit trail for AI-assisted feature definition"
Write-Host ""

$Release = Invoke-RestMethod "https://api.github.com/repos/$Repo/releases/latest"
$Version = $Release.tag_name -replace '^v', ''
Write-Info "Latest version: v$Version"

# Download
$Archive = "upfront_${Version}_windows_amd64.zip"
$Url = "https://github.com/$Repo/releases/download/v${Version}/$Archive"
$TmpDir = Join-Path $env:TEMP "upfront-install"

if (Test-Path $TmpDir) { Remove-Item $TmpDir -Recurse -Force }
New-Item -ItemType Directory -Path $TmpDir | Out-Null

Write-Info "Downloading $Archive..."
Invoke-WebRequest -Uri $Url -OutFile (Join-Path $TmpDir $Archive)

# Extract
Expand-Archive -Path (Join-Path $TmpDir $Archive) -DestinationPath $TmpDir -Force

# Install
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}
Copy-Item (Join-Path $TmpDir "upfront.exe") (Join-Path $InstallDir "upfront.exe") -Force
Write-Info "Installed to $InstallDir\upfront.exe"

# Add to PATH if needed
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    Write-Info "Added $InstallDir to user PATH (restart your terminal)"
}

# Register Claude Code hook
$SettingsFile = Join-Path $env:USERPROFILE ".claude\settings.json"
$SettingsDir = Split-Path $SettingsFile

if (-not (Test-Path $SettingsDir)) {
    New-Item -ItemType Directory -Path $SettingsDir | Out-Null
}

if (-not (Test-Path $SettingsFile)) {
    '{}' | Out-File -FilePath $SettingsFile -Encoding utf8
}

$Settings = Get-Content $SettingsFile -Raw | ConvertFrom-Json

$HookExists = $false
if ($Settings.hooks -and $Settings.hooks.PostToolUse) {
    foreach ($entry in $Settings.hooks.PostToolUse) {
        if ($entry.matcher -eq "Skill") {
            foreach ($h in $entry.hooks) {
                if ($h.command -eq "upfront hook") {
                    $HookExists = $true
                }
            }
        }
    }
}

if (-not $HookExists) {
    # Back up settings
    $Backup = "$SettingsFile.backup.$(Get-Date -Format 'yyyyMMddHHmmss')"
    Copy-Item $SettingsFile $Backup
    Write-Info "Backed up settings to $Backup"

    # Add hook (manual JSON manipulation to avoid dependency on jq)
    $Hook = @{
        matcher = "Skill"
        hooks = @(@{ type = "command"; command = "upfront hook" })
    }

    if (-not $Settings.hooks) {
        $Settings | Add-Member -NotePropertyName "hooks" -NotePropertyValue @{} -Force
    }
    if (-not $Settings.hooks.PostToolUse) {
        $Settings.hooks | Add-Member -NotePropertyName "PostToolUse" -NotePropertyValue @() -Force
    }
    $Settings.hooks.PostToolUse += $Hook

    $Settings | ConvertTo-Json -Depth 10 | Out-File -FilePath $SettingsFile -Encoding utf8
    Write-Info "Registered PostToolUse hook in Claude Code"
} else {
    Write-Info "PostToolUse hook already registered"
}

# Cleanup
Remove-Item $TmpDir -Recurse -Force

Write-Host ""
Write-Info "Done! Run 'upfront status' to verify."
Write-Host ""
Write-Host "  Docs:  https://github.com/$Repo"
Write-Host "  Book:  https://upfront.dev/book"
Write-Host ""
