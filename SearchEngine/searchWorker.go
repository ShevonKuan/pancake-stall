package searchengine

import (
	"database/sql"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	_ "github.com/mattn/go-sqlite3"

	"github.com/tidwall/gjson"
)

var (
	re = regexp.MustCompile(`[\pP\pS]+`)
)

func GenerateRandomNumber() float32 {
	offset := rand.Float32()
	return offset
}

func SearchHanzi(k rune, enableSplit bool, resultCounts int, db *sql.DB) SearchResults {
	start := time.Now()
	var result []*SearchResult
	count := 0
	rows, err := db.Query("SELECT * FROM poetry_info WHERE paragraphs LIKE '%' || ? || '%'", string(k))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var title, author, paragraphs string
		var index int
		var para string
		err := rows.Scan(&id, &title, &author, &paragraphs)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("ID: %d, Title: %s, Author: %s, Paragraphs: %s\n", id, title, author, paragraphs)
		paras := gjson.Parse(paragraphs).Array()
		for i := range paras { // i 为单句诗
			index = strings.IndexRune(paras[i].Str, k)
			if index != (-1) { // 搜索到则跳出
				para = paras[i].Str
				break
			}
		}
		if enableSplit {
			parts := strings.Split(para, "，") // 使用逗号作为分隔符

			for _, part := range parts {
				index = strings.IndexRune(part, k)
				if index != (-1) { // 搜索到则跳出
					para = part
					break
				}
			}
		}

		// 使用正则表达式替换中文文本中的标点符号为空字符串
		para = re.ReplaceAllString(para, "")

		r := SearchResult{
			Paragraph: para,
			Title:     title,
			Author:    author,
			PoetryID:  id,
			Index:     utf8.RuneCountInString(para[:strings.Index(para, string(k))]),
			Length:    utf8.RuneCountInString(para),
		}
		result = append(result, &r)
		// log.Printf("%+v\n", r)
		count += 1
		if count == resultCounts {
			break
		}
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return SearchResults{result, time.Since(start)}
}
