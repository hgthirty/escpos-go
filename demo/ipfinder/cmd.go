package main

import (
	"fmt"
	"github.com/gw123/escpos-go/printer/conntion/scanner"
	"github.com/gw123/escpos-go/util"
)



func main() {
	fmt.Println("本机地址:")
	hostAddrs := util.GetHostNet()
	fmt.Println(hostAddrs)
	fmt.Println()

	fmt.Println("USB打印机:")
	usbPrinterList := scanner.GetUsbPrinter("/dev/usb")
	fmt.Println(usbPrinterList)
	fmt.Println()

	fmt.Println("网络打印机:")
	printerList := scanner.GetNetPrinter()
	fmt.Println(printerList)
}
