# Changelog

## 0.59.2 (2025-03-04)

* fix permission #686 (yseto)
* added container registry GHCR #684 (yseto)
* added container registry ECR Public #683 (yseto)
* Bump mackerelio/workflows from 1.2.0 to 1.4.0 #682 (dependabot[bot])


## 0.59.1 (2025-01-27)

* Bump the stable-packages group across 1 directory with 3 updates #672 (dependabot[bot])
* Fix CI build #671 (ne-sachirou)
* Bump golang.org/x/net from 0.23.0 to 0.33.0 #669 (dependabot[bot])
* host-status option doesn't receive File in alerts command #668 (do-su-0805)
* Use strings.Replacer and strings.ReplaceAll where appropriate #667 (itchyny)
* use mackerelio/workflows@v1.2.0 #666 (yseto)
* Bump alpine from 3.19.1 to 3.20.3 #661 (dependabot[bot])
* Bump golang from 1.22-alpine to 1.23-alpine #660 (dependabot[bot])
* Bump docker/build-push-action from 5 to 6 #654 (dependabot[bot])


## 0.59.0 (2024-11-26)

* Support query monitorings #664 (rmatsuoka)
* fix unintentionally return in validateRules #663 (rmatsuoka)


## 0.58.0 (2024-06-12)

* Bump the stable-packages group with 3 updates #651 (dependabot[bot])
* Update Go and dependencies #650 (lufia)
* Add `--jq` option to `org` sub command #646 (kmuto)
* Bump github.com/urfave/cli from 1.22.14 to 1.22.15 #641 (dependabot[bot])


## 0.57.1 (2024-05-01)

* Bump github.com/mackerelio/mackerel-agent from 0.80.0 to 0.81.0 in the stable-packages group #642 (dependabot[bot])
* Bump golang.org/x/net from 0.22.0 to 0.23.0 #640 (dependabot[bot])
* Bump the stable-packages group with 4 updates #639 (dependabot[bot])
* Bump github.com/itchyny/gojq from 0.12.14 to 0.12.15 #638 (dependabot[bot])
* Bump google.golang.org/protobuf from 1.31.0 to 1.33.0 #634 (dependabot[bot])
* Bump peter-evans/repository-dispatch from 2 to 3 #623 (dependabot[bot])
* Bump actions/download-artifact from 3 to 4 #615 (dependabot[bot])
* Bump actions/upload-artifact from 3 to 4 #614 (dependabot[bot])
* Bump docker/build-push-action from 4 to 5 #591 (dependabot[bot])


## 0.57.0 (2024-03-18)

* Bump the stable-packages group with 4 updates #635 (dependabot[bot])
* Bump the dev-dependencies group with 1 update #629 (dependabot[bot])
* Bump mackerelio/workflows from 1.0.2 to 1.1.0 #627 (dependabot[bot])
* Bump golang from 1.21-alpine to 1.22-alpine #625 (dependabot[bot])
* Bump alpine from 3.17.3 to 3.19.1 #624 (dependabot[bot])
* Bump github.com/itchyny/gojq from 0.12.13 to 0.12.14 #609 (dependabot[bot])


## 0.56.0 (2024-03-07)

* Always set CGO_ENABLED=0 #630 (fujiwara)


## 0.55.0 (2024-02-27)

* accept UTF8-BOM when reading dashboard or monitor from JSON file #622 (kmuto)


## 0.54.0 (2023-12-22)

* Bump the stable-packages group with 1 update #617 (dependabot[bot])
* Bump golang.org/x/crypto from 0.13.0 to 0.17.0 #616 (dependabot[bot])
* Bump actions/setup-go from 4 to 5 #612 (dependabot[bot])
* Bump the stable-packages group with 5 updates #610 (dependabot[bot])
* update Go version to 1.21 and 1.20 by using reusable workflow #604 (lufia)


## 0.53.0 (2023-09-22)

* Bump docker/login-action from 2 to 3 #593 (dependabot[bot])
* Bump docker/setup-buildx-action from 2 to 3 #592 (dependabot[bot])
* Bump docker/setup-qemu-action from 2 to 3 #590 (dependabot[bot])
* Bump actions/checkout from 3 to 4 #589 (dependabot[bot])
* Bump golang.org/x/oauth2 from 0.7.0 to 0.12.0 #588 (dependabot[bot])
* Remove old rpm packaging #585 (yseto)


