#!/bin/sh

set -xeu

NEW_VERSION=0.9.2

# create pull request
git checkout $TRAVIS_BRANCH
~/bin/hub pull-request -m "Release version $NEW_VERSION" -b stanaka/mkr:master
