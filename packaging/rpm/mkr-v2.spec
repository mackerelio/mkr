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
