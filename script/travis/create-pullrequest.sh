#!/bin/sh

set -xeu

# create pull request
echo $PATH
~/bin/hub pull-request -m 'test'
