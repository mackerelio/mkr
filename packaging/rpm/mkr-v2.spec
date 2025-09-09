Name:      mkr
Version:   %{_version}
Release:   1%{?dist}
License:   ASL 2.0
Summary:   mackerel.io api client tool
URL:       https://mackerel.io
Group:     Application/System
Packager:  Hatena
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root

%description
mackerel.io api client tool

%prep

%build

%install
rm -rf %{buildroot}

install -d -m 755 %{buildroot}/%{_bindir}
install    -m 655 %{_builddir}/%{name}  %{buildroot}/%{_bindir}

%clean
rm -f %{buildroot}%{_bindir}/%{name}

%pre

%post

%preun

%files
%defattr(-,root,root)
%{_bindir}/%{name}

%changelog
* Tue Sep 9 2025 <mackerel-developers@hatena.ne.jp> - 0.62.1
- fix tarball version (by yseto)
- Revert "Release version 0.62.1" (by yseto)
- Release version 0.62.1 (by mackerelbot)
- Make plugin install safety to avoid text busy error (by kazeburo)

* Mon Aug 25 2025 <mackerel-developers@hatena.ne.jp> - 0.62.0
- Get plugin's releaseTag from the releases URL instead of using api (by kazeburo)
- Bump golang from 1.24-alpine to 1.25-alpine (by dependabot[bot])
- Bump actions/checkout from 4 to 5 (by dependabot[bot])
- Bump actions/download-artifact from 4 to 5 (by dependabot[bot])
- Bump alpine from 3.20.3 to 3.22.1 (by dependabot[bot])
- Bump golang.org/x/oauth2 from 0.25.0 to 0.27.0 (by dependabot[bot])
- Bump github.com/mholt/archives from 0.1.1 to 0.1.3 (by dependabot[bot])
- Bump the stable-packages group across 1 directory with 5 updates (by dependabot[bot])
- Bump github.com/urfave/cli from 1.22.15 to 1.22.17 (by dependabot[bot])
- Bump mackerelio/workflows from 1.4.0 to 1.5.0 (by dependabot[bot])
- Bump golang.org/x/net from 0.36.0 to 0.38.0 (by dependabot[bot])
- Bump github.com/itchyny/gojq from 0.12.16 to 0.12.17 (by dependabot[bot])
- Bump github.com/stretchr/testify from 1.9.0 to 1.10.0 in the dev-dependencies group (by dependabot[bot])

* Fri May 16 2025 <mackerel-developers@hatena.ne.jp> - 0.61.0
- Remove rewrite some files on every releases (by yseto)
- replace to use github.com/mholt/archives (by yseto)
- use Go 1.24 (by yseto)
- Implement find alert logs command (by appare45)
- Updated mackerel-client-go to v0.37.0 (by appare45)

* Mon Mar 31 2025 <mackerel-developers@hatena.ne.jp> - 0.60.0
- update mackerel-client-go to v0.36.0 (by fujiwara)
- replace to newer runner-images (by yseto)
- Bump golang.org/x/net from 0.33.0 to 0.36.0 (by dependabot[bot])

* Tue Mar 4 2025 <mackerel-developers@hatena.ne.jp> - 0.59.2
- fix permission (by yseto)
- added container registry GHCR (by yseto)
- added container registry ECR Public (by yseto)
- Bump mackerelio/workflows from 1.2.0 to 1.4.0 (by dependabot[bot])

* Mon Jan 27 2025 <mackerel-developers@hatena.ne.jp> - 0.59.1
- Bump the stable-packages group across 1 directory with 3 updates (by dependabot[bot])
- Fix CI build (by ne-sachirou)
- Bump golang.org/x/net from 0.23.0 to 0.33.0 (by dependabot[bot])
- host-status option doesn't receive File in alerts command (by do-su-0805)
- Use strings.Replacer and strings.ReplaceAll where appropriate (by itchyny)
- use mackerelio/workflows@v1.2.0 (by yseto)
- Bump alpine from 3.19.1 to 3.20.3 (by dependabot[bot])
- Bump golang from 1.22-alpine to 1.23-alpine (by dependabot[bot])
- Bump docker/build-push-action from 5 to 6 (by dependabot[bot])

* Tue Nov 26 2024 <mackerel-developers@hatena.ne.jp> - 0.59.0
- Support query monitorings (by rmatsuoka)
- fix unintentionally return in validateRules (by rmatsuoka)

