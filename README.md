# escpos-go
go escpos
# 使用用例参考 /demo/demo.go
`
    
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
 	//conn, err := conntion.NewUsbConntion("/dev/usb/lp2")
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
`
