package main

import (
	"main/gdi32"
	"main/memory"
	"main/prcocess"
	"main/view"
	"syscall"
	"time"
)

const (
	rolePtrArrayOffset = 0x4DA2E74
	rolePtrArrayPtrSize = 0x10

	// Game Property
	gameWindowName = "Counter-Strike: Global Offensive"
	gameMaxPlayer = 20
	viewMatrixOffset = 0x4D94774

	// Role Property
	roleHpOffset = 0x100
	roleTeamOffset = 0xF4
	roleXOffset = 0x138
	roleYOffset = 0x13C
	roleZOffset = 0x140

	// View Property
	rolePitchOffset = 0x12C
	roleYawOffset = 0x130
	roleRollOffset = 0x134
)

type Role struct{
	Hp       uint32
	Team     uint32
	Position view.RolePosition
}

var (
	gamePid          uint
	clientAddress    uint
	rolePtrArrayPtr  uint
	gameHandle       syscall.Handle
	viewMatrix       [4][4]float32
	gameWindowHandle syscall.Handle
	gameWindowRect   gdi32.WindowRect
	gameWindowDC     syscall.Handle

	roles            [gameMaxPlayer]Role
)

func init() {
	InitialMemoryAddress()
	InitialDrawing()
}

func main() {
	for {
		for i := 0; i < gameMaxPlayer; i++ {
			roleUpdate(i)
			if roles[i].Team == roles[0].Team {
				continue
			}
			if roles[i].Hp < 1 {
				continue
			}
			screenPosition := view.WorldToScreen(roles[i].Position, viewMatrix, gameWindowRect)
			DrawESP(screenPosition)
		}
	}
}

/*
	绘制方块
 */
func DrawESP(p view.ScreenPosition) {
	width := float32(p.TopY - p.BottomY)
	width *= 0.516515151552
	gdi32.DrawBorderBox(int(float32(p.X) - width / 2), int(p.TopY), int(width), int(p.BottomY - p.TopY), 3)
}

/**
	更新全局角色数据， 三个err处理主要是为了防止游戏在一局结束后读取失败造成程序退出
 */
func roleUpdate(i int){
	var err error
	viewMatrix, err = memory.ReadMemoryViewMatrix(gameHandle, clientAddress + viewMatrixOffset)
	role, err := memory.ReadMemoryUint32(gameHandle, rolePtrArrayPtr + uint(i * rolePtrArrayPtrSize))
	if err != nil {
		time.Sleep(time.Duration(1) * time.Second)
		return
	}
	roles[i].Hp, err = memory.ReadMemoryUint32(gameHandle, uint(role+roleHpOffset))
	if err != nil {
		return
	}
	roles[i].Team, err = memory.ReadMemoryUint32(gameHandle, uint(role+roleTeamOffset))
	roles[i].Position.X, err = memory.ReadMemoryFloat32(gameHandle, uint(role+roleXOffset))
	roles[i].Position.Y, err = memory.ReadMemoryFloat32(gameHandle, uint(role+roleYOffset))
	roles[i].Position.Z, err = memory.ReadMemoryFloat32(gameHandle, uint(role+roleZOffset))
	if err != nil {
		time.Sleep(time.Duration(1) * time.Second)
		return
	}
}

/*
	初始化绘制
 */
func InitialDrawing() {
	for {
		gameWindowHandle, err := gdi32.FindWindow(gameWindowName)
		if err != nil {
			time.Sleep(time.Duration(1) * time.Second)
			continue
		}
		gameWindowRect = gdi32.GetWindowRect(gameWindowHandle)
		gameWindowDC = gdi32.GetDC(gameWindowHandle)
		gdi32.SetGameHdc(gameWindowDC)
		gdi32.SetLineColor(34, 139, 34)
		gdi32.SetEnemyBrush(43, 244, 61)
		break
	}
}

/*
	内存读取初始化
 */
func InitialMemoryAddress() {
	for {
		var err error
		gamePid, err = prcocess.FindProcessIdByName("csgo.exe")
		if err != nil {
			time.Sleep(time.Duration(1) * time.Second)
			continue
		}
		gameHandle, err = prcocess.GetProcessHandleByPid(gamePid)
		if err != nil {
			time.Sleep(time.Duration(1) * time.Second)
			continue
		}
		clientAddress, err = prcocess.GetModuleHandleByDllNameWithProcessId(gamePid, "client.dll")
		if err != nil {
			time.Sleep(time.Duration(1) * time.Second)
			continue
		}
		rolePtrArrayPtr = clientAddress + rolePtrArrayOffset
		break
	}
}

