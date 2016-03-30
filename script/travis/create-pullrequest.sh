#!/bin/sh

set -xeu

# create pull request
pwd
git branch
git checkout $TRAVIS_BRANCH
~/bin/hub pull-request -m 'test' -b stanaka/mkr:master -h stanaka/mkr:$TRAVIS_BRANCH
