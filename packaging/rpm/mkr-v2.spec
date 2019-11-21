Name:      mkr
Version:   %{_version}
Release:   1%{?dist}
License:   ASL 2.0
Summary:   macekrel.io api client tool
URL:       https://mackerel.io
Group:     Application/System
Packager:  Hatena
BuildArch: %{buildarch}
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root

%description
macekrel.io api client tool

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
* Thu Nov 21 2019 <mackerel-developers@hatena.ne.jp> - 0.39.4
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
