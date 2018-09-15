package scanner

import (
	"fmt"
	"strings"
	"time"
	"strconv"
	"net"
	"sync"
	"github.com/gw123/escpos-go/util"
)

var printerList []string = make([]string, 0)

// 主机所在网络内escpos打印机
func GetNetPrinter() []string {
	ips := util.GetHostNet()
	for _, ip := range ips {
		scanNetwork(ip, "24")
	}
	return printerList
}

// 扫描网络内打印机
func scanNetwork(net string, netMask string) {
	sp := strings.Split(net, ".")
	waitGroup := sync.WaitGroup{}
	mutex := sync.Mutex{}
	for index := 1; index < 255; index++ {
		ipaddr := fmt.Sprintf("%s.%s.%s.%s", sp[0], sp[1], sp[2], strconv.Itoa(index))
		//fmt.Println("Addr ", ipaddr)
		waitGroup.Add(1)
		go func() {
			if IsPrinterAvailable(ipaddr, 9100, time.Second*3) && net != ipaddr {
				//fmt.Println(ipaddr, "是打印机")
				mutex.Lock()
				printerList = append(printerList, ipaddr)
				mutex.Unlock()
			} else {
				//fmt.Println(ipaddr, "不是打印机")
			}
			waitGroup.Done()
		}()
	}
	waitGroup.Wait()

}

// IsPrinterAvailable 检查打印机端口是否开放
func IsPrinterAvailable(ip string, port int, timeout time.Duration) bool {
	remoteAddr := ip + ":" + strconv.Itoa(port)
	conn, _ := net.DialTimeout("tcp", remoteAddr, timeout)

	if conn == nil {
		return false
	}

	defer conn.Close()

	return true
}
