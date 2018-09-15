package main

import (
	"gobot.io/x/gobot/platforms/mqtt"
	"gobot.io/x/gobot"
	"fmt"
	"github.com/gw123/escpos-go/printer/conntion/scanner"
		"github.com/gw123/escpos-go/printer/conntion"
	"github.com/gw123/escpos-go/printer/driver/escpos"
)

/***
mqttClientId: clientId+"|securemode=3,signmethod=hmacsha1,timestamp=132323232|"
mqttUsername: deviceName+"&"+productKey
mqttPassword: sign_hmac(deviceSecret,content)
*/

func main() {

	fmt.Println("USB打印机:")
	usbPrinterList := scanner.GetUsbPrinter("/dev/usb")
	fmt.Println(usbPrinterList)
	fmt.Println()


	host := "a1LrZ29YEJ9.iot-as-mqtt.cn-shanghai.aliyuncs.com:1883"
	clientId := "router01|securemode=3,signmethod=hmacsha1,timestamp=132323232|"
	username := "router01&a1LrZ29YEJ9"
	password := "d250d030eb54bbb345f655783d629d5156c69459"
	mqttAdaptor := mqtt.NewAdaptorWithAuth(host, clientId, username, password)

	conn, err := conntion.NewUsbConntion(usbPrinterList[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	printDirver := escpos.NewEscpos(conn)

	work := func() {
		mqttAdaptor.On("/a1LrZ29YEJ9/router01/get", func(msg mqtt.Message) {
			fmt.Println(string(msg.Payload()))

			/***
	 * 测试xml ,测试解析xml文件
	 */

			root, err := escpos.ParseString(string(msg.Payload()))
			if err != nil {
				fmt.Println(err)
				return
			}

			/***
			 * 开始打印
			 */
			printDirver.Linefeed()
			printDirver.WriteXml(root)
			printDirver.Cut()
			printDirver.End()
			printDirver.Linefeed()

		})
		//mqttAdaptor.On("hola", func(msg mqtt.Message) {
		//	fmt.Println(msg)
		//})
		//data := []byte("o")
		//gobot.Every(1*time.Second, func() {
		//	mqttAdaptor.Publish("/a1LrZ29YEJ9/router01/update", data)
		//})
		//gobot.Every(5*time.Second, func() {
		//	mqttAdaptor.Publish("hola", data)
		//})
	}

	robot := gobot.NewRobot("mqttBot",
		[]gobot.Connection{mqttAdaptor},
		work,
	)

	robot.Start()
}