* Wed Jun 12 2024 <mackerel-developers@hatena.ne.jp> - 0.58.0
- Bump the stable-packages group with 3 updates (by dependabot[bot])
- Update Go and dependencies (by lufia)
- Add `--jq` option to `org` sub command (by kmuto)
- Bump github.com/urfave/cli from 1.22.14 to 1.22.15 (by dependabot[bot])

* Wed May 1 2024 <mackerel-developers@hatena.ne.jp> - 0.57.1
- Bump github.com/mackerelio/mackerel-agent from 0.80.0 to 0.81.0 in the stable-packages group (by dependabot[bot])
- Bump golang.org/x/net from 0.22.0 to 0.23.0 (by dependabot[bot])
- Bump the stable-packages group with 4 updates (by dependabot[bot])
- Bump github.com/itchyny/gojq from 0.12.14 to 0.12.15 (by dependabot[bot])
- Bump google.golang.org/protobuf from 1.31.0 to 1.33.0 (by dependabot[bot])
- Bump peter-evans/repository-dispatch from 2 to 3 (by dependabot[bot])
- Bump actions/download-artifact from 3 to 4 (by dependabot[bot])
- Bump actions/upload-artifact from 3 to 4 (by dependabot[bot])
- Bump docker/build-push-action from 4 to 5 (by dependabot[bot])

* Mon Mar 18 2024 <mackerel-developers@hatena.ne.jp> - 0.57.0
- Bump the stable-packages group with 4 updates (by dependabot[bot])
- Bump the dev-dependencies group with 1 update (by dependabot[bot])
- Bump mackerelio/workflows from 1.0.2 to 1.1.0 (by dependabot[bot])
- Bump golang from 1.21-alpine to 1.22-alpine (by dependabot[bot])
- Bump alpine from 3.17.3 to 3.19.1 (by dependabot[bot])
- Bump github.com/itchyny/gojq from 0.12.13 to 0.12.14 (by dependabot[bot])

* Thu Mar 7 2024 <mackerel-developers@hatena.ne.jp> - 0.56.0
- Always set CGO_ENABLED=0 (by fujiwara)

* Tue Feb 27 2024 <mackerel-developers@hatena.ne.jp> - 0.55.0
- accept UTF8-BOM when reading dashboard or monitor from JSON file (by kmuto)

* Fri Dec 22 2023 <mackerel-developers@hatena.ne.jp> - 0.54.0
- Bump the stable-packages group with 1 update (by dependabot[bot])
- Bump golang.org/x/crypto from 0.13.0 to 0.17.0 (by dependabot[bot])
- Bump actions/setup-go from 4 to 5 (by dependabot[bot])
- Bump the stable-packages group with 5 updates (by dependabot[bot])
- update Go version to 1.21 and 1.20 by using reusable workflow (by lufia)

* Fri Sep 22 2023 <mackerel-developers@hatena.ne.jp> - 0.53.0
- Bump docker/login-action from 2 to 3 (by dependabot[bot])
- Bump docker/setup-buildx-action from 2 to 3 (by dependabot[bot])
- Bump docker/setup-qemu-action from 2 to 3 (by dependabot[bot])
- Bump actions/checkout from 3 to 4 (by dependabot[bot])
- Bump golang.org/x/oauth2 from 0.7.0 to 0.12.0 (by dependabot[bot])
- Remove old rpm packaging (by yseto)

* Wed Jul 26 2023 <mackerel-developers@hatena.ne.jp> - 0.52.0
- read annotation's description from file or stdin (by Arthur1)
- Bump github.com/mackerelio/mackerel-agent from 0.77.0 to 0.77.1 (by dependabot[bot])
- Bump golang.org/x/sync from 0.1.0 to 0.3.0 (by dependabot[bot])
- Bump github.com/urfave/cli from 1.22.12 to 1.22.14 (by dependabot[bot])
- Bump github.com/itchyny/gojq from 0.12.12 to 0.12.13 (by dependabot[bot])
- Bump github.com/stretchr/testify from 1.8.2 to 1.8.4 (by dependabot[bot])
- Bump actions/setup-go from 3 to 4 (by dependabot[bot])

* Thu Jul 13 2023 <mackerel-developers@hatena.ne.jp> - 0.51.1
- Bump github.com/mackerelio/mackerel-agent from 0.75.1 to 0.77.0 (by dependabot[bot])

* Wed May 31 2023 <mackerel-developers@hatena.ne.jp> - 0.51.0
- Bump github.com/mackerelio/mackerel-client-go from 0.25.0 to 0.26.0 (by dependabot[bot])

* Wed May 17 2023 <mackerel-developers@hatena.ne.jp> - 0.50.0
- update macos,windows Actions Runner Image. (by yseto)
- Add metric-names subcommand. (by fujiwara)

