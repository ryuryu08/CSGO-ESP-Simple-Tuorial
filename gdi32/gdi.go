package gdi32

/*
#cgo LDFLAGS: -lgdi32
#include <windows.h>
#include <wingdi.h>
#include <espgdi32.h>
*/
import "C"
import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	user32             = syscall.MustLoadDLL("user32.dll")
	procEnumWindows    = user32.MustFindProc("EnumWindows")
	procGetWindowTextW = user32.MustFindProc("GetWindowTextW")
	procGetWindowRect  = user32.MustFindProc("GetWindowRect")
	procGetDC		   = user32.MustFindProc("GetDC")
)

type WindowRect struct {
	Left int32
	Top int32
	Right int32
	Bottom int32
}

func GetDC(windowHandle syscall.Handle) syscall.Handle {
	ret, _, err := procGetDC.Call(uintptr(windowHandle))
	length := int32(ret)
	if length == 0 {
		if err != nil {
			panic(err.Error())
		} else {
			err = syscall.EINVAL
		}
	}
	return syscall.Handle(ret)
}

func EnumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumWindows.Addr(), 2, enumFunc, lparam, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetWindowText(hwnd syscall.Handle, str *uint16, maxCount int32) (len int32, err error) {
	r0, _, e1 := syscall.Syscall(procGetWindowTextW.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	len = int32(r0)
	if len == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func FindWindow(title string) (syscall.Handle, error) {
	var hwnd syscall.Handle
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		b := make([]uint16, 200)
		_, err := GetWindowText(h, &b[0], int32(len(b)))
		if err != nil {
			// ignore the error
			return 1 // continue enumeration
		}
		if syscall.UTF16ToString(b) == title {
			// note the window
			hwnd = h
			return 0 // stop enumeration
		}
		return 1 // continue enumeration
	})
	_ = EnumWindows(cb, 0)
	if hwnd == 0 {
		return 0, fmt.Errorf("no window with title '%s' found", title)
	}
	return hwnd, nil
}

func GetWindowRect(handle syscall.Handle) WindowRect {
	var windowRect WindowRect
	ret, _, err := procGetWindowRect.Call(uintptr(handle), uintptr(unsafe.Pointer(&windowRect)))
	length := int32(ret)
	if length == 0 {
		if err != nil {
			panic(err.Error())
		} else {
			err = syscall.EINVAL
		}
	}
	return windowRect
}


func SetLineColor(r, g, b int) {
	C.SetLineColor(C.int(r), C.int(g) ,C.int(b))
}

func DrawBorderBox(x, y, w, h, thickness int){
	C.DrawBorderBox(C.int(x), C.int(y), C.int(w), C.int(h), C.int(thickness))
}

func SetEnemyBrush(r, g, b int) {
	C.SetEnemyBrush(C.int(r), C.int(g) ,C.int(b))
}

func SetGameHdc(hdc syscall.Handle) {
	C.SetGameHdc(C.HDC(unsafe.Pointer(hdc)))
}