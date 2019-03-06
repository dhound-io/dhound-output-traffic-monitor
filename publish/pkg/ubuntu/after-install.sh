#!/bin/sh

chown -R dhound-agent:dhound-agent /opt/dhound-output-traffic-monitor
chown dhound-agent /var/log/dhound-agent
chown dhound-agent:dhound-agent /var/lib/dhound-output-traffic-monitor

update-rc.d dhound-agent defaults

echo "Logs for dhound-output-traffic-monitor will be in /var/log/dhound-agent/"
cd d