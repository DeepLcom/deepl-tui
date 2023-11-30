#!/bin/sh

get_git_dir() {
    git rev-parse --show-toplevel
}

get_tag() {
    git describe --exact-match 2>/dev/null
}

get_commit_sha() {
    git rev-parse --verify ${1:-HEAD}
}

get_commit_date() {
    TZ=UTC git show --no-patch --format='%cd' --date='format-local:%Y-%m-%dT%H:%M:%SZ' ${1:-HEAD}
}

get_pseudo_version() {
    ref=${1:-"HEAD"}
    prefix=${2:-"v0.0.0"}

    latest_tag=$(git describe --tags --abbrev=0 2>/dev/null || true)
    if [ -n "${latest_tag}" ]; then
        prefix=${latest_tag}-$(git rev-list --count ${latest_tag}..${ref})
    fi

    # UTC time the revision was created (yyyymmddhhmmss).
    timestamp=$(TZ=UTC git show --no-patch --format='%cd' --date='format-local:%Y%m%d%H%M%S' ${ref})

    # 12-character prefix of the commit hash
    revision=$(git rev-parse --short=12 --verify ${ref})

    echo "${prefix}-${timestamp}-${revision}"
}

get_version() {
    tag=$(git describe --exact-match 2>/dev/null || true)
    if [ -n "${tag}" ]; then
        echo ${tag}
    else
        get_pseudo_version
    fi
}

is_dirty() {
    ! git diff --quiet
}

is_tagged() {
    test -n "$(get_tag)"
}

get_changes() {
    from=${1:-""}
    to=${2:-HEAD}

    if [ -z "${from}" ]; then
        from=$(git describe --tags --abbrev=0 2>/dev/null)
    fi

    if [ "${from}" == "$(git describe --tags --exact-match ${to} 2>/dev/null)" ]; then
        from=$(git describe --tags --abbrev=0 --exclude=${from} 2>/dev/null)
    fi

    if [ -n "${from}" ]; then
        git log --oneline --no-decorate ${from}..${to}
    else
        git log --oneline --no-decorate ${to}
    fi
}
