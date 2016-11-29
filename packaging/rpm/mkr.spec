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
