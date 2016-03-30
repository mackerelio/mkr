#!/bin/sh

set -xeu

# create pull request
pwd
~/bin/hub pull-request -m 'test' -b stanaka/mkr:master -h stanaka/mkr:$TRAVIS_BRANCH
