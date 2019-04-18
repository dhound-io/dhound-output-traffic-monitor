#!/bin/sh

REPOFOLDER=/var/dhound.io/repository-dotm

# open packages directory
cd packages/

# publish debian packages
find -name \*.deb -exec reprepro --ignore=undefinedtarget -Vb $REPOFOLDER/deb includedeb dhound-output-traffic-monitor {} \;

# public rpm packages
cp *.rpm $REPOFOLDER/rpm/
createrepo --outputdir=$REPOFOLDER/rpm/ . --update

sudo chmod -R ugo+rX $REPOFOLDER