#!/bin/sh
# Based on flyctl install.sh and https://github.com/curquiza/MeiliSearch/blob/e291d9954a18413c99105def1a2d32e63e5715be/download-latest.sh
# Based on Deno installer: Copyright 2019 the Deno authors. All rights reserved. MIT license.
# TODO(everyone): Keep this script simple and easily auditable.

os=$(uname -s)
arch=$(uname -m)
#version=${1:-latest}

# Literally with quotes: "v[number].[number].[number]"
# This makes sure we ignore pre-releases that are in the format "v[number].[number].[number]-pre-[number]"
releases_regexp='\"v\([0-9]*\)[.]\([0-9]*\)[.]\([0-9]*\)\"'

# semverParseInto and semverLT from https://github.com/cloudflare/semver_bash/blob/master/semver.sh
# usage: semverParseInto version major minor patch special
# version: the string version
# major, minor, patch, special: will be assigned by the function
semverParseInto() {
    local RE='[^0-9]*\([0-9]*\)[.]\([0-9]*\)[.]\([0-9]*\)\([0-9A-Za-z-]*\)'
    #MAJOR
    eval $2=`echo $1 | sed -e "s#$RE#\1#"`
    #MINOR
    eval $3=`echo $1 | sed -e "s#$RE#\2#"`
    #MINOR
    eval $4=`echo $1 | sed -e "s#$RE#\3#"`
    #SPECIAL
    eval $5=`echo $1 | sed -e "s#$RE#\4#"`
}

# usage: semverLT version1 version2
semverLT() {
    local MAJOR_A=0
    local MINOR_A=0
    local PATCH_A=0
    local SPECIAL_A=0

    local MAJOR_B=0
    local MINOR_B=0
    local PATCH_B=0
    local SPECIAL_B=0

    semverParseInto $1 MAJOR_A MINOR_A PATCH_A SPECIAL_A
    semverParseInto $2 MAJOR_B MINOR_B PATCH_B SPECIAL_B

    if [ $MAJOR_A -lt $MAJOR_B ]; then
        return 0
    fi
    if [ $MAJOR_A -le $MAJOR_B ] && [ $MINOR_A -lt $MINOR_B ]; then
        return 0
    fi
    if [ $MAJOR_A -le $MAJOR_B ] && [ $MINOR_A -le $MINOR_B ] && [ $PATCH_A -lt $PATCH_B ]; then
        return 0
    fi
    if [ "_$SPECIAL_A"  == "_" ] && [ "_$SPECIAL_B"  == "_" ] ; then
        return 1
    fi
    if [ "_$SPECIAL_A"  == "_" ] && [ "_$SPECIAL_B"  != "_" ] ; then
        return 1
    fi
    if [ "_$SPECIAL_A"  != "_" ] && [ "_$SPECIAL_B"  == "_" ] ; then
        return 0
    fi
    if [ "_$SPECIAL_A" < "_$SPECIAL_B" ]; then
        return 0
    fi

    return 1
}

# Get all tag releases.
# Grep edits:
# "name": "v0.0.1",
# name: v0.0.1,
# name: v0.0.1
# 0.0.1
tags=$(curl -s https://api.github.com/repos/devbookhq/devbookctl/tags \
  | grep "$releases_regexp" \
  | grep 'name' \
  | tr -d '"' \
  | tr -d ',' \
  | cut -d 'v' -f2)

# Sort the tags
latest=""
for tag in $tags; do
  echo "t $tag"
  if [ "$latest" = "" ]; then
    latest="$tag"
  else
    semverLT $tag $latest
    if [ $? -eq 1 ]; then
      latest="$tag"
    fi
  fi
done

if [ ! "$latest" ]; then
  echo "No releases found"
  exit 1
fi

dbk_uri="https://github.com/devbookhq/devbookctl/releases/$version/download/dbk_${os}_${arch}.tar.gz"


#dbk_install="${DBK_INSTALL:-$HOME/.dbk}"
dbk_install="/usr/local"

bin_dir="$dbk_install/bin"
exe="$bin_dir/dbk"
simexe="$bin_dir/devbookctl"

if [ ! -d "$bin_dir" ]; then
 	mkdir -p "$bin_dir"
fi

curl --fail --location --progress-bar --output "$exe.tar.gz" "$dbk_uri"
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

echo "dbk was installed successfully to $exe"
if command -v dbk >/dev/null; then
	echo "Run 'dbk --help' to get started"
else
	case $SHELL in
	/bin/zsh) shell_profile=".zshrc" ;;
  /bin/fish) shell_profile=".config/fish/config.fish" ;;
	*) shell_profile=".bash_profile" ;;
	esac
	echo "Manually add the directory to your \$HOME/$shell_profile (or similar)"
	echo "  export DBK_INSTALL=\"$dbk_install\""
	echo "  export PATH=\"\$DBK_INSTALL/bin:\$PATH\""
	echo "Run '$exe --help' to get started"
fi


