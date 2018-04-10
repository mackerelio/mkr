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