## 0.52.0 (2023-07-26)

* read annotation's description from file or stdin #582 (Arthur1)
* Bump github.com/mackerelio/mackerel-agent from 0.77.0 to 0.77.1 #581 (dependabot[bot])
* Bump golang.org/x/sync from 0.1.0 to 0.3.0 #577 (dependabot[bot])
* Bump github.com/urfave/cli from 1.22.12 to 1.22.14 #574 (dependabot[bot])
* Bump github.com/itchyny/gojq from 0.12.12 to 0.12.13 #573 (dependabot[bot])
* Bump github.com/stretchr/testify from 1.8.2 to 1.8.4 #570 (dependabot[bot])
* Bump actions/setup-go from 3 to 4 #555 (dependabot[bot])


## 0.51.1 (2023-07-13)

* Bump github.com/mackerelio/mackerel-agent from 0.75.1 to 0.77.0 #576 (dependabot[bot])


## 0.51.0 (2023-05-31)

* Bump github.com/mackerelio/mackerel-client-go from 0.25.0 to 0.26.0 #568 (dependabot[bot])


## 0.50.0 (2023-05-17)

* update macos,windows Actions Runner Image. #566 (yseto)
* Add metric-names subcommand. #561 (fujiwara)


## 0.49.3 (2023-04-12)

* Bump github.com/mackerelio/mackerel-client-go from 0.24.0 to 0.25.0 #558 (dependabot[bot])
* Bump golang.org/x/oauth2 from 0.5.0 to 0.7.0 #557 (dependabot[bot])
* Bump alpine from 3.17.2 to 3.17.3 #556 (dependabot[bot])
* Bump github.com/fatih/color from 1.14.1 to 1.15.0 #554 (dependabot[bot])
* Bump github.com/itchyny/gojq from 0.12.11 to 0.12.12 #552 (dependabot[bot])


## 0.49.2 (2023-02-27)

* Bump golang.org/x/crypto from 0.0.0-20210817164053-32db794688a5 to 0.1.0 #550 (dependabot[bot])
* Bump golang.org/x/net from 0.6.0 to 0.7.0 #549 (dependabot[bot])
* Bump actions/checkout from 2 to 3 #548 (dependabot[bot])
* Bump github.com/stretchr/testify from 1.8.1 to 1.8.2 #547 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.75.0 to 0.75.1 #546 (dependabot[bot])
* Bump docker/build-push-action from 3 to 4 #542 (dependabot[bot])


## 0.49.1 (2023-02-15)

* Bump alpine from 3.17.1 to 3.17.2 #544 (dependabot[bot])
* Bump golang.org/x/oauth2 from 0.4.0 to 0.5.0 #543 (dependabot[bot])
* Bump docker/setup-qemu-action from 1 to 2 #541 (dependabot[bot])
* Bump actions/download-artifact from 2 to 3 #540 (dependabot[bot])
* Bump peter-evans/repository-dispatch from 1 to 2 #539 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.74.1 to 0.75.0 #537 (dependabot[bot])


## 0.49.0 (2023-02-01)

* Bump docker/build-push-action from 2 to 3 #535 (dependabot[bot])
* Bump docker/setup-buildx-action from 1 to 2 #534 (dependabot[bot])
* Bump actions/upload-artifact from 2 to 3 #533 (dependabot[bot])
* Bump docker/login-action from 1 to 2 #532 (dependabot[bot])
* Bump actions/cache from 2 to 3 #531 (dependabot[bot])
* Enables Dependabot version updates for GitHub Actions #530 (Arthur1)
* Remove debian package v1 process. #529 (yseto)
* remove useless packaging script #528 (yseto)
* Bump github.com/fatih/color from 1.13.0 to 1.14.1 #527 (dependabot[bot])
* Bump github.com/urfave/cli from 1.22.10 to 1.22.12 #526 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.73.1 to 0.74.1 #525 (dependabot[bot])
* Bump github.com/google/go-github/v49 from 49.0.0 to 49.1.0 #522 (dependabot[bot])
* Bump alpine from 3.16.2 to 3.17.1 #521 (dependabot[bot])
* Bump golang.org/x/oauth2 from 0.3.0 to 0.4.0 #520 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-client-go from 0.23.0 to 0.24.0 #519 (dependabot[bot])
* Bump github.com/itchyny/gojq from 0.12.10 to 0.12.11 #516 (dependabot[bot])


