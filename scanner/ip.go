package scanner

import (
	"bytes"
	"math"
	"net"
	"strconv"
	"strings"
)

// DigitalIP ip地址数字形式
type DigitalIP uint32

// 网络号
func (ip DigitalIP) isBegin() bool {
	return ip&0x000000ff == 0
}

// String 将 IP(uint32) 转换成 可读性IP字符串
func (ip DigitalIP) String() string {
	var bf bytes.Buffer
	for i := 1; i <= 4; i++ {
		bf.WriteString(strconv.Itoa(int((ip >> ((4 - uint(i)) * 8)) & 0xff)))
		if i != 4 {
			bf.WriteByte('.')
		}
	}
	return bf.String()
}
func (ip DigitalIP) toIP() net.IP {
	return net.ParseIP(ip.String())
}

// Table 根据IP和mask换算内网IP范围
func Table(ipNet *net.IPNet) []DigitalIP {

	var ips []DigitalIP

	beginIP := getBeginIP(ipNet)
	broadcastIP := getBroadcastIP(ipNet)

	for i := beginIP; i < broadcastIP; i++ {
		if i.isBegin() {
			continue
		}
		ips = append(ips, i)
	}

	return ips
}

// ParseIP []byte --> IP
func ParseIP(b []byte) DigitalIP {
	return DigitalIP(DigitalIP(b[0])<<24 + DigitalIP(b[1])<<16 + DigitalIP(b[2])<<8 + DigitalIP(b[3]))
}

// ParseIPString string --> IP
func ParseIPString(s string) DigitalIP {
	var b []byte
	for _, i := range strings.Split(s, ".") {
		v, _ := strconv.Atoi(i)
		b = append(b, uint8(v))
	}
	return ParseIP(b)
}

func getBroadcastIP(ipNet *net.IPNet) DigitalIP {
	currentLength, _ := ipNet.Mask.Size()
	beginIP := getBeginIP(ipNet)
	return beginIP | DigitalIP(math.Pow(2, float64(32-currentLength))-1)
}

func getBeginIP(ipNet *net.IPNet) DigitalIP {
	ip := ipNet.IP.To4()
	var min DigitalIP

	for i := 0; i < 4; i++ {
		b := DigitalIP(ip[i] & ipNet.Mask[i])
		min += b << ((3 - uint(i)) * 8)
	}

	return min
}

// IsIPV4 判断是否ipv4
func IsIPV4(ip *net.IPNet) bool {
	return ip.IP.To4() != nil
}
