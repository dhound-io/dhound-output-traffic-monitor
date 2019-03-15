# create  dhound group
if ! getent group dhound >/dev/null; then
  groupadd -r dhound
fi

# create  dhound user
if ! getent passwd dhound >/dev/null; then
  useradd -r -g  dhound -d /opt/dhound-output-traffic-monitor \
    -s /sbin/nologin -c "dhound Service user" dhound
fi