## 0.48.0 (2023-01-18)

* Update some libraries #518 (yseto)
* improve `mkr status -v` to show host metrics #517 (kmuto)
* added compile option, fix packaging format #515 (yseto)
* Bump github.com/itchyny/gojq from 0.12.9 to 0.12.10 #513 (dependabot[bot])


## 0.47.2 (2022-11-04)

* Bump github.com/stretchr/testify from 1.8.0 to 1.8.1 #505 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-client-go from 0.21.2 to 0.22.0 #503 (dependabot[bot])
* Bump github.com/golangci/golangci-lint from 1.47.1 to 1.50.0 #502 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.72.14 to 0.73.1 #501 (dependabot[bot])
* Bump github.com/urfave/cli from 1.22.9 to 1.22.10 #498 (dependabot[bot])
* Bump github.com/itchyny/gojq from 0.12.8 to 0.12.9 #497 (dependabot[bot])
* Bump github.com/Songmu/goxz from 0.8.2 to 0.9.1 #492 (dependabot[bot])
* Bump alpine from 3.16.1 to 3.16.2 #488 (dependabot[bot])


## 0.47.1 (2022-09-14)

* Bump github.com/mackerelio/mackerel-client-go from 0.21.1 to 0.21.2 #494 (dependabot[bot])


## 0.47.0 (2022-09-06)

* Go 1.17 -> 1.19 #496 (yseto)
* added filter by gojq #490 (yseto)
* Organize flat packages in layers #489 (yseto)


## 0.46.9 (2022-07-20)

* Bump github.com/golangci/golangci-lint from 1.46.2 to 1.47.1 #479 (dependabot[bot])
* Bump alpine from 3.16.0 to 3.16.1 #478 (dependabot[bot])
* Bump github.com/stretchr/testify from 1.7.2 to 1.8.0 #475 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-client-go from 0.21.0 to 0.21.1 #473 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.72.13 to 0.72.14 #472 (dependabot[bot])
* Bump alpine from 3.15.1 to 3.16.0 #461 (dependabot[bot])


## 0.46.8 (2022-06-22)

* Bump github.com/mackerelio/mackerel-agent from 0.72.12 to 0.72.13 #468 (dependabot[bot])


## 0.46.7 (2022-06-08)

* Bump github.com/stretchr/testify from 1.7.1 to 1.7.2 #466 (dependabot[bot])
* Bump github.com/Songmu/prompter from 0.5.0 to 0.5.1 #464 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.72.11 to 0.72.12 #463 (dependabot[bot])


## 0.46.6 (2022-05-26)

* Bump github.com/golangci/golangci-lint from 1.45.2 to 1.46.2 #460 (dependabot[bot])
* Bump github.com/Songmu/goxz from 0.8.1 to 0.8.2 #456 (dependabot[bot])
* Bump github.com/urfave/cli from 1.22.5 to 1.22.9 #455 (dependabot[bot])
* [ci] Fix: Input 'job-number' has been deprecated with message: use flag-name instead #452 (ne-sachirou)
* Bump github.com/mackerelio/mackerel-agent from 0.72.8 to 0.72.11 #451 (dependabot[bot])


## 0.46.5 (2022-03-30)

* Bump github.com/golangci/golangci-lint from 1.45.0 to 1.45.2 #445 (dependabot[bot])
* Bump github.com/golangci/golangci-lint from 1.44.2 to 1.45.0 #443 (dependabot[bot])
* Bump alpine from 3.15.0 to 3.15.1 #441 (dependabot[bot])
* Bump github.com/stretchr/testify from 1.7.0 to 1.7.1 #439 (dependabot[bot])


## 0.46.4 (2022-03-15)

* refine README and others #437 (lufia)
* Bump github.com/golangci/golangci-lint from 1.44.0 to 1.44.2 #436 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.72.7 to 0.72.8 #434 (dependabot[bot])


## 0.46.3 (2022-02-16)

