#!/bin/sh

# check official debian ports https://www.debian.org/ports/

cd ..

env CC=gcc GOOS=linux GOARCH=386 CGO_ENABLED=1 GOGCCFLAGS="-fPIC -m32 -fmessage-length=0" CGO_LDFLAGS="-L/dhound-build-server/dhound-output-traffic-monitor/publish/libpcap/i386"  ARCHITECTURE=i386 make deb
env CC=gcc GOOS=linux GOARCH=amd64 CGO_ENABLED=1 GOGCCFLAGS="-fPIC -m64 -fmessage-length=0" ARCHITECTURE=amd64  make deb
env CC=gcc GOOS=linux GOARCH=amd64 CGO_ENABLED=1 GOGCCFLAGS="-fPIC -m64 -fmessage-length=0" ARCHITECTURE=ia64 make deb

mv *.deb publish/pakages
rm dhound-output-traffic-monitor