# Agent Extensions Install (Windows PowerShell)
# Creates global symlinks and local copies of the repo's folders
# Local copies avoid repo-specific symlinks checked into git

$ErrorActionPreference = "Stop"

function Write-Info { param($Message) Write-Host "[INFO] $Message" -ForegroundColor Blue }
function Write-Success { param($Message) Write-Host "[OK] $Message" -ForegroundColor Green }
function Write-Warn { param($Message) Write-Host "[WARN] $Message" -ForegroundColor Yellow }
function Write-Err { param($Message) Write-Host "[ERROR] $Message" -ForegroundColor Red; exit 1 }

function Confirm-YesNo {
    param([string]$Prompt)
    $answer = Read-Host "$Prompt [y/N]"
    return $answer -match "^[Yy]"
}

function Find-GitRoot {
    $dir = (Get-Location).Path
    while ($true) {
        if (Test-Path (Join-Path $dir ".git")) { return $dir }
        $parent = Split-Path $dir -Parent
        if ($parent -eq $dir) { break }
        $dir = $parent
    }
    return $null
}

function Install-Symlinks {
    param(
        [string]$TargetRoot,
        [string]$PayloadDir,
        [string]$Label
    )

    if (-not (Test-Path $PayloadDir)) {
        Write-Warn "Payload directory '$PayloadDir' not found, skipping $Label"
        return $false
    }

    Write-Info "Installing $Label to: $TargetRoot"

    $files = Get-ChildItem -Path $PayloadDir -Recurse -File |
        Where-Object { $_.Name -notin @(".DS_Store", "Thumbs.db") } |
        ForEach-Object { $_.FullName.Substring($PayloadDir.Length + 1) }

    $conflicts = @()
    foreach ($file in $files) {
        $dest = Join-Path $TargetRoot $file
        if (Test-Path $dest) {
            if ((Get-Item $dest -ErrorAction SilentlyContinue).LinkType) {
                $existingTarget = (Get-Item $dest).Target
                $src = Join-Path $PayloadDir $file
                if ($existingTarget -eq $src) { continue }
            }
            $conflicts += $file
        }
    }

    if ($conflicts.Count -gt 0) {
        Write-Host ""
        Write-Warn "Found $($conflicts.Count) conflicting file(s)/symlink(s) for $Label:"
        Write-Host ""
        $conflicts | Select-Object -First 20 | ForEach-Object { Write-Host $_ }
        if ($conflicts.Count -gt 20) {
            Write-Host "  ... and $($conflicts.Count - 20) more"
        }
        Write-Host ""
        Write-Warn "These will be replaced with symlinks to the repo."
        Write-Host ""
        if (-not (Confirm-YesNo "Replace ALL conflicting files with symlinks for $Label?")) {
            Write-Warn "Skipping $Label install"
            return $false
        }
        Write-Host ""
    }

    Write-Info "Creating symlinks for $Label..."

    foreach ($file in $files) {
        $src = Join-Path $PayloadDir $file
        $dest = Join-Path $TargetRoot $file
        $destDir = Split-Path $dest -Parent

        if (-not (Test-Path $destDir)) {
            New-Item -ItemType Directory -Path $destDir -Force | Out-Null
        }

        if (Test-Path $dest) {
            Remove-Item $dest -Force
        }

        New-Item -ItemType SymbolicLink -Path $dest -Target $src | Out-Null
    }

    Write-Success "Symlinks created for $Label at: $TargetRoot"
    return $true
}

function Install-Copies {
    param(
        [string]$TargetRoot,
        [string]$PayloadDir,
        [string]$Label
    )

    if (-not (Test-Path $PayloadDir)) {
        Write-Warn "Payload directory '$PayloadDir' not found, skipping $Label"
        return $false
    }

    Write-Info "Installing $Label to: $TargetRoot"

    $files = Get-ChildItem -Path $PayloadDir -Recurse -File |
        Where-Object { $_.Name -notin @(".DS_Store", "Thumbs.db") } |
        ForEach-Object { $_.FullName.Substring($PayloadDir.Length + 1) }

    $conflicts = @()
    foreach ($file in $files) {
        $dest = Join-Path $TargetRoot $file
        if (Test-Path $dest) {
            $conflicts += $file
        }
    }

    if ($conflicts.Count -gt 0) {
        Write-Host ""
        Write-Warn "Found $($conflicts.Count) existing file(s) that will be overwritten:"
        Write-Host ""
        $conflicts | Select-Object -First 20 | ForEach-Object { Write-Host $_ }
        if ($conflicts.Count -gt 20) {
            Write-Host "  ... and $($conflicts.Count - 20) more"
        }
        Write-Host ""
        if (-not (Confirm-YesNo "Overwrite ALL conflicting files for $Label?")) {
            Write-Warn "Skipping $Label install"
            return $false
        }
        Write-Host ""
    }

    Write-Info "Copying files for $Label..."

    foreach ($file in $files) {
        $src = Join-Path $PayloadDir $file
        $dest = Join-Path $TargetRoot $file
        $destDir = Split-Path $dest -Parent

        if (-not (Test-Path $destDir)) {
            New-Item -ItemType Directory -Path $destDir -Force | Out-Null
        }

        if (Test-Path $dest) {
            Remove-Item $dest -Force
        }

        Copy-Item -Path $src -Destination $dest -Force
    }

    Write-Success "Files copied for $Label at: $TargetRoot"
    return $true
}