* upgrade Go version: 1.16 -> 1.17 #432 (lufia)
* Bump github.com/mackerelio/mackerel-agent from 0.72.6 to 0.72.7 #431 (dependabot[bot])
* Bump alpine from 3.14.2 to 3.15.0 #417 (dependabot[bot])


## 0.46.2 (2022-02-02)

* Bump github.com/golangci/golangci-lint from 1.43.0 to 1.44.0 #429 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.72.4 to 0.72.6 #428 (dependabot[bot])
* Bump github.com/Songmu/goxz from 0.7.0 to 0.8.1 #424 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-client-go from 0.20.0 to 0.21.0 #423 (dependabot[bot])


## 0.46.1 (2022-01-12)

* Add a job to build docker images and push to DockerHub #425 (Krout0n)


## 0.46.0 (2021-12-01)

* added memo on host #418 (yseto)
* Bump github.com/fatih/color from 1.12.0 to 1.13.0 #416 (dependabot[bot])
* Bump github.com/golangci/golangci-lint from 1.41.1 to 1.43.0 #415 (dependabot[bot])
* Bump mackerel-agent from 0.72.1 to 0.72.4 #414 (susisu)
* Add arm64/darwin build artifacts to GitHub release #413 (astj)


## 0.45.3 (2021-10-14)

* implement list cloud integration settings subcommand #407 (Gompei)
* Bump alpine from 3.14.0 to 3.14.2 #405 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.72.0 to 0.72.1 #396 (dependabot[bot])


## 0.45.2 (2021-06-23)

* Bump alpine from 3.13.5 to 3.14.0 #391 (dependabot[bot])
* Bump github.com/golangci/golangci-lint from 1.40.1 to 1.41.1 #394 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.71.2 to 0.72.0 #393 (dependabot[bot])


## 0.45.1 (2021-06-03)

* Bump github.com/fatih/color from 1.11.0 to 1.12.0 #388 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.71.1 to 0.71.2 #389 (dependabot[bot])
* Bump github.com/fatih/color from 1.10.0 to 1.11.0 #386 (dependabot[bot])
* Bump github.com/golangci/golangci-lint from 1.39.0 to 1.40.1 #387 (dependabot[bot])
* Bump alpine from 3.13.2 to 3.13.5 #379 (dependabot[bot])
* Bump github.com/Songmu/goxz from 0.6.0 to 0.7.0 #382 (dependabot[bot])


## 0.45.0 (2021-04-26)

* Bump github.com/mackerelio/mackerel-client-go from 0.16.0 to 0.17.0 #383 (dependabot[bot])
* [Breaking changes] Remove `mkr dashboards generate` and `mkr dashboards migrate` #380 (shibayu36)
* Bump github.com/golangci/golangci-lint from 1.38.0 to 1.39.0 #377 (dependabot[bot])
* Bump github.com/Songmu/prompter from 0.4.0 to 0.5.0 #375 (dependabot[bot])
* update mackerel-client-go 0.16.0 #374 (yseto)


## 0.44.2 (2021-03-05)

* Bump github.com/golangci/golangci-lint from 1.37.1 to 1.38.0 #369 (dependabot[bot])
* Bump alpine from 3.13.0 to 3.13.2 #362 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-client-go from 0.14.0 to 0.15.0 #370 (dependabot[bot])
* fix CI build stage. #367 (yseto)
* refactor: CI and dependency management #366 (lufia)
* build / test with Go 1.16 #365 (astj)
* added repository_dispatch to homebrew-mackerel-agent #363 (yseto)
* replace token #364 (yseto)


## 0.44.1 (2021-02-18)

* fix condition to delete monitors #360 (lufia)


## 0.44.0 (2021-02-16)

* replace mackerel-github-release #357 (yseto)
* add dashboard list, pull, push, and migrate commands for current dashboards. #353 (fujiwara)
* Bump github.com/mackerelio/mackerel-agent from 0.71.0 to 0.71.1 #351 (dependabot[bot])
* Update environments #348 (lufia)


## 0.43.1 (2021-01-21)

* Bump alpine from 3.12.2 to 3.13.0 #347 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-client-go from 0.12.0 to 0.13.0 #346 (dependabot[bot])
* Bump github.com/stretchr/testify from 1.6.1 to 1.7.0 #345 (dependabot[bot])


## 0.43.0 (2020-12-17)

