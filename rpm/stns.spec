Summary: SimpleTomlNameService is Linux User,Group Name Service
Name:             stns-v2
Version:          2.2.11
Release:          1
License:          GPLv3
URL:              https://github.com/STNS/STNS
Group:            System Environment/Base
Packager:         pyama86 <www.kazu.com@gmail.com>
Source:           %{name}-%{version}.tar.gz
BuildRequires:    make
BuildRoot:        %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)
BuildArch:        i386, x86_64

%ifarch x86_64
%global gohostarch  amd64
%endif
%ifarch %{ix86}
%global gohostarch  386
%endif
%ifarch %{arm}
%global gohostarch  arm
%endif
%ifarch aarch64
%global gohostarch  arm64
%endif
%define debug_package %{nil}

%description
This server can easily manage Linux user group with TOML format setting file.

%prep
%setup -q -n %{name}-%{version}

%build
export GOOS=linux
export GOARCH=%{gohostarch}
make

%install
%{__rm} -rf %{buildroot}
mkdir -p %{buildroot}/usr/sbin
mkdir -p %{buildroot}%{_sysconfdir}/stns/server
make PREFIX=%{buildroot}/usr/ install
install -m 644 package/stns-v2.conf %{buildroot}%{_sysconfdir}/stns/server/stns.conf

%if 0%{?rhel} < 7
mkdir -p %{buildroot}%{_sysconfdir}/init.d
install -m 755 package/stns-v2.initd  %{buildroot}%{_sysconfdir}/init.d/stns
%else
mkdir -p %{buildroot}%{_sysconfdir}/systemd/system/
install -m 644 package/stns-v2.systemd %{buildroot}%{_sysconfdir}/systemd/system/stns.service
%endif

mkdir -p %{buildroot}%{_sysconfdir}/logrotate.d
install -m 644 package/stns-v2.logrotate %{buildroot}%{_sysconfdir}/logrotate.d/stns

%clean
%{__rm} -rf %{buildroot}

%post

%preun

%postun

%files
%defattr(-, root, root)
/usr/sbin/stns
%config(noreplace) /etc/stns/server/stns.conf
/usr/local/stns/modules.d/mod_stns_etcd.so
/usr/local/stns/modules.d/mod_stns_dynamodb.so
/etc/logrotate.d/stns

%if 0%{?rhel} < 7
/etc/init.d/stns
%else
/etc/systemd/system/stns.service
%endif

%changelog
* Tue Oct 12 2021 pyama86 <www.kazu.com@gmail.com> - 2.2.11-1
- Fix typo
* Fri Dec 18 2020 pyama86 <www.kazu.com@gmail.com> - 2.2.10-1
- Add healthcheck endpoint to v1
* Tue Nov 24 2020 pyama86 <www.kazu.com@gmail.com> - 2.2.9-1
- Add Parameter IP Filter
* Tue Oct 27 2020 pyama86 <www.kazu.com@gmail.com> - 2.2.8-1
- FIX should shutdown when STNS recived signals
* Tue Jul 21 2020 pyama86 <www.kazu.com@gmail.com> - 2.2.7-1
- delete fatal methods
* Tue Jul 14 2020 pyama86 <www.kazu.com@gmail.com> - 2.2.6-1
- use Go 1.14.4
- support debian8/9,CentOS8
* Thu Apr 5 2019 pyama86 <www.kazu.com@gmail.com> - 2.2.5-1
- #108 Not check when user password empty
* Thu Apr 4 2019 pyama86 <www.kazu.com@gmail.com> - 2.2.4-1
- #107 Config may be empty
* Thu Apr 4 2019 pyama86 <www.kazu.com@gmail.com> - 2.2.3-1
- #105 check diff before update
* Thu Apr 4 2019 pyama86 <www.kazu.com@gmail.com> - 2.2.2-1
- #104 Add password envs
* Thu Apr 4 2019 pyama86 <www.kazu.com@gmail.com> - 2.2.1-1
- #103 Listen port
* Mon Apr 1 2019 pyama86 <www.kazu.com@gmail.com> - 2.2.0-1
- #92 Support LDAP Interface
- #93 Update password
- #95 HTTP(S) request log も stns.logへ出力する
- #96 Fix the problem of password update
- #97 support ldaps
- #98 support read config from s3
- #99 Allow only a single backend
- #100 support dynamodb
- #101 add redis backend
* Wed Feb 27 2019 pyama86 <www.kazu.com@gmail.com> - 2.1.1-1
- #89 typo
* Thu Feb 21 2019 pyama86 <www.kazu.com@gmail.com> - 2.1.0-1
- #88 Support TLS Authentication
* Thu Nov 29 2018 pyama86 <www.kazu.com@gmail.com> - 2.0.3-1
- #84 add config check with systemd
* Sat Nov 10 2018 pyama86 <www.kazu.com@gmail.com> - 2.0.2-1
- #81 add checkconfig command
* Wed Oct 3 2018 pyama86 <www.kazu.com@gmail.com> - 2.0.1-1
- #77 add modole to package
* Wed Oct 3 2018 pyama86 <www.kazu.com@gmail.com> - 2.0.0-1
- #75 Support etcd Backend
* Thu Sep 20 2018 pyama86 <www.kazu.com@gmail.com> - 1.0.0-3
- #74 forget } at the end.
* Thu Sep 11 2018 pyama86 <www.kazu.com@gmail.com> - 1.0.0-2
- #70 Logger aggregates into gommon
* Mon Sep 3 2018 pyama86 <www.kazu.com@gmail.com> - 1.0.0-1
- Release
* Sun Aug 26 2018 pyama86 <www.kazu.com@gmail.com> - 0.1.0-1
- Initial packaging
