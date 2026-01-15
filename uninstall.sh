#!/bin/sh
set -e

# Agent Extensions Uninstaller (macOS/Linux)
# Removes the Spec-Driven Development process for OpenCode, Augment, and/or Codex

# Colors (if terminal supports them)
if [ -t 1 ]; then
  RED='\033[0;31m'
  GREEN='\033[0;32m'
  YELLOW='\033[0;33m'
  BLUE='\033[0;34m'
  NC='\033[0m'
else
  RED=''
  GREEN=''
  YELLOW=''
  BLUE=''
  NC=''
fi

info() { printf "${BLUE}[INFO]${NC} %s\n" "$1"; }
success() { printf "${GREEN}[OK]${NC} %s\n" "$1"; }
warn() { printf "${YELLOW}[WARN]${NC} %s\n" "$1"; }
error() { printf "${RED}[ERROR]${NC} %s\n" "$1"; exit 1; }

# Prompt for yes/no
confirm() {
  printf "%s [y/N] " "$1"
  read -r answer
  case "$answer" in
    [Yy]|[Yy][Ee][Ss]) return 0 ;;
    *) return 1 ;;
  esac
}

# Find git root by walking up directories
find_git_root() {
  dir="$PWD"
  while [ "$dir" != "/" ]; do
    if [ -d "$dir/.git" ] || [ -f "$dir/.git" ]; then
      echo "$dir"
      return 0
    fi
    dir="$(dirname "$dir")"
  done
  return 1
}

# Remove a file or symlink if it exists
remove_file() {
  path="$1"
  if [ -e "$path" ] || [ -L "$path" ]; then
    rm -f "$path"
    return 0
  fi
  return 1
}

# Remove a directory only if it's empty
remove_empty_dir() {
  dir="$1"
  if [ -d "$dir" ] && [ -z "$(ls -A "$dir" 2>/dev/null)" ]; then
    rmdir "$dir" 2>/dev/null || true
  fi
}

# Remove OpenCode SDD files from a target directory
remove_opencode_files() {
  target="$1"
  removed=0

  # OpenCode file list - exact files installed (keep in sync with repo contents)
  files="
    agent/sdd/plan.md
    agent/sdd/build.md
    agent/librarian.md
    agent/chat.md
    agent/designer.md
    command/sdd/init.md
    command/sdd/plan.md
    command/sdd/implement.md
    command/sdd/proposal.md
    command/sdd/specs.md
    command/sdd/discovery.md
    command/sdd/tasks.md
    command/sdd/reconcile.md
    command/sdd/finish.md
    command/sdd/status.md
    command/sdd/brainstorm.md
    command/sdd/explain.md
    command/sdd/fast/vibe.md
    command/sdd/fast/bug.md
    command/sdd/tools/critique.md
    command/sdd/tools/scenario-test.md
    command/sdd/tools/taxonomy-map.md
    command/sdd/tools/prime-specs.md
    command/tool/commit.md
    command/create/agent.md
    command/create/command.md
    plugin/spec-validate.ts
    skill/architecture-fit-check/SKILL.md
    skill/architecture-workshop/SKILL.md
    skill/bun-shell-commands/SKILL.md
    skill/design-case-study-generator/SKILL.md
    skill/design-case-study-generator/references/case-study-template.md
    skill/design-case-study-generator/references/demo-skeleton.md
    skill/design-case-study-generator/references/pitch-card-template.md
    skill/design-case-study-generator/references/token-schema-guidance.md
    skill/design-case-study-generator/references/tokens-css-emitter.md
    skill/design-case-study-generator/scripts/copy-version.ts
    skill/keep-current/SKILL.md
    skill/merge-change-specs/SKILL.md
    skill/merge-change-specs/references/delta-merge-rules.md
    skill/merge-change-specs/scripts/merge-change-specs.ts
    skill/research/SKILL.md
    skill/sdd-state-management/SKILL.md
    skill/skill-creator/SKILL.md
    skill/skill-creator/references/references-guide.md
    skill/skill-creator/references/scripts-overview.md
    skill/skill-creator/references/scripts-runtime-node-first.md
    skill/spec-format/SKILL.md
    skill/spec-format/scripts/validate-change-spec.ts
  "

  for file in $files; do
    if remove_file "$target/$file"; then
      removed=$((removed + 1))
    fi
  done

  # Clean up empty directories (leaf to root)
  dirs="
    skill/design-case-study-generator/references
    skill/design-case-study-generator/scripts
    skill/design-case-study-generator
    skill/merge-change-specs/references
    skill/merge-change-specs/scripts
    skill/merge-change-specs
    skill/skill-creator/references
    skill/skill-creator
    skill/spec-format/scripts
    skill/spec-format
    skill/sdd-state-management
    skill/research
    skill/architecture-fit-check
    skill/architecture-workshop
    skill/bun-shell-commands
    skill/keep-current
    skill
    plugin
    command/tool
    command/sdd/tools
    command/sdd/fast
    command/sdd
    command/create
    command
    agent/sdd
    agent
  "

  for dir in $dirs; do
    remove_empty_dir "$target/$dir"
  done

  echo "$removed"
}

