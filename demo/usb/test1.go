package main

import (
	"syscall"
	"unsafe"
	"log"
	"fmt"
)

func main() {
	calltest3()
}

func calltest3() {

	const (
		DIGCF_DEFAULT         = 0x00000001
		DIGCF_PRESENT         = 0x00000002
		DIGCF_ALLCLASSES      = 0x00000004
		DIGCF_PROFILE         = 0x00000008
		DIGCF_DEVICEINTERFACE = 0x00000010
	)

	h := syscall.MustLoadDLL("SetupAPI.dll")
	SetupDiGetClassDevs := h.MustFindProc("SetupDiGetClassDevsW")

	handel , _ ,err := SetupDiGetClassDevs.Call(uintptr(0), uintptr(0), uintptr(0), uintptr(DIGCF_ALLCLASSES|DIGCF_PRESENT))
	if err != nil{
		fmt.Println(err)
	}else {
		fmt.Println( err ,handel)
	}

	SetupDiEnumDeviceInterfaces := h.MustFindProc("SetupDiEnumDeviceInterfaces")
	icount :=0
	SetupDiEnumDeviceInterfaces.Call( handel ,0 , 0 ,uintptr(icount) , )
	//icount := 0
	//for  {
	//
	//}

	//while ( SetupDiEnumDeviceInterfaces(hDevInfo, NULL,&USB_Device_GUID, icount, &DeviceInterfaceData) )
	//{
	//
	//}
}

func calltest1() {
	h := syscall.MustLoadDLL("kernel32.dll")
	c := h.MustFindProc("CreateFileW")

	//lpTotalNumberOfBytes := 0x00000000

	r2, _, err := c.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("D:\\tmp.data"))),
		uintptr(0x40000000),
		uintptr(0x00000000),
		uintptr(0),
		uintptr(1),
		uintptr(0x00000080),
		uintptr(0),
	)
	if r2 != 0 {
		log.Println(r2, err)
	} else {
		log.Println(r2, err)
	}
}
func calltest() {
	h := syscall.MustLoadDLL("kernel32.dll")
	c := h.MustFindProc("GetDiskFreeSpaceExW")
	lpFreeBytesAvailable := int64(0)
	lpTotalNumberOfBytes := int64(0)
	lpTotalNumberOfFreeBytes := int64(0)
	r2, _, err := c.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("C:"))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)))
	if r2 != 0 {
		log.Println(r2, err, lpFreeBytesAvailable/1024/1024, "MB")
	} else {
		log.Println(r2, err)
	}
}

func calltest2() {
	//首先,准备输入参数, GetDiskFreeSpaceEx需要4个参数, 可查MSDN
	dir := "D:"
	lpFreeBytesAvailable := int64(0) //注意类型需要跟API的类型相符 ,类型失败有时候会编译成功但是会发生内容不全的问题
	lpTotalNumberOfBytes := int64(0)
	lpTotalNumberOfFreeBytes := int64(0)

	//获取方法的引用
	kernel32, err := syscall.LoadLibrary("Kernel32.dll") // 严格来说需要加上
	defer syscall.FreeLibrary(kernel32)
	if err != nil {
		fmt.Println(err)
	}
	GetDiskFreeSpaceEx, err := syscall.GetProcAddress(syscall.Handle(kernel32), "GetDiskFreeSpaceExW")

	//执行之. 因为有4个参数,故取Syscall6才能放得下. 最后2个参数,自然就是0了
	r, _, _ := syscall.Syscall6(uintptr(GetDiskFreeSpaceEx), 4,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(dir))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)), 0, 0)

	// 注意, errno并非error接口的, 不可能是nil
	// 而且,根据MSDN的说明,返回值为0就fail, 不为0就是成功
	if r != 0 {
		log.Printf("Free %dmb ,Free %dmb ,Free %dmb ",
			lpFreeBytesAvailable/1024/1024,
			lpTotalNumberOfBytes/1024/1024,
			lpTotalNumberOfFreeBytes/1024/1024)
	}
}
