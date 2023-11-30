#!/usr/bin/env bash

GITHUB_OWNER=${GITHUB_OWNER}
GITHUB_REPO=${GITHUB_REPO}
GITHUB_TOKEN=${GITHUB_TOKEN}

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

owner=${GITHUB_OWNER}
repo=${GITHUB_REPO}
token=${GITHUB_TOKEN}

#
# Sanity checks
#

if is_dirty; then
    echo "Working directory is dirty"
    exit 1
fi

tag=$(get_tag || true)
test -n "${tag}" || (echo "No tag exactly matches current commit" && exit 1)

test -n "${GITHUB_TOKEN}" || (echo "Missing required github token" && exit 1)

#
# Build release assets
#

${verbose} && echo "Creating binary distributions..."
sh -c "${scriptsdir}/dist.sh -n '${bin_name}' -d '${dist_dir}' -p '${package}'"

assets=$(find ${dist_dir}/ -type f \
    -name "${bin_name}_${tag}_*"
)

#
# Create release changelog
#

changes=$(get_changes)
commits_url="https://github.com/${GITHUB_OWNER}/${GITHUB_REPO}/commit/"
changelog=$(awk '{ printf "%s\\n", $0}' <<-EOF
## Changelog
$(awk '{ print "  - [`" $1 "`][" $1 "] " substr($0, index($0, $2)) }' <<< $changes)

<!-- Link -->
$(awk -v commits_url=${commits_url} '{ print "[" $1 "]: " commits_url $1 }' <<< $changes)
EOF
)

#
# Create release
#

release_data() {
    cat <<-EOF
{
    "tag_name": "${tag}",
    "name": "${tag}",
    "body": "${changelog}"
}
EOF
}

${verbose} && echo "Creating release for ${tag}..."
release_url="https://api.github.com/repos/${GITHUB_OWNER}/${GITHUB_REPO}/releases"
response=$(curl -L -s \
    -H "Accept: application/vnd.github+json" \
    -H "Authorization: Bearer ${GITHUB_TOKEN}" \
    -H "X-GitHub-Api-Version 2022-11-28" \
    -d "$(release_data | awk '{ printf "%s ", $0}')" \
    ${release_url} \
    2>/dev/null
)

error=$(jq -r '.errors[0].code // empty' <<< $response)
[ -n "$error" ] && {
    echo ${error}
    exit 1
}

release_id=$(jq -r '.id // empty' <<< $response)
[ -n "$release_id" ] || {
    echo "No release with tag: ${tag}"
    echo $response
    exit 1
}

#
# Upload release assets
#

upload_url="https://uploads.github.com/repos/${owner}/${repo}/releases/${release_id}/assets"
for asset in ${assets}; do
    name=$(basename ${asset})
    ${verbose} && echo "uploading asset: ${name}"
    response=$(curl -L -s \
        -H "Accept: application/vnd.github+json" \
        -H "Authorization: Bearer ${GITHUB_TOKEN}" \
        -H "X-GitHub-Api-Version 2022-11-28" \
        -H "Content-Type: application/octet-stream" \
        ${upload_url}?name=${name} \
        --data-binary "@${asset}"
    )
done
