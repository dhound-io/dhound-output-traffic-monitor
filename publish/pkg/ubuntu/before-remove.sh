#!/bin/sh

if [ $1 = "remove" ]; then
  service dhound-output-traffic-monitor stop >/dev/null 2>&1 || true
  update-rc.d -f dhound-output-traffic-monitor remove
fi
