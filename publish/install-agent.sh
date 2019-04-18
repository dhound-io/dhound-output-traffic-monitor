#!/usr/bin/env bash
set -o pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
LGRAY='\033[0;37m'
DGRAY='\033[0;30m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

### Based on scripts developed by boundary
### Copyright 2011-2013, Boundary
### Copyright 2016-2017, DHound
### Licensed under the Apache License, Version 2.0 (the "License");
### you may not use this file except in compliance with the License.
### You may obtain a copy of the License at
###
###     http://www.apache.org/licenses/LICENSE-2.0
###
### Unless required by applicable law or agreed to in writing, software
### distributed under the License is distributed on an "AS IS" BASIS,
### WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
### See the License for the specific language governing permissions and
### limitations under the License.
###

PLATFORMS=("Ubuntu" "Debian" "CentOS" "RHEL" "CloudLinux" "Fedora" "Amazon" "Oracle")

# Put minimum version numbers here.
Ubuntu_VERSION_MIN="12"
Debian_VERSION_MIN="7"
CentOS_VERSION_MIN="5"
RHEL_VERSION_MIN="5"
CloudLinux_VERSION_MIN="5"
Oracle_VERSION_MIN="6"
Fedora_VERSION_MIN="21"
Amazon_VERSION_MIN=""

APIHOST="https://gate.dhound.io"
DEBREPOSITORY="https://repository.dhound.io/deb"
RPMREPOSITORY="https://repository.dhound.io/rpm"
SIGNKEY="https://repository.dhound.io/dhound-agent-packaging-public.key"

TARGET_DIR="/etc/dhound-output-traffic-monitor"

SUPPORTED_ARCH=1
SUPPORTED_PLATFORM=0
UPDATEAGENT=0

APT_CMD="apt-get -q -y"
YUM_CMD="yum -d0 -e0 -y --nogpgcheck"


trap "exit" INT TERM EXIT

function print_supported_platforms() {
    echo "Supported platforms are:"
    for d in ${PLATFORMS[*]}
    do
        echo -n " * $d: > "
        eval echo "\${${d}_VERSION_MIN}"
    done
}

function print_help() {
    echo "  Arguments:"
    echo "      -u: update agent to latest version"
    echo "      -h: help"
    echo "  Examples:"
    echo "      install-agent.sh -a 1234567890 -s 1234567890"
    echo "      install-agent.sh -u"
    exit 0
}

