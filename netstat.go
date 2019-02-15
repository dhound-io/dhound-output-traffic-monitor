package main

import (
	"time"
)

type NetStatInfo struct {
	Pid                int32
	LocalIp            string
	LocalPort          uint32
	EventTimeUtcNumber int64
}

type NetStatManager struct {
	_cache []*NetStatInfo
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
