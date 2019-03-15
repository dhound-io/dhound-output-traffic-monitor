#!/bin/sh

chown -R dhound:dhound /opt/dhound-output-traffic-monitor
chown dhound-agent /var/log/dhound
chown dhound:dhound /var/lib/dhound-output-traffic-monitor

update-rc.d dhound-output-traffic-monitor defaults

echo "Logs for dhound-output-traffic-monitor will be in /var/log/dhound/"
