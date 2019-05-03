package main

import (
	"time"

	"github.com/shirou/gopsutil/net"
)

type NetStatInfo struct {
	Pid                int32
	LocalIp            string
	LocalPort          uint32
	EventTimeUtcNumber int64
}

type NetStatManager struct {
	_cache  []*NetStatInfo
	Options *Options
}

func (netstat *NetStatManager) Init() {
	netstat._cache = make([]*NetStatInfo, 0)
}

func (netstat *NetStatManager) Run() {

}

func (manager *NetStatManager) SyncPortList(protocol string) {
	connections, err := net.Connections(protocol)
	if err != nil {
		return
	}

	currentTimeNumber := time.Now().UTC().Unix()
	// keep info about pid only particular seconds
	minTime := currentTimeNumber - 60*3

	var list []*NetStatInfo
	for _, stat := range manager._cache {
		if stat.EventTimeUtcNumber > minTime {
			list = append(list, stat)
		}
	}

	if list == nil {
		list = manager._cache
	}

	for _, connection := range connections {
		if connection.Laddr.Port > 0 && connection.Laddr.IP != "127.0.0.1" && connection.Laddr.IP != "::1" {
			isFound := false
			for _, stat := range list {
				if stat.LocalIp == connection.Laddr.IP && stat.LocalPort == connection.Laddr.Port {
					stat.Pid = connection.Pid
					isFound = true
				}
			}

			if !isFound {
				stat := &NetStatInfo{
					Pid:                connection.Pid,
					LocalIp:            connection.Laddr.IP,
					LocalPort:          connection.Laddr.Port,
					EventTimeUtcNumber: currentTimeNumber,
				}

				list = append(list, stat)
				debugJson(stat)
			}
		}
	}

	manager._cache = list
}

func (netstat *NetStatManager) FindNetstatInfoByLocalPort(localIp string, localPort uint32, protocol NetworkProtocol) *NetStatInfo {
	var string prot
	isIpV6 := IsIPv6(localIp)

	if protocol == TCP && !isIpV6 {
		prot = "tcp4"

	} else if protocol == TCP && isIpV6 {
		prot = "tcp6"
	} else if protocol == UDP && !isIpV6 {
		prot = "udp4"
	} else if protocol == UDP && isIpV6 {
		prot = "udp6"
	}

	result := netstat.InternalFindNetstatInfoByLocalPort(localIp, localPort)
	if result == nil {
		netstat.SyncPortList(prot)
		result := netstat.InternalFindNetstatInfoByLocalPort(localIp, localPort)
	}

	return result
}

func (netstat *NetStatManager) InternalFindNetstatInfoByLocalPort(localIp string, localPort uint32) *NetStatInfo {
	cache := netstat._cache
	if len(cache) > 0 {
		// check by ip and port
		for _, info := range cache {
			if info.LocalIp == localIp && info.LocalPort == localPort {
				return info
			}
		}

		// check only by port
		for _, info := range cache {
			if info.LocalPort == localPort {
				return info
			}
		}
	}
	return nil
}
