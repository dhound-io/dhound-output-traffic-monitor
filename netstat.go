package main

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

