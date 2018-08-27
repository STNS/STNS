Summary: SimpleTomlNameService is Linux User,Group Name Service
Name:             stns-v2
Version:          0.0.1
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
install -m 644 package/stns-v2.initd  %{buildroot}%{_sysconfdir}/init.d/stns
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
/etc/stns/server/stns.conf
/etc/logrotate.d/stns

%if 0%{?rhel} < 7
/etc/init.d/stns
%else
/etc/systemd/system/stns.service
%endif

%changelog
* Sun Aug 26 2018 pyama86 <www.kazu.com@gmail.com> - 0.1.0-1
- Initial packaging
