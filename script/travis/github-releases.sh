#!/bin/bash

go get github.com/aktau/github-release


#VERSION=$(cat ruby-version)
VERSION=0.9.1
USER="stanaka"
REPO="mkr"

# create description

# create release
github-release release \
  --user $USER \
  --repo $REPO \
  --tag $VERSION \
  --name "$REPO-$VERSION" \
  --description "not release"

# upload files
echo "Use at your own risk!" >> description.md
echo "" >> description.md
for i in $(ls -1 *.rpm)
do
  echo "* $i" >> description.md
  echo "  * $(openssl sha256 $i)" >> description.md
  github-release upload --user $USER \
    --repo $REPO \
    --tag $VERSION \
    --name "$i" \
    --file $i
done

# edit description
github-release edit \
  --user $USER \
  --repo $REPO \
  --tag $VERSION \
  --description "$(cat description.md)"