# Remove Augment SDD files from a target directory
remove_augment_files() {
  target="$1"
  removed=0

  # Augment file list - exact files installed (keep in sync with repo contents)
  files="
    agents/librarian.md
    commands/sdd/init.md
    commands/sdd/plan.md
    commands/sdd/implement.md
    commands/sdd/proposal.md
    commands/sdd/specs.md
    commands/sdd/discovery.md
    commands/sdd/tasks.md
    commands/sdd/reconcile.md
    commands/sdd/finish.md
    commands/sdd/status.md
    commands/sdd/brainstorm.md
    commands/sdd/explain.md
    commands/sdd/fast/vibe.md
    commands/sdd/fast/bug.md
    commands/sdd/tools/critique.md
    commands/sdd/tools/scenario-test.md
    commands/sdd/tools/taxonomy-map.md
    commands/sdd/tools/prime-specs.md
    commands/sdd/skill-discovery.md
    commands/tool/commit.md
    commands/create/agent.md
    commands/create/command.md
    skills/sdd-state-management.md
    skills/spec-format.md
    skills/research.md
    skills/architecture-fit-check.md
    skills/architecture-workshop.md
    skills/keep-current.md
  "

  for file in $files; do
    if remove_file "$target/$file"; then
      removed=$((removed + 1))
    fi
  done

  # Clean up empty directories (leaf to root)
  dirs="
    skills
    commands/tool
    commands/sdd/tools
    commands/sdd/fast
    commands/sdd
    commands/create
    commands
    agents
  "

  for dir in $dirs; do
    remove_empty_dir "$target/$dir"
  done

  echo "$removed"
}

# Remove Codex SDD files from a target directory
remove_codex_files() {
  target="$1"
  removed=0

  # Codex file list - exact files installed (keep in sync with repo contents)
  files="
    prompts/create-command.md
    prompts/sdd-brainstorm.md
    prompts/sdd-discovery.md
    prompts/sdd-explain.md
    prompts/sdd-fast-bug.md
    prompts/sdd-fast-vibe.md
    prompts/sdd-finish.md
    prompts/sdd-implement.md
    prompts/sdd-init.md
    prompts/sdd-plan.md
    prompts/sdd-proposal.md
    prompts/sdd-reconcile.md
    prompts/sdd-specs.md
    prompts/sdd-status.md
    prompts/sdd-tasks.md
    prompts/sdd-tools-critique.md
    prompts/sdd-tools-prime-specs.md
    prompts/sdd-tools-scenario-test.md
    prompts/sdd-tools-taxonomy-map.md
    prompts/tool-commit.md
    skills/architecture-fit-check/SKILL.md
    skills/architecture-workshop/SKILL.md
    skills/bun-shell-commands/SKILL.md
    skills/design-case-study-generator/SKILL.md
    skills/design-case-study-generator/references/case-study-template.md
    skills/design-case-study-generator/references/demo-skeleton.md
    skills/design-case-study-generator/references/pitch-card-template.md
    skills/design-case-study-generator/references/token-schema-guidance.md
    skills/design-case-study-generator/references/tokens-css-emitter.md
    skills/design-case-study-generator/scripts/copy-version.ts
    skills/keep-current/SKILL.md
    skills/merge-change-specs/SKILL.md
    skills/merge-change-specs/references/delta-merge-rules.md
    skills/merge-change-specs/scripts/merge-change-specs.ts
    skills/research/SKILL.md
    skills/sdd-state-management/SKILL.md
    skills/skill-creator/SKILL.md
    skills/skill-creator/references/references-guide.md
    skills/skill-creator/references/scripts-overview.md
    skills/skill-creator/references/scripts-runtime-node-first.md
    skills/spec-format/SKILL.md
    skills/spec-format/scripts/validate-change-spec.ts
  "

  for file in $files; do
    if remove_file "$target/$file"; then
      removed=$((removed + 1))
    fi
  done

  # Clean up empty directories (leaf to root)
  dirs="
    skills/design-case-study-generator/references
    skills/design-case-study-generator/scripts
    skills/design-case-study-generator
    skills/merge-change-specs/references
    skills/merge-change-specs/scripts
    skills/merge-change-specs
    skills/skill-creator/references
    skills/skill-creator
    skills/spec-format/scripts
    skills/spec-format
    skills/sdd-state-management
    skills/research
    skills/architecture-fit-check
    skills/architecture-workshop
    skills/bun-shell-commands
    skills/keep-current
    skills
    prompts
  "

  for dir in $dirs; do
    remove_empty_dir "$target/$dir"
  done

  echo "$removed"
}

