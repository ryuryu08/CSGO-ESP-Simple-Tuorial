package helper

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	user32             = syscall.MustLoadDLL("user32.dll")
	gdi32              = syscall.MustLoadDLL("gdi32.dll")
	procEnumWindows    = user32.MustFindProc("EnumWindows")
	procGetWindowTextW = user32.MustFindProc("GetWindowTextW")
	procGetWindowRect  = user32.MustFindProc("GetWindowRect")
	procGetDC		   = user32.MustFindProc("GetDC")
	procCreateSolidBrush  = gdi32.MustFindProc("CreateSolidBrush")
	procCreatePen  = gdi32.MustFindProc("CreatePen")
	procMoveToEx     	  = gdi32.MustFindProc("MoveToEx")
	procLineTo        	  = gdi32.MustFindProc("LineTo")
	procDeleteObject      = gdi32.MustFindProc("DeleteObject")
	procSelectObject      = gdi32.MustFindProc("SelectObject")
)


var (
	HDC syscall.Handle
	brush syscall.Handle
	PS_SOLID = 0
)

type WindowRect struct {
	Left int32
	Top int32
	Right int32
	Bottom int32
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

func SetupDrawing(hdc syscall.Handle) {
	HDC = hdc
	color := 0x32CD32
	ret, _, err := procCreateSolidBrush.Call(uintptr(unsafe.Pointer(&color)))
	length := int32(ret)
	if length == 0 {
		if err != nil {
			panic(err.Error())
		} else {
			err = syscall.EINVAL
		}
	}
	brush = syscall.Handle(ret)
}

func DrawLine(startX float32, startY float32, endX float32, endY float32, color int32) {
	pen, _, err := procCreatePen.Call(uintptr(unsafe.Pointer(&PS_SOLID)), 2, uintptr(brush))
	length := int32(pen)
	if length == 0 {
		if err != nil {
			panic(err.Error())
		} else {
			err = syscall.EINVAL
		}
	}
	ret, _, err := procMoveToEx.Call(
		uintptr(HDC),
		uintptr(unsafe.Pointer(&startX)),
		uintptr(unsafe.Pointer(&startY)),
		uintptr(unsafe.Pointer(nil)))
	length = int32(ret)
	if length == 0 {
		if err != nil {
			panic(err.Error())
		} else {
			err = syscall.EINVAL
		}
	}
	ret, _, err = procLineTo.Call(
		uintptr(HDC),
		uintptr(unsafe.Pointer(&endX)),
		uintptr(unsafe.Pointer(&endY)))
	object, _, _ := procSelectObject.Call(
		uintptr(HDC),
		pen)
	length = int32(ret)
	if length == 0 {
		if err != nil {
			panic(err.Error())
		} else {
			err = syscall.EINVAL
		}
	}
	ret, _, err = procDeleteObject.Call(
		object)
	length = int32(ret)
	if length == 0 {
		if err != nil {
			panic(err.Error())
		} else {
			err = syscall.EINVAL
		}
	}
}
