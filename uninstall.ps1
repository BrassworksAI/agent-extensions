# Agent Extensions Uninstaller (Windows PowerShell)
# Removes the Spec-Driven Development process for OpenCode, Augment, and/or Codex

$ErrorActionPreference = "Stop"

function Write-Info { param($Message) Write-Host "[INFO] $Message" -ForegroundColor Blue }
function Write-Success { param($Message) Write-Host "[OK] $Message" -ForegroundColor Green }
function Write-Warn { param($Message) Write-Host "[WARN] $Message" -ForegroundColor Yellow }
function Write-Err { param($Message) Write-Host "[ERROR] $Message" -ForegroundColor Red; exit 1 }

# Find git root by walking up directories
function Find-GitRoot {
    $dir = Get-Location
    while ($dir -ne $null -and $dir.Path -ne "") {
        $gitDir = Join-Path $dir.Path ".git"
        if (Test-Path $gitDir) {
            return $dir.Path
        }
        $parent = Split-Path $dir.Path -Parent
        if ($parent -eq $dir.Path) { break }
        $dir = Get-Item $parent -ErrorAction SilentlyContinue
    }
    return $null
}

# Remove a file if it exists
function Remove-FileIfExists {
    param($Path)
    if (Test-Path $Path) {
        Remove-Item -Path $Path -Force
        return $true
    }
    return $false
}

# Remove a directory only if empty
function Remove-EmptyDirectory {
    param($Path)
    if ((Test-Path $Path) -and ((Get-ChildItem $Path -Force | Measure-Object).Count -eq 0)) {
        Remove-Item -Path $Path -Force -ErrorAction SilentlyContinue
    }
}

# Remove OpenCode SDD files from a target directory
function Remove-OpenCodeFiles {
    param($Target)
    $removed = 0

    # OpenCode file list - exact files installed (keep in sync with repo contents)
    $files = @(
        "agent\sdd\plan.md",
        "agent\sdd\build.md",
        "agent\librarian.md",
        "agent\chat.md",
        "agent\designer.md",
        "command\sdd\init.md",
        "command\sdd\plan.md",
        "command\sdd\implement.md",
        "command\sdd\proposal.md",
        "command\sdd\specs.md",
        "command\sdd\discovery.md",
        "command\sdd\tasks.md",
        "command\sdd\reconcile.md",
        "command\sdd\finish.md",
        "command\sdd\status.md",
        "command\sdd\brainstorm.md",
        "command\sdd\explain.md",
        "command\sdd\fast\vibe.md",
        "command\sdd\fast\bug.md",
        "command\sdd\tools\critique.md",
        "command\sdd\tools\scenario-test.md",
        "command\sdd\tools\taxonomy-map.md",
        "command\sdd\tools\prime-specs.md",
        "command\tool\commit.md",
        "command\create\agent.md",
        "command\create\command.md",
        "plugin\spec-validate.ts",
        "skill\architecture-fit-check\SKILL.md",
        "skill\architecture-workshop\SKILL.md",
        "skill\bun-shell-commands\SKILL.md",
        "skill\design-case-study-generator\SKILL.md",
        "skill\design-case-study-generator\references\case-study-template.md",
        "skill\design-case-study-generator\references\demo-skeleton.md",
        "skill\design-case-study-generator\references\pitch-card-template.md",
        "skill\design-case-study-generator\references\token-schema-guidance.md",
        "skill\design-case-study-generator\references\tokens-css-emitter.md",
        "skill\design-case-study-generator\scripts\copy-version.ts",
        "skill\keep-current\SKILL.md",
        "skill\merge-change-specs\SKILL.md",
        "skill\merge-change-specs\references\delta-merge-rules.md",
        "skill\merge-change-specs\scripts\merge-change-specs.ts",
        "skill\research\SKILL.md",
        "skill\sdd-state-management\SKILL.md",
        "skill\skill-creator\SKILL.md",
        "skill\skill-creator\references\references-guide.md",
        "skill\skill-creator\references\scripts-overview.md",
        "skill\skill-creator\references\scripts-runtime-node-first.md",
        "skill\spec-format\SKILL.md",
        "skill\spec-format\scripts\validate-change-spec.ts"
    )

    foreach ($file in $files) {
        $path = Join-Path $Target $file
        if (Remove-FileIfExists -Path $path) {
            $removed++
        }
    }

    # Clean up empty directories (leaf to root)
    $dirs = @(
        "skill\design-case-study-generator\references",
        "skill\design-case-study-generator\scripts",
        "skill\design-case-study-generator",
        "skill\merge-change-specs\references",
        "skill\merge-change-specs\scripts",
        "skill\merge-change-specs",
        "skill\skill-creator\references",
        "skill\skill-creator",
        "skill\spec-format\scripts",
        "skill\spec-format",
        "skill\sdd-state-management",
        "skill\research",
        "skill\architecture-fit-check",
        "skill\architecture-workshop",
        "skill\bun-shell-commands",
        "skill\keep-current",
        "skill",
        "plugin",
        "command\tool",
        "command\sdd\tools",
        "command\sdd\fast",
        "command\sdd",
        "command\create",
        "command",
        "agent\sdd",
        "agent"
    )

    foreach ($dir in $dirs) {
        Remove-EmptyDirectory -Path (Join-Path $Target $dir)
    }

    return $removed
}

