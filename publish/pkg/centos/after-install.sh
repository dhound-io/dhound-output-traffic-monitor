/sbin/chkconfig --add dhound-output-traffic-monitor

chown -R dhound-agent: dhound-agent /opt/dhound-output-traffic-monitor
chown dhound-agent /var/log/dhound
chown  dhound-agent: dhound-agent /var/lib/dhound-output-traffic-monitor

echo "Logs for dhound-output-traffic-monitor will be in /var/log/dhound-logs/"
