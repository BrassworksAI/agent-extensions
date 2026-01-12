#!/bin/sh
set -e

# Agent Extensions Install (macOS/Linux)
# Creates global symlinks and local copies of the repo's folders
# Local copies avoid repo-specific symlinks checked into git

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

  if [ ! -d "$PAYLOAD_DIR" ]; then
    warn "Payload directory '$PAYLOAD_DIR' not found, skipping $LABEL"
    return 1
  fi

  info "Installing $LABEL to: $TARGET_ROOT"

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
    warn "Found $CONFLICT_COUNT conflicting file(s)/symlink(s) for $LABEL:"
    echo ""
    printf "$CONFLICTS" | head -20
    if [ $CONFLICT_COUNT -gt 20 ]; then
      echo "  ... and $((CONFLICT_COUNT - 20)) more"
    fi
    echo ""
    warn "These will be replaced with symlinks to the repo."
    echo ""
    if ! confirm "Replace ALL conflicting files with symlinks for $LABEL?"; then
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

  success "Symlinks created for $LABEL at: $TARGET_ROOT"
  return 0
}

# Install copies to a target directory
# Arguments: $1 = target root, $2 = payload dir, $3 = label for messages
install_copies() {
  TARGET_ROOT="$1"
  PAYLOAD_DIR="$2"
  LABEL="$3"

  if [ ! -d "$PAYLOAD_DIR" ]; then
    warn "Payload directory '$PAYLOAD_DIR' not found, skipping $LABEL"
    return 1
  fi

  info "Installing $LABEL to: $TARGET_ROOT"

  # Build file list and detect conflicts
  CONFLICTS=""
  CONFLICT_COUNT=0

  cd "$PAYLOAD_DIR"
  FILES="$(find . -type f | sed 's|^\./||')"
  cd - > /dev/null

  for file in $FILES; do
    dest="$TARGET_ROOT/$file"
    if [ -e "$dest" ] || [ -L "$dest" ]; then
      CONFLICTS="$CONFLICTS$file\n"
      CONFLICT_COUNT=$((CONFLICT_COUNT + 1))
    fi
  done

  if [ "$CONFLICT_COUNT" -gt 0 ]; then
    warn "Found $CONFLICT_COUNT existing file(s) that will be overwritten:"
    printf "%b" "$CONFLICTS"
    echo ""
    if ! confirm "Continue and overwrite these files?"; then
      warn "Installation cancelled."
      return 1
    fi
  fi

  # Create copies
  info "Copying files for $LABEL..."

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

    # Copy file
    cp -p "$src" "$dest" || error "Failed to copy $file"
  done

  success "Files copied for $LABEL at: $TARGET_ROOT"
  return 0
}

# Main
main() {
  echo ""
  echo "  Agent Extensions Install"
  echo "  ========================"
  echo ""

  # Determine script location to find repo root
  SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

  info "Source directory: $SCRIPT_DIR"
  echo ""

  # Choose tool
  echo "Which tool(s) would you like to install extensions for?"
  echo "  1) OpenCode"
  echo "  2) Augment (Auggie)"
  echo "  3) Both"
  echo ""
  printf "Enter choice [1/2/3]: "
  read -r tool_choice

  INSTALL_OPENCODE=""
  INSTALL_AUGMENT=""

  case "$tool_choice" in
    1)
      INSTALL_OPENCODE="yes"
      ;;
    2)
      INSTALL_AUGMENT="yes"
      ;;
    3)
      INSTALL_OPENCODE="yes"
      INSTALL_AUGMENT="yes"
      ;;
    *)
      error "Invalid choice. Please enter 1, 2, or 3."
      ;;
  esac

  # Choose install mode
  echo ""
  echo "Where would you like to install?"
  echo "  1) Global (user config directory)"
  echo "  2) Local (current repo)"
  echo "  3) Both global and local"
  echo ""
  printf "Enter choice [1/2/3]: "
  read -r scope_choice

  INSTALL_GLOBAL=""
  INSTALL_LOCAL=""

  case "$scope_choice" in
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
  fi

  echo ""

  INSTALLED_COUNT=0

  # Install OpenCode globally if requested
  if [ -n "$INSTALL_OPENCODE" ] && [ -n "$INSTALL_GLOBAL" ]; then
    OPENCODE_PAYLOAD="$SCRIPT_DIR/opencode"
    OPENCODE_TARGET="$HOME/.config/opencode"
    if install_symlinks "$OPENCODE_TARGET" "$OPENCODE_PAYLOAD" "OpenCode (global)"; then
      INSTALLED_COUNT=$((INSTALLED_COUNT + 1))
    fi
    echo ""
  fi

  # Install OpenCode locally if requested (copy)
  if [ -n "$INSTALL_OPENCODE" ] && [ -n "$INSTALL_LOCAL" ]; then
    OPENCODE_PAYLOAD="$SCRIPT_DIR/opencode"
    OPENCODE_TARGET="$GIT_ROOT/.opencode"
    if install_copies "$OPENCODE_TARGET" "$OPENCODE_PAYLOAD" "OpenCode (local)"; then
      INSTALLED_COUNT=$((INSTALLED_COUNT + 1))
    fi
    echo ""
  fi

  # Install Augment globally if requested
  if [ -n "$INSTALL_AUGMENT" ] && [ -n "$INSTALL_GLOBAL" ]; then
    AUGMENT_PAYLOAD="$SCRIPT_DIR/augment"
    AUGMENT_TARGET="$HOME/.augment"
    if install_symlinks "$AUGMENT_TARGET" "$AUGMENT_PAYLOAD" "Augment (global)"; then
      INSTALLED_COUNT=$((INSTALLED_COUNT + 1))
    fi
    echo ""
  fi

  # Install Augment locally if requested (copy)
  if [ -n "$INSTALL_AUGMENT" ] && [ -n "$INSTALL_LOCAL" ]; then
    AUGMENT_PAYLOAD="$SCRIPT_DIR/augment"
    AUGMENT_TARGET="$GIT_ROOT/.augment"
    if install_copies "$AUGMENT_TARGET" "$AUGMENT_PAYLOAD" "Augment (local)"; then
      INSTALLED_COUNT=$((INSTALLED_COUNT + 1))
    fi
    echo ""
  fi

  if [ $INSTALLED_COUNT -eq 0 ]; then
    error "No installations completed"
  fi

  success "Install complete!"
  echo ""
  info "Source files at: $SCRIPT_DIR"
  echo ""
  echo "Global installs use symlinks to the repo."
  echo "Local installs copy files into the repo."
  echo ""
}

main "$@"
