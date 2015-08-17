#!/bin/sh
set -ex

deploykey=~/.ssh/deploy.key

echo "
Host github.com
    StrictHostKeyChecking no
    IdentityFile $deploykey
" >> ~/.ssh/config
openssl aes-256-cbc -K $encrypted_e192d40ccb57_key -iv $encrypted_e192d40ccb57_iv -in script/travis/mkr.pem.enc -out $deploykey -d
chmod 600 $deploykey
git config --global user.email "mackerel-developers@hatena.ne.jp"
git config --global user.name  "mackerel"
git remote set-url origin git@github.com:mackerelio/mkr.git
tool/autotag
