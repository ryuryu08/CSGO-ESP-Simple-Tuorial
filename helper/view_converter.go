package helper

type ScreenPosition struct {
	X float32
	Y float32
	IsShow bool
}

type RolePosition struct {
	X float32
	Y float32
	Z float32
}

type ViewPosition struct {
	Pitch float32
	Yaw float32
	Roll float32
}

func WorldToScreen(
	p RolePosition,
	viewMatrix [4][4]float32,
	gameWindow WindowRect,
) (screenPosition ScreenPosition){
	width := gameWindow.Right - gameWindow.Left
	height := gameWindow.Bottom - gameWindow.Top
	width /= 2
	height /= 2

	w := p.X * viewMatrix[0][2] + p.Y * viewMatrix[1][2] + p.Z * viewMatrix[2][2] + viewMatrix[3][2]
	bili := 1/ w
	if bili < 0 {
		screenPosition.IsShow = false
		return
	}

	floatX := viewMatrix[0][0] * p.X + viewMatrix[1][0] * p.Y + viewMatrix[2][0] * p.Z + viewMatrix[3][0]
	floatY := viewMatrix[0][1] * p.X + viewMatrix[1][1] * p.Y + viewMatrix[2][1] * p.Z + viewMatrix[3][1]
	floatW := viewMatrix[0][3] * p.X + viewMatrix[1][3] * p.Y + viewMatrix[2][3] * p.Z + viewMatrix[3][3]


	screenPosition.IsShow = true

	floatX /= floatW
	floatY /= floatW
	floatX /= 2
	floatY /= 2

	floatX += 0.5
	floatY += 0.5

	floatX += float32(width / 2)
	floatY = float32(height / 2) - floatY

	screenPosition.X = floatX
	screenPosition.Y = floatY

	return
}