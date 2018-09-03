Summary:          SimpleTomlNameService Nss Module
Name:             libnss-stns-v2 
Version:          1.0.0
Release:          1
License:          GPLv3
URL:              https://github.com/STNS/STNS
Source:           %{name}-%{version}.tar.gz
Group:            System Environment/Base
Packager:         pyama86 <www.kazu.com@gmail.com>
%if 0%{?rhel} < 6
Requires:         glibc curl-devel jansson-devel
%else
Requires:         glibc libcurl-devel jansson-devel
%endif
BuildRequires:    gcc make
BuildRoot:        %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)
BuildArch:        i386, x86_64

%define debug_package %{nil}

%description
We provide name resolution of Linux user group using STNS.

%prep
%setup -q -n %{name}-%{version}

%build
make

%install
%{__rm} -rf %{buildroot}
mkdir -p %{buildroot}/usr/{lib64,bin}
mkdir -p %{buildroot}%{_sysconfdir}
make PREFIX=%{buildroot}/usr install
install -d -m 777 %{buildroot}/var/cache/stns
install -d -m 744 %{buildroot}%{_sysconfdir}/stns/client/
install -m 644 stns.conf.example %{buildroot}%{_sysconfdir}/stns/client/stns.conf

%clean
%{__rm} -rf %{buildroot}

%post

%preun

%postun

%files
%defattr(-, root, root)
/usr/lib64/libnss_stns.so
/usr/lib64/libnss_stns.so.2
/usr/lib64/libnss_stns.so.2.0
/usr/lib/stns/stns-key-wrapper
/var/cache/stns
/etc/stns/client/stns.conf

%changelog
* Mon Sep 3 2018 pyama86 <www.kazu.com@gmail.com> - 1.0.0-1
- Release
* Mon Aug 27 2018 pyama86 <www.kazu.com@gmail.com> - 0.0.1-1
- Initial packaging