* Bump github.com/mackerelio/mackerel-agent from 0.70.2 to 0.71.0 #342 (dependabot[bot])
* Bump gopkg.in/yaml.v2 from 2.3.0 to 2.4.0 #338 (dependabot[bot])
* Bump alpine from 3.12.1 to 3.12.2 #341 (dependabot[bot])
* migrate to GitHub Actions #339 (lufia)
* Bump github.com/fatih/color from 1.9.0 to 1.10.0 #330 (dependabot-preview[bot])


## 0.42.0 (2020-11-25)

* Add AlertStatusOnGone by updating mackerel-client-go to v0.12.0 #335 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.69.3 to 0.70.2 #334 (dependabot[bot])
* Bump github.com/urfave/cli from 1.22.4 to 1.22.5 #331 (dependabot-preview[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.69.2 to 0.69.3 #329 (dependabot-preview[bot])
* Update Dependabot config file #332 (dependabot-preview[bot])


## 0.41.0 (2020-10-28)

* Print empty array if hosts not found #326 (kazeburo)
* Add armhf Debian package to GitHub release #323 (hnw)
* Bump alpine from 3.12.0 to 3.12.1 #325 (dependabot-preview[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.69.1 to 0.69.2 #321 (dependabot-preview[bot])


## 0.40.4 (2020-10-01)

* Bump github.com/mackerelio/mackerel-client-go from 0.10.1 to 0.11.0 #319 (dependabot-preview[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.68.2 to 0.69.1 #317 (dependabot-preview[bot])


## 0.40.3 (2020-09-15)

* revert go1.15-alpine #315 (lufia)
* revert changing filename #314 (lufia)
* Bump golang from 1.15.0-alpine to 1.15.2-alpine #313 (dependabot-preview[bot])
* Bump github.com/Songmu/prompter from 0.3.0 to 0.4.0 #311 (dependabot-preview[bot])
* Bump golang from 1.14-alpine to 1.15.0-alpine #309 (dependabot-preview[bot])
* add arm64 packages, and fix Architecture field of deb #310 (lufia)
* Bump github.com/mackerelio/mackerel-agent from 0.68.0 to 0.68.2 #308 (dependabot-preview[bot])


## 0.40.2 (2020-07-20)

* Bump github.com/mackerelio/mackerel-client-go from 0.10.0 to 0.10.1 #305 (dependabot-preview[bot])
* define default file name while doMonitorsPull() at once #304 (hgsgtk)
* Update `Installation` section of README  #303 (astj)
* Bump github.com/stretchr/testify from 1.6.0 to 1.6.1 #300 (dependabot-preview[bot])
* Update go.sum #301 (shibayu36)
* Bump alpine from 3.11.6 to 3.12.0 #299 (dependabot-preview[bot])
* Bump github.com/mackerelio/mackerel-client-go from 0.9.1 to 0.10.0 #296 (dependabot-preview[bot])
* Bump github.com/stretchr/testify from 1.5.1 to 1.6.0 #298 (dependabot-preview[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.67.1 to 0.68.0 #295 (dependabot-preview[bot])


## 0.40.1 (2020-05-14)

* Bump gopkg.in/yaml.v2 from 2.2.8 to 2.3.0 #293 (dependabot-preview[bot])
* Bump alpine from 3.11.5 to 3.11.6 #292 (dependabot-preview[bot])
* Build with Go 1.14 #286 (lufia)
* Bump github.com/urfave/cli from 1.22.3 to 1.22.4 #287 (dependabot-preview[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.67.0 to 0.67.1 #289 (dependabot-preview[bot])


## 0.40.0 (2020-04-03)

* Bump alpine from 3.11.3 to 3.11.5 #285 (dependabot-preview[bot])
* Implement `mkr channels pull` command #270 (stefafafan)
* Bump github.com/urfave/cli from 1.22.2 to 1.22.3 #283 (dependabot-preview[bot])
* Stop building 32bit Darwin artifacts #282 (astj)
* Bump github.com/mackerelio/mackerel-client-go from 0.8.0 to 0.9.1 #276 (dependabot-preview[bot])
* Bump github.com/stretchr/testify from 1.4.0 to 1.5.1 #280 (dependabot-preview[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.65.0 to 0.67.0 #275 (dependabot-preview[bot])
* Bump gopkg.in/yaml.v2 from 2.2.7 to 2.2.8 #267 (dependabot-preview[bot])


