%define _localbindir /usr/local/bin
%define _binaries_in_noarch_packages_terminate_build   0
Summary: SimpleTomlNameService is Linux User,Group Name Service
Name: stns
Group: SipmleTomlNameService
URL: https://github.com/pyama86/SimpleTomlNameService
Version: 0.1
Release: 1
License: MIT
Source0:   %{name}.initd
Source2:   %{name}.logrotate
Source3:   %{name}.conf
Packager:  stns
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root
Requires(post): /sbin/chkconfig
Requires(preun): /sbin/chkconfig, /sbin/service
Requires(postun): /sbin/service

%description
Api of SimpleTomlNameService generate json from toml config file

%install
rm -rf %{buildroot}
install -d -m 755 %{buildroot}/%{_localbindir}
install    -m 655 %{_builddir}/%{name}  %{buildroot}/%{_localbindir}

install -d -m 755 %{buildroot}/%{_localstatedir}/log/

install -d -m 755 %{buildroot}/%{_initrddir}
install    -m 755 %{_sourcedir}/%{name}.initd        %{buildroot}/%{_initrddir}/%{name}


install -d -m 755 %{buildroot}/%{_sysconfdir}/logrotate.d/
install    -m 644 %{_sourcedir}/%{name}.logrotate %{buildroot}/%{_sysconfdir}/logrotate.d/%{name}

install -d -m 755 %{buildroot}/%{_sysconfdir}/%{name}/
install    -m 644 %{_sourcedir}/%{name}.conf %{buildroot}/%{_sysconfdir}/%{name}/%{name}.conf

%clean
rm -rf %{_builddir}/*
rm -rf %{buildroot}

%post
chkconfig --add %{name}

%preun
if [ $1 = 0 ]; then
  service %{name} stop > /dev/null 2>&1
  chkconfig --del %{name}
fi

%files
%defattr(-,root,root)
%{_initrddir}/%{name}
%{_localbindir}/%{name}
%config(noreplace) %{_sysconfdir}/%{name}/%{name}.conf
%{_sysconfdir}/logrotate.d/%{name}
