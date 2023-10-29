package pancake

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	searchengine "github.com/ShevonKuan/pancake-stall/SearchEngine"
	_ "github.com/mattn/go-sqlite3"
)

func StallPancake(
	xLength int, // 行数
	yLength int, // 列数
	sqlLitePath, // sqlite路径
	startW string, //以什么字开始
	enableSplit bool, // 是否允许断句
	resultCounts int, // 搜索结果目标值越多随机性越大
	trialCounts int, // 注入位置（选字）次数，设定停止条件，当尝试多少个注入位置依然无法注入时停止随机选择注入位置
	enableLog bool, // 开启日志
	poetryCounts int, // 注入诗句条数 小于等于0 则表示填入至满只用trialCounts终止条件
) *Canvas {
	// 启动数据库
	db, _ := sql.Open("sqlite3", "poetry.db")
	canvas := NewCanvas(xLength, yLength)
	// 写入第一个字
	canvas.Input(PoetryInput{canvas.GetCenterCoordinate(): []rune(startW)[0]})
	log.Println("开始摊煎饼")
	// 开始摊煎饼
	trial := 0 // 寻找注入位置次数
	counter := 0
	for {
		poetryInputNow := canvas.ExportToPoetryInput()
		// 从已填中选取一个字
		// 注入多次失败后处理
		var start Coordinate
		if trial == trialCounts {
			break
		} else {
			for { // 随机抽取一个字，计算其周围含字的密度，密度越大这个抽取出来的字被丢弃的概率越大
				start = poetryInputNow.getRandomKey()
				xLeftTop := start.X - 2
				xRightBottom := start.X + 2
				yLeftTop := start.Y - 2
				yRightBottom := start.Y + 2
				if xLeftTop < 0 {
					xLeftTop = 0
				}
				if xRightBottom > (xLength - 1) {
					xRightBottom = xLength - 1
				}
				if yLeftTop < 0 {
					yLeftTop = 0
				}
				if yRightBottom > (yLength - 1) {
					yRightBottom = yLength - 1
				}
				var runeCounts int
				for i := xLeftTop; i <= xRightBottom; i++ {
					for j := yLeftTop; j <= yRightBottom; j++ {
						if (*canvas)[i][j] != '\x00' { // 有字
							runeCounts += 1
						}
					}
				}
				s := (xRightBottom - xLeftTop + 1) * (yRightBottom - yLeftTop + 1)
				log.Println(s, xLeftTop, xRightBottom, yLeftTop, yRightBottom)
				if rand.Intn(s) >= runeCounts {
					break
				}
			}
		}
		//log.Println(poetryInputNow)
		log.Println("搜索位于 ", start, " 的\"", string((*poetryInputNow)[start]), "\"")
		// 搜索诗句
		result := searchengine.SearchInterface((*poetryInputNow)[start], enableSplit, resultCounts, db)
		// 获得所有可用注入
		var inject []*PoetryInput
		for i := range result.Results {

			h, err := canvas.CheckHorizonAvailable(start, result.Results[i])
			if err == nil {
				inject = append(inject, h)
			}
			v, err := canvas.CheckVerticalAvailable(start, result.Results[i])
			if err == nil {
				inject = append(inject, v)
			}
		}
		log.Println("共有", len(result.Results), "个搜索结果", len(inject), "个注入可能性")
		// 判断注入是否非空并随机选出注入
		if len(inject) != 0 {
			injectIndex := rand.Intn(len(inject))
			canvas.Input(*inject[injectIndex])
			log.Println("注入 ", *inject[injectIndex])
			trial = 0
			counter += 1
		} else {
			trial += 1
			log.Println("无法注入")
		}
		if enableLog {
			fmt.Println("")
			canvas.StdOut()
		}
		if poetryCounts > 0 && counter == poetryCounts {
			break
		}
	}
	// 引出 poetryInputNow和对应的canvas

	return canvas
}

// 随机选一个出来
func (m PoetryInput) getRandomKey() Coordinate {
	keys := make([]Coordinate, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}

	randomIndex := rand.Intn(len(keys))
	return keys[randomIndex]
}
