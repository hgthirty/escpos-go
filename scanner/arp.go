package scanner

import (
	"errors"
	"log"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

const (
	// OptionSetResponse 0x0002 arp response
	OptionSetResponse = uint16(2)
)

// ARP arp请求
type ARP struct {
	RequsetTimeout  time.Duration
	ResponseTimeout time.Duration
	Channel         chan IPMac
}

// SendRequest 发送arp请求
func (arp *ARP) SendRequest(localIP *net.IPNet, localAddr net.HardwareAddr, iface string, targetIP *net.IPNet) error {
	srcIP := localIP.IP.To4()
	dstIP := targetIP.IP.To4()

	if srcIP == nil || dstIP == nil {
		return errors.New("ip 解析出问题")
	}

	buffer := gopacket.NewSerializeBuffer()

	var opt gopacket.SerializeOptions
	gopacket.SerializeLayers(buffer, opt, initEnternet(localAddr), initARP(srcIP, localAddr, dstIP))
	outgoingPacket := buffer.Bytes()

	handle, err := pcap.OpenLive(iface, 2048, false, arp.RequsetTimeout)
	if err != nil {
		return err
	}
	defer handle.Close()

	err = handle.WritePacketData(outgoingPacket)
	if err != nil {
		return err
	}

	return nil
}

// ListenResponse 监听arp响应
func (arp *ARP) ListenResponse(iface string, ips []net.IP) {

	ipMac := InitIPMacWithIPs(ips)

	handle, err := sendARPRequest(iface, arp.RequsetTimeout)
	if err != nil {
		arp.Channel <- ipMac
	}

	ticker := time.NewTicker(arp.ResponseTimeout)
	defer ticker.Stop()
	defer handle.Close()
	ps := gopacket.NewPacketSource(handle, handle.LinkType())
	for {
		select {
		case <-ticker.C:
			log.Print("listenARP timeout")
			arp.Channel <- ipMac
			return
		case p := <-ps.Packets():
			arpResponse := p.Layer(layers.LayerTypeARP).(*layers.ARP)
			if arpResponse.Operation != OptionSetResponse {
				continue
			}

			ip := net.IP(arpResponse.SourceProtAddress).String()

			if ipMac.Exist(ip) {
				ipMac.Set(ip, net.HardwareAddr(arpResponse.SourceHwAddress))
			}

			if ipMac.Full() {
				arp.Channel <- ipMac
				return
			}
		}
	}
}

func sendARPRequest(iface string, time time.Duration) (*pcap.Handle, error) {
	handle, err := pcap.OpenLive(iface, 1024, false, time)
	if err != nil {
		log.Fatal("pcap打开失败:", err)
		return nil, err
	}

	handle.SetBPFFilter("arp")
	return handle, nil
}

func initEnternet(addr net.HardwareAddr) *layers.Ethernet {
	return &layers.Ethernet{
		SrcMAC:       addr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}
}

func initARP(ip net.IP, addr net.HardwareAddr, remoteIP net.IP) *layers.ARP {
	return &layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     uint8(6),
		ProtAddressSize:   uint8(4),
		Operation:         uint16(1), // request
		SourceHwAddress:   addr,
		SourceProtAddress: ip,
		DstHwAddress:      net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		DstProtAddress:    remoteIP,
	}
}
