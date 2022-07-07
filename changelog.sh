#!/bin/bash

IFS=$'\n'
# obtain a list of all the tags in the repository, sort by number after last decimal point, then reverse the list
tags=$(git tag --sort=-v:refname --list | sed -e 's/^v//')
#tags=$(git tag -l)
# obtain the initial commit hash
initial_commit=$(git rev-list --max-parents=0 HEAD)
# obtain the first tag after the initial commit
first_tag=$(echo "$tags" | head -n 1)
# generate a changelog from the initial commit to the very first tag
entry="$(git log --oneline $initial_commit..$first_tag | sed -r '/^.{,40}$/d')"

author=$2
email=$3
packagename=$1

changelogentry(){
    version=$1
    #newentry=$2
    echo "$packagename ($version-1) UNRELEASED; urgency=medium"
    echo ""
    echo "  * tag $version"
    for line in $entry; do
        echo "  * $line"
    done
    echo ""
    echo " -- $author <$email> $(date -R)"
    echo ""
}

#changelogentry "$first_tag" $entry

for tag in ${tags}; do
    # obtain the tag after this one
    next_tag=$(echo "$tags" | grep -A 1 $tag | tail -n 1)
    echo "Generating changelog for $tag, $next_tag" 1>&2
    # if there is no next tag, quit
    # sleep 20s
    if [ "$tag" = "$next_tag" ]; then
        next_tag=$initial_commit
    fi
    # generate a changelog from the commit after the tag to the next tag
    entry="$(git log --oneline $tag..$first_tag | sed -r '/^.{,30}$/d')"
    changelogentry "$tag" "$entry"
    echo ""
done