## 0.39.7 (2020-02-05)

* show mandantory flags usage in modern CLI manner #273 (aereal)
* rename: github.com/motemen/gobump -> github.com/x-motemen/gobump #268 (lufia)


## 0.39.6 (2020-01-22)

* Bump github.com/pkg/errors from 0.8.1 to 0.9.1 #262 (dependabot-preview[bot])
* Bump alpine from 3.11.2 to 3.11.3 #264 (dependabot-preview[bot])
* Stop using alias imports for mackerel-client-go for simplicity #263 (stefafafan)
* Bump github.com/fatih/color from 1.8.0 to 1.9.0 #258 (dependabot-preview[bot])
* Bump github.com/fatih/color from 1.7.0 to 1.8.0 #257 (dependabot-preview[bot])
* Bump alpine from 3.11.0 to 3.11.2 #255 (dependabot-preview[bot])
* Bump alpine from 3.10.3 to 3.11.0 #254 (dependabot-preview[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.64.1 to 0.65.0 #252 (dependabot-preview[bot])


## 0.39.5 (2019-12-05)

* Bump gopkg.in/yaml.v2 from 2.2.5 to 2.2.7 #244 (dependabot-preview[bot])
* Bump github.com/Songmu/prompter from 0.2.0 to 0.3.0 #250 (dependabot-preview[bot])
* Bump github.com/urfave/cli from 1.22.1 to 1.22.2 #248 (dependabot-preview[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.64.0 to 0.64.1 #249 (dependabot-preview[bot])


## 0.39.4 (2019-11-21)

* always set GO111MODULE=on #246 (lufia)
* Bump alpine from 3.9 to 3.10.3 #238 (dependabot-preview[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.63.0 to 0.64.0 #240 (dependabot-preview[bot])
* Bump gopkg.in/yaml.v2 from 2.2.4 to 2.2.5 #242 (dependabot-preview[bot])


## 0.39.3 (2019-10-24)

* Build with Go 1.12.12
* Bump gopkg.in/yaml.v2 from 2.2.3 to 2.2.4 #236 (dependabot-preview[bot])


## 0.39.2 (2019-10-02)

* Bump gopkg.in/yaml.v2 from 2.2.2 to 2.2.3 #234 (dependabot-preview[bot])
* Bump github.com/urfave/cli from 1.22.0 to 1.22.1 #230 (dependabot-preview[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.62.1 to 0.63.0 #227 (dependabot-preview[bot])
* Bump github.com/mackerelio/mackerel-client-go from 0.7.0 to 0.8.0 #229 (dependabot-preview[bot])
* Bump github.com/mackerelio/mackerel-client-go from 0.6.0 to 0.7.0 #224 (dependabot-preview[bot])
* Bump github.com/urfave/cli from 1.21.0 to 1.22.0 #225 (dependabot-preview[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.62.0 to 0.62.1 #223 (dependabot-preview[bot])
* Bump github.com/stretchr/testify from 1.3.0 to 1.4.0 #221 (dependabot-preview[bot])


## 0.39.1 (2019-08-29)

* rename gopkg.in/urfave/cli.v1 -> github.com/urfave/cli #220 (lufia)
* Bump github.com/mackerelio/mackerel-agent from 0.61.1 to 0.62.0 #218 (dependabot-preview[bot])
* Bump github.com/mackerelio/mackerel-agent from 0.60.0 to 0.61.1 #216 (dependabot-preview[bot])


## 0.39.0 (2019-07-22)

* Implement anomaly detection monitor #212 (syou6162)
* Bump github.com/mackerelio/mackerel-agent from 0.59.3 to 0.60.0 #213 (dependabot-preview[bot])
* add fakeroot to build dependencies #211 (astj)
* upgrade builder image and enable go modules #210 (lufia)


## 0.38.0 (2019-06-11)

* Move implementation of `mkr create` into mkr/hosts #208 (astj)
* Build With Go 1.12 #207 (astj)
* Update go module dependencies #206 (astj)
* support go modules #205 (lufia)


## 0.37.0 (2019-05-08)

* [wrap] Improve message truncation algorithm #203 (itchyny)


## 0.36.0 (2019-03-27)

