
#include <Windows.h>
#include <wingdi.h>

HFONT Font;
HDC gameHDC;
HBRUSH EnemyBrush;
COLORREF LineColor;
COLORREF TextColor;

void SetLineColor(int r, int g, int b) {
	LineColor = RGB(r, g, b);
}
void SetTexTColor(int r, int g, int b) {
	TextColor = RGB(r, g, b);
}

void SetGameHdc(HDC hdc) {
	gameHDC = hdc;
}

void SetGameHDCByWindowHandle(HWND handle) {
	gameHDC = GetDC(handle);
}

void SetEnemyBrush(int r, int g, int b) {
	EnemyBrush = CreateSolidBrush(RGB(r, g, b));
}

/*Drawing function*/
void DrawFilledRect(int x, int y, int w, int h) {
	RECT rect = { x, y, x + w, y + h };
	FillRect(gameHDC, &rect, EnemyBrush);
}
void DrawBorderBox(int x, int y, int w, int h, int thickness)
{
	DrawFilledRect(x, y, w, thickness);
	DrawFilledRect(x, y, thickness, h);
	DrawFilledRect((x + w), y, thickness, h);
	DrawFilledRect(x, y + h, w + thickness, thickness);
}

void DrawLine(float StartX, float StartY, float EndX, float EndY)
{
	HPEN hOPen;
	HPEN hNPen = CreatePen(PS_SOLID, 2, LineColor);
	MoveToEx(gameHDC, StartX, StartY, NULL);
	LineTo(gameHDC, EndX, EndY);
	DeleteObject(SelectObject(gameHDC, hNPen));
}


