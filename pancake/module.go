package pancake

import (
	"errors"

	searchengine "github.com/ShevonKuan/pancake-stall/SearchEngine"
)

// 画布
type Canvas [][]rune      // 使用时注意传递指针
type CanvasColors [][]int // 画布上色使用时注意传递指针
// 坐标
type Coordinate struct {
	X int
	Y int
}

// 新建画布

// 颜色容器
var colorsContainer = make(map[*Canvas]*CanvasColors)
var statusContainer = make(map[*Canvas]int)

func NewCanvas(x int, y int) *Canvas {
	canvas := make(Canvas, x)
	canvasColors := make(CanvasColors, x)
	for i := range canvas {
		canvas[i] = make([]rune, y)
		canvasColors[i] = make([]int, y)
	}

	colorsContainer[&canvas] = &canvasColors
	statusContainer[&canvas] = 0
	return &canvas
}

func (c *Canvas) GetCanvasColors() *CanvasColors {
	return colorsContainer[c]
}

// 诗句注入信息
type PoetryInput map[Coordinate]rune

// 画布中心点
func (c *Canvas) GetCenterCoordinate() Coordinate {
	xLength, yLength := len(*c), len((*c)[0])
	return Coordinate{xLength / 2, yLength / 2}
}

// 诗句填入与可行性检测（横向）
func (c *Canvas) CheckHorizonAvailable(
	sharedHanzi Coordinate, //共用字坐标
	readyInput *searchengine.SearchResult, // 待写入诗句
) (*PoetryInput, error) {
	// 横向占位
	horizonOccupation := PoetryInput{}
	for i := 0; i < readyInput.Length; i++ {
		x := sharedHanzi.X
		y := sharedHanzi.Y - readyInput.Index + i
		w := []rune((*readyInput).Paragraph)[i] // 待注入字符
		w1, _ := searchengine.SC2TC.Convert(string(w))
		w2, _ := searchengine.TC2SC.Convert(string(w))
		if (y < 0) || (y >= len((*c)[0])) { // 撞击边缘检测
			return nil, errors.New("不支持的注入方式：注入位置超出画布范围")
		}
		if ((*c)[x][y] != '\x00') && (((*c)[x][y] != []rune(w1)[0]) && ((*c)[x][y] != []rune(w2)[0])) { // 字不一样并且占用了已有空位跳出
			return nil, errors.New("不支持的注入方式：注入位置已有不同注入")
		}

		horizonOccupation[Coordinate{x, y}] = w // 创建注入信息
	}
	for i := range horizonOccupation { // 如果每个字都一样就不要填了
		if (*c)[i.X][i.Y] != horizonOccupation[i] {
			return &horizonOccupation, nil
		}
	}
	return nil, errors.New("已存在该注入")
}

// 诗句填入与可行性检测（纵向）
func (c *Canvas) CheckVerticalAvailable(
	sharedHanzi Coordinate, //共用字坐标
	readyInput *searchengine.SearchResult, // 待写入诗句
) (*PoetryInput, error) {

	// 纵向占位
	verticalOccupation := PoetryInput{}
	for i := 0; i < readyInput.Length; i++ {
		x := sharedHanzi.X - readyInput.Index + i
		y := sharedHanzi.Y
		w := []rune((*readyInput).Paragraph)[i] // 待注入字符
		w1, _ := searchengine.SC2TC.Convert(string(w))
		w2, _ := searchengine.TC2SC.Convert(string(w))

		if (x < 0) || (x >= len(*c)) { // 撞击边缘检测
			return nil, errors.New("不支持的注入方式：注入位置超出画布范围")
		}
		if ((*c)[x][y] != '\x00') && (((*c)[x][y] != []rune(w1)[0]) && ((*c)[x][y] != []rune(w2)[0])) { // 字不一样并且占用了已有空位跳出
			return nil, errors.New("不支持的注入方式：注入位置已有不同注入")
		}

		verticalOccupation[Coordinate{x, y}] = w // 创建注入信息
	}
	for i := range verticalOccupation { // 如果每个字都一样就不要填了
		if (*c)[i.X][i.Y] != verticalOccupation[i] {
			return &verticalOccupation, nil
		}
	}
	return nil, errors.New("已存在该注入")
}

// 诗句填入
func (c *Canvas) Input(I PoetryInput) {
	cc := c.GetCanvasColors()
	statusContainer[c] += 1
	for i := range I {
		(*c)[i.X][i.Y] = I[i]
		(*cc)[i.X][i.Y] = statusContainer[c]
	}
	if statusContainer[c] == 3 {
		statusContainer[c] = 0
	}
}

// 扩展画布
func (c *Canvas) EnLarge(x int, y int) (*Canvas, error) {
	c = c.Tight()
	xLength, yLength := len(*c), len((*c)[0])
	if x == xLength && y == yLength {
		return c, nil
	}
	if x < xLength || y < yLength {
		return nil, errors.New("新画布过小无法存放原有画布信息")
	}

	dw := x - xLength // 行数的变化量
	dh := y - yLength // 列数的变化量

	offsetX := dw / 2 // 行数的偏移量
	offsetY := dh / 2 // 列数的偏移量

	c2 := NewCanvas(x, y)
	newPoetryInput := *(*c).ExportToPoetryInput()
	for c := range newPoetryInput {
		x := c.X + offsetX
		y := c.Y + offsetY
		(*c2)[x][y] = newPoetryInput[c]
	}
	return c2, nil
}

// 画布导出为 *PoetryInput类型
func (c *Canvas) ExportToPoetryInput() *PoetryInput {
	xLength, yLength := len(*c), len((*c)[0])
	output := PoetryInput{}
	for i := 0; i < xLength; i++ {
		for j := 0; j < yLength; j++ {
			if (*c)[i][j] != '\x00' {
				output[Coordinate{i, j}] = (*c)[i][j]
			}
		}
	}
	return &output
}

// 修剪画布切除白边
func (c *Canvas) Tight() *Canvas {
	dataCoordinates := c.ExportToPoetryInput()
	xLength, yLength := len(*c), len((*c)[0])
	var xMin, xMax, yMin, yMax = xLength, 0, yLength, 0

	for c := range *dataCoordinates {
		if xMin > c.X {
			xMin = c.X
		}
		if xMax < c.X {
			xMax = c.X
		}
		if yMin > c.Y {
			yMin = c.Y
		}
		if yMax < c.Y {
			yMax = c.Y
		}
	}
	xLength = xMax - xMin + 1
	yLength = yMax - yMin + 1
	output := NewCanvas(xLength, yLength)
	for c := range *dataCoordinates {
		x := c.X - xLength
		y := c.Y - yLength
		(*output)[x][y] = (*dataCoordinates)[c]

	}
	return output
}
