#!/usr/bin/env bash

set -euo pipefail

usage() { echo "usage: $(basename -- $0) [-h] [-v] [-n NAME] [-d DIR] [-p PKG]" 1>&2; }

package=""
dist_dir=""
bin_name=""
verbose=false
while getopts ":hd:n:p:v" opt; do
    case "${opt}" in
        h)
            usage
            exit 0
            ;;
        d)
            dist_dir=${OPTARG}
            ;;
        n)
            bin_name=${OPTARG}
            ;;
        p)
            package=${OPTARG}
            ;;
        v)
            verbose=true
            ;;
        *)
            usage
            exit 1
            ;;
    esac
done

# load common helpers
scriptsdir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
source "${scriptsdir}/functions.sh"

if [ -z "${package}" ]; then
    package="."
fi

if [ -z "${dist_dir}" ]; then
    git_dir=$(get_git_dir)
    dist_dir="${git_dir}/dist"
fi

if [ -z "${bin_name}" ]; then
    bin_name=$(basename $(realpath "${package}"))
fi

version=$(get_version)

declare -A OSARCHMAP=(
    [linux]="amd64,arm,arm64"
    [darwin]="amd64,arm64"
    [windows]="amd64,arm,arm64"
)

${verbose} && echo "Building binaries..."
for os in ${!OSARCHMAP[@]}; do
    for arch in ${OSARCHMAP[$os]//,/ }; do
        out="${dist_dir}/${bin_name}_${version}_${os}_${arch}"

        ${verbose} && echo "  for ${os}/${arch}"
        GOOS=${os} GOARCH=${arch} ${scriptsdir}/build.sh -o "${out}" -p "${package}"
    done
done
${verbose} && echo "Building binaries... done"

exit 0
