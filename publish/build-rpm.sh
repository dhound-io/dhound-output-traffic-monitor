#!/bin/sh

cd ..

env CC=gcc GOOS=linux GOARCH=386 CGO_ENABLED=1 GOGCCFLAGS="-fPIC -m32 -fmessage-length=0" ARCHITECTURE=i386 CGO_LDFLAGS="-L/dhound-build-server/dhound-output-traffic-monitor/publish/libpcap/i386"  make rpm
env CC=gcc GOOS=linux GOARCH=amd64 CGO_ENABLED=1 GOGCCFLAGS="-fPIC -m64 -fmessage-length=0" ARCHITECTURE=amd64 make rpm
env CC=gcc GOOS=linux GOARCH=amd64 CGO_ENABLED=1 GOGCCFLAGS="-fPIC -m64 -fmessage-length=0" ARCHITECTURE=ia64 make rpm

mv *.rpm publish
cd publish/

echo 'Sign rpm packages'
for line in $(find . -iname '*.rpm'); do
   sign-rpm.sh ""  "$line"
   echo $line
   rpm --checksig "$line"
done

mv *.rpm packages/
