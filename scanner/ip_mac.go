package scanner

import "net"

// IPMac ip和mac对应关系
type IPMac map[string]net.HardwareAddr

// InitIPMacWithIPs 通过IP初始化
func InitIPMacWithIPs(ips []net.IP) IPMac {
	ipMac := IPMac{}
	for _, ip := range ips {
		ipMac[ip.String()] = nil
	}

	return ipMac
}

// Set 设置mac地址
func (ipmac IPMac) Set(key string, mac net.HardwareAddr) bool {
	ipmac[key] = mac
	return true
}

// Exist ip存在
func (ipmac IPMac) Exist(key string) bool {
	_, ok := ipmac[key]
	return ok
}

// Full 判断是否所有的IP都获取到mac地址了
func (ipmac IPMac) Full() bool {
	for _, one := range ipmac {
		if one == nil {
			return false
		}
	}

	return true
}
