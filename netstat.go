package main

import (
	"time"
	net "github.com/shirou/gopsutil/net"
)

type NetStatInfo struct {
	Pid                int32
	LocalIp            string
	LocalPort          uint32
	EventTimeUtcNumber int64
}

type NetStatManager struct {
	_cache []*NetStatInfo
	Options *Options
}

func (netstat *NetStatManager) Init() {
	netstat._cache = make([]*NetStatInfo, 0)
}

func (netstat *NetStatManager) Run() {
	go func() {
		for {
			netstat.SyncPortList()
			time.Sleep(2 * time.Second)
		}
	}()
}

func (netstat *NetStatManager) FindNetstatInfoByLocalPort(localIp string, localPort uint32) *NetStatInfo {
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

func (manager *NetStatManager) SyncPortList() {
	connections, err := net.Connections("all")
	if err != nil {
		return
	}

	currentTimeNumber := time.Now().UTC().Unix()
	// keep info about pid only 60 seconds
	minTime := currentTimeNumber - 60

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
		if connection.Laddr.Port > 0 && connection.Pid > 0 && connection.Laddr.IP != "127.0.0.1" && connection.Laddr.IP != "::1" {
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