function Main {
    Write-Host ""
    Write-Host "  Agent Extensions Install"
    Write-Host "  ========================"
    Write-Host ""

    $scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

    Write-Info "Source directory: $scriptDir"
    Write-Host ""

    Write-Host "Which tool(s) would you like to install extensions for?"
    Write-Host "  1) OpenCode"
    Write-Host "  2) Augment (Auggie)"
    Write-Host "  3) Codex"
    Write-Host "  4) OpenCode + Augment"
    Write-Host "  5) OpenCode + Codex"
    Write-Host "  6) Augment + Codex"
    Write-Host "  7) All (OpenCode + Augment + Codex)"
    Write-Host ""
    $toolChoice = Read-Host "Enter choice [1/2/3/4/5/6/7]"

    $InstallOpenCode = $false
    $InstallAugment = $false
    $InstallCodex = $false

    switch ($toolChoice) {
        "1" { $InstallOpenCode = $true }
        "2" { $InstallAugment = $true }
        "3" { $InstallCodex = $true }
        "4" { $InstallOpenCode = $true; $InstallAugment = $true }
        "5" { $InstallOpenCode = $true; $InstallCodex = $true }
        "6" { $InstallAugment = $true; $InstallCodex = $true }
        "7" { $InstallOpenCode = $true; $InstallAugment = $true; $InstallCodex = $true }
        default { Write-Err "Invalid choice. Please enter 1, 2, 3, 4, 5, 6, or 7." }
    }

    Write-Host ""
    Write-Host "Where would you like to install?"
    Write-Host "  1) Global (user config directory)"
    Write-Host "  2) Local (current repo)"
    Write-Host "  3) Both global and local"
    Write-Host ""
    $scopeChoice = Read-Host "Enter choice [1/2/3]"

    $InstallGlobal = $false
    $InstallLocal = $false

    switch ($scopeChoice) {
        "1" { $InstallGlobal = $true }
        "2" { $InstallLocal = $true }
        "3" { $InstallGlobal = $true; $InstallLocal = $true }
        default { Write-Err "Invalid choice. Please enter 1, 2, or 3." }
    }

    if ($InstallLocal) {
        $gitRoot = Find-GitRoot
        if (-not $gitRoot) { Write-Err "Not inside a git repository. Cannot determine repo root for local install." }
    }

    Write-Host ""

    $installedCount = 0

    if ($InstallOpenCode -and $InstallGlobal) {
        $opencodePayload = Join-Path $scriptDir "opencode"
        $opencodeTarget = Join-Path $HOME ".config\opencode"
        if (Install-Symlinks -TargetRoot $opencodeTarget -PayloadDir $opencodePayload -Label "OpenCode (global)") {
            $installedCount++
        }
        Write-Host ""
    }

    if ($InstallOpenCode -and $InstallLocal) {
        $opencodePayload = Join-Path $scriptDir "opencode"
        $opencodeTarget = Join-Path $gitRoot ".opencode"
        if (Install-Copies -TargetRoot $opencodeTarget -PayloadDir $opencodePayload -Label "OpenCode (local)") {
            $installedCount++
        }
        Write-Host ""
    }

    if ($InstallAugment -and $InstallGlobal) {
        $augmentPayload = Join-Path $scriptDir "augment"
        $augmentTarget = Join-Path $HOME ".augment"
        if (Install-Symlinks -TargetRoot $augmentTarget -PayloadDir $augmentPayload -Label "Augment (global)") {
            $installedCount++
        }
        Write-Host ""
    }

    if ($InstallAugment -and $InstallLocal) {
        $augmentPayload = Join-Path $scriptDir "augment"
        $augmentTarget = Join-Path $gitRoot ".augment"
        if (Install-Copies -TargetRoot $augmentTarget -PayloadDir $augmentPayload -Label "Augment (local)") {
            $installedCount++
        }
        Write-Host ""
    }

    if ($InstallCodex -and $InstallGlobal) {
        $codexPayload = Join-Path $scriptDir "codex"
        $codexTarget = Join-Path $HOME ".codex"
        if (Install-Symlinks -TargetRoot $codexTarget -PayloadDir $codexPayload -Label "Codex (global)") {
            $installedCount++
        }
        Write-Host ""
    }

    if ($InstallCodex -and $InstallLocal) {
        $codexPayload = Join-Path $scriptDir "codex"
        $codexTarget = Join-Path $gitRoot ".codex"
        if (Install-Copies -TargetRoot $codexTarget -PayloadDir $codexPayload -Label "Codex (local)") {
            $installedCount++
        }
        Write-Host ""
    }

    if ($installedCount -eq 0) { Write-Err "No installations completed" }

    Write-Success "Install complete!"
    Write-Host ""
    Write-Info "Source files at: $scriptDir"
    Write-Host ""
    Write-Host "Global installs use symlinks to the repo."
    Write-Host "Local installs copy files into the repo."
    Write-Host ""
}

Main
