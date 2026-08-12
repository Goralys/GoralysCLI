#!/bin/sh

# This script was generated using AI
# The model used was: Claude Sonnet 5 [reasoning: low]
# Might need a human rewrite one day
# Do not hesitate if you have the knowledge and faith to do it :)

set -eu

REPO_BASE="https://cli.goralys.fr/release"
BINARY_NAME="goralys-cli"
INSTALL_DIR="${GORALYS_INSTALL_DIR:-$HOME/.goralys/bin}"
VERSION="${1:-latest}"

# --- detect OS/arch ---
os="$(uname -s)"
arch="$(uname -m)"

case "$os" in
    Linux)  os="linux" ;;
    Darwin) os="darwin" ;;
    *) echo "Unsupported OS: $os" >&2; exit 1 ;;
esac

case "$arch" in
    x86_64|amd64) arch="amd64" ;;
    aarch64|arm64) arch="arm64" ;;
    *) echo "Unsupported architecture: $arch" >&2; exit 1 ;;
esac

asset="${BINARY_NAME}-${os}-${arch}"
url="${REPO_BASE}/${VERSION}/${asset}"

# --- download ---
mkdir -p "$INSTALL_DIR"
tmpfile="$(mktemp)"

echo "Downloading ${asset} (${VERSION})..."
if ! curl -fsSL "$url" -o "$tmpfile"; then
    echo "Failed to download from $url" >&2
    echo "Check that this platform/architecture is supported and the version exists." >&2
    rm -f "$tmpfile"
    exit 1
fi

chmod +x "$tmpfile"
mv "$tmpfile" "$INSTALL_DIR/$BINARY_NAME"

echo "Installed to $INSTALL_DIR/$BINARY_NAME"

# --- add to PATH ---
add_to_path() {
    rc_file="$1"

    if [ -f "$rc_file" ] && grep -qF "$INSTALL_DIR" "$rc_file" 2>/dev/null; then
        return
    fi

    printf '\n# Added by GoralysCLI installer\nexport PATH="$PATH:%s"\n' "$INSTALL_DIR" >> "$rc_file"
    echo "Added $INSTALL_DIR to PATH in $rc_file"
}

case "$(basename "${SHELL:-}")" in
    zsh)  add_to_path "$HOME/.zshrc" ;;
    bash) add_to_path "$HOME/.bashrc" ;;
    fish)
        mkdir -p "$HOME/.config/fish"
        if ! grep -qF "$INSTALL_DIR" "$HOME/.config/fish/config.fish" 2>/dev/null; then
            echo "fish_add_path $INSTALL_DIR" >> "$HOME/.config/fish/config.fish"
        fi
        ;;
    *) add_to_path "$HOME/.profile" ;;
esac

echo ""
echo "Done! Restart your shell or run:"
echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
echo "Then verify with: $BINARY_NAME --version"