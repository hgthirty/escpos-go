package main

import (
	"fmt"
	"io/ioutil"
	"encoding/xml"
	"os"
	"time"
	"github.com/gw123/escpos-go/printer/conntion"
	"github.com/gw123/escpos-go/printer/driver/escpos"
)

func main() {

	/***
	 * 测试xml ,测试解析xml文件
	 */
	root, err := escpos.ParseLocalXml("test.xml")
	if err != nil {
		fmt.Println(err)
		return
	}

	/***
	 * 建立连接打印机的conntion
	 */

	//方式1  usb连接
	//conn, err := conntion.NewUsbConntion("/dev/usb/lp0")
	//方式2  网络连接
	conn, err := conntion.NewNetConntion("10.0.1.150:9100")
	defer conn.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	/***
	 * 开始打印
	 */
	printDirver := escpos.NewEscpos(conn)
	printDirver.WriteXml(root)
	printDirver.Cut()
	printDirver.End()

	return
}


/***
 * 通过golang 的 ioutil 读取文件打印
 */
func test1() {

	fileContent, err := ioutil.ReadFile("test.xml")
	if err != nil {
		fmt.Println(err)
	}
	root := escpos.Root{}
	err = xml.Unmarshal(fileContent, &root)

	file, err := os.OpenFile("/dev/usb/lp2", os.O_WRONLY, 0660)
	if err != nil {
		fmt.Println(err)
	}
	printDirver := escpos.NewEscpos(file)

	//fmt.Print(root)
	for _, line := range root.Lines {
		printDirver.WriteLine(line)
	}
	printDirver.Linefeed()
	printDirver.Linefeed()
	printDirver.Cut()
	return
}

/***
 * 模仿抖音 连续打印图片
 */
func test2() {
	file, err := os.OpenFile("/dev/usb/lp2", os.O_WRONLY, 0660)
	if err != nil {
		fmt.Println(err)
	}
	printDirver := escpos.NewEscpos(file)
	printDirver.Init()
	rootPath := "/home/gw/data/doubi/"
	printDirver.WriteLocalImage(rootPath + "1.jpeg")
	time.Sleep(time.Second * 1)
	printDirver.WriteLocalImage(rootPath + "2.jpeg")
	time.Sleep(time.Second * 1)
	printDirver.WriteLocalImage(rootPath + "3.jpeg")
	time.Sleep(time.Second * 1)
	printDirver.WriteLocalImage(rootPath + "4.jpeg")
	time.Sleep(time.Second * 1)
	printDirver.WriteLocalImage(rootPath + "5.jpeg")
	time.Sleep(time.Second * 1)
	printDirver.WriteLocalImage(rootPath + "6.png")
	time.Sleep(time.Second * 1)
	printDirver.WriteLocalImage(rootPath + "7.jpeg")
	printDirver.Linefeed()
	printDirver.Cut()
	return
}

/***
 * 基础功能测试
 */
func test3() {

	client ,err := conntion.NewNetConntion("10.0.1.150:9100")
	defer client.Close()
	if err != nil{
		fmt.Println(err)
		return
	}

	printDirver := escpos.NewEscpos(client)
	printDirver.Init()
	printDirver.SetFont("A")
	printDirver.SetFontSize(2, 2)
	printDirver.SetAlign("center")
	printDirver.WriteGbk("*乐惠扫码点餐*")
	printDirver.Linefeed()
	printDirver.Linefeed()
	printDirver.SetAlign("left")
	printDirver.WriteLRLine("唐食", "#0001")
	printDirver.Linefeed()

	printDirver.WriteALine('*')

	/***测试打印一行**/
	line := escpos.NewLine(1, "", "left")
	cell := escpos.NewCell(0.25, "品名", "left")
	line.AppendCell(*cell)
	cell = escpos.NewCell(0.25, "数量", "center")
	line.AppendCell(*cell)
	cell = escpos.NewCell(0.25, "单价", "center")
	line.AppendCell(*cell)
	cell = escpos.NewCell(0.25, "总额", "right")
	line.AppendCell(*cell)
	printDirver.WriteLine(*line)
	printDirver.WriteALine('*')

	line = escpos.NewLine(2, "测试一行", "left")
	printDirver.WriteLine(*line)


}

/***
 * 基础功能测试
 */
func test4() {

	client ,err := conntion.NewNetConntion("10.0.1.150:9100")
	defer client.Close()
	if err != nil{
		fmt.Println(err)
		return
	}

	printDirver := escpos.NewEscpos(client)

	//测试xml
	fileContent, err := ioutil.ReadFile("/home/gw/data/testimg.data")
	if err != nil {
		fmt.Println(err)
	}
	printDirver.WriteRaw(fileContent)
	printDirver.Linefeed()
	printDirver.Cut()

}

