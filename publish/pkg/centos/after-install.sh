/sbin/chkconfig --add dhound-output-traffic-monitor

chown -R dhound:dhound /opt/dhound-output-traffic-monitor
chown dhound /var/log/dhound
chown  dhound:dhound /var/lib/dhound-output-traffic-monitor

echo "Logs for dhound-output-traffic-monitor will be in /var/log/dhound/"
