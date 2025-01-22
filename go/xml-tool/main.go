package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// 定义游戏数据结构
type Game struct {
	Path        string  `xml:"path"`
	Name        string  `xml:"name"`
	Desc        string  `xml:"desc"`
	Rating      float32 `xml:"rating"`
	ReleaseDate string  `xml:"releasedate"`
	Developer   string  `xml:"developer"`
	Publisher   string  `xml:"publisher"`
	Genre       string  `xml:"genre"`
	Players     string  `xml:"players"`
}

// 定义游戏列表数据结构
type GameList struct {
	XMLName xml.Name `xml:"gameList"`
	Games   []Game   `xml:"game"`
}

func main() {
	if len(os.Args) < 2 {
        fmt.Println("Usage: xml-tool <input.xml>")
        return
    }

    inputFile := os.Args[1]

	// 读取 XML 文件内容
	xmlContent, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading XML file:", err)
		return
	}

	// 解析 XML 内容
	var gameList GameList
	err = xml.Unmarshal(xmlContent, &gameList)
	if err != nil {
		fmt.Println("Error unmarshalling XML:", err)
		return
	}

	// 处理每个游戏
	for i := range gameList.Games {
		// 修改 path 和 name 字段
		gameList.Games[i].Name = modifyName(gameList.Games[i].Path, gameList.Games[i].Name)
	}

	// 将修改后的数据重新编码为 XML
	output, err := xml.MarshalIndent(gameList, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling XML:", err)
		return
	}

	// 为输出结果创建一个新的 XML 文件
	outputFile := "output.xml"
	err = os.WriteFile(outputFile, output, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	// 提示用户文件已保存
	fmt.Printf("Modified XML has been saved to %s\n", outputFile)
}

// 修改 name 字段，通过从 path 提取中文和 [中] 来替换，并替换 name 中的 '.' 为 '-'
func modifyName(path, originalName string) string {
	// 正则表达式：提取 path 中 [] 中的内容，包括中文和 [中]
	re := regexp.MustCompile(`\[(.*?)\]`)
	matches := re.FindAllStringSubmatch(path, -1)

	// 处理中文和 [中] 部分
	var chinesePart string
	var langPart string
	if len(matches) > 0 {
		// 如果第一个 [] 中的部分是中文内容，提取并更新 name
		if !strings.Contains(matches[0][1], "中") {
			chinesePart = matches[0][1] // 提取中文部分，去掉 [ ]
			originalName = chinesePart
		}
		// 如果有 [中] 部分，追加到 name 后
		for _, match := range matches {
			if match[1] == "中" {
				langPart = match[0] // 提取 [中] 部分
				originalName += langPart
			}
		}
	}

	// 将 name 中的 '.' 替换成 '-'
	originalName = strings.ReplaceAll(originalName, ".", "-")

	return originalName
}
