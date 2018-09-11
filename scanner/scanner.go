package scanner

import (
	"net"
	"strconv"
	"time"
)

const (
	// Timeout 超时时间
	Timeout = 2 * time.Second
	// PrintPort 打印端口
	PrintPort int = 9100
)

// NetPort 网络接口，包含接口信息，地址信息，ip地址集合
type NetPort struct {
	Interface net.Interface
	IPNet     *net.IPNet
	IPs       []DigitalIP
}

// IsPrinterAvailable 检查打印机端口是否开放
func IsPrinterAvailable(ip DigitalIP, port int, timeout time.Duration) bool {

	remoteAddr := ip.String() + ":" + strconv.Itoa(port)
	conn, _ := net.DialTimeout("tcp", remoteAddr, timeout)

	if conn == nil {
		return false
	}

	defer conn.Close()

	return true
}

// GetAvailableNetPorts 获取可用的网络接口集合
func GetAvailableNetPorts() []*NetPort {
	interfaces, _ := net.Interfaces()

	var netPorts []*NetPort

	for _, i := range interfaces {
		addrs, _ := i.Addrs()

		for _, j := range addrs {
			ipnet, ok := j.(*net.IPNet)
			if !ok {
				continue
			}
			if ipnet.IP.IsLoopback() {
				continue
			}
			if !IsIPV4(ipnet) {
				continue
			}

			netPorts = append(netPorts, &NetPort{Interface: i, IPNet: ipnet, IPs: Table(ipnet)})
		}
	}

	return netPorts
}

// FilterPrintIPs 过滤出是打印服务器的ip地址
func FilterPrintIPs(netPorts []*NetPort) []*NetPort {
	for index, i := range netPorts {
		channel := make(chan DigitalIP)
		defer close(channel)

		for _, ip := range i.IPs {
			go func(ip DigitalIP) {
				if IsPrinterAvailable(ip, PrintPort, Timeout) {
					channel <- ip
				}
			}(ip)
		}

		var printIPs []DigitalIP

		select {
		case ip := <-channel:
			printIPs = append(printIPs, ip)
		}

		netPorts[index].IPs = printIPs
	}

	return netPorts
}
