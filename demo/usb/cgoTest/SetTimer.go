package main

import (
	"syscall"
	//~ "time"
	"unsafe"
	"golang.org/x/sys/windows"
	"fmt"
)

// String returns a human-friendly display name of the hotkey
// such as "Hotkey[Id: 1, Alt+Ctrl+O]"
var (
	user32                  = windows.NewLazySystemDLL("user32.dll")
	procSetWindowsHookEx    = user32.NewProc("SetWindowsHookExW")
	procLowLevelKeyboard    = user32.NewProc("LowLevelKeyboardProc")
	procCallNextHookEx      = user32.NewProc("CallNextHookEx")
	procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	procGetMessage          = user32.NewProc("GetMessageW")
	procTranslateMessage    = user32.NewProc("TranslateMessage")
	procDispatchMessage     = user32.NewProc("DispatchMessageW")
	setTimer                = user32.NewProc("SetTimer")
	killTimer               = user32.NewProc("KillTimer")
	keybd_event             = user32.NewProc("keybd_event")
	keyboardHook            HHOOK
)

const (
	WH_KEYBOARD_LL = 13
	WH_KEYBOARD    = 2
	WM_KEYDOWN     = 256
	WM_SYSKEYDOWN  = 260
	WM_KEYUP       = 257
	WM_SYSKEYUP    = 261
	WM_KEYFIRST    = 256
	WM_KEYLAST     = 264
	PM_NOREMOVE    = 0x000
	PM_REMOVE      = 0x001
	PM_NOYIELD     = 0x002
	WM_LBUTTONDOWN = 513
	WM_RBUTTONDOWN = 516
	NULL           = 0
)

type (
	DWORD uint32
	WPARAM uintptr
	LPARAM uintptr
	LRESULT uintptr
	HANDLE uintptr
	HINSTANCE HANDLE
	HHOOK HANDLE
	HWND HANDLE
	UINT uint32
)

type HOOKPROC func(int, WPARAM, LPARAM) LRESULT
type TIMERPROC func(HWND, UINT, uintptr, DWORD) LRESULT

func SetTimer(hwnd HWND, uintptr2 uintptr, uint2 UINT, timerproc TIMERPROC) uintptr {
	ret, _, err := setTimer.Call(uintptr(hwnd), uintptr2, uintptr(uint2), uintptr(syscall.NewCallback(timerproc)))
	fmt.Println(err)
	return uintptr(ret)
}

func KillTimer(hwnd HWND, timer uintptr) bool {
	ret, _, _ := killTimer.Call(uintptr(hwnd), timer)
	return ret != 0
}

func Keybd_Event(bVk, bScan byte, dwFlags DWORD, dwExtraInfo uintptr) {
	keybd_event.Call(uintptr(bVk), uintptr(bScan), uintptr(dwFlags), uintptr(dwExtraInfo));
	return
}
// http://msdn.microsoft.com/en-us/library/windows/desktop/dd162805.aspx
type POINT struct {
	X, Y int32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms644958.aspx
type MSG struct {
	Hwnd    HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

func GetMessage(msg *MSG, hwnd HWND, msgFilterMin uint32, msgFilterMax uint32) int {
	ret, _, _ := procGetMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax))
	return int(ret)
}

//定时器没有起作用
func main()  {

	 SetTimer(0,0,80, (TIMERPROC)(func(hwnd HWND, i UINT, u uintptr, dword DWORD) LRESULT  {

		 fmt.Printf("timeout \n")
		return 0
	}))

	var msg MSG
	for {
		GetMessage(&msg, 0, 0, 0)
	}
}