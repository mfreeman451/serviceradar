Name:           serviceradar-cloud
Version:        %{version}
Release:        %{release}%{?dist}
Summary:        ServiceRadar cloud service with web interface
License:        Proprietary

BuildRequires:  systemd
Requires:       systemd
%{?systemd_requires}

Source: systemd/serviceradar-cloud.service  # Corrected: Added Source tag

%description
Provides centralized monitoring and web dashboard for ServiceRadar.

%install
mkdir -p %{buildroot}/usr/local/bin
mkdir -p %{buildroot}/etc/serviceradar
mkdir -p %{buildroot}/lib/systemd/system
mkdir -p %{buildroot}/var/lib/serviceradar


install -m 755 %{_builddir}/serviceradar-cloud %{buildroot}/usr/local/bin/
install -m 644 %{_sourcedir}/systemd/serviceradar-cloud.service %{buildroot}/lib/systemd/system/serviceradar-cloud.service
install -m 644 %{_sourcedir}/config/cloud.json %{buildroot}/etc/serviceradar/


%files
%attr(0755, root, root) /usr/local/bin/serviceradar-cloud
%config(noreplace) %attr(0644, serviceradar, serviceradar) /etc/serviceradar/cloud.json
%attr(0644, root, root) /lib/systemd/system/serviceradar-cloud.service
%dir %attr(0755, root, root) /etc/serviceradar
%dir %attr(0755, serviceradar, serviceradar) /var/lib/serviceradar

%pre
# Create serviceradar user if it doesn't exist
if ! id -u serviceradar >/dev/null 2>&1; then
    useradd --system --no-create-home --shell /usr/sbin/nologin serviceradar
fi

%post
%systemd_post serviceradar-cloud.service

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
    restorecon -Rv /usr/local/bin/serviceradar-cloud # Adjust path
    logger "Configured SELinux for ServiceRadar (port 8090/tcp)."
fi

%preun
%systemd_preun serviceradar-cloud.service

%postun
%systemd_postun_with_restart serviceradar-cloud.service