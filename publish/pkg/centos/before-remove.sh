if [ $1 -eq 0 ]; then
  /sbin/service dhound-output-traffic-monitor stop >/dev/null 2>&1 || true
  /sbin/chkconfig --del dhound-output-traffic-monitor
fi