* [monitors] support interruption monitoring thresholds of service metric monitors #201 (itchyny)
* Improve Makefile #200 (itchyny)
* Create services package and add tests for services subcommand #198 (itchyny)
* Create org package and add tests for org subcommand #196 (itchyny)
* Add tests for hosts package #195 (itchyny)
* Introduce Client interface for Mackerel API #194 (itchyny)
* Add out stream argument to format.PrettyPrintJSON #193 (itchyny)


## 0.35.1 (2019-03-06)

* Fix status option of hosts subcommand #191 (itchyny)


## 0.35.0 (2019-02-13)

* separate hosts package from main and reorganize package structure #189 (Songmu)
* add wrap subcommand to monitor batch jobs to run with cron etc #186 (Songmu)
* add `checks run` subcommand to confirm check plugin settings #187 (Songmu)


## 0.34.2 (2018-11-27)

* Fix the default limit of mkr alerts #184 (itchyny)


## 0.34.1 (2018-11-26)

* avoid out of range exception in mkr alerts #181 (astj)


## 0.34.0 (2018-11-26)

* Fixed issue where nextID was not inherited #179 (yaminoma)
* Implemented according to new specification of Alerts API. #178 (yaminoma)


## 0.33.0 (2018-11-12)

* Follow https://github.com/mholt/archiver's change #175 (astj)
* update subcommand should keep the interface infromation #174 (itchyny)
* show alert.Message as monitor message in mkr alerts list #173 (astj)
* implement org subcommand #172 (itchyny)


## 0.32.1 (2018-10-17)

* Build with Go 1.11 #169 (astj)


## 0.32.0 (2018-08-30)

* logs to Stderr #167 (Songmu)
* Omit symbol table and debug information from the executable #166 (itchyny)
* fix Dockerfile to create alpine based docker image #163 (hayajo)


## 0.31.1 (2018-07-04)

* Fix mkr throw --retry not working #161 (astj)


## 0.31.0 (2018-07-04)

* add Retry feature and --retry option to `mkr throw` #159 (astj)


## 0.30.0 (2018-06-20)

* Build with Go 1.10 #157 (astj)


## 0.29.0 (2018-04-10)

* Change createdAt of hosts subcommand output to ISO 8601 extended format #154 (hayajo)


## 0.28.0 (2018-03-28)

* Add --upgrade option to plugin install. #150 (fujiwara)


## 0.27.1 (2018-03-15)

* Add <direct_url> help #151 (shibayu36)


## 0.27.0 (2018-03-01)

* Support empty threshold for monitors #148 (itchyny)


## 0.26.0 (2018-01-23)

* Fix copying plugin archives while installing a plugin on Windows #144 (itchyny)
* update rpm-v2 task for building Amazon Linux 2 package #143 (hayajo)


## 0.25.0 (2018-01-10)

* Add plugin document links #141 (shibayu36)
* introduce goxz and adjust deps #139 (Songmu)
* add appveyor.yml and adjust tests for windows #140 (Songmu)
* Define plugin default installation path in windows environment #138 (Songmu)


## 0.24.1 (2017-12-13)

* Rebuild to avoid panic when action of check-plugin was not specified #136 (astj)


## 0.24.0 (2017-12-12)

* Support maxCheckAttempts for host metric and service metric monitors #134 (itchyny)


## 0.23.0 (2017-11-28)

* [plugin.install] support direct URL target #130 (Songmu)
* [plugin.install] support tarball archives #131 (Songmu)
* fix hostId option flag in the command help of throw and metrics commands #129 (astj)
* Refactor mkr plugin install implementation #127 (shibayu36)


## 0.22.0 (2017-10-26)

* Release mkr plugin install #125 (shibayu36)
* Add metrics command #119 (edangelion)


## 0.21.0 (2017-10-04)

* Use new API BaseURL #116 (astj)


## 0.20.0 (2017-09-27)

* build with Go 1.9 #114 (astj)


## 0.19.0 (2017-09-20)

* Support fetch command to retrieve many hosts #112 (itchyny)
* Prefer apibase in mackerel-agent confFile #108 (astj)


## 0.18.0 (2017-09-12)

* add --customIdentifier option to mkr create #110 (astj)


## 0.17.0 (2017-08-23)

* [dashboards] Add unit to expression graph #106 (edangelion)
* [dashboards] Add title param to expression graph #104 (edangelion)


