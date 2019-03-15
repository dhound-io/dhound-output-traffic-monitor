#!/bin/sh

# create dhound group
if ! getent group dhound >/dev/null; then
  groupadd -r dhound
fi

# create dhound user
if ! getent passwd dhound >/dev/null; then
  useradd -M -r -g dhound -d /var/lib/dhound-output-traffic-monitor \
    -s /usr/sbin/nologin -c "dhound.io Service User" dhound
fi
