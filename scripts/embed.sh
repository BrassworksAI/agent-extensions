#!/bin/sh
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
EMBED_DIR="$ROOT_DIR/internal/embedded/content"

rm -rf "$EMBED_DIR"
mkdir -p "$EMBED_DIR"

cp "$ROOT_DIR/tools.yaml" "$EMBED_DIR/"
cp -r "$ROOT_DIR/repository" "$EMBED_DIR/"

echo "Embedded content prepared in $EMBED_DIR"
