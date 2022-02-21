#!/bin/sh
# Based on Deno installer: Copyright 2019 the Deno authors. All rights reserved. MIT license.
# TODO(everyone): Keep this script simple and easily auditable.

set -e

os=$(uname -s)
arch=$(uname -m)
version=${1:-latest}
devbookctl_uri="https://github.com/devbookhq/devbookctl/releases/$version/download/devbookctl_${os}_${arch}.tar.gz"

if [ ! "$devbookctl_uri" ]; then
  # TODO
	echo "Error: Unable to find a devbookctl release for $os/$arch/$version - see github.com/devbookhq/devbookctl/releases for all versions" 1>&2
	exit 1
fi

#devbookctl_install="${DEVBOOKCTL_INSTALL:-$HOME/.dbk}"
devbookctl_install="/usr/local"

bin_dir="$devbookctl_install/bin"
exe="$bin_dir/dbk"
simexe="$bin_dir/devbookctl"

if [ ! -d "$bin_dir" ]; then
 	mkdir -p "$bin_dir"
fi

curl --fail --location --progress-bar --output "$exe.tar.gz" "$devbookctl_uri"
cd "$bin_dir"
tar xzf "$exe.tar.gz"
chmod +x "$exe"
rm "$exe.tar.gz"

ln -sf $exe $simexe

# TODO: We don't support `dbk version` yet.
#if [ "${2}" = "prerel" ] || [ "${1}" = "pre" ]; then
#	"$exe" version -s "shell-prerel"
#else
#	"$exe" version -s "shell"
#fi

echo "devbookctl was installed successfully to $exe"
if command -v devbookctl >/dev/null; then
	echo "Run 'devbookctl --help' to get started"
else
	case $SHELL in
	/bin/zsh) shell_profile=".zshrc" ;;
  /bin/fish) shell_profile=".config/fish/config.fish" ;;
	*) shell_profile=".bash_profile" ;;
	esac
	echo "Manually add the directory to your \$HOME/$shell_profile (or similar)"
	echo "  export DEVBOOKCTL_INSTALL=\"$devbookctl_install\""
	echo "  export PATH=\"\$DEVBOOKCTL_INSTALL/bin:\$PATH\""
	echo "Run '$exe --help' to get started"
fi