function do_install() {

    check_url_status "$SIGNKEY"
echo
    if [ "$DISTRO" = "Ubuntu" ] || [ $DISTRO = "Debian" ]; then

        echo -e "${GREEN}Adding repository $DEBREPOSITORY ${NC}"
        #sh -c "echo \"deb $DEBREPOSITORY dhound.io main\" | sudo tee /etc/apt/sources.list.d/dhound.list"
         echo -e "deb $DEBREPOSITORY dhound-agent main\ndeb $DEBREPOSITORY dhound-output-traffic-monitor main" | sudo tee /etc/apt/sources.list.d/dhound.list

        $CURL -Ls $SIGNKEY | sudo apt-key add - > /dev/null
        if [ $? -gt 0 ]; then
            echo "Error downloading GPG key from $SIGNKEY"
            exit 1
        fi

        echo -e "${GREEN}Updating apt repository cache...${NC}"
        $APT_CMD update -o Dir::Etc::sourcelist=/etc/apt/sources.list.d/dhound.list -o Dir::Etc::sourceparts="-" -o APT::Get::List-Cleanup="0"
        echo -e "${GREEN}Updating apt repositories...${NC}"
        apt update > /dev/null
        if [ $UPDATEAGENT -eq 0 ]; then
            echo -e "${GREEN}Installing dhound-output-traffic-monitor...${NC}"
            $APT_CMD install dhound-output-traffic-monitor > /dev/null
            echo -e "${GREEN}Installation finished${NC}"
        else
            echo -e "${GREEN}Updating dhound-output-traffic-monitor...${NC}"
            $APT_CMD install --only-upgrade dhound-output-traffic-monitor
            echo -e "${GREEN}Update finished${NC}"
        fi

    elif [ "$DISTRO" = "CentOS" ] || [ $DISTRO = "Amazon" ] || [ $DISTRO = "RHEL" ] || [ $DISTRO = "Oracle" ]; then
        if [ $UPDATEAGENT -eq 0 ]; then
            GPG_KEY_LOCATION=/etc/pki/rpm-gpg/RPM-GPG-KEY-DHound

            echo "Adding repository ${RPMREPOSITORY}"

            sh -c "cat - > /etc/yum.repos.d/dhound.repo <<EOF
[repository.dhound.io]
name=dhound-agent
baseurl=https://repository.dhound.io/rpm
failovermethod=priority
enabled=1
gpgcheck=1
metadata_expire=300
gpgkey=file://$GPG_KEY_LOCATION

[repository.dhound.io]
name=dhound-output-traffic-monitor
baseurl=https://repository.dhound.io/rpm
failovermethod=priority
enabled=1
gpgcheck=1
metadata_expire=300
gpgkey=file://$GPG_KEY_LOCATION
EOF"
            $CURL -s "$SIGNKEY"  | tee "$GPG_KEY_LOCATION" > /dev/null
            if [ $? -gt 0 ]; then
                echo "Error downloading GPG key from $SIGNKEY!"
                exit 1
            fi

            rpm --import "$GPG_KEY_LOCATION"
            $YUM_CMD install dhound-output-traffic-monitor

        else
            $YUM_CMD install dhound-output-traffic-monitor
        fi
    fi

    echo -e "${BLUE}Restarting dhound-output-traffic-monitor...${NC}"
    service dhound-output-traffic-monitor restart

    if [ $? -gt 0 ]; then
      echo -e "${RED}dhound-output-traffic-monitor installation failed.${NC}"
      exit 1
    fi

    echo ""
    if [ $UPDATEAGENT -eq 1 ]; then
       echo -e "${GREEN}dhound-output-traffic-monitor has been updated successfully!${NC}"
    else
       echo -e "${GREEN}dhound-output-traffic-monitor has been installed successfully!${NC}"
    fi
    echo -e "${BLUE}dhound-output-traffic-monitor output information can be found in the file: /var/log/dhound/dhound-output-traffic-monitor.log${NC}"
}

# 1st parameter - url, 2nd - error
function check_url_status()
{
    url="$1"
    status=$(curl --write-out %{http_code} -silent --output /dev/null $url)
    error="$2"
    if test $status -ne 200; then
        echo -e "${RED}Failed loading url: $url (status: $status)${NC}"
        echo "$error"
        echo -e "${RED}Installation failed${NC}"
        exit 1
    fi
}

function pre_install_sanity() {
    which curl > /dev/null
    if [ $? -gt 0 ]; then
                echo -e "${GREEN}Installing curl ...${NC}"

                if [ $DISTRO = "Ubuntu" ] || [ $DISTRO = "Debian" ]; then
                        echo "Updating apt repository cache..."
                        $APT_CMD update > /dev/null
                        $APT_CMD install curl

                elif [ $DISTRO = "CentOS" ] || [ $DISTRO = "Amazon" ] || [ $DISTRO = "RHEL" ] || [ $DISTRO = "Oracle" ]; then
                        if [ "$MACHINE" = "i686" ]; then
                                $YUM_CMD install curl.i686
                        fi

                        if [ "$MACHINE" = "x86_64" ]; then
                                $YUM_CMD install curl.x86_64
                        fi

                elif [ $DISTRO = "FreeBSD" ]; then
                        pkg_add -r curl
                fi
    fi

    CURL="`which curl`"

    if [ $DISTRO = "Ubuntu" ] || [ $DISTRO = "Debian" ]; then
        test -f /usr/lib/apt/methods/https
        if [ $? -gt 0 ];then
            echo "apt-transport-https is not installed to access DHound Gate HTTPS based APT repository ..."
                        echo "Updating apt repository cache..."
                        $APT_CMD update > /dev/null
                        echo "Installing apt-transport-https ..."
                        $APT_CMD install apt-transport-https
        fi
    fi
}

