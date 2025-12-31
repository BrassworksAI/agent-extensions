#!/bin/sh
set -e

# OpenCode SDD Dev Install (macOS/Linux only)
# Creates symlinks to the local repo's opencode/ folder
# This allows you to edit files and push changes back to the repo

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

# Find git root by walking up directories (no git required)
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

# Install symlinks to a target directory
# Arguments: $1 = target root, $2 = payload dir, $3 = label for messages
install_symlinks() {
  TARGET_ROOT="$1"
  PAYLOAD_DIR="$2"
  LABEL="$3"

  info "Installing to $LABEL: $TARGET_ROOT"

  # Build file list and detect conflicts
  CONFLICTS=""
  CONFLICT_COUNT=0

  cd "$PAYLOAD_DIR"
  FILES="$(find . -type f | sed 's|^\./||')"
  cd - > /dev/null

  for file in $FILES; do
    dest="$TARGET_ROOT/$file"
    if [ -e "$dest" ] || [ -L "$dest" ]; then
      # Check if it's already a symlink to our source
      src="$PAYLOAD_DIR/$file"
      if [ -L "$dest" ]; then
        existing_target="$(readlink "$dest")"
        if [ "$existing_target" = "$src" ]; then
          # Already linked to us, not a conflict
          continue
        fi
      fi
      CONFLICTS="$CONFLICTS$file\n"
      CONFLICT_COUNT=$((CONFLICT_COUNT + 1))
    fi
  done

  # Handle conflicts
  if [ $CONFLICT_COUNT -gt 0 ]; then
    echo ""
    warn "Found $CONFLICT_COUNT conflicting file(s)/symlink(s) in $LABEL:"
    echo ""
    printf "$CONFLICTS" | head -20
    if [ $CONFLICT_COUNT -gt 20 ]; then
      echo "  ... and $((CONFLICT_COUNT - 20)) more"
    fi
    echo ""
    warn "These will be replaced with symlinks to the repo."
    echo ""
    if ! confirm "Replace ALL conflicting files with symlinks?"; then
      warn "Skipping $LABEL install"
      return 1
    fi
    echo ""
  fi

  # Create symlinks
  info "Creating symlinks for $LABEL..."

  for file in $FILES; do
    src="$PAYLOAD_DIR/$file"
    dest="$TARGET_ROOT/$file"
    dest_dir="$(dirname "$dest")"

    # Create parent directories if needed
    if [ ! -d "$dest_dir" ]; then
      mkdir -p "$dest_dir"
    fi

    # Remove existing file/symlink if present
    if [ -e "$dest" ] || [ -L "$dest" ]; then
      rm "$dest"
    fi

    # Create symlink
    ln -s "$src" "$dest" || error "Failed to create symlink for $file"
  done

  success "Symlinks created at: $TARGET_ROOT"
  return 0
}

# Main
main() {
  echo ""
  echo "  OpenCode SDD Dev Install (Symlink Mode)"
  echo "  ========================================"
  echo ""

  # Determine script location to find repo root
  SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
  PAYLOAD_DIR="$SCRIPT_DIR/opencode"

  if [ ! -d "$PAYLOAD_DIR" ]; then
    error "Cannot find 'opencode/' directory at $SCRIPT_DIR"
  fi

  info "Source: $PAYLOAD_DIR"
  echo ""

  # Choose install mode
  echo "Where would you like to install symlinks?"
  echo "  1) Global (~/.config/opencode)"
  echo "  2) Local (current repo's .opencode folder)"
  echo "  3) Both global and local"
  echo ""
  printf "Enter choice [1/2/3]: "
  read -r choice

  INSTALL_GLOBAL=""
  INSTALL_LOCAL=""

  case "$choice" in
    1)
      INSTALL_GLOBAL="yes"
      ;;
    2)
      INSTALL_LOCAL="yes"
      ;;
    3)
      INSTALL_GLOBAL="yes"
      INSTALL_LOCAL="yes"
      ;;
    *)
      error "Invalid choice. Please enter 1, 2, or 3."
      ;;
  esac

  # Validate local install is possible
  if [ -n "$INSTALL_LOCAL" ]; then
    GIT_ROOT="$(find_git_root)" || error "Not inside a git repository. Cannot determine repo root for local install."
    LOCAL_TARGET="$GIT_ROOT/.opencode"
  fi

  echo ""

  INSTALLED_COUNT=0

  # Install globally if requested
  if [ -n "$INSTALL_GLOBAL" ]; then
    GLOBAL_TARGET="$HOME/.config/opencode"
    if install_symlinks "$GLOBAL_TARGET" "$PAYLOAD_DIR" "global"; then
      INSTALLED_COUNT=$((INSTALLED_COUNT + 1))
    fi
    echo ""
  fi

  # Install locally if requested
  if [ -n "$INSTALL_LOCAL" ]; then
    if install_symlinks "$LOCAL_TARGET" "$PAYLOAD_DIR" "local"; then
      INSTALLED_COUNT=$((INSTALLED_COUNT + 1))
    fi
    echo ""
  fi

  if [ $INSTALLED_COUNT -eq 0 ]; then
    error "No installations completed"
  fi

  success "Dev install complete!"
  echo ""
  info "Source files at: $PAYLOAD_DIR"
  echo ""
  echo "Any edits you make to the symlinked files will modify the repo files."
  echo "You can commit and push changes directly."
  echo ""
}

main "$@"
