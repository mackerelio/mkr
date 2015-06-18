RELEASE
=======

1. git checkout -b release-v$VERSION
1. Update CHANGELOG.md. etc.
1. git push origin release-v$VERSION
1. Pull Request.
1. Merge the PR.
1. git checkout master
1. git tag v$VERSION
1. git push --tag