# Check if OpenCode installation exists
check_opencode_exists() {
  target="$1"
  [ -d "$target/agent/sdd" ] || \
    [ -d "$target/command/sdd" ] || \
    [ -f "$target/agent/librarian.md" ] || \
    [ -f "$target/plugin/spec-validate.ts" ] || \
    [ -d "$target/skill/merge-change-specs" ]
}

# Check if Augment installation exists
check_augment_exists() {
  target="$1"
  [ -d "$target/commands/sdd" ] || [ -f "$target/agents/librarian.md" ]
}

# Check if Codex installation exists
check_codex_exists() {
  target="$1"
  [ -d "$target/prompts" ] || [ -d "$target/skills" ]
}

# Main
main() {
  echo ""
  echo "  Agent Extensions Uninstaller"
  echo "  ============================"
  echo ""

  # Determine possible locations
  OPENCODE_GLOBAL="$HOME/.config/opencode"
  AUGMENT_GLOBAL="$HOME/.augment"
  CODEX_GLOBAL="$HOME/.codex"
  LOCAL_ROOT=""
  
  if GIT_ROOT="$(find_git_root 2>/dev/null)"; then
    LOCAL_ROOT="$GIT_ROOT"
  fi

  # Check what installations exist
  OPENCODE_GLOBAL_EXISTS=false
  OPENCODE_LOCAL_EXISTS=false
  AUGMENT_GLOBAL_EXISTS=false
  AUGMENT_LOCAL_EXISTS=false
  CODEX_GLOBAL_EXISTS=false
  CODEX_LOCAL_EXISTS=false

  if check_opencode_exists "$OPENCODE_GLOBAL"; then
    OPENCODE_GLOBAL_EXISTS=true
  fi

  if [ -n "$LOCAL_ROOT" ] && check_opencode_exists "$LOCAL_ROOT/.opencode"; then
    OPENCODE_LOCAL_EXISTS=true
  fi

  if check_augment_exists "$AUGMENT_GLOBAL"; then
    AUGMENT_GLOBAL_EXISTS=true
  fi

  if [ -n "$LOCAL_ROOT" ] && check_augment_exists "$LOCAL_ROOT/.augment"; then
    AUGMENT_LOCAL_EXISTS=true
  fi

  if check_codex_exists "$CODEX_GLOBAL"; then
    CODEX_GLOBAL_EXISTS=true
  fi

  if [ -n "$LOCAL_ROOT" ] && check_codex_exists "$LOCAL_ROOT/.codex"; then
    CODEX_LOCAL_EXISTS=true
  fi

  # Build found list
  FOUND_ANY=false
  echo "Found installations:"
  echo ""
  
  if [ "$OPENCODE_GLOBAL_EXISTS" = true ]; then
    echo "  - OpenCode (global): $OPENCODE_GLOBAL"
    FOUND_ANY=true
  fi
  if [ "$OPENCODE_LOCAL_EXISTS" = true ]; then
    echo "  - OpenCode (local): $LOCAL_ROOT/.opencode"
    FOUND_ANY=true
  fi
  if [ "$AUGMENT_GLOBAL_EXISTS" = true ]; then
    echo "  - Augment (global): $AUGMENT_GLOBAL"
    FOUND_ANY=true
  fi
  if [ "$AUGMENT_LOCAL_EXISTS" = true ]; then
    echo "  - Augment (local): $LOCAL_ROOT/.augment"
    FOUND_ANY=true
  fi
  if [ "$CODEX_GLOBAL_EXISTS" = true ]; then
    echo "  - Codex (global): $CODEX_GLOBAL"
    FOUND_ANY=true
  fi
  if [ "$CODEX_LOCAL_EXISTS" = true ]; then
    echo "  - Codex (local): $LOCAL_ROOT/.codex"
    FOUND_ANY=true
  fi

  if [ "$FOUND_ANY" = false ]; then
    echo "  (none)"
    echo ""
    exit 0
  fi

  echo ""

  # Ask which tool to uninstall
  echo "Which tool(s) would you like to uninstall?"
  echo "  1) OpenCode"
  echo "  2) Augment"
  echo "  3) Codex"
  echo "  4) OpenCode + Augment"
  echo "  5) OpenCode + Codex"
  echo "  6) Augment + Codex"
  echo "  7) All"
  echo ""
  printf "Enter choice [1/2/3/4/5/6/7]: "
  read -r tool_choice

  UNINSTALL_OPENCODE=""
  UNINSTALL_AUGMENT=""
  UNINSTALL_CODEX=""

  case "$tool_choice" in
    1) UNINSTALL_OPENCODE="yes" ;;
    2) UNINSTALL_AUGMENT="yes" ;;
    3) UNINSTALL_CODEX="yes" ;;
    4) UNINSTALL_OPENCODE="yes"; UNINSTALL_AUGMENT="yes" ;;
    5) UNINSTALL_OPENCODE="yes"; UNINSTALL_CODEX="yes" ;;
    6) UNINSTALL_AUGMENT="yes"; UNINSTALL_CODEX="yes" ;;
    7) UNINSTALL_OPENCODE="yes"; UNINSTALL_AUGMENT="yes"; UNINSTALL_CODEX="yes" ;;
    *) error "Invalid choice" ;;
  esac

  # Ask which scope to uninstall
  echo ""
  echo "Which installation scope?"
  echo "  1) Global only"
  echo "  2) Local only"
  echo "  3) Both"
  echo ""
  printf "Enter choice [1/2/3]: "
  read -r scope_choice

  UNINSTALL_GLOBAL=""
  UNINSTALL_LOCAL=""

  case "$scope_choice" in
    1) UNINSTALL_GLOBAL="yes" ;;
    2) UNINSTALL_LOCAL="yes" ;;
    3) UNINSTALL_GLOBAL="yes"; UNINSTALL_LOCAL="yes" ;;
    *) error "Invalid choice" ;;
  esac

  echo ""
  warn "This will remove SDD agents, commands, and skills."
  warn "Your project files (changes/, specs/, etc.) will NOT be affected."
  echo ""

  if ! confirm "Proceed with uninstall?"; then
    echo "Aborted."
    exit 0
  fi

  echo ""

  TOTAL_REMOVED=0

  # Uninstall OpenCode
  if [ -n "$UNINSTALL_OPENCODE" ]; then
    if [ -n "$UNINSTALL_GLOBAL" ] && [ "$OPENCODE_GLOBAL_EXISTS" = true ]; then
      info "Removing OpenCode (global)..."
      count=$(remove_opencode_files "$OPENCODE_GLOBAL")
      success "Removed $count files from $OPENCODE_GLOBAL"
      TOTAL_REMOVED=$((TOTAL_REMOVED + count))
    fi

    if [ -n "$UNINSTALL_LOCAL" ] && [ "$OPENCODE_LOCAL_EXISTS" = true ]; then
      info "Removing OpenCode (local)..."
      count=$(remove_opencode_files "$LOCAL_ROOT/.opencode")
      success "Removed $count files from $LOCAL_ROOT/.opencode"
      TOTAL_REMOVED=$((TOTAL_REMOVED + count))
    fi
  fi

  # Uninstall Augment
  if [ -n "$UNINSTALL_AUGMENT" ]; then
    if [ -n "$UNINSTALL_GLOBAL" ] && [ "$AUGMENT_GLOBAL_EXISTS" = true ]; then
      info "Removing Augment (global)..."
      count=$(remove_augment_files "$AUGMENT_GLOBAL")
      success "Removed $count files from $AUGMENT_GLOBAL"
      TOTAL_REMOVED=$((TOTAL_REMOVED + count))
    fi

    if [ -n "$UNINSTALL_LOCAL" ] && [ "$AUGMENT_LOCAL_EXISTS" = true ]; then
      info "Removing Augment (local)..."
      count=$(remove_augment_files "$LOCAL_ROOT/.augment")
      success "Removed $count files from $LOCAL_ROOT/.augment"
      TOTAL_REMOVED=$((TOTAL_REMOVED + count))
    fi
  fi

  # Uninstall Codex
  if [ -n "$UNINSTALL_CODEX" ]; then
    if [ -n "$UNINSTALL_GLOBAL" ] && [ "$CODEX_GLOBAL_EXISTS" = true ]; then
      info "Removing Codex (global)..."
      count=$(remove_codex_files "$CODEX_GLOBAL")
      success "Removed $count files from $CODEX_GLOBAL"
      TOTAL_REMOVED=$((TOTAL_REMOVED + count))
    fi

    if [ -n "$UNINSTALL_LOCAL" ] && [ "$CODEX_LOCAL_EXISTS" = true ]; then
      info "Removing Codex (local)..."
      count=$(remove_codex_files "$LOCAL_ROOT/.codex")
      success "Removed $count files from $LOCAL_ROOT/.codex"
      TOTAL_REMOVED=$((TOTAL_REMOVED + count))
    fi
  fi

  echo ""
  if [ $TOTAL_REMOVED -gt 0 ]; then
    success "Uninstall complete! Removed $TOTAL_REMOVED files."
  else
    warn "No files were removed (installations may not exist at selected locations)."
  fi
  echo ""
}

main "$@"
