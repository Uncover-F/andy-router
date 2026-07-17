#!/bin/sh

# MIT License
#
# Copyright (c) 2026 Uncover-F
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

VERSION="1.2.1-beta.2"
REPO="https://github.com/Uncover-F/andy-router/releases/download/$VERSION"

INSTALL_DIR="$HOME/.local/bin"
BIN="$INSTALL_DIR/andy-router"

die() {
    printf "%s\n" "$@" >&2
    exit 1
}

check_bin() {
    command -v "$1" >/dev/null 2>&1
}

check_path() {
    case ":$1:" in
        (*":$HOME/.local/bin:"*) return 0 ;;
    esac
    return 1
}

download() {
    printf "Downloading andy-router...\n"

    curl -fsSL "$1" -o "$BIN.tmp" ||
        die "Failed to download binary."

    chmod +x "$BIN.tmp" ||
        die "Failed to make binary executable."

    mv "$BIN.tmp" "$BIN" ||
        die "Failed to install binary."
}

main() {
    printf "Installing andy-router...\n"

    [ -n "$HOME" ] || die "HOME is not set."
    check_bin curl || die "Please install curl."

    case "$(uname -s)" in
        Linux)  OS=linux ;;
        Darwin) OS=darwin ;;
        *) die "Unsupported operating system: $(uname -s)" ;;
    esac

    case "$(uname -m)" in
        x86_64|amd64) ARCH=amd64 ;;
        arm64|aarch64) ARCH=arm64 ;;
        *) die "Unsupported architecture: $(uname -m)" ;;
    esac

    printf "Detected OS: %s\n" "$OS"
    printf "Detected Architecture: %s\n" "$ARCH"

    URL="$REPO/andy-router-$OS-$ARCH"

    mkdir -p "$INSTALL_DIR" ||
        die "Couldn't create $INSTALL_DIR."

    download "$URL"

    printf "\nInstallation completed successfully.\n\n"

    if check_path "$PATH"; then
        cat <<EOF
Run the following command:

  andy-router

EOF
        exit 0
    fi

    LOGIN_SHELL="${SHELL:-/bin/sh}"
    LOGIN_PATH=$("$LOGIN_SHELL" -l -c 'printf %s "$PATH"' 2>/dev/null)

    if check_path "$LOGIN_PATH"; then
        cat <<EOF
Please restart your terminal or open a new shell.

Then run:

  andy-router

EOF
        exit 0
    fi

    RC_FILE=

    case "${SHELL##*/}" in
        bash) RC_FILE=".bash_profile" ;;
        zsh)  RC_FILE=".zprofile" ;;
    esac

    if [ -n "$RC_FILE" ]; then
        cat <<EOF
~/.local/bin is not on your PATH.

To add it permanently, run:

  echo 'export PATH="\$HOME/.local/bin:\$PATH"' >> ~/$RC_FILE

Then restart your terminal.

EOF
    else
        cat <<EOF
Add the following line to your shell profile:

  export PATH="\$HOME/.local/bin:\$PATH"

Then restart your shell.

EOF
    fi

    cat <<EOF
Or run andy-router directly:

  ~/.local/bin/andy-router

EOF
}

main "$@"