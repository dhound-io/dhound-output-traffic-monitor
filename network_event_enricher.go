package main

import (
	"fmt"
	"time"
	/*
		"strings"

	*/)

type NetworkProtocol int

const (
	TCP NetworkProtocol = iota
	UDP
)

type ProcessNetworkEvent struct {
	Protocol NetworkProtocol

	LocalIpAddress  string
	LocalPort       uint32
	RemoteIpAddress string
	RemotePort      uint32
	Success         bool

	// process info
	ProcessId   int32
	ProcessInfo *ProcessInfo

	Domains []string

	Size uint16

	Hash               string
	EventTimeUtcNumber int64
}

type NetworkEventEnricher struct {
	Input      chan *NetworkEvent
	Output     chan []string
	SysManager *SysProcessManager
	NetStat    *NetStatManager
	_cache     []*ProcessNetworkEvent
	_dnsCache  []*DnsAnswer
}

func (enricher *NetworkEventEnricher) Init() {
	enricher._cache = make([]*ProcessNetworkEvent, 0)
	enricher._dnsCache = make([]*DnsAnswer, 0)
}

func (enricher *NetworkEventEnricher) Run() {
	go func() {
		for {
			time.Sleep(2 * time.Second)
			enricher._sync()
		}
	}()

	go func() {
		for networkEvent := range enricher.Input {
			enricher._processInput(networkEvent)
		}
	}()
}

func (enricher *NetworkEventEnricher) _processInput(networkEvent *NetworkEvent) {
	if networkEvent == nil {
		return
	}

	connection := networkEvent.Connection
	if networkEvent.Type == TcpConnectionInitiatedByHost { // host initiated TCP connection

		event := &ProcessNetworkEvent{
			Protocol:           TCP,
			EventTimeUtcNumber: connection.EventTimeUtcNumber,
			LocalIpAddress:     connection.LocalIpAddress,
			LocalPort:          connection.LocalPort,
			RemoteIpAddress:    connection.RemoteIpAddress,
			RemotePort:         connection.RemotePort,
			Success:            false,
			ProcessId:          -1,
		}

		event.Hash = enricher.CalculateHashForTCP(connection.LocalIpAddress, connection.LocalPort, connection.RemoteIpAddress, connection.RemotePort, connection.Sequence)
		enricher._cache = append(enricher._cache, event)

	} else if networkEvent.Type == TcpConnectionSetUp { // host received TCP connection acknowledge from external source
		hash := enricher.CalculateHashForTCP(connection.LocalIpAddress, connection.LocalPort, connection.RemoteIpAddress, connection.RemotePort, connection.Sequence-1)
		for _, event := range enricher._cache {
			if event != nil && event.Hash == hash {
				event.Success = true
			}
		}
	} else if networkEvent.Type == UdpSendByHost { // host initiated UDP connection
		event := &ProcessNetworkEvent{
			Protocol:           UDP,
			EventTimeUtcNumber: connection.EventTimeUtcNumber,
			LocalIpAddress:     connection.LocalIpAddress,
			LocalPort:          connection.LocalPort,
			RemoteIpAddress:    connection.RemoteIpAddress,
			RemotePort:         connection.RemotePort,
			Size:               connection.Size,
			Success:            false,
			ProcessId:          -1,
		}
		event.Hash = enricher.CalculateHashForUDP(connection.LocalIpAddress, connection.LocalPort, connection.RemoteIpAddress, connection.RemotePort)
		enricher._cache = append(enricher._cache, event)

	} else if networkEvent.Type == UdpResponse { // host received UDP response from external source
		hash := enricher.CalculateHashForUDP(connection.LocalIpAddress, connection.LocalPort, connection.RemoteIpAddress, connection.RemotePort)
		for _, event := range enricher._cache {
			if event != nil && event.Hash == hash {
				event.Success = true
				event.Size += connection.Size
			}
		}
	} else if networkEvent.Type == DnsResponseReceived && networkEvent.Dns != nil { // dns received
		enricher._dnsCache = append(enricher._dnsCache, networkEvent.Dns)
	}
}

func (enricher *NetworkEventEnricher) CalculateHashForTCP(
	localIpAddress string, localPort uint32,
	remoteIpAddress string, remotePort uint32, sequence uint32) string {

	var hash = fmt.Sprintf("tcp_%s:%d->%s:%d_%d", localIpAddress, localPort, remoteIpAddress, remotePort, sequence)
	return hash
}

func (enricher *NetworkEventEnricher) CalculateHashForUDP(
	localIpAddress string, localPort uint32,
	remoteIpAddress string, remotePort uint32) string {

	var hash = fmt.Sprintf("udp_%s:%d->%s:%d", localIpAddress, localPort, remoteIpAddress, remotePort)
	return hash
}

