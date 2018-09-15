package conntion

import (
	"net"
	"fmt"
)

type NetConntion struct {
	Conn   net.Conn
	Status Status
}

func NewNetConntion(addr string) (netConntion *NetConntion, err error) {
	tempConn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("net.Dial 连接"+addr+"失败", err)
		return netConntion, err
	}
	netConntion = new(NetConntion)
	netConntion.Conn = tempConn
	return
}

func (this *NetConntion) Write(data []byte) (count int, err error) {
	count, err = this.Conn.Write(data)
	return
}

func (this *NetConntion) Read(data []byte) (count int, err error) {
	count, err = this.Conn.Read(data)
	return
}

func (this *NetConntion) Info() Status {
	return this.Status
}

func (this *NetConntion) Close() {
	this.Conn.Close()
}

