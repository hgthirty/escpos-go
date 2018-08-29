package conntion

import (
	"os"
	"fmt"
)

type UsbConntion struct {
	Conn *os.File
	Status     Status
}

func NewUsbConntion(filename string) (usb *UsbConntion, err error) {
	usb = new(UsbConntion)
	file, err := os.OpenFile(filename, os.O_RDWR, 0660)
	if err != nil {
		fmt.Print(err)
		return
	}
	usb.Conn = file
	return
}

func (this *UsbConntion) Info() Status {
	return this.Status
}

func (this *UsbConntion) Write(data []byte) (count int, err error) {
	count, err = this.Conn.Write(data)
	return
}

func (this *UsbConntion) Read(data []byte) (count int, err error) {
	count, err = this.Conn.Read(data)
	return
}

func (this *UsbConntion) Close()  {
	this.Conn.Close()
}