func (enricher *NetworkEventEnricher) _sync() {

	if len(enricher._cache) > 0 {
		eventsToPublish := make([]*ProcessNetworkEvent, 0)
		for index, networkEvent := range enricher._cache {
			if networkEvent == nil {
				continue
			}

			isToPublish := false

			// check dns
			if networkEvent.Domains == nil {
				domains := make([]string, 0)
				for _, dnsAnswer := range enricher._dnsCache {
					for _, ip := range *dnsAnswer.IpAddresses {
						if networkEvent.RemoteIpAddress == ip {
							domains = append(domains, dnsAnswer.DomainName)
						}
					}
				}
				networkEvent.Domains = domains
			}

			// try to find info about process
			if networkEvent.ProcessId < 0 {
				netStatInfo := enricher.NetStat.FindNetstatInfoByLocalPort(networkEvent.LocalIpAddress, networkEvent.LocalPort)
				if netStatInfo != nil {
					networkEvent.ProcessId = netStatInfo.Pid
				}
			}

			if networkEvent.ProcessId >= 0 && networkEvent.ProcessInfo == nil {
				networkEvent.ProcessInfo = enricher.SysManager.FindProcessInfoByPid(networkEvent.ProcessId)
			}

			difference := time.Now().Sub(time.Unix(networkEvent.EventTimeUtcNumber, 0).UTC())
			// max time for setting up connection - we give only 1 minute

			if difference.Minutes() > 1 || enricher._isNetworkEventProcessCompleted(networkEvent) {
				isToPublish = true
			}

			/**/
			if isToPublish {
				eventsToPublish = append(eventsToPublish, networkEvent)
				enricher._cache = enricher.RemoveIndex(enricher._cache, index)
			}
		}

		if len(eventsToPublish) > 0 {
			debugJson(eventsToPublish)
			// we can publish events
			// linesToPublish := make([]string, len(eventsToPublish))
			// debugJson(linesToPublish)
			/*
				for index, event := range eventsToPublish {
					output := ""
					switch eventType := event.Type; eventType {
					case TcpConnectionInitiatedByHost, TcpConnectionSetUp:
						{
							output = fmt.Sprintf("[%s]: TCP %s:%s -> %s:%s success:%t", time.Unix(event.EventTimeUtcNumber, 0).Format(time.RFC3339), event.Connection.LocalIpAddress, fmt.Sprint(event.Connection.LocalPort), event.Connection.RemoteIpAddress, fmt.Sprint(event.Connection.RemotePort), event.Success)
							if event.NetStatInfo != nil {
								output = output + fmt.Sprintf(" pid: %d", event.NetStatInfo.Pid)
								if event.ProcessInfo != nil {
									output = output + fmt.Sprintf(" process: %s commandline: %s", event.ProcessInfo.Name, event.ProcessInfo.CommandLine)
								}
							}
						}

					case UdpSendByHost:
						{
							output = fmt.Sprintf("[%s]: UDP %s:%s -> %s:%s", time.Unix(event.EventTimeUtcNumber, 0).Format(time.RFC3339), event.Connection.LocalIpAddress, fmt.Sprint(event.Connection.LocalPort), event.Connection.RemoteIpAddress, fmt.Sprint(event.Connection.RemotePort))
							if event.NetStatInfo != nil {
								output = output + fmt.Sprintf(" pid: %d", event.NetStatInfo.Pid)
								if event.ProcessInfo != nil {
									output = output + fmt.Sprintf(" process: %s commandline: %s", event.ProcessInfo.Name, event.ProcessInfo.CommandLine)
								}
							}
						}
						// @TODO Write logs for DNS types
					case DnsResponseReceived:
						{
							var ips = strings.Join(*event.Dns.IpAddresses, ", ")
							output = fmt.Sprintf("[%s]: DNS domain:%s, ip: [%s]", time.Unix(event.EventTimeUtcNumber, 0).Format(time.RFC3339), event.Dns.DomainName, ips)
						}
					}
					if output != "" {
						linesToPublish[index] = output
					}
				}
				enricher.Output <- linesToPublish
			*/
		}
	}

	// debug("Sync end: %d", len(enricher._cache))
}

func (enricher *NetworkEventEnricher) _isNetworkEventProcessCompleted(event *ProcessNetworkEvent) bool {
	if event == nil {
		return false
	}

	if event.ProcessId > -1 && event.Success == true && event.ProcessInfo != nil {
		return true
	}

	return false
}

func (enricher *NetworkEventEnricher) RemoveIndex(array []*ProcessNetworkEvent, index int) []*ProcessNetworkEvent {
	array[index] = array[len(array)-1] // Copy last element to index i.
	array[len(array)-1] = nil          // Erase last element (write zero value).
	array = array[:len(array)-1]       // Truncate slice.

	return array
}
