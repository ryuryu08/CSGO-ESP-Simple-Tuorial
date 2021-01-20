# CSGO-ESP-Tutorial

[TOC]

## 概述

ESP即Extrasensory perception的缩写， 通常被称作方框透视， 通过方框将游戏中其他目标角色框出来， 来得知额外的位置信息来作弊。

本文将指导你如何使用Go语言制作一个CSGO

由于代码也是比较少，简单易懂， 所以直接进行数据查找的说明， 代码可以clone下来看一看

[![sR1yTI.png](https://s3.ax1x.com/2021/01/20/sR1yTI.png)](https://imgchr.com/i/sR1yTI)

## 开发环境

+ 截至 2020年1月19日 CSGO Steam最新版
+ Windows 10 专业版 1903 
+ Go 1.56.6
+ Goland
+ Cheat Engine 7.2

## 制作方框透视需要什么

+ 玩家列表
+ 玩家信息（坐标、生命值、视角）
+ 视图矩阵
+ WorldToScreen函数
+ 绘制图像到游戏中
+ 数据结构和操作系统 的基础知识

## 数据查找

### 生命值与玩家列表

首先从最容易观察到的生命值进行搜索
在回合开始时， 我们的生命值为 100， 在CE中搜索 4Byte 精确数值 100。

随后我们使用任意方式减少自己的血量， 我们的血量减少为 46， 在CE中继续搜索 46， 此时CE中将剩下 21 组数据。

[![sR3k9K.png](https://s3.ax1x.com/2021/01/20/sR3k9K.png)](https://imgchr.com/i/sR3k9K)

我们可以继续变动自己的HP， 并进行搜索， 可以发现最终剩下7组数据（不确定数值）， 不过这并不重要， 只要少到可操作范围即可

随后我们从上至下依次点击右键，选择 Find out what access this address，首先观察是否存在 cmp 汇编指令来讲一个指针的值 与 0进行对比， 游戏中逻辑往往通过 if 语句对血量进行判断， 当然也不一定， 在这款游戏中， 他就是不断的判断自己血量是否为 0， 为 0 则触发死亡函数

[![sR3NHs.png](https://s3.ax1x.com/2021/01/20/sR3NHs.png)](https://imgchr.com/i/sR3NHs)

可以发现在一个地址中符合上面所说的情况， 随便选择一个 cmp 代码行， 双击查看信息

[![sR3BCV.png](https://s3.ax1x.com/2021/01/20/sR3BCV.png)](https://imgchr.com/i/sR3BCV)

我们可以发现他从 ecx 这个寄存器 + 230来访问的血量。复制ECX寄存器的值， 在CE中选择新的搜索， 并勾选 16进制， 对ECX的值进行搜索， 来查找一下究竟是谁存放了这个地址值

[![sR3cDJ.png](https://s3.ax1x.com/2021/01/20/sR3cDJ.png)](https://imgchr.com/i/sR3cDJ)

首先观察绿色的基地址， 发现他从 server.dll 模块来的， 因为我们开的是人机， 所以服务器在本地创建， 但是玩竞技游戏时，我们只有客户端， 所以忽略这个地址。 从上面血量的列表中继续依次查找

[![sR3vPP.md.png](https://s3.ax1x.com/2021/01/20/sR3vPP.md.png)](https://imgchr.com/i/sR3vPP)

可以观察到这个也符合之前的要求。同样的步骤， 对ECX寄存器的值进行查找

[![sR3z28.png](https://s3.ax1x.com/2021/01/20/sR3z28.png)](https://imgchr.com/i/sR3z28)

可以看到这次有4个绿色的基地址

> 在游戏中， 往往是通过全局变量来存储这些信息。例如虚幻4引擎通过链表， Unity3D则使用指针数组。 在Value自家游戏引擎中使用的是一组结构体数组，其中有一个指针指向了角色实体所在位置。

我们继续选择 Find out what access this address，一般数组访问可以观察
**mov 寄存器A，地址 + 寄存器B**
可以结合写代码中访问数组的代码来理解这段汇编代码

[![sR8yJP.png](https://s3.ax1x.com/2021/01/20/sR8yJP.png)](https://imgchr.com/i/sR8yJP)

很明显在我们选中的部分与访问数组的代码非常相似， 我们复制 Instruction Address
打开Memory View
跳转到 这条指令的地址 

[![sR8LyF.md.png](https://s3.ax1x.com/2021/01/20/sR8LyF.md.png)](https://imgchr.com/i/sR8LyF)

右键选择Set Breakpoint，设置断点并不断执行

[![sR8vw9.png](https://s3.ax1x.com/2021/01/20/sR8vw9.png)](https://imgchr.com/i/sR8vw9)

[![sR8xoR.png](https://s3.ax1x.com/2021/01/20/sR8xoR.png)](https://imgchr.com/i/sR8xoR)

我们可以发现 EAX 寄存器的值， 每次上升 10， 非常符合访问数组的代码的特征
通过对地址的+10 * n 进行访问， 再加上之前访问HP的偏移 100， 我们可以确定上一个地址就是存放游戏角色数组的地址

### 坐标

回到之前的HP的地址， 打开Memory View， 选择工具中的 Dissect structrues

[![sRGPSK.png](https://s3.ax1x.com/2021/01/20/sRGPSK.png)](https://imgchr.com/i/sRGPSK)

通过观察这组数据附近的，应该存在着3个连续的float类型数据， 一般为 X Y Z坐标， 通过移动角色可以观察出这三个数据即为坐标数据

[![sRGFyD.png](https://s3.ax1x.com/2021/01/20/sRGFyD.png)](https://imgchr.com/i/sRGFyD)

### 视角

同样观察附近， 应该有 两个 连续并随着转动视角变化的值。即是视角的三个值， 他们分别是 Pitch Yaw Roll 控制着相机的转动。下面的0则是Roll表示镜头不需要旋转

[![sRGl6S.png](https://s3.ax1x.com/2021/01/20/sRGl6S.png)](https://imgchr.com/i/sRGl6S)

### 视图矩阵

视图矩阵是用来将游戏中的 世界坐标 投影转化为 屏幕坐标 的一组 4x4 的矩阵
它随着视角的转动而变化， 当鼠标不动视图矩阵不动，狙击枪开镜与关镜则会改变矩阵头
具体相关内容不做介绍， 详细了解可以查看相关计算机图形学的内容

我们可以通过算法来计算矩阵， 来帮助我们找到游戏中的视图矩阵

http://andre-gaschler.com/rotationconverter/

[![sRGJTs.png](https://s3.ax1x.com/2021/01/20/sRGJTs.png)](https://imgchr.com/i/sRGJTs)

视图矩阵也是存放在基地址当中， 我们对几个绿色的基地址分别根据内存段进行右键查看浏览内存区域， 通过观察附近， 找到一组 4x4 在移动中只有最后一列会变化，其他不变化。 晃动视角则是4x4范围内变动的一组数据。

## WorldToScreen

世界坐标转屏幕坐标算法如下所示， 不同游戏大同小异，位置可以根据调整算法的参数来进行微调

~~~go
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
~~~

## 绘制

绘制则采用的是Windows提供GDI对游戏的窗口进行绘制， 使用什么绘制都可以进行绘图， 这个不进行详细说明