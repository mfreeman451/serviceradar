Name:           serviceradar-core
Version:        %{version}
Release:        %{release}%{?dist}
Summary:        ServiceRadar core service with web interface
License:        Proprietary

BuildRequires:  systemd
Requires:       systemd
%{?systemd_requires}

Source: systemd/serviceradar-core.service

%description
Provides centralized monitoring and web dashboard for ServiceRadar.

%install
mkdir -p %{buildroot}/usr/local/bin
mkdir -p %{buildroot}/etc/serviceradar
mkdir -p %{buildroot}/etc/serviceradar/checkers/sweep
mkdir -p %{buildroot}/lib/systemd/system
mkdir -p %{buildroot}/var/lib/serviceradar


install -m 755 %{_builddir}/serviceradar-core %{buildroot}/usr/local/bin/
install -m 644 %{_sourcedir}/systemd/serviceradar-core.service %{buildroot}/lib/systemd/system/serviceradar-core.service
install -m 644 %{_sourcedir}/config/core.json %{buildroot}/etc/serviceradar/
install -m 644 %{_sourcedir}/config/checkers/sweep/sweep.json %{buildroot}/etc/serviceradar/checkers/sweep/


%files
%attr(0755, root, root) /usr/local/bin/serviceradar-core
%config(noreplace) %attr(0644, serviceradar, serviceradar) /etc/serviceradar/core.json
%config(noreplace) %attr(0644, serviceradar, serviceradar) /etc/serviceradar/checkers/sweep/sweep.json
%attr(0644, root, root) /lib/systemd/system/serviceradar-core.service
%dir %attr(0755, root, root) /etc/serviceradar
%dir %attr(0755, serviceradar, serviceradar) /var/lib/serviceradar

%pre
# Create serviceradar user if it doesn't exist
if ! id -u serviceradar >/dev/null 2>&1; then
    useradd --system --no-create-home --shell /usr/sbin/nologin serviceradar
fi

%post
%systemd_post serviceradar-core.service

# Check if firewalld is running
if systemctl is-active firewalld.service >/dev/null 2>&1; then
    firewall-cmd --permanent --add-port=8090/tcp --zone=public >/dev/null 2>&1 # Adjust zone
    firewall-cmd --reload >/dev/null 2>&1
    logger "Configured firewalld for ServiceRadar (port 8090/tcp, zone public)."
else
    logger "Firewalld is not running. Skipping firewall configuration."
fi

# SELinux (if needed)
if getenforce | grep -q "Enforcing"; then
    semanage -a -t http_port_t -p tcp 8090 >/dev/null 2>&1
    restorecon -Rv /usr/local/bin/serviceradar-core
    logger "Configured SELinux for ServiceRadar (port 8090/tcp)."
fi

%preun
%systemd_preun serviceradar-core.service

%postun
%systemd_postun_with_restart serviceradar-core.service