package main

import (
	"syscall"
	//~ "time"
	"unsafe"
	"golang.org/x/sys/windows"
	"sync"
	"fmt"
	"time"
)

// String returns a human-friendly display name of the hotkey
// such as "Hotkey[Id: 1, Alt+Ctrl+O]"
var (
	user32                  = windows.NewLazySystemDLL("user32.dll")
	kernel32                = windows.NewLazySystemDLL("Kernel32.dll")
	procSetWindowsHookEx    = user32.NewProc("SetWindowsHookExA")
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
	getModuleHandle         = kernel32.NewProc("GetModuleHandleW")
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

type KBDLLHOOKSTRUCT struct {
	VkCode      DWORD
	ScanCode    DWORD
	Flags       DWORD
	Time        DWORD
	DwExtraInfo uintptr
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

func SetWindowsHookEx(idHook int, lpfn HOOKPROC, hMod HINSTANCE, dwThreadId DWORD) HHOOK {
	ret, _, _ := procSetWindowsHookEx.Call(
		uintptr(idHook),
		uintptr(syscall.NewCallback(lpfn)),
		uintptr(hMod),
		uintptr(dwThreadId),
	)
	return HHOOK(ret)
}

func GetModuleHandle(lpModuleName uintptr ) HINSTANCE  {
	ret ,_,_ := getModuleHandle.Call(lpModuleName)
	return  HINSTANCE(ret)
}

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

func CallNextHookEx(hhk HHOOK, nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
	ret, _, _ := procCallNextHookEx.Call(
		uintptr(hhk),
		uintptr(nCode),
		uintptr(wParam),
		uintptr(lParam),
	)
	return LRESULT(ret)
}

func UnhookWindowsHookEx(hhk HHOOK) bool {
	ret, _, _ := procUnhookWindowsHookEx.Call(
		uintptr(hhk),
	)
	return ret != 0
}

func GetMessage(msg *MSG, hwnd HWND, msgFilterMin uint32, msgFilterMax uint32) int {
	ret, _, _ := procGetMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax))
	return int(ret)
}

func TranslateMessage(msg *MSG) bool {
	ret, _, _ := procTranslateMessage.Call(
		uintptr(unsafe.Pointer(msg)))
	return ret != 0
}

func DispatchMessage(msg *MSG) uintptr {
	ret, _, _ := procDispatchMessage.Call(
		uintptr(unsafe.Pointer(msg)))
	return ret
}

func LowLevelKeyboardProc(nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
	ret, _, _ := procLowLevelKeyboard.Call(
		uintptr(nCode),
		uintptr(wParam),
		uintptr(lParam),
	)
	return LRESULT(ret)
}

const BufLen = 64

type KeyBorardBuffer struct {
	KeyBuffer           [BufLen]byte
	KeyDetailBuffer     [64]KBDLLHOOKSTRUCT
	Pos                 uint32
	timer               uintptr
	LastKeyDownTime     DWORD
	LastKeyDownFullTime int64
}
var resetLock = sync.Mutex{}
func (this *KeyBorardBuffer) Reset(isOutput bool) {
	resetLock.Lock()
	var i uint32 = 0
	for ; i < this.Pos; i++ {
		// 输出
		if isOutput {
			fmt.Printf("%c",this.KeyBuffer[i])
			Keybd_Event(byte(this.KeyDetailBuffer[i].VkCode),
				byte(this.KeyDetailBuffer[i].ScanCode),
				this.KeyDetailBuffer[i].Flags,
				this.KeyDetailBuffer[i].DwExtraInfo)

			Keybd_Event(byte(this.KeyDetailBuffer[i].VkCode),
				byte(this.KeyDetailBuffer[i].ScanCode),
				0x0002,
				this.KeyDetailBuffer[i].DwExtraInfo)

		}
		this.KeyBuffer[i] = 0
	}

	this.Pos = 0
	this.timer = 0
	this.LastKeyDownTime = 0
	resetLock.Unlock()
}

func (this *KeyBorardBuffer) Append(temp KBDLLHOOKSTRUCT) uint32 {
	if this.Pos >= BufLen {
		return this.Pos
	}
	resetLock.Lock()
	this.KeyBuffer[this.Pos] = byte(temp.VkCode)
	this.KeyDetailBuffer[this.Pos] = temp
	this.Pos++
	//this.LastKeyDownTime = temp.Time
	this.LastKeyDownTime = temp.Time
	this.LastKeyDownFullTime = time.Now().UnixNano()
	resetLock.Unlock()
	return this.Pos
}

func (this *KeyBorardBuffer) isValidCode() bool {
	if this.Pos > 16 {
		return true;
	}
	return false;
}

var keyBorardBuffer KeyBorardBuffer = KeyBorardBuffer{}

func Start(wait sync.WaitGroup) {
	// defer user32.Release()
	keyboardHook = SetWindowsHookEx(WH_KEYBOARD_LL,
		(HOOKPROC)(func(nCode int, wparam WPARAM, lparam LPARAM) LRESULT {
			kbdstruct := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lparam))
			//加上标记防止消息重复处理
			if (kbdstruct.DwExtraInfo == 1) {
				return CallNextHookEx(keyboardHook, nCode, wparam, lparam)
			}
			//
			if (keyBorardBuffer.timer != 0) {
				KillTimer(0, keyBorardBuffer.timer);
				keyBorardBuffer.timer = 0;
			}

			kbdstruct.DwExtraInfo = 1;
			if nCode == 0 && wparam == WM_KEYDOWN {
				if kbdstruct.VkCode >= '0' && kbdstruct.VkCode <= '9' && (
					keyBorardBuffer.LastKeyDownTime == 0 || kbdstruct.Time-keyBorardBuffer.LastKeyDownTime < 80) {
					keyBorardBuffer.Append(*kbdstruct)
					return 1;
				} else if kbdstruct.VkCode == 0x0D {
					//判断是否为二维码
					if (keyBorardBuffer.isValidCode()) {
						fmt.Printf("Code: %s \n", keyBorardBuffer.KeyBuffer[0: keyBorardBuffer.Pos]);
						keyBorardBuffer.Reset(false);
						return 1;
					} else {
						keyBorardBuffer.Reset(true);
					}
				} else {
					keyBorardBuffer.Reset(true);
				}

			}
			return CallNextHookEx(keyboardHook, nCode, wparam, lparam)
		}), GetModuleHandle(0), 0)
	var msg MSG
	for {
		GetMessage(&msg, 0, 0, 0)
	}

	//for {
	//
	// GetMessage(uintptr(unsafe.Pointer(msg)), 0, 0, 0)
	// TranslateMessage(msg)
	// DispatchMessage(msg)
	// // fmt.Println("key pressed:")
	//
	// }

	UnhookWindowsHookEx(keyboardHook)
	keyboardHook = 0
	//fmt.Println("end")
	wait.Done()
}

func main() {
	wait := sync.WaitGroup{}
	wait.Add(1)
	go Start(wait)

	tick := time.Tick(1000 * time.Millisecond)
	for {
		select {
		case <-tick:

			now := time.Now().UnixNano()
			temp := (now - keyBorardBuffer.LastKeyDownFullTime) / int64(time.Millisecond)
			if temp > 80 && keyBorardBuffer.Pos >0 {
				keyBorardBuffer.Reset(true)
			}
			//fmt.Printf("temp time : %d  last %d   now %d \n", temp, keyBorardBuffer.LastKeyDownFullTime, now)
		}
	}
	wait.Wait()

}
