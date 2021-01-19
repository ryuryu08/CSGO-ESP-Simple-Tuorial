package view

import (
	"main/gdi32"
)

type ScreenPosition struct {
	X 			int32
	TopY 		int32
	BottomY 	int32
	IsShow 		bool
}

type RolePosition struct {
	X float32
	Y float32
	Z float32
}

type Angles struct {
	Pitch 	float32
	Yaw 	float32
	Roll	float32
}

func WorldToScreen(
	p RolePosition,
	viewMatrix [4][4]float32,
	gameWindow gdi32.WindowRect,
) (screenPosition ScreenPosition){
	width := gameWindow.Right - gameWindow.Left
	height := gameWindow.Bottom - gameWindow.Top
	width /= 2
	height /= 2
	w := viewMatrix[2][0] * p.X + viewMatrix[2][1] * p.Y + viewMatrix[2][2] * p.Z + viewMatrix[2][3]

	bili := 1 / w
	if bili < 0 {
		screenPosition.IsShow = false
		return
	}

	x := float32(width)
	x += (viewMatrix[0][0] * p.X + viewMatrix[0][1] * p.Y + viewMatrix[0][2] * p.Z + viewMatrix[0][3]) * float32(width) * bili

	topY := float32(height) - (viewMatrix[1][0] * p.X + viewMatrix[1][1] * p.Y + viewMatrix[1][2] * (p.Z + 8) + viewMatrix[1][3]) * float32(height) * bili

	bottomY := float32(height) - (viewMatrix[1][0] * p.X + viewMatrix[1][1] * p.Y + viewMatrix[1][2] * (p.Z + 78) + viewMatrix[1][3]) * float32(height) * bili

	screenPosition.X = int32(x)
	screenPosition.TopY = int32(topY)
	screenPosition.BottomY = int32(bottomY)

	return
}