# Remove Augment SDD files from a target directory
function Remove-AugmentFiles {
    param($Target)
    $removed = 0

    # Augment file list - exact files installed (keep in sync with repo contents)
    $files = @(
        "agents\librarian.md",
        "commands\sdd\init.md",
        "commands\sdd\plan.md",
        "commands\sdd\implement.md",
        "commands\sdd\proposal.md",
        "commands\sdd\specs.md",
        "commands\sdd\discovery.md",
        "commands\sdd\tasks.md",
        "commands\sdd\reconcile.md",
        "commands\sdd\finish.md",
        "commands\sdd\status.md",
        "commands\sdd\brainstorm.md",
        "commands\sdd\explain.md",
        "commands\sdd\fast\vibe.md",
        "commands\sdd\fast\bug.md",
        "commands\sdd\tools\critique.md",
        "commands\sdd\tools\scenario-test.md",
        "commands\sdd\tools\taxonomy-map.md",
        "commands\sdd\tools\prime-specs.md",
        "commands\sdd\skill-discovery.md",
        "commands\tool\commit.md",
        "commands\create\agent.md",
        "commands\create\command.md",
        "skills\sdd-state-management.md",
        "skills\spec-format.md",
        "skills\research.md",
        "skills\architecture-fit-check.md",
        "skills\architecture-workshop.md",
        "skills\keep-current.md"
    )

    foreach ($file in $files) {
        $path = Join-Path $Target $file
        if (Remove-FileIfExists -Path $path) {
            $removed++
        }
    }

    # Clean up empty directories (leaf to root)
    $dirs = @(
        "skills",
        "commands\tool",
        "commands\sdd\tools",
        "commands\sdd\fast",
        "commands\sdd",
        "commands\create",
        "commands",
        "agents"
    )

    foreach ($dir in $dirs) {
        Remove-EmptyDirectory -Path (Join-Path $Target $dir)
    }

    return $removed
}

# Remove Codex SDD files from a target directory
function Remove-CodexFiles {
    param($Target)
    $removed = 0

    # Codex file list - exact files installed (keep in sync with repo contents)
    $files = @(
        "prompts\create-command.md",
        "prompts\sdd-brainstorm.md",
        "prompts\sdd-discovery.md",
        "prompts\sdd-explain.md",
        "prompts\sdd-fast-bug.md",
        "prompts\sdd-fast-vibe.md",
        "prompts\sdd-finish.md",
        "prompts\sdd-implement.md",
        "prompts\sdd-init.md",
        "prompts\sdd-plan.md",
        "prompts\sdd-proposal.md",
        "prompts\sdd-reconcile.md",
        "prompts\sdd-specs.md",
        "prompts\sdd-status.md",
        "prompts\sdd-tasks.md",
        "prompts\sdd-tools-critique.md",
        "prompts\sdd-tools-prime-specs.md",
        "prompts\sdd-tools-scenario-test.md",
        "prompts\sdd-tools-taxonomy-map.md",
        "prompts\tool-commit.md",
        "skills\architecture-fit-check\SKILL.md",
        "skills\architecture-workshop\SKILL.md",
        "skills\bun-shell-commands\SKILL.md",
        "skills\design-case-study-generator\SKILL.md",
        "skills\design-case-study-generator\references\case-study-template.md",
        "skills\design-case-study-generator\references\demo-skeleton.md",
        "skills\design-case-study-generator\references\pitch-card-template.md",
        "skills\design-case-study-generator\references\token-schema-guidance.md",
        "skills\design-case-study-generator\references\tokens-css-emitter.md",
        "skills\design-case-study-generator\scripts\copy-version.ts",
        "skills\keep-current\SKILL.md",
        "skills\merge-change-specs\SKILL.md",
        "skills\merge-change-specs\references\delta-merge-rules.md",
        "skills\merge-change-specs\scripts\merge-change-specs.ts",
        "skills\research\SKILL.md",
        "skills\sdd-state-management\SKILL.md",
        "skills\skill-creator\SKILL.md",
        "skills\skill-creator\references\references-guide.md",
        "skills\skill-creator\references\scripts-overview.md",
        "skills\skill-creator\references\scripts-runtime-node-first.md",
        "skills\spec-format\SKILL.md",
        "skills\spec-format\scripts\validate-change-spec.ts"
    )

    foreach ($file in $files) {
        $path = Join-Path $Target $file
        if (Remove-FileIfExists -Path $path) {
            $removed++
        }
    }

    # Clean up empty directories (leaf to root)
    $dirs = @(
        "skills\design-case-study-generator\references",
        "skills\design-case-study-generator\scripts",
        "skills\design-case-study-generator",
        "skills\merge-change-specs\references",
        "skills\merge-change-specs\scripts",
        "skills\merge-change-specs",
        "skills\skill-creator\references",
        "skills\skill-creator",
        "skills\spec-format\scripts",
        "skills\spec-format",
        "skills\sdd-state-management",
        "skills\research",
        "skills\architecture-fit-check",
        "skills\architecture-workshop",
        "skills\bun-shell-commands",
        "skills\keep-current",
        "skills",
        "prompts"
    )

    foreach ($dir in $dirs) {
        Remove-EmptyDirectory -Path (Join-Path $Target $dir)
    }

    return $removed
}

