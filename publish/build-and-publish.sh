#!/bin/sh

# Remove old packages
./clean.sh
# Build .deb packages
./build-debian.sh
# Build .rpm packages
./build-rpm.sh
# Publish and reconfigure repository
./publish-into-repo.sh