# Grab some system information
if [ -f /etc/redhat-release ] ; then
    PLATFORM=`cat /etc/redhat-release`
    DISTRO=`echo $PLATFORM | awk '{print $1}'`
    if [ "$DISTRO" = "Fedora" ]; then
       DISTRO="RHEL"
       VERSION="6"
    else
       if [ "$DISTRO" != "CentOS" ]; then
           if [ "$DISTRO" = "Enterprise" ] || [ -f /etc/oracle-release ]; then
                # Oracle "Enterprise Linux"/"Linux"
                DISTRO="Oracle"
                VERSION=`echo $PLATFORM | awk '{print $7}'`
           elif [ "$DISTRO" = "Red" ]; then
                DISTRO="RHEL"
                VERSION=`echo $PLATFORM | awk '{print $7}'`
           else
                DISTRO="unknown"
                PLATFORM="unknown"
                VERSION="unknown"
           fi
       elif [ "$DISTRO" = "CentOS" ]; then
           VERSION=`echo $PLATFORM | awk '{print $3}'`
           if [ "$VERSION" = "release" ]; then
             VERSION=`echo $PLATFORM | awk '{print $4}'`
           fi
       fi
    fi
    MACHINE=`uname -m`
elif [ -f /etc/system-release ]; then
    PLATFORM=`cat /etc/system-release | head -n 1`
    DISTRO=`echo $PLATFORM | awk '{print $1}'`
    VERSION=`echo $PLATFORM | awk '{print $5}'`
    MACHINE=`uname -m`
elif [ -f /etc/lsb-release ] ; then
    #Ubuntu version lsb-release - https://help.ubuntu.com/community/CheckingYourUbuntuVersion
    . /etc/lsb-release
    PLATFORM=$DISTRIB_DESCRIPTION
    DISTRO=$DISTRIB_ID
    VERSION=$DISTRIB_RELEASE
    MACHINE=`uname -m`
    if [ "$DISTRO" = "LinuxMint" ]; then
       DISTRO="Ubuntu"
       VERSION="12.04"
    fi
elif [ -f /etc/debian_version ] ; then
    #Debian Version /etc/debian_version - Source: http://www.debian.org/doc/manuals/debian-faq/ch-software.en.html#s-isitdebian
    DISTRO="Debian"
    VERSION=`cat /etc/debian_version`
    INFO="$DISTRO $VERSION"
    PLATFORM=$INFO
    MACHINE=`uname -m`
elif [ -f /etc/os-release ] ; then
    . /etc/os-release
    PLATFORM=$PRETTY_NAME
    DISTRO=$NAME
    VERSION=$VERSION_ID
    MACHINE=`uname -m`
elif [ -f /etc/gentoo-release ] ; then
    PLATFORM="Gentoo"
    DISTRO="Gentoo"
    VERSION=`cat /etc/gentoo-release | cut -d ' ' -f 5`
    MACHINE=`uname -m`
else
    PLATFORM=`uname -sv | grep 'SunOS joyent'` > /dev/null
    if [ "$?" = "0" ]; then
      PLATFORM="SmartOS"
      DISTRO="SmartOS"
      MACHINE="i386"
      VERSION=13
      if [ -f /etc/product ]; then
        grep "base64" /etc/product > /dev/null
        if [ "$?" = "0" ]; then
            MACHINE="x86_64"
        fi
        VERSION=`grep 'Image' /etc/product | awk '{ print $3}' | awk -F. '{print $1}'`
      fi
    elif [ "$?" != "0" ]; then
        uname -sv | grep 'FreeBSD' > /dev/null
        if [ "$?" = "0" ]; then
            PLATFORM="FreeBSD"
            DISTRO="FreeBSD"
            VERSION=`uname -r`
            MACHINE=`uname -m`
        else
            uname -sv | grep 'Darwin' > /dev/null
            if [ "$?" = "0" ]; then
                PLATFORM="Darwin"
                DISTRO="OS X"
                VERSION=`uname -r`
                MACHINE=`uname -m`
            fi
        fi
    fi
fi