# Check if OpenCode installation exists
function Test-OpenCodeExists {
    param($Target)
    return (Test-Path "$Target\agent\sdd") -or
        (Test-Path "$Target\command\sdd") -or
        (Test-Path "$Target\agent\librarian.md") -or
        (Test-Path "$Target\plugin\spec-validate.ts") -or
        (Test-Path "$Target\skill\merge-change-specs")
}

# Check if Augment installation exists
function Test-AugmentExists {
    param($Target)
    return (Test-Path "$Target\commands\sdd") -or (Test-Path "$Target\agents\librarian.md")
}

# Check if Codex installation exists
function Test-CodexExists {
    param($Target)
    return (Test-Path "$Target\prompts") -or (Test-Path "$Target\skills")
}

# Main
function Main {
    Write-Host ""
    Write-Host "  Agent Extensions Uninstaller"
    Write-Host "  ============================"
    Write-Host ""

    # Determine possible locations
    $OpenCodeGlobal = Join-Path $HOME ".config\opencode"
    $AugmentGlobal = Join-Path $HOME ".augment"
    $CodexGlobal = Join-Path $HOME ".codex"
    $LocalRoot = Find-GitRoot

    # Check what installations exist
    $OpenCodeGlobalExists = Test-OpenCodeExists -Target $OpenCodeGlobal
    $OpenCodeLocalExists = $false
    $AugmentGlobalExists = Test-AugmentExists -Target $AugmentGlobal
    $AugmentLocalExists = $false
    $CodexGlobalExists = Test-CodexExists -Target $CodexGlobal
    $CodexLocalExists = $false

    if ($LocalRoot) {
        $OpenCodeLocalExists = Test-OpenCodeExists -Target (Join-Path $LocalRoot ".opencode")
        $AugmentLocalExists = Test-AugmentExists -Target (Join-Path $LocalRoot ".augment")
        $CodexLocalExists = Test-CodexExists -Target (Join-Path $LocalRoot ".codex")
    }

    # Build found list
    $FoundAny = $false
    Write-Host "Found installations:"
    Write-Host ""

    if ($OpenCodeGlobalExists) {
        Write-Host "  - OpenCode (global): $OpenCodeGlobal"
        $FoundAny = $true
    }
    if ($OpenCodeLocalExists) {
        Write-Host "  - OpenCode (local): $LocalRoot\.opencode"
        $FoundAny = $true
    }
    if ($AugmentGlobalExists) {
        Write-Host "  - Augment (global): $AugmentGlobal"
        $FoundAny = $true
    }
    if ($AugmentLocalExists) {
        Write-Host "  - Augment (local): $LocalRoot\.augment"
        $FoundAny = $true
    }
    if ($CodexGlobalExists) {
        Write-Host "  - Codex (global): $CodexGlobal"
        $FoundAny = $true
    }
    if ($CodexLocalExists) {
        Write-Host "  - Codex (local): $LocalRoot\.codex"
        $FoundAny = $true
    }

    if (-not $FoundAny) {
        Write-Host "  (none)"
        Write-Host ""
        exit 0
    }

    Write-Host ""

    # Ask which tool to uninstall
    Write-Host "Which tool(s) would you like to uninstall?"
    Write-Host "  1) OpenCode"
    Write-Host "  2) Augment"
    Write-Host "  3) Codex"
    Write-Host "  4) OpenCode + Augment"
    Write-Host "  5) OpenCode + Codex"
    Write-Host "  6) Augment + Codex"
    Write-Host "  7) All"
    Write-Host ""
    $toolChoice = Read-Host "Enter choice [1/2/3/4/5/6/7]"

    $UninstallOpenCode = $false
    $UninstallAugment = $false
    $UninstallCodex = $false

    switch ($toolChoice) {
        "1" { $UninstallOpenCode = $true }
        "2" { $UninstallAugment = $true }
        "3" { $UninstallCodex = $true }
        "4" { $UninstallOpenCode = $true; $UninstallAugment = $true }
        "5" { $UninstallOpenCode = $true; $UninstallCodex = $true }
        "6" { $UninstallAugment = $true; $UninstallCodex = $true }
        "7" { $UninstallOpenCode = $true; $UninstallAugment = $true; $UninstallCodex = $true }
        default { Write-Err "Invalid choice" }
    }

    # Ask which scope to uninstall
    Write-Host ""
    Write-Host "Which installation scope?"
    Write-Host "  1) Global only"
    Write-Host "  2) Local only"
    Write-Host "  3) Both"
    Write-Host ""
    $scopeChoice = Read-Host "Enter choice [1/2/3]"

    $UninstallGlobal = $false
    $UninstallLocal = $false

    switch ($scopeChoice) {
        "1" { $UninstallGlobal = $true }
        "2" { $UninstallLocal = $true }
        "3" { $UninstallGlobal = $true; $UninstallLocal = $true }
        default { Write-Err "Invalid choice" }
    }

    Write-Host ""
    Write-Warn "This will remove SDD agents, commands, and skills."
    Write-Warn "Your project files (changes/, specs/, etc.) will NOT be affected."
    Write-Host ""

    $confirm = Read-Host "Proceed with uninstall? [y/N]"
    if ($confirm -notmatch "^[Yy]") {
        Write-Host "Aborted."
        exit 0
    }

    Write-Host ""

    $TotalRemoved = 0

    # Uninstall OpenCode
    if ($UninstallOpenCode) {
        if ($UninstallGlobal -and $OpenCodeGlobalExists) {
            Write-Info "Removing OpenCode (global)..."
            $count = Remove-OpenCodeFiles -Target $OpenCodeGlobal
            Write-Success "Removed $count files from $OpenCodeGlobal"
            $TotalRemoved += $count
        }

        if ($UninstallLocal -and $OpenCodeLocalExists) {
            Write-Info "Removing OpenCode (local)..."
            $localPath = Join-Path $LocalRoot ".opencode"
            $count = Remove-OpenCodeFiles -Target $localPath
            Write-Success "Removed $count files from $localPath"
            $TotalRemoved += $count
        }
    }

    # Uninstall Augment
    if ($UninstallAugment) {
        if ($UninstallGlobal -and $AugmentGlobalExists) {
            Write-Info "Removing Augment (global)..."
            $count = Remove-AugmentFiles -Target $AugmentGlobal
            Write-Success "Removed $count files from $AugmentGlobal"
            $TotalRemoved += $count
        }

        if ($UninstallLocal -and $AugmentLocalExists) {
            Write-Info "Removing Augment (local)..."
            $localPath = Join-Path $LocalRoot ".augment"
            $count = Remove-AugmentFiles -Target $localPath
            Write-Success "Removed $count files from $localPath"
            $TotalRemoved += $count
        }
    }

    # Uninstall Codex
    if ($UninstallCodex) {
        if ($UninstallGlobal -and $CodexGlobalExists) {
            Write-Info "Removing Codex (global)..."
            $count = Remove-CodexFiles -Target $CodexGlobal
            Write-Success "Removed $count files from $CodexGlobal"
            $TotalRemoved += $count
        }

        if ($UninstallLocal -and $CodexLocalExists) {
            Write-Info "Removing Codex (local)..."
            $localPath = Join-Path $LocalRoot ".codex"
            $count = Remove-CodexFiles -Target $localPath
            Write-Success "Removed $count files from $localPath"
            $TotalRemoved += $count
        }
    }

    Write-Host ""
    if ($TotalRemoved -gt 0) {
        Write-Success "Uninstall complete! Removed $TotalRemoved files."
    } else {
        Write-Warn "No files were removed (installations may not exist at selected locations)."
    }
    Write-Host ""
}

Main