## 0.16.1 (2017-06-07)

* v2 packages #102 (Songmu)


## 0.16.0 (2017-05-09)

* Add services subcommand #97 (yuuki)


## 0.15.0 (2017-04-06)

* bump Go to 1.8 #95 (astj)


## 0.14.5 (2017-03-27)

* Colors on Windows #93 (mattn)


## 0.14.4 (2017-03-22)

* use new bot token #88 (daiksy)
* use new bot token #89 (daiksy)
* Workaround for git fetch failure #90 (daiksy)
* Apply git fetch workaround #91 (astj)


## 0.14.3 (2017-02-16)

* Support annotations command for graph annotation #83 (syou6162)
* Improve help management and fix usage help for command #85 (haya14busa)
* remove unused functions #86 (haya14busa)


## 0.14.2 (2017-02-08)

* [monitors diff] Add the "to-remote" bool flag #82 (yoheimuta)


## 0.14.1 (2017-01-11)

* formatter.NewAsciiFormatter now needs config #80 (astj)


## 0.14.0 (2016-12-21)

* Support expression monitor alerts in mkr alerts list #78 (itchyny)


## 0.13.0 (2016-11-29)

* remove unreachable code: monitor type cannot be "check" #72 (haya14busa)
* Fix the links to the api documents #73 (itchyny)
* catch up monitor interface changes of mackerel-client-go #74 (haya14busa)
* Introduce yudai/gojsondiff for `mkr monitors diff` #75 (haya14busa)
* fix test according to mackerel-client-go changes #76 (haya14busa)


## 0.12.0 (2016-10-27)

* Rename a dependent package #68 (usk81)
* Support `-apibase` option #69 (astj)
* [breaking change] Prepend `custom.` prefix to host metric name by default #70 (astj)


## 0.11.3 (2016-07-14)

* fix `validateRules()`,  when monitor has rule of "expression". #66 (daiksy)


## 0.11.2 (2016-06-23)

* replace angle brackets for json #63 (daiksy)


## 0.11.1 (2016-06-10)

* fix version number #61 (stanaka)


## 0.11.0 (2016-06-09)

* add dashboard generator #56 (daiksy)
* Add flag to overwrite host's roles  #58 (haya14busa)


## 0.10.1 (2016-05-25)

* fix signnatures. codegangsta/cli #54 (tknzk)


## 0.10.0 (2016-05-10)

* support `isMute` field of monitors #49 (Songmu)
* support boolean at isEmpty #51 (stanaka)
* bump up go version to 1.6.2 #52 (stanaka)


## 0.9.1 (2016-03-25)

* use GOARCH=amd64 for now #41 (Songmu)


## 0.9.0 (2016-02-18)

* Support displayName of host's json #39 (stanaka)


## 0.8.1 (2016-01-07)

* fix handling host-status option #37 (stanaka)


## 0.8.0 (2016-01-06)

* support alerts subcommand #31 (stanaka)
* Fix README example about mkr throw #32 (yuuki1)
* Build with Go 1.5 #33 (itchyny)
* Fixed the english used in the command descriptions #35 (stefafafan)


## 0.7.1 (2015-11-12)

* support `notificationIntervai` field in monitors (stanaka)
* [bug] fix json parameter s/hostID/hostId/g (Songmu)

## 0.7.0 (2015-10-26)

* append newline to the end of monitors.json #23 (Songmu)
* fix printMonitor #24 (Songmu)
* fix diff output between slices #25 (Songmu)

## 0.6.0 (2015-10-15)

* Fix update command bug about overwriting hostname #17 (y_uuki)
* Stop the parallel request sending temporarily #18 (y_uuki)
* Suppress to display empty fields when mkr monitors diff #20 (by stanaka)

## 0.5.0 (2015-09-14)

 * add fields for external URL monitors (stanaka)

## 0.4.1 (2015-08-28)

* Create deb/rpm package for Linux release #11 (Sixeight)


## 0.3.0 (2015-07-05)

* [feature] add --conf option to specify conf file path #4 (Sixeight)
* [fix] Fix update command as firstaid #7 (Sixeight)

## 0.2.0 (2015-06-18)

* [feature] add -f flag to hosts command to format the output #2 (motemen)
