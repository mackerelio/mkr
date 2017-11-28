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
