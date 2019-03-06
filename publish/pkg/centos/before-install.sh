# create  dhound-agent group
if ! getent group dhound-agent >/dev/null; then
  groupadd -r dhound-agent
fi

# create  dhound-agent user
if ! getent passwd dhound-agent >/dev/null; then
  useradd -r -g  dhound-agent -d /opt/dhound-output-traffic-monitor \
    -s /sbin/nologin -c "dhound-agent Service user" dhound-agent
fi
