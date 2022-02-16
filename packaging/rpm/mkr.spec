# sudo yum -y install rpmdevtools go && rpmdev-setuptree
# rpmbuild -ba ~/rpmbuild/SPECS/mkr.spec

%define _binaries_in_noarch_packages_terminate_build   0
%define _localbindir /usr/local/bin

Name:      mkr
Version:   %{_version}
Release:   1
License:   Apache-2.0
Summary:   mackerel.io api client tool
URL:       https://mackerel.io
Group:     Hatena
Packager:  Hatena
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root

%description
mackerel.io api client tool

%prep

%build

%install
rm -rf %{buildroot}
install -d -m 755 %{buildroot}/%{_localbindir}
install    -m 655 %{_builddir}/%{name}  %{buildroot}/%{_localbindir}
install -d -m 755 %{buildroot}/%{_localstatedir}/log/

%clean
rm -f %{buildroot}%{_bindir}/${name}

%pre

%post

%preun

%files
%defattr(-,root,root)
%{_localbindir}/%{name}

%changelog
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

* Thu Apr 06 2017 <mackerel-developers@hatena.ne.jp> - 0.15.0-1
- bump Go to 1.8 (by astj)

* Mon Mar 27 2017 <mackerel-developers@hatena.ne.jp> - 0.14.5-1
- Colors on Windows (by mattn)

* Wed Mar 22 2017 <mackerel-developers@hatena.ne.jp> - 0.14.4-1
- use new bot token (by daiksy)
- use new bot token (by daiksy)
- Workaround for git fetch failure (by daiksy)
- Apply git fetch workaround (by astj)

* Thu Feb 16 2017 <mackerel-developers@hatena.ne.jp> - 0.14.3-1
- Support annotations command for graph annotation (by syou6162)
- Improve help management and fix usage help for command (by haya14busa)
- remove unused functions (by haya14busa)

* Wed Feb 08 2017 <mackerel-developers@hatena.ne.jp> - 0.14.2-1
- [monitors diff] Add the "to-remote" bool flag (by yoheimuta)

* Wed Jan 11 2017 <mackerel-developers@hatena.ne.jp> - 0.14.1-1
- formatter.NewAsciiFormatter now needs config (by astj)

* Wed Dec 21 2016 <mackerel-developers@hatena.ne.jp> - 0.14.0-1
- Support expression monitor alerts in mkr alerts list (by itchyny)

* Tue Nov 29 2016 <mackerel-developers@hatena.ne.jp> - 0.13.0-1
- remove unreachable code: monitor type cannot be "check" (by haya14busa)
- Fix the links to the api documents (by itchyny)
- catch up monitor interface changes of mackerel-client-go (by haya14busa)
- Introduce yudai/gojsondiff for `mkr monitors diff` (by haya14busa)
- fix test according to mackerel-client-go changes (by haya14busa)

* Thu Oct 27 2016 <mackerel-developers@hatena.ne.jp> - 0.12.0-1
- Rename a dependent package (by usk81)
- Support `-apibase` option (by astj)
- [breaking change] Prepend `custom.` prefix to host metric name by default (by astj)

* Thu Jul 14 2016 <mackerel-developers@hatena.ne.jp> - 0.11.3-1
- fix `validateRules()`,  when monitor has rule of "expression". (by daiksy)

* Thu Jun 23 2016 <mackerel-developers@hatena.ne.jp> - 0.11.2-1
- replace angle brackets for json (by daiksy)

* Fri Jun 10 2016 <mackerel-developers@hatena.ne.jp> - 0.11.1-1
- fix version number (by stanaka)

* Thu Jun 09 2016 <mackerel-developers@hatena.ne.jp> - 0.11.0-1
- add dashboard generator (by daiksy)
- Add flag to overwrite host's roles  (by haya14busa)

* Wed May 25 2016 <mackerel-developers@hatena.ne.jp> - 0.10.1-1
- fix signnatures. codegangsta/cli (by tknzk)

* Tue May 10 2016 <mackerel-developers@hatena.ne.jp> - 0.10.0-1
- support `isMute` field of monitors (by Songmu)
- support boolean at isEmpty (by stanaka)
- bump up go version to 1.6.2 (by stanaka)

* Fri Mar 25 2016 <y.songmu@gmail.com> - 0.9.1-1
- use GOARCH=amd64 for now (by Songmu)

* Thu Feb 18 2016 <stefafafan@hatena.ne.jp> - 0.9.0-1
- Support displayName of host's json (by stanaka)

* Thu Jan 07 2016 <y.songmu@gmail.com> - 0.8.1-1
- fix handling host-status option (by stanaka)

* Wed Jan 06 2016 <y.songmu@gmail.com> - 0.8.0-1
- support alerts subcommand (by stanaka)
- Fix README example about mkr throw (by yuuki1)
- Build with Go 1.5 (by itchyny)
- Fixed the english used in the command descriptions (by stefafafan)

* Thu Nov 12 2015 <y.songmu@gmail.com> - 0.7.1-1
- support `notificationIntervai` field in monitors (stanaka)
- [bug] fix json parameter s/hostID/hostId/g (Songmu)

* Mon Oct 26 2015 <daiksy@hatena.ne.jp> - 0.7.0-1
- append newline to the end of monitors.json (by Songmu)
- fix printMonitor (by Songmu)
- fix diff output between slices (by Songmu)

* Thu Oct 15 2015 <itchyny@hatena.ne.jp> - 0.6.0-1
- Fix update command bug about overwriting hostname (by y_uuki)
- Stop the parallel request sending temporarily (by y_uuki)
- Suppress to display empty fields when mkr monitors diff (by stanaka)

* Mon Sep 14 2015 <itchyny@hatena.ne.jp> - 0.5.0-1
- add fields for external URL monitors (by stanaka)

* Fri Aug 28 2015 <tomohiro68@gmail.com> - 0.4.1-1
- Create deb/rpm package for Linux release (by Sixeight)

* Fri Aug 14 2015 <sixeight@hatena.ne.jp> - 0.4.0-1
- first release for rpm
