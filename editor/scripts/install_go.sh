#!/usr/bin/env bash



set -e

BASHRC="$HOME/.bashrc"
GO_VERSION="${1:-1.24.5}"
GO_TARBALL="go${GO_VERSION}.linux-amd64.tar.gz"
GO_URL="https://go.dev/dl/${GO_TARBALL}"
GO_INSTALL_DIR="/usr/local/go"

source "$BASHRC"
echo "[*] Checking for Go ${GO_VERSION}..."

if command -v go >/dev/null 2>&1; then
    CUR_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    if [ "$CUR_VERSION" = "$GO_VERSION" ]; then
        echo "[✓] Go ${GO_VERSION} is already installed."
        exit 0
    else
        echo "[!] Found Go ${CUR_VERSION} — replacing with ${GO_VERSION}..."
    fi
else
    echo "[*] Go not found — installing ${GO_VERSION}..."
fi

# Download and install Go
wget "$GO_URL"
sudo rm -rf "$GO_INSTALL_DIR"
sudo tar -C /usr/local -xzf "$GO_TARBALL"
rm "$GO_TARBALL"
echo "[✓] Installed Go ${GO_VERSION}"

# Ensure PATH and GOPATH are set in ~/.bashrc
echo "[*] Updating ~/.bashrc..."



ensure_line() {
    grep -qxF "$1" "$BASHRC" || echo "$1" >> "$BASHRC"
}

ensure_line 'export PATH=$PATH:/usr/local/go/bin'
ensure_line 'export GOPATH=$HOME/go'
ensure_line 'export PATH=$PATH:$GOPATH/bin'

echo "[✓] ~/.bashrc updated. Run 'source ~/.bashrc' or restart your shell."
