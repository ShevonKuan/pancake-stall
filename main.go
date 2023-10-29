package main

import (
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/ShevonKuan/pancake-stall/pancake"
	"github.com/spf13/cobra"
)

//go:embed "全唐诗.db"
var database embed.FS

func main() {
	rootCmd := &cobra.Command{
		Use:   "摊煎饼",
		Short: "这是一个摊煎饼小程序",
		Run: func(cmd *cobra.Command, args []string) {
			// Execute StallPancake function with provided arguments
			xLength, _ := cmd.Flags().GetInt("xlength")
			yLength, _ := cmd.Flags().GetInt("ylength")
			sqlLitePath, _ := cmd.Flags().GetString("sql")
			startW, _ := cmd.Flags().GetString("start")
			enableSplit, _ := cmd.Flags().GetBool("split")
			resultCounts, _ := cmd.Flags().GetInt("resultCounts")
			trialCounts, _ := cmd.Flags().GetInt("trialCounts")
			enableLog, _ := cmd.Flags().GetBool("enableLog")
			outputfile, _ := cmd.Flags().GetString("output")
			poetryCounts, _ := cmd.Flags().GetInt("poetryCounts")
			//check flags
			if startW == "" {
				log.Fatal("😟 请指定第一个汉字")
			}
			if sqlLitePath == "" { //未指定时使用嵌入数据库
				dbData, _ := database.ReadFile("全唐诗.db")
				// Write the embedded data to a temporary file
				tmpFile, err := os.CreateTemp("", "temp.db")
				if err != nil {
					log.Fatalf("failed to create temporary file %v", err)
				}
				defer os.Remove(tmpFile.Name())

				if _, err := tmpFile.WriteString(string(dbData)); err != nil {
					log.Fatalf("failed to write SQLite data to temporary file: %v", err)
				}

				// Use tmpFile.Name() as the SQLite database path
				sqlLitePath = tmpFile.Name()
			}
			o := pancake.StallPancake(
				xLength, yLength, sqlLitePath, startW, enableSplit, resultCounts, trialCounts, enableLog, poetryCounts,
			)
			err := os.WriteFile(outputfile, []byte(o.GenerateLatexCode()), 0644)
			if err != nil {
				fmt.Println("写入文件时发生错误:", err)
				return
			}
		},
	}

	// Define command-line flags
	rootCmd.Flags().Int("xlength", 20, "画布行数")
	rootCmd.Flags().Int("ylength", 20, "画布列数")
	rootCmd.Flags().Int("poetryCounts", 0, "注入诗句条数，<=0表示使用注入位置（选字）次数为终止条件")
	rootCmd.Flags().String("sql", "", "sqlite数据库文件路径，不指定时使用自带的全唐诗数据库")
	rootCmd.Flags().String("start", "", "用于开始的第一个字")
	rootCmd.Flags().Bool("split", true, "是否允许切断诗句")
	rootCmd.Flags().String("output", "output.tex", "tex文件输出路径")
	rootCmd.Flags().Int("resultCounts", 5000, "[高级选项] 数据库查询结果返回计数")
	rootCmd.Flags().Int("trialCounts", 10, "[高级选项] 注入位置（选字）次数，设定停止条件，当尝试多少个注入位置依然无法注入时停止随机选择注入位置")
	rootCmd.Flags().Bool("enableLog", true, "[高级选项] 在stdIO输出画布")

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
