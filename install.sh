#!/bin/sh
set -eu

repo="raphael-p/datashard"
install_dir="${HOME}/.local/bin"
data_dir="${HOME}/.dash"
tmp_dir="$(mktemp -d)"

trap 'rm -rf "$tmp_dir"' EXIT

case "$(uname -s)-$(uname -m)" in
Darwin-arm64)
    archive="dash_Darwin_arm64.tar.gz"
    ;;
Darwin-x86_64)
    archive="dash_Darwin_x86_64.tar.gz"
    ;;
Linux-aarch64|Linux-arm64)
    archive="dash_Linux_arm64.tar.gz"
    ;;
Linux-x86_64)
    archive="dash_Linux_x86_64.tar.gz"
    ;;
*)
    echo "Unsupported platform: $(uname -s)-$(uname -m)" >&2
    exit 1
    ;;
esac

mkdir -p "$install_dir" "$data_dir"

curl -fsSL \
-o "$tmp_dir/$archive" \
"https://github.com/$repo/releases/latest/download/$archive"

tar -xzf "$tmp_dir/$archive" -C "$tmp_dir"

install -m 755 "$tmp_dir/dash" "$install_dir/dash"
if [ ! -f "$data_dir/config.json" ]; then
    install -m 644 "$tmp_dir/config.json" "$data_dir/config.json"
fi

echo "dash binary installed at: $install_dir/dash"
echo "dash data directory: $data_dir"