while getopts ":h:u" opt; do
    case $opt in
        u)
            UPDATEAGENT=1
        ;;
        h)
                print_help
                ;;
        [?])
                exit 1
                ;;
        \?)
                exit 1
                ;;
        :)
            echo -e "${RED}Option -$OPTARG requires an argument.${NC}" >&2
            exit 1
            ;;
    esac
done

echo -e "${GREEN}===DHound Output Traffic Monitor===${NC}"
which /opt/dhound-output-traffic-monitor/bin/dhound-output-traffic-monitor > /dev/null
if [ $? -eq 0 ]; then
    DHOUND_INSTALLED=1
else
    DHOUND_INSTALLED=0
fi

if [ $DHOUND_INSTALLED -eq 1 ]; then
    if  [ $UPDATEAGENT -eq 0 ]; then
            echo -e "{$BLUE}Dhound Output Traffic Monitor already installed into the system. Use -u option for the script to upgrade dhound-output-traffic-monitor to the latest version.${NC}"
        print_help
        exit 1
    fi

    echo -e "${GREEN}DHound Output Traffic Monitor already installed. The script will upgrade dhound-output-traffic-monitor to the latest version.${NC}"
    UPDATEAGENT=1
else
    UPDATEAGENT=0
fi

if [ "$MACHINE" = "i686" ] ||
   [ "$MACHINE" = "i586" ] ||
   [ "$MACHINE" = "i386" ] ; then
    ARCH="32"
    SUPPORTED_ARCH=1
fi

#determine hard vs. soft float using readelf
if [[ "$MACHINE" == arm* ]] ; then
        if [ -x /usr/bin/readelf ] ; then
                HARDFLOAT=`readelf -a /proc/self/exe | grep armhf`
                if [ -z "$HARDFLOAT" ]; then
                        if [ "$MACHINE" = "armv7l" ] ||
                           [ "$MACHINE" = "armv6l" ] ||
                           [ "$MACHINE" = "armv5tel" ] ||
                           [ "$MACHINE" = "armv5tejl" ] ; then
                                ARCH="32"
                                SUPPORTED_ARCH=1
                                echo "Detected $MACHINE running armel"
                        fi
                else
                        ARCH="32"
                        SUPPORTED_ARCH=1
                        echo "Detected $MACHINE running armhf"
                fi
        else
                echo -e "${RED}Cannot determine ARM ABI, please install the 'binutils' package${NC}"
        fi
fi

if [ "$MACHINE" = "x86_64" ] || [ "$MACHINE" = "amd64" ]; then
    ARCH="64"
    SUPPORTED_ARCH=1
fi

if [ $SUPPORTED_ARCH -eq 0 ]; then
    echo -e "${RED}Unsupported architecture ($MACHINE) ...${NC}"
    echo -e "${RED}This is an unsupported platform for the dhound.${NC}"
    exit 1
fi

# ${BLUE}
# ${NC}
echo -e "${GREEN}Supplied paramenters:${NC} \n${BLUE}OS:${NC} $DISTRO $VERSION..."

DISTROMAJORVERSION=$(echo "$VERSION" | grep -oP "[0-9]+" | head -1)

# Check the distribution and version
for d in ${PLATFORMS[*]} ; do
    if [[ "$DISTRO" = "$d" ]]; then
        TEMP="\${${DISTRO}_VERSION_MIN}"
        MIN_VERSION=$(eval echo "$TEMP")
        if [[ "$DISTROMAJORVERSION" -ge "$MIN_VERSION" ]]; then
            SUPPORTED_PLATFORM=1
            break
        fi
    fi
done

if [ $SUPPORTED_PLATFORM -eq 0 ]; then
    echo -e "${RED}Your platform is not supported by this script at this time.${NC}"
    print_supported_platforms
    exit 1
fi

# If this script is being run by root for some reason, don't use sudo.
if [ "$(id -u)" != "0" ]; then
    echo -e "${RED}This script must be executed as the 'root' user or with sudo${NC}"
    echo -e "${RED}Please install sudo or run again as the 'root' user.${NC}"
    exit 1
fi

# At this point, we think we have a supported OS.
pre_install_sanity $d

do_install
