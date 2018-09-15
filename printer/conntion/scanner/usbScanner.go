package scanner

import (
	"path/filepath"
	"os"
	"fmt"
	"strings"
)

// 主机的usb打印机
func GetUsbPrinter(path string) []string {
	var usbPrinterList []string = make([]string, 0)
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if (f == nil) {
			return err
		}
		if f.IsDir() {
			return nil
		}

		if strings.HasPrefix(path,"/dev/usb/lp"){
			usbPrinterList = append(usbPrinterList, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
		return nil
	}
	return usbPrinterList
}
