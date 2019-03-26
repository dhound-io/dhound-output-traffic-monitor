dhound-output-traffic-monitor
==================
Monitor and log UDP and TCP connections initiated by host to Internet. Connections are associated with process info.

Utility to show network traffic (both TCP and UDP v4 and v6) split by process and remote host.

All of this functionality is fully configurable.

## Documentation
These instructions will get you install and configure Dhound Output Traffic Monitor on your server.

## Install
To install login to the server using ssh and run the following command:
```
curl https://raw.githubusercontent.com/dhound-io/dhound-output-traffic-monitor/master/publish/install-agent.sh 2>/dev/null | sudo bash -s -- -u
```
After executing this command the installer will be downloaded and started.

## Running
Start Dhound Output Traffic Monitor service
```
service dhound-output-traffic-monitor start
```
Stop Dhound Output Traffic Monitor service
```
service dhound-output-traffic-monitor stop
```
Show service status
```
service dhound-output-traffic-monitor status
```

### Options
```
    -log-file
```
network events output: syslog, console, <path to a custom file>; default: console

```
    -eth
```
listen to a particular network interface; default: listen to all active network interfaces
```
    -verbose
```
log more detailed and debug information; default: false
```
    -version
```
Display dhound output traffic monitor version

```
    -pprof
```
