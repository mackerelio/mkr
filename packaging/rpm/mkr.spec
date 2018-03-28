# sudo yum -y install rpmdevtools go && rpmdev-setuptree
# rpmbuild -ba ~/rpmbuild/SPECS/mkr.spec

%define _binaries_in_noarch_packages_terminate_build   0
%define _localbindir /usr/local/bin

Name:      mkr
Version:   %{_version}
Release:   1
License:   Apache-2.0
Summary:   macekrel.io api client tool
URL:       https://mackerel.io
Group:     Hatena
Packager:  Hatena
BuildArch: %{buildarch}
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root

%description
macekrel.io api client tool

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
