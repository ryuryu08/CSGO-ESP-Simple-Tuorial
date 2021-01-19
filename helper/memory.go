package helper

import (
	"encoding/binary"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

var (
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procReadProcessMemory = modkernel32.NewProc("ReadProcessMemory")
)

func ReadMemoryUint32(handle syscall.Handle, address uint) uint32 {
	var data [4]byte
	_, err := ReadProcessMemory(handle, uintptr(address), data[:])
	if err != nil {
		panic(err.Error())
	}
	return binary.LittleEndian.Uint32(data[:])
}


func ReadMemoryViewMatrix(handle syscall.Handle, address uint) (matrix [4][4]float32) {
	var data [64]byte
	_, err := ReadProcessMemory(handle, uintptr(address), data[:])
	if err != nil {
		panic(err.Error())
	}
	array := Float32SliceFromBytes(data[:])
	for i, val := range array {
		matrix[i/4][i%4] = val
	}
	return
}

func ReadMemoryFloat(handle syscall.Handle, address uint) float32 {
	var data float32
	var length uint32
	_, _, _ = procReadProcessMemory.Call(
		uintptr(handle),
		uintptr(address),
		uintptr(unsafe.Pointer(&data)), 4, uintptr(unsafe.Pointer(&length)))
	return data
}



func _ReadProcessMemory(handle syscall.Handle, baseAddress uintptr, buffer uintptr, size uintptr, numRead *uintptr) (err error) {
	r1, _, e1 := syscall.Syscall6(procReadProcessMemory.Addr(), 5, uintptr(handle), baseAddress, buffer, size, uintptr(unsafe.Pointer(numRead)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = e1
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func ReadProcessMemory(handle syscall.Handle, baseAddress uintptr, dest []byte) (numRead uintptr, err error) {
	n := len(dest)
	if n == 0 {
		return 0, nil
	}
	if err = _ReadProcessMemory(handle, baseAddress, uintptr(unsafe.Pointer(&dest[0])), uintptr(n), &numRead); err != nil {
		return 0, err
	}
	return numRead, nil
}
