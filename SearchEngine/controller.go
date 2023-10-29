package searchengine

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/longbridgeapp/opencc"
	_ "github.com/mattn/go-sqlite3"
)

var (
	DatabaseAvailableStatus = make(map[string]bool)
	DatabasePath            = make(map[string][]string)
	SC2TC, _                = opencc.New("s2t")
	TC2SC, _                = opencc.New("t2s")
)

type SearchResult struct {
	Paragraph string // 搜索到符合条件的诗句
	Title     string // 题目
	Author    string //作者
	PoetryID  int    //诗歌id
	Index     int    // 字出现的位置
	Length    int    // 长度
}

type SearchResults struct {
	Results []*SearchResult
	T       time.Duration
}

type SRs []*SearchResults

// 私有查询接口，将单字转换为简繁体两种形式进行搜索并concatenate搜索结果
func searchChannel(
	k rune,
	enableSplit bool,
	resultCounts int,
	db *sql.DB,
	ch chan<- *SearchResults) {
	// 搜索结果通过channel返回
	r := SearchHanzi(k, enableSplit, resultCounts, db)
	ch <- &r
}

func (sr SRs) SearchResultsConcatenate() *SearchResults {
	var results []*SearchResult
	var t time.Duration
	for _, s := range sr {
		results = append(results, s.Results...)
		if s.T > t {
			t = s.T
		}
	}
	return &SearchResults{Results: results,
		T: t,
	}
}

func SearchInterface(k rune, enableSplit bool, resultCounts int, db *sql.DB) *SearchResults {
	// 简繁体转换

	sk, err := SC2TC.Convert(string(k))
	if err != nil {
		fmt.Println(err)
	}
	tk, err := TC2SC.Convert(string(k))
	if err != nil {
		fmt.Println(err)
	}
	// 并发查询
	ch := make(chan *SearchResults)
	for _, w := range []rune{[]rune(sk)[0], []rune(tk)[0]} {
		go searchChannel(w, enableSplit, resultCounts, db, ch)
	}
	r1 := <-ch
	r2 := <-ch
	return SRs{r1, r2}.SearchResultsConcatenate()
}

// 打印结果
func (o *SearchResults) Print() {
	spew.Dump(o)
}
