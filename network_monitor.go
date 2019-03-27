package main

import (
	"fmt"
	gopacket "github.com/google/gopacket"
	layers "github.com/google/gopacket/layers"
	pcap "github.com/google/gopacket/pcap"
	"regexp"
	"strings"
	"time"
)

type NetworkEventType int

const (
	TcpConnectionInitiatedByHost NetworkEventType = iota
	TcpConnectionSetUp
	UdpSendByHost
	UdpResponse
	DnsResponseReceived
)

type NetworkEvent struct {
	Type       NetworkEventType
	Connection *NetworkConnectionData
	Dns        *DnsAnswer
}

type DnsAnswer struct {
	EventTimeUtcNumber int64
	DomainName         string
	IpAddresses        *[]string
}

type NetworkConnectionData struct {
	LocalIpAddress     string
	LocalPort          uint32
	RemoteIpAddress    string
	RemotePort         uint32
	Size               uint16
	Sequence           uint32
	EventTimeUtcNumber int64
}

type NetworkMonitor struct {
	Options           *Options
	Output            chan *NetworkEvent
	SysProcessManager *SysProcessManager
	_firstRun         bool
}

type NetInterfaceInfo struct {
	Addresses    []string
	Name         string
	Active       bool
	NetInterface *pcap.Interface
}

var (
	Ipv4Validator     = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	localNetworkRegex = regexp.MustCompile(`^((0\.)|(127\.0\.0\.1)|(192\.168\.)|(10\.)|(172\.(1[6-9]|2[0-9]|3[0-1])\.)|(fc00:)|(fe80:))`)
)

func (monitor *NetworkMonitor) Run() {
	netInterfaces := monitor.FindActiveNetInterfaces()

	for _, dev := range netInterfaces {
		go monitor._monitorInterfaceTraffic(dev)
	}
}

func (monitor *NetworkMonitor) FindActiveNetInterfaces() []*NetInterfaceInfo {
	activeInterfaces := make([]*NetInterfaceInfo, 0)
	netInterfaces := monitor.FindAllNetInterfaces()
	var targetInterface *NetInterfaceInfo

	for _, dev := range netInterfaces {
		if dev.Active {
			activeInterfaces = append(activeInterfaces, dev)

			if monitor.Options.NetworkInterface == dev.Name {
				targetInterface = dev
			}
		}
	}

	if len(monitor.Options.NetworkInterface) > 0 {
		if targetInterface == nil {
			emitLine(logLevel.important, "specified network interface '%s' not found among active network interfaces.", monitor.Options.NetworkInterface)
		} else {
			activeInterfaces = make([]*NetInterfaceInfo, 0)
			activeInterfaces = append(activeInterfaces, targetInterface)
		}
	}

	return activeInterfaces

}

func (monitor *NetworkMonitor) FindAllNetInterfaces() []*NetInterfaceInfo {
	interfaces := make([]*NetInterfaceInfo, 0)

	devices, err := pcap.FindAllDevs()
	if err != nil {
		emitLine(logLevel.important, "Failed find network devices", err.Error())
		return interfaces
	}

	for _, device := range devices {
		intf := &NetInterfaceInfo{}
		intf.NetInterface = &device
		intf.Name = device.Name
		intf.Active = false
		interfaces = append(interfaces, intf)

		for _, address := range device.Addresses {
			ip := address.IP.String()

			if ip != "127.0.0.1" && ip != "0.0.0.0" {
				intf.Addresses = append(intf.Addresses, ip)

				if Ipv4Validator.MatchString(ip) {
					intf.Active = true
				}
			}
		}
	}

	return interfaces
}