* Wed Apr 12 2023 <mackerel-developers@hatena.ne.jp> - 0.49.3
- Bump github.com/mackerelio/mackerel-client-go from 0.24.0 to 0.25.0 (by dependabot[bot])
- Bump golang.org/x/oauth2 from 0.5.0 to 0.7.0 (by dependabot[bot])
- Bump alpine from 3.17.2 to 3.17.3 (by dependabot[bot])
- Bump github.com/fatih/color from 1.14.1 to 1.15.0 (by dependabot[bot])
- Bump github.com/itchyny/gojq from 0.12.11 to 0.12.12 (by dependabot[bot])

* Mon Feb 27 2023 <mackerel-developers@hatena.ne.jp> - 0.49.2
- Bump golang.org/x/crypto from 0.0.0-20210817164053-32db794688a5 to 0.1.0 (by dependabot[bot])
- Bump golang.org/x/net from 0.6.0 to 0.7.0 (by dependabot[bot])
- Bump actions/checkout from 2 to 3 (by dependabot[bot])
- Bump github.com/stretchr/testify from 1.8.1 to 1.8.2 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.75.0 to 0.75.1 (by dependabot[bot])
- Bump docker/build-push-action from 3 to 4 (by dependabot[bot])

* Wed Feb 15 2023 <mackerel-developers@hatena.ne.jp> - 0.49.1
- Bump alpine from 3.17.1 to 3.17.2 (by dependabot[bot])
- Bump golang.org/x/oauth2 from 0.4.0 to 0.5.0 (by dependabot[bot])
- Bump docker/setup-qemu-action from 1 to 2 (by dependabot[bot])
- Bump actions/download-artifact from 2 to 3 (by dependabot[bot])
- Bump peter-evans/repository-dispatch from 1 to 2 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.74.1 to 0.75.0 (by dependabot[bot])

* Wed Feb 1 2023 <mackerel-developers@hatena.ne.jp> - 0.49.0
- Bump docker/build-push-action from 2 to 3 (by dependabot[bot])
- Bump docker/setup-buildx-action from 1 to 2 (by dependabot[bot])
- Bump actions/upload-artifact from 2 to 3 (by dependabot[bot])
- Bump docker/login-action from 1 to 2 (by dependabot[bot])
- Bump actions/cache from 2 to 3 (by dependabot[bot])
- Enables Dependabot version updates for GitHub Actions (by Arthur1)
- Remove debian package v1 process. (by yseto)
- remove useless packaging script (by yseto)
- Bump github.com/fatih/color from 1.13.0 to 1.14.1 (by dependabot[bot])
- Bump github.com/urfave/cli from 1.22.10 to 1.22.12 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.73.1 to 0.74.1 (by dependabot[bot])
- Bump github.com/google/go-github/v49 from 49.0.0 to 49.1.0 (by dependabot[bot])
- Bump alpine from 3.16.2 to 3.17.1 (by dependabot[bot])
- Bump golang.org/x/oauth2 from 0.3.0 to 0.4.0 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-client-go from 0.23.0 to 0.24.0 (by dependabot[bot])
- Bump github.com/itchyny/gojq from 0.12.10 to 0.12.11 (by dependabot[bot])

* Wed Jan 18 2023 <mackerel-developers@hatena.ne.jp> - 0.48.0
- Update some libraries (by yseto)
- improve `mkr status -v` to show host metrics (by kmuto)
- added compile option, fix packaging format (by yseto)
- Bump github.com/itchyny/gojq from 0.12.9 to 0.12.10 (by dependabot[bot])

* Fri Nov 4 2022 <mackerel-developers@hatena.ne.jp> - 0.47.2
- Bump github.com/stretchr/testify from 1.8.0 to 1.8.1 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-client-go from 0.21.2 to 0.22.0 (by dependabot[bot])
- Bump github.com/golangci/golangci-lint from 1.47.1 to 1.50.0 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.72.14 to 0.73.1 (by dependabot[bot])
- Bump github.com/urfave/cli from 1.22.9 to 1.22.10 (by dependabot[bot])
- Bump github.com/itchyny/gojq from 0.12.8 to 0.12.9 (by dependabot[bot])
- Bump github.com/Songmu/goxz from 0.8.2 to 0.9.1 (by dependabot[bot])
- Bump alpine from 3.16.1 to 3.16.2 (by dependabot[bot])

* Wed Sep 14 2022 <mackerel-developers@hatena.ne.jp> - 0.47.1
- Bump github.com/mackerelio/mackerel-client-go from 0.21.1 to 0.21.2 (by dependabot[bot])

