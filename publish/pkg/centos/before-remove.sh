if [ $1 -eq 0 ]; then
  /sbin/service dhound-agent stop >/dev/null 2>&1 || true
  /sbin/chkconfig --del dhound-agent
fi
