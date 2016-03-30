#!/bin/sh

set -xeu

#VERSION=$(cat ruby-version)
VERSION=v0.9.1
USER="stanaka"
REPO="mkr"
pwd
# create description

# create release
github-release release \
  --user $USER \
  --repo $REPO \
  --tag $VERSION \
  --name "$REPO-$VERSION" \
  --description "not release" \
  --pre-release

# upload files
echo "Use at your own risk!" >> description.md
echo "" >> description.md

for i in $(ls -1 ~/rpmbuild/RPMS/noarch/*.rpm) $(ls -1 packaging/*.deb) $(ls -1 snapshot/mkr_*)
do
  name=$(basename "$path" ".php")
  echo "* $name" >> description.md
  echo "  * $(openssl sha256 $i)" >> description.md
  github-release upload --user $USER \
    --repo $REPO \
    --tag $VERSION \
    --name "$name" \
    --file $i
done

# edit description
github-release edit \
  --user $USER \
  --repo $REPO \
  --tag $VERSION \
  --description "$(cat description.md)"
