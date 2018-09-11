package main

import (
	"fmt"
	"net"
	"time"
	"github.com/gw123/escpos-go/scanner"
)

func main() {

	netPorts := scanner.GetAvailableNetPorts()
	netPorts = scanner.FilterPrintIPs(netPorts)

	arp := scanner.ARP{
		RequsetTimeout:  10 * time.Second,
		ResponseTimeout: 10 * time.Second,
		Channel:         make(chan scanner.IPMac),
	}
	defer close(arp.Channel)

	for _, netPort := range netPorts {
		interfaceName := netPort.Interface.Name
		localAddress := netPort.Interface.HardwareAddr

		var printerIPs []net.IP
		for _, ip := range netPort.IPs {
			printerIPs = append(printerIPs, net.ParseIP(ip.String()))
		}

		go arp.ListenResponse(interfaceName, printerIPs)

		for _, ip := range netPort.IPs {
			targetIP := &net.IPNet{IP: net.ParseIP(ip.String()), Mask: netPort.IPNet.Mask}
			go arp.SendRequest(netPort.IPNet, localAddress, interfaceName, targetIP)
		}
	}

	timer := time.NewTicker(10 * time.Second)
	defer timer.Stop()

	interfaceLen := len(netPorts)

	receivedInterface := 0
	select {
	case <-timer.C:
		return
	case ipmac := <-arp.Channel:
		fmt.Println(ipmac, ipmac.Full())

		receivedInterface++
		if receivedInterface == interfaceLen {
			return
		}
	}
}
