#!/bin/sh

RED='\033[0;31m'
GREEN='\033[0;32m'
LGRAY='\033[0;37m'
DGRAY='\033[0;30m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

chown -R dhound:dhound /opt/dhound-output-traffic-monitor
chown dhound-agent /var/log/dhound
chown dhound:dhound /var/lib/dhound-output-traffic-monitor

update-rc.d dhound-output-traffic-monitor defaults

echo "Logs for dhound-output-traffic-monitor will be in ${GREEN}/var/log/dhound/${NC}"
echo "${BLUE}How to${NC}"
echo "${LGRAY}sudo service dhound-output-traffic-monitor status${NC} - application status"
echo "${GREEN}sudo service dhound-output-traffic-monitor start${NC} - start application"
echo "${RED}sudo service dhound-output-traffic-monitor status${NC} - stop application"

service dhound-output-traffic-monitor start

echo "${GREEN}support@dhound.io${NC}"