* Tue Sep 6 2022 <mackerel-developers@hatena.ne.jp> - 0.47.0
- Go 1.17 -> 1.19 (by yseto)
- added filter by gojq (by yseto)
- Organize flat packages in layers (by yseto)

* Wed Jul 20 2022 <mackerel-developers@hatena.ne.jp> - 0.46.9
- Bump github.com/golangci/golangci-lint from 1.46.2 to 1.47.1 (by dependabot[bot])
- Bump alpine from 3.16.0 to 3.16.1 (by dependabot[bot])
- Bump github.com/stretchr/testify from 1.7.2 to 1.8.0 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-client-go from 0.21.0 to 0.21.1 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.72.13 to 0.72.14 (by dependabot[bot])
- Bump alpine from 3.15.1 to 3.16.0 (by dependabot[bot])

* Wed Jun 22 2022 <mackerel-developers@hatena.ne.jp> - 0.46.8
- Bump github.com/mackerelio/mackerel-agent from 0.72.12 to 0.72.13 (by dependabot[bot])

* Wed Jun 8 2022 <mackerel-developers@hatena.ne.jp> - 0.46.7
- Bump github.com/stretchr/testify from 1.7.1 to 1.7.2 (by dependabot[bot])
- Bump github.com/Songmu/prompter from 0.5.0 to 0.5.1 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.72.11 to 0.72.12 (by dependabot[bot])

* Thu May 26 2022 <mackerel-developers@hatena.ne.jp> - 0.46.6
- Bump github.com/golangci/golangci-lint from 1.45.2 to 1.46.2 (by dependabot[bot])
- Bump github.com/Songmu/goxz from 0.8.1 to 0.8.2 (by dependabot[bot])
- Bump github.com/urfave/cli from 1.22.5 to 1.22.9 (by dependabot[bot])
- [ci] Fix: Input 'job-number' has been deprecated with message: use flag-name instead (by ne-sachirou)
- Bump github.com/mackerelio/mackerel-agent from 0.72.8 to 0.72.11 (by dependabot[bot])

* Wed Mar 30 2022 <mackerel-developers@hatena.ne.jp> - 0.46.5
- Bump github.com/golangci/golangci-lint from 1.45.0 to 1.45.2 (by dependabot[bot])
- Bump github.com/golangci/golangci-lint from 1.44.2 to 1.45.0 (by dependabot[bot])
- Bump alpine from 3.15.0 to 3.15.1 (by dependabot[bot])
- Bump github.com/stretchr/testify from 1.7.0 to 1.7.1 (by dependabot[bot])

* Tue Mar 15 2022 <mackerel-developers@hatena.ne.jp> - 0.46.4
- refine README and others (by lufia)
- Bump github.com/golangci/golangci-lint from 1.44.0 to 1.44.2 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.72.7 to 0.72.8 (by dependabot[bot])

* Wed Feb 16 2022 <mackerel-developers@hatena.ne.jp> - 0.46.3
- upgrade Go version: 1.16 -> 1.17 (by lufia)
- Bump github.com/mackerelio/mackerel-agent from 0.72.6 to 0.72.7 (by dependabot[bot])
- Bump alpine from 3.14.2 to 3.15.0 (by dependabot[bot])

* Wed Feb 2 2022 <mackerel-developers@hatena.ne.jp> - 0.46.2
- Bump github.com/golangci/golangci-lint from 1.43.0 to 1.44.0 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.72.4 to 0.72.6 (by dependabot[bot])
- Bump github.com/Songmu/goxz from 0.7.0 to 0.8.1 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-client-go from 0.20.0 to 0.21.0 (by dependabot[bot])

* Wed Jan 12 2022 <mackerel-developers@hatena.ne.jp> - 0.46.1
- Add a job to build docker images and push to DockerHub (by Krout0n)

* Wed Dec 1 2021 <mackerel-developers@hatena.ne.jp> - 0.46.0
- added memo on host (by yseto)
- Bump github.com/fatih/color from 1.12.0 to 1.13.0 (by dependabot[bot])
- Bump github.com/golangci/golangci-lint from 1.41.1 to 1.43.0 (by dependabot[bot])
- Bump mackerel-agent from 0.72.1 to 0.72.4 (by susisu)
- Add arm64/darwin build artifacts to GitHub release (by astj)

* Thu Oct 14 2021 <mackerel-developers@hatena.ne.jp> - 0.45.3
- implement list cloud integration settings subcommand (by Gompei)
- Bump alpine from 3.14.0 to 3.14.2 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.72.0 to 0.72.1 (by dependabot[bot])

* Wed Jun 23 2021 <mackerel-developers@hatena.ne.jp> - 0.45.2
- Bump alpine from 3.13.5 to 3.14.0 (by dependabot[bot])
- Bump github.com/golangci/golangci-lint from 1.40.1 to 1.41.1 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.71.2 to 0.72.0 (by dependabot[bot])

* Thu Jun 03 2021 <mackerel-developers@hatena.ne.jp> - 0.45.1
- Bump github.com/fatih/color from 1.11.0 to 1.12.0 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.71.1 to 0.71.2 (by dependabot[bot])
- Bump github.com/fatih/color from 1.10.0 to 1.11.0 (by dependabot[bot])
- Bump github.com/golangci/golangci-lint from 1.39.0 to 1.40.1 (by dependabot[bot])
- Bump alpine from 3.13.2 to 3.13.5 (by dependabot[bot])
- Bump github.com/Songmu/goxz from 0.6.0 to 0.7.0 (by dependabot[bot])

* Mon Apr 26 2021 <mackerel-developers@hatena.ne.jp> - 0.45.0
- Bump github.com/mackerelio/mackerel-client-go from 0.16.0 to 0.17.0 (by dependabot[bot])
- [Breaking changes] Remove `mkr dashboards generate` and `mkr dashboards migrate` (by shibayu36)
- Bump github.com/golangci/golangci-lint from 1.38.0 to 1.39.0 (by dependabot[bot])
- Bump github.com/Songmu/prompter from 0.4.0 to 0.5.0 (by dependabot[bot])
- update mackerel-client-go 0.16.0 (by yseto)

* Fri Mar 05 2021 <mackerel-developers@hatena.ne.jp> - 0.44.2
- Bump github.com/golangci/golangci-lint from 1.37.1 to 1.38.0 (by dependabot[bot])
- Bump alpine from 3.13.0 to 3.13.2 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-client-go from 0.14.0 to 0.15.0 (by dependabot[bot])
- fix CI build stage. (by yseto)
- refactor: CI and dependency management (by lufia)
- build / test with Go 1.16 (by astj)
- added repository_dispatch to homebrew-mackerel-agent (by yseto)
- replace token (by yseto)

* Thu Feb 18 2021 <mackerel-developers@hatena.ne.jp> - 0.44.1
- fix condition to delete monitors (by lufia)

* Tue Feb 16 2021 <mackerel-developers@hatena.ne.jp> - 0.44.0
- replace mackerel-github-release (by yseto)
- add dashboard list, pull, push, and migrate commands for current dashboards. (by fujiwara)
- Bump github.com/mackerelio/mackerel-agent from 0.71.0 to 0.71.1 (by dependabot[bot])
- Update environments (by lufia)

* Thu Jan 21 2021 <mackerel-developers@hatena.ne.jp> - 0.43.1
- Bump alpine from 3.12.2 to 3.13.0 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-client-go from 0.12.0 to 0.13.0 (by dependabot[bot])
- Bump github.com/stretchr/testify from 1.6.1 to 1.7.0 (by dependabot[bot])

* Thu Dec 17 2020 <mackerel-developers@hatena.ne.jp> - 0.43.0
- Bump github.com/mackerelio/mackerel-agent from 0.70.2 to 0.71.0 (by dependabot[bot])
- Bump gopkg.in/yaml.v2 from 2.3.0 to 2.4.0 (by dependabot[bot])
- Bump alpine from 3.12.1 to 3.12.2 (by dependabot[bot])
- migrate to GitHub Actions (by lufia)
- Bump github.com/fatih/color from 1.9.0 to 1.10.0 (by dependabot-preview[bot])

