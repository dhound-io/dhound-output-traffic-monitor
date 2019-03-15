RED='\033[0;31m'
GREEN='\033[0;32m'
LGRAY='\033[0;37m'
DGRAY='\033[0;30m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color
/sbin/chkconfig --add dhound-output-traffic-monitor

chown -R dhound:dhound /opt/dhound-output-traffic-monitor
chown dhound /var/log/dhound
chown  dhound:dhound /var/lib/dhound-output-traffic-monitor

echo -e "Logs for dhound-output-traffic-monitor will be in ${GREEN}/var/log/dhound/${NC}"
echo -e "${BLUE}How to${NC}"
echo -e "${LGRAY}sudo service dhound-output-traffic-monitor status${NC} - application status"
echo -e "${GREEN}sudo service dhound-output-traffic-monitor start${NC} - start application"
echo -e "${RED}sudo service dhound-output-traffic-monitor status${NC} - stop application"

service dhound-output-traffic-monitor start

echo -e "${GREEN}support@dhound.io${NC}"