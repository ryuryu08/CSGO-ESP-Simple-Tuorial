package helper

import (
	"fmt"
	"github.com/winlabs/gowin32"
	"github.com/winlabs/gowin32/wrappers"
	"golang.org/x/sys/windows"
	"strings"
	"syscall"
	"unsafe"
)

func FindProcessIdByName( name string) (pid uint, err error) {
	// 获取进程列表
	processes, err := gowin32.GetProcesses()
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	// 遍历查找进程
	for _, p := range processes {
		if strings.ToLower(p.ExeFile) == strings.ToLower(name) {
			pid = p.ProcessID
			return
		}
	}
	return
}

func GetProcessHandleByPid(pid uint) syscall.Handle {
	kernel32, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		fmt.Println(err)
	}
	defer kernel32.Release()
	handle, err := syscall.OpenProcess(windows.PROCESS_VM_READ, false, uint32(pid))
	if err != nil {
		fmt.Println(err)
	}
	return handle
}

func GetModuleHandleByDllNameWithProcessId(pid uint, moduleName string) uint{
	moduleInfos, err := GetProcessModules32(uint32(pid))
	if err != nil {
		panic(err.Error())
	}
	for _, moduleInfo := range moduleInfos {
		if moduleInfo.ModuleName == moduleName {
			return *(*uint)(unsafe.Pointer(&moduleInfo.ModuleBaseAddress))
		}
	}
	panic("Can not find Game Module")
}

func GetProcessModules32(pid uint32) ([]gowin32.ModuleInfo, error) {
	hSnapshot, err := wrappers.CreateToolhelp32Snapshot(wrappers.TH32CS_SNAPMODULE | wrappers.TH32CS_SNAPMODULE32, pid)
	if err != nil {
		return nil, gowin32.NewWindowsError("CreateToolhelp32Snapshot", err)
	}
	defer wrappers.CloseHandle(hSnapshot)
	me := wrappers.MODULEENTRY32{}
	me.Size = uint32(unsafe.Sizeof(me))
	if err := wrappers.Module32First(hSnapshot, &me); err != nil {
		return nil, gowin32.NewWindowsError("Module32First", err)
	}
	var mi []gowin32.ModuleInfo
	for {
		mi = append(mi, gowin32.ModuleInfo{
			ProcessID:         uint(me.ProcessID),
			ModuleBaseAddress: me.ModBaseAddr,
			ModuleBaseSize:    uint(me.ModBaseSize),
			ModuleHandle:      me.Module,
			ModuleName:        syscall.UTF16ToString((&me.ModuleName)[:]),
			ExePath:           syscall.UTF16ToString((&me.ExePath)[:]),
		})
		err := wrappers.Module32Next(hSnapshot, &me)
		if err == wrappers.ERROR_NO_MORE_FILES {
			return mi, nil
		} else if err != nil {
			return nil, gowin32.NewWindowsError("Module32Next", err)
		}
	}
}