func (monitor *NetworkMonitor) _monitorInterfaceTraffic(dev *NetInterfaceInfo) {

	deviceName := dev.Name
	handle, err := pcap.OpenLive(deviceName, 1600, false, pcap.BlockForever)
	if err != nil {
		emitLine(logLevel.important, "Failed listening network interface %s. Error: %s.", deviceName, err.Error())
		// monitor.ShowDevices()
		return
	}
	defer handle.Close()

	// Set filter
	filterLines := make([]string, 0)

	for _, address := range dev.Addresses {
		if monitor.Options.Protocol == "all" || monitor.Options.Protocol == "tcp"{
			// track TCP SYN that host tries to make to initiate connection
			filterLines = append(filterLines, fmt.Sprintf("((tcp[tcpflags] == tcp-syn) && src %s)", address))
			// track TCP SYN-ACK sent to host (TCP connection is opened)
			filterLines = append(filterLines, fmt.Sprintf("(tcp[13] = 18 and dst host %s)", address))
		}
		if monitor.Options.Protocol == "all" || monitor.Options.Protocol == "udp"{
			// track all UDP traffic initiated by host
			filterLines = append(filterLines, fmt.Sprintf("(udp && src host %s)", address))
		}

		// track DNS traffic that comes on the host
		filterLines = append(filterLines, fmt.Sprintf("(udp && port 53 && dst host %s)", address))
	}

	trafficFilter := strings.Join(filterLines, " || ")
	// this filter can be used for tcpdump to observe traffic
	// debug(trafficFilter)

	err = handle.SetBPFFilter(trafficFilter)
	if err != nil {
		emitLine(logLevel.important, "Failed set BPF Filter on device %s. Error: %s. Filter: %s", deviceName, err.Error(), trafficFilter)
		return
	}

	emitLine(logLevel.important, "Start listening to network device: %s (%s). Ip Addresses: %s.", deviceName, dev.NetInterface.Description, strings.Join(dev.Addresses, ", "))

	var eth layers.Ethernet
	var ip4 layers.IPv4
	var ip6 layers.IPv6
	var tcp layers.TCP
	var udp layers.UDP
	var dns layers.DNS

	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4, &ip6, &tcp, &udp, &dns)

	decodedLayers := make([]gopacket.LayerType, 0, 10)
	for {
		data, _, err := handle.ReadPacketData()
		if err != nil {
			continue
		}

		err = parser.DecodeLayers(data, &decodedLayers)
		if err != nil {
			continue
		}

		// define source and destination addresses
		srcIp := ""
		dstIp := ""

		for _, typ := range decodedLayers {
			switch typ {
			case layers.LayerTypeIPv4:
				srcIp = ip4.SrcIP.String()
				dstIp = ip4.DstIP.String()
			case layers.LayerTypeIPv6:
				srcIp = ip6.SrcIP.String()
				dstIp = ip6.DstIP.String()
			case layers.LayerTypeTCP:
				// server initiate external connection to Internet
				if tcp.SYN && !tcp.ACK {
					if localNetworkRegex.MatchString(dstIp) == false && Contains(dev.Addresses, srcIp) {
						// emitLine(logLevel.verbose, "TCP INIT: %s:%d->%s:%d, seq: %d, ack: %d.", srcIp, tcp.SrcPort, dstIp, tcp.DstPort, tcp.Seq, tcp.Ack)
						monitor.Output <- &NetworkEvent{
							Type: TcpConnectionInitiatedByHost,
							Connection: &NetworkConnectionData{
								LocalIpAddress:     srcIp,
								RemoteIpAddress:    dstIp,
								LocalPort:          uint32(tcp.SrcPort),
								Sequence:           tcp.Seq,
								RemotePort:         uint32(tcp.DstPort),
								EventTimeUtcNumber: time.Now().UTC().Unix(),
							},
						}
					}
				}
				// external source responded on initiated connection
				if tcp.SYN && tcp.ACK {
					if localNetworkRegex.MatchString(srcIp) == false && Contains(dev.Addresses, dstIp) {
						// emitLine(logLevel.verbose, "TCP ACCEPTED: %s:%d->%s:%d, seq: %d, ack: %d.", srcIp, tcp.SrcPort, dstIp, tcp.DstPort, tcp.Seq, tcp.Ack)
						monitor.Output <- &NetworkEvent{
							Type: TcpConnectionSetUp,
							Connection: &NetworkConnectionData{
								LocalIpAddress:     dstIp,
								RemoteIpAddress:    srcIp,
								LocalPort:          uint32(tcp.DstPort),
								Sequence:           tcp.Ack,
								RemotePort:         uint32(tcp.SrcPort),
								EventTimeUtcNumber: time.Now().UTC().Unix(),
							},
						}
					}
				}

			case layers.LayerTypeUDP:
			//	emitLine(logLevel.verbose, "UDP: %s:%d->%s:%d, checksum: %d, len: %d.", srcIp, udp.SrcPort, dstIp, udp.DstPort, udp.Checksum, udp.Length)

				// host sent UDP request to Internet
				if localNetworkRegex.MatchString(dstIp) == false && Contains(dev.Addresses, srcIp) {
						// emitLine(logLevel.verbose, "UDP INIT: %s:%d->%s:%d, checksum: %d, len: %d.", srcIp, udp.SrcPort, dstIp, udp.DstPort, udp.Checksum, udp.Length)

						monitor.Output <- &NetworkEvent{
							Type: UdpSendByHost,
							Connection: &NetworkConnectionData{
								LocalIpAddress:     srcIp,
								RemoteIpAddress:    dstIp,
								Size:               udp.Length,
								LocalPort:          uint32(udp.SrcPort),
								RemotePort:         uint32(udp.DstPort),
								EventTimeUtcNumber: time.Now().UTC().Unix(),
							},
						}
					}

				// Internet source sent an UDP request back
				if localNetworkRegex.MatchString(srcIp) == false && Contains(dev.Addresses, dstIp) {
					// emitLine(logLevel.verbose, "UDP RESPONSE: %s:%d->%s:%d, checksum: %d, len: %d.", srcIp, udp.SrcPort, dstIp, udp.DstPort, udp.Checksum, udp.Length)
					monitor.Output <- &NetworkEvent{
						Type: UdpResponse,
						Connection: &NetworkConnectionData{
							LocalIpAddress:     dstIp,
							RemoteIpAddress:    srcIp,
							Size:               udp.Length,
							LocalPort:          uint32(udp.DstPort),
							RemotePort:         uint32(udp.SrcPort),
							EventTimeUtcNumber: time.Now().UTC().Unix(),
						},
					}
				}

			case layers.LayerTypeDNS:
				if int(dns.ANCount) > 0 {
					for _, dnsQuestion := range dns.Questions {
						var ips []string
						dnsName := dnsQuestion.Name

						// debugJson(dns)

						for _, dnsAnswer := range dns.Answers {
							if dnsAnswer.Type == layers.DNSTypeA || dnsAnswer.Type == layers.DNSTypeAAAA {

								ip := dnsAnswer.IP.String()

								if ip != "" && ip != "<nil>" && localNetworkRegex.MatchString(ip) == false {
									ips = append(ips, ip)
								}
							}
						}

						if len(ips) > 0 {
							// debug("DNS Response: %s: %s", dnsName, strings.Join(ips, ","))
							monitor.Output <- &NetworkEvent{
								Type: DnsResponseReceived,
								Dns: &DnsAnswer{
									DomainName:         fmt.Sprintf("%s", dnsName),
									IpAddresses:        &ips,
									EventTimeUtcNumber: time.Now().UTC().Unix(),
								},
							}
						}
					}
				}
			default:
				// debugJson(typ)
			}
		}
	}
}