* Wed Nov 25 2020 <mackerel-developers@hatena.ne.jp> - 0.42.0
- Add AlertStatusOnGone by updating mackerel-client-go to v0.12.0 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.69.3 to 0.70.2 (by dependabot[bot])
- Bump github.com/urfave/cli from 1.22.4 to 1.22.5 (by dependabot-preview[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.69.2 to 0.69.3 (by dependabot-preview[bot])
- Update Dependabot config file (by dependabot-preview[bot])

* Wed Oct 28 2020 <mackerel-developers@hatena.ne.jp> - 0.41.0
- Print empty array if hosts not found (by kazeburo)
- Add armhf Debian package to GitHub release (by hnw)
- Bump alpine from 3.12.0 to 3.12.1 (by dependabot-preview[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.69.1 to 0.69.2 (by dependabot-preview[bot])

* Thu Oct 01 2020 <mackerel-developers@hatena.ne.jp> - 0.40.4
- Bump github.com/mackerelio/mackerel-client-go from 0.10.1 to 0.11.0 (by dependabot-preview[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.68.2 to 0.69.1 (by dependabot-preview[bot])

* Tue Sep 15 2020 <mackerel-developers@hatena.ne.jp> - 0.40.3
- revert go1.15-alpine (by lufia)
- revert changing filename (by lufia)
- Bump golang from 1.15.0-alpine to 1.15.2-alpine (by dependabot-preview[bot])
- Bump github.com/Songmu/prompter from 0.3.0 to 0.4.0 (by dependabot-preview[bot])
- Bump golang from 1.14-alpine to 1.15.0-alpine (by dependabot-preview[bot])
- add arm64 packages, and fix Architecture field of deb (by lufia)
- Bump github.com/mackerelio/mackerel-agent from 0.68.0 to 0.68.2 (by dependabot-preview[bot])

* Mon Jul 20 2020 <mackerel-developers@hatena.ne.jp> - 0.40.2
- Bump github.com/mackerelio/mackerel-client-go from 0.10.0 to 0.10.1 (by dependabot-preview[bot])
- define default file name while doMonitorsPull() at once (by hgsgtk)
- Update `Installation` section of README  (by astj)
- Bump github.com/stretchr/testify from 1.6.0 to 1.6.1 (by dependabot-preview[bot])
- Update go.sum (by shibayu36)
- Bump alpine from 3.11.6 to 3.12.0 (by dependabot-preview[bot])
- Bump github.com/mackerelio/mackerel-client-go from 0.9.1 to 0.10.0 (by dependabot-preview[bot])
- Bump github.com/stretchr/testify from 1.5.1 to 1.6.0 (by dependabot-preview[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.67.1 to 0.68.0 (by dependabot-preview[bot])

* Thu May 14 2020 <mackerel-developers@hatena.ne.jp> - 0.40.1
- Bump gopkg.in/yaml.v2 from 2.2.8 to 2.3.0 (by dependabot-preview[bot])
- Bump alpine from 3.11.5 to 3.11.6 (by dependabot-preview[bot])
- Build with Go 1.14 (by lufia)
- Bump github.com/urfave/cli from 1.22.3 to 1.22.4 (by dependabot-preview[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.67.0 to 0.67.1 (by dependabot-preview[bot])

* Fri Apr 03 2020 <mackerel-developers@hatena.ne.jp> - 0.40.0
- Bump alpine from 3.11.3 to 3.11.5 (by dependabot-preview[bot])
- Implement `mkr channels pull` command (by stefafafan)
- Bump github.com/urfave/cli from 1.22.2 to 1.22.3 (by dependabot-preview[bot])
- Stop building 32bit Darwin artifacts (by astj)
- Bump github.com/mackerelio/mackerel-client-go from 0.8.0 to 0.9.1 (by dependabot-preview[bot])
- Bump github.com/stretchr/testify from 1.4.0 to 1.5.1 (by dependabot-preview[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.65.0 to 0.67.0 (by dependabot-preview[bot])
- Bump gopkg.in/yaml.v2 from 2.2.7 to 2.2.8 (by dependabot-preview[bot])

* Wed Feb 05 2020 <mackerel-developers@hatena.ne.jp> - 0.39.7
- show mandantory flags usage in modern CLI manner (by aereal)
- rename: github.com/motemen/gobump -> github.com/x-motemen/gobump (by lufia)

* Wed Jan 22 2020 <mackerel-developers@hatena.ne.jp> - 0.39.6
- Bump github.com/pkg/errors from 0.8.1 to 0.9.1 (by dependabot-preview[bot])
- Bump alpine from 3.11.2 to 3.11.3 (by dependabot-preview[bot])
- Stop using alias imports for mackerel-client-go for simplicity (by stefafafan)
- Bump github.com/fatih/color from 1.8.0 to 1.9.0 (by dependabot-preview[bot])
- Bump github.com/fatih/color from 1.7.0 to 1.8.0 (by dependabot-preview[bot])
- Bump alpine from 3.11.0 to 3.11.2 (by dependabot-preview[bot])
- Bump alpine from 3.10.3 to 3.11.0 (by dependabot-preview[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.64.1 to 0.65.0 (by dependabot-preview[bot])

* Thu Dec 05 2019 <mackerel-developers@hatena.ne.jp> - 0.39.5
- Bump gopkg.in/yaml.v2 from 2.2.5 to 2.2.7 (by dependabot-preview[bot])
- Bump github.com/Songmu/prompter from 0.2.0 to 0.3.0 (by dependabot-preview[bot])
- Bump github.com/urfave/cli from 1.22.1 to 1.22.2 (by dependabot-preview[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.64.0 to 0.64.1 (by dependabot-preview[bot])

* Thu Nov 21 2019 <mackerel-developers@hatena.ne.jp> - 0.39.4
- always set GO111MODULE=on (by lufia)
- Bump alpine from 3.9 to 3.10.3 (by dependabot-preview[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.63.0 to 0.64.0 (by dependabot-preview[bot])
- Bump gopkg.in/yaml.v2 from 2.2.4 to 2.2.5 (by dependabot-preview[bot])

* Thu Oct 24 2019 <mackerel-developers@hatena.ne.jp> - 0.39.3
- Build with Go 1.12.12
- Bump gopkg.in/yaml.v2 from 2.2.3 to 2.2.4 (by dependabot-preview[bot])

* Wed Oct 02 2019 <mackerel-developers@hatena.ne.jp> - 0.39.2
- Bump gopkg.in/yaml.v2 from 2.2.2 to 2.2.3 (by dependabot-preview[bot])
- Bump github.com/urfave/cli from 1.22.0 to 1.22.1 (by dependabot-preview[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.62.1 to 0.63.0 (by dependabot-preview[bot])
- Bump github.com/mackerelio/mackerel-client-go from 0.7.0 to 0.8.0 (by dependabot-preview[bot])
- Bump github.com/mackerelio/mackerel-client-go from 0.6.0 to 0.7.0 (by dependabot-preview[bot])
- Bump github.com/urfave/cli from 1.21.0 to 1.22.0 (by dependabot-preview[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.62.0 to 0.62.1 (by dependabot-preview[bot])
- Bump github.com/stretchr/testify from 1.3.0 to 1.4.0 (by dependabot-preview[bot])

* Thu Aug 29 2019 <mackerel-developers@hatena.ne.jp> - 0.39.1
- rename gopkg.in/urfave/cli.v1 -> github.com/urfave/cli (by lufia)
- Bump github.com/mackerelio/mackerel-agent from 0.61.1 to 0.62.0 (by dependabot-preview[bot])
- Bump github.com/mackerelio/mackerel-agent from 0.60.0 to 0.61.1 (by dependabot-preview[bot])

* Mon Jul 22 2019 <mackerel-developers@hatena.ne.jp> - 0.39.0
- Implement anomaly detection monitor (by syou6162)
- Bump github.com/mackerelio/mackerel-agent from 0.59.3 to 0.60.0 (by dependabot-preview[bot])
- add fakeroot to build dependencies (by astj)
- upgrade builder image and enable go modules (by lufia)

* Tue Jun 11 2019 <mackerel-developers@hatena.ne.jp> - 0.38.0
- Move implementation of `mkr create` into mkr/hosts (by astj)
- Build With Go 1.12 (by astj)
- Update go module dependencies (by astj)
- support go modules (by lufia)

* Wed May 08 2019 <mackerel-developers@hatena.ne.jp> - 0.37.0
- [wrap] Improve message truncation algorithm (by itchyny)

* Wed Mar 27 2019 <mackerel-developers@hatena.ne.jp> - 0.36.0
- [monitors] support interruption monitoring thresholds of service metric monitors (by itchyny)
- Improve Makefile (by itchyny)
- Create services package and add tests for services subcommand (by itchyny)
- Create org package and add tests for org subcommand (by itchyny)
- Add tests for hosts package (by itchyny)
- Introduce Client interface for Mackerel API (by itchyny)
- Add out stream argument to format.PrettyPrintJSON (by itchyny)

* Wed Mar 06 2019 <mackerel-developers@hatena.ne.jp> - 0.35.1
- Fix status option of hosts subcommand (by itchyny)

* Wed Feb 13 2019 <mackerel-developers@hatena.ne.jp> - 0.35.0
- separate hosts package from main and reorganize package structure (by Songmu)
- add wrap subcommand to monitor batch jobs to run with cron etc (by Songmu)
- add `checks run` subcommand to confirm check plugin settings (by Songmu)

* Tue Nov 27 2018 <mackerel-developers@hatena.ne.jp> - 0.34.2
- Fix the default limit of mkr alerts (by itchyny)

* Mon Nov 26 2018 <mackerel-developers@hatena.ne.jp> - 0.34.1
- avoid out of range exception in mkr alerts (by astj)

* Mon Nov 26 2018 <mackerel-developers@hatena.ne.jp> - 0.34.0
- Fixed issue where nextID was not inherited (by yaminoma)
- Implemented according to new specification of Alerts API. (by yaminoma)

* Mon Nov 12 2018 <mackerel-developers@hatena.ne.jp> - 0.33.0
- Follow https://github.com/mholt/archiver's change (by astj)
- update subcommand should keep the interface infromation (by itchyny)
- show alert.Message as monitor message in mkr alerts list (by astj)
- implement org subcommand (by itchyny)

* Wed Oct 17 2018 <mackerel-developers@hatena.ne.jp> - 0.32.1
- Build with Go 1.11 (by astj)

* Thu Aug 30 2018 <mackerel-developers@hatena.ne.jp> - 0.32.0
- logs to Stderr (by Songmu)
- Omit symbol table and debug information from the executable (by itchyny)
- fix Dockerfile to create alpine based docker image (by hayajo)

* Wed Jul 04 2018 <mackerel-developers@hatena.ne.jp> - 0.31.1
- Fix mkr throw --retry not working (by astj)

* Wed Jul 04 2018 <mackerel-developers@hatena.ne.jp> - 0.31.0
- add Retry feature and --retry option to `mkr throw` (by astj)

* Wed Jun 20 2018 <mackerel-developers@hatena.ne.jp> - 0.30.0
- Build with Go 1.10 (by astj)

* Tue Apr 10 2018 <mackerel-developers@hatena.ne.jp> - 0.29.0
- Change createdAt of hosts subcommand output to ISO 8601 extended format (by hayajo)

* Wed Mar 28 2018 <mackerel-developers@hatena.ne.jp> - 0.28.0
- Add --upgrade option to plugin install. (by fujiwara)

* Thu Mar 15 2018 <mackerel-developers@hatena.ne.jp> - 0.27.1
- Add <direct_url> help (by shibayu36)

* Thu Mar 01 2018 <mackerel-developers@hatena.ne.jp> - 0.27.0
- Support empty threshold for monitors (by itchyny)

* Tue Jan 23 2018 <mackerel-developers@hatena.ne.jp> - 0.26.0
- Fix copying plugin archives while installing a plugin on Windows (by itchyny)
- update rpm-v2 task for building Amazon Linux 2 package (by hayajo)

* Wed Jan 10 2018 <mackerel-developers@hatena.ne.jp> - 0.25.0
- Add plugin document links (by shibayu36)
- introduce goxz and adjust deps (by Songmu)
- add appveyor.yml and adjust tests for windows (by Songmu)
- Define plugin default installation path in windows environment (by Songmu)

* Wed Dec 13 2017 <mackerel-developers@hatena.ne.jp> - 0.24.1
- Rebuild to avoid panic when action of check-plugin was not specified (by astj)

* Tue Dec 12 2017 <mackerel-developers@hatena.ne.jp> - 0.24.0
- Support maxCheckAttempts for host metric and service metric monitors (by itchyny)

* Tue Nov 28 2017 <mackerel-developers@hatena.ne.jp> - 0.23.0
- [plugin.install] support direct URL target (by Songmu)
- [plugin.install] support tarball archives (by Songmu)
- fix hostId option flag in the command help of throw and metrics commands (by astj)
- Refactor mkr plugin install implementation (by shibayu36)

* Thu Oct 26 2017 <mackerel-developers@hatena.ne.jp> - 0.22.0
- Release mkr plugin install (by shibayu36)
- Add metrics command (by edangelion)

* Wed Oct 04 2017 <mackerel-developers@hatena.ne.jp> - 0.21.0
- Use new API BaseURL (by astj)

* Wed Sep 27 2017 <mackerel-developers@hatena.ne.jp> - 0.20.0
- build with Go 1.9 (by astj)

* Wed Sep 20 2017 <mackerel-developers@hatena.ne.jp> - 0.19.0
- Support fetch command to retrieve many hosts (by itchyny)
- Prefer apibase in mackerel-agent confFile (by astj)

* Tue Sep 12 2017 <mackerel-developers@hatena.ne.jp> - 0.18.0
- add --customIdentifier option to mkr create (by astj)

* Wed Aug 23 2017 <mackerel-developers@hatena.ne.jp> - 0.17.0
- [dashboards] Add unit to expression graph (by edangelion)
- [dashboards] Add title param to expression graph (by edangelion)

* Wed Jun 07 2017 <mackerel-developers@hatena.ne.jp> - 0.16.1
- v2 packages (by Songmu)

* Tue May 09 2017 <mackerel-developers@hatena.ne.jp> - 0.16.0-1
- Add services subcommand (by yuuki)
