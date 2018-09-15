package main

import (
	"fmt"
	"github.com/gw123/escpos-go/printer/conntion/scanner"
	"github.com/gw123/escpos-go/util"
	"io"
	"github.com/gw123/escpos-go/printer/driver/escpos"
	"github.com/gw123/escpos-go/printer/conntion"
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


	/***
	 * 建立连接打印机的conntion
	 */
	//方式1  usb连接
	for _, usbpath := range usbPrinterList {
		conn, err := conntion.NewUsbConntion(usbpath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("usb打印机"+usbpath)
		printerRun(conn, "usb打印机"+usbpath)
	}

	return
	//方式2  网络连接
	fmt.Println("网络打印机:")
	printerList := scanner.GetNetPrinter()
	fmt.Println(printerList)
	for _, addr := range printerList {
		conn, err := conntion.NewNetConntion(addr + ":9100")
		if err != nil {
			fmt.Println(err)
			continue
		}
		printerRun(conn, "网络打印机"+addr)
	}

}

func printerRun(conn io.ReadWriter, title string) {
	/***
	 * 测试xml ,测试解析xml文件
	 */
	root, err := escpos.ParseLocalXml("test.xml")
	if err != nil {
		fmt.Println(err)
		return
	}

	/***
	 * 开始打印
	 */
	printDirver := escpos.NewEscpos(conn)
	printDirver.Linefeed()
	printDirver.WriteGbk(title)

	//printDirver.WriteGbk("adc")
	printDirver.Linefeed()
	printDirver.WriteXml(root)
	printDirver.Cut()
	printDirver.End()
	printDirver.Linefeed()
}
