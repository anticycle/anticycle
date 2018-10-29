#!/usr/bin/env bash
# Output full Anticycle version based on git history.
# If repository has no tags then use v0.0.0 as default
# and last revision set to HEAD for short commit hash.
git_tag=$(git describe --abbrev=0 --tags 2> /dev/null)
if [[ ${git_tag} == "" ]]; then
    git_tag="v0.0.0"
    rev_for="HEAD"
else
    rev_for=${git_tag}
fi

git_rev=$(git rev-parse --short --verify ${rev_for})
echo "${git_tag} ${git_rev}"
