package main

import (
	"bufio"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
	"strings"
	"time"
)

func main() {
	argv := len(os.Args)
	if 2 > argv {
		fmt.Printf("%s\n", usage())
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "drop":
		if len(os.Args) < 3 {
			usage()
		} else {
			drop()
		}

	case "take":
		take()

	default:
		fmt.Printf("%s\n", usage())
	}
}

func usage() string {
	return "使い方出ます"
}

func init() {
	// ディレクトリ確認
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	confDir := path.Join(home, ".candy")
	os.Mkdir(confDir, 0777)
}

func drop() {
	historyWrite()
	candyWrite()
}

func historyWrite() {
	history := getHistoryPath()
	current := getCandyPath()

	if !containHistory(history, current) {
		fileWriter(history, current)
	}
}

func containHistory(historyPath string, path string) bool {
	fp, err := os.OpenFile(historyPath, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	var line string
	scan := bufio.NewScanner(fp)
	for scan.Scan() {
		line = scan.Text()
		if err = scan.Err(); err != nil {
			panic(err)
		}

		if string(line) == path {
			return true
		}
	}
	return false
}

func candyWrite() {
	str := nowtime() + " " + strings.Join(os.Args[2:], " ")
	path := getCandyPath()

	fileWriter(path, str)
}

func take() {
	// カレントディレクトリcandyを取る
	// todo: 引数にpathの文字列配列を入れて、空配列なら
	//       カレントディレクトリの.candyを読みだす。
	//       普段は配列を読み取る。無ければスルー
	filepath := path.Join(getCurrentPath(), ".candy")
	fp, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Errorf("%s\n", err)
	}

	fmt.Printf("=== %s\n", filepath)

	var line string
	scan := bufio.NewScanner(fp)
	for scan.Scan() {
		line = scan.Text()
		if err = scan.Err(); err != nil {
			panic(err)
		}

		fmt.Printf("  - %s\n", line)
	}
}

func takes() {
	// .candy\droped から取得する
}

func nowtime() string {
	const timeformat string = "2006-01-02 03:04:05"

	return time.Now().Format(timeformat)
}

func getHistoryPath() string {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	return path.Join(home, ".candy", "droped")
}

func getCandyPath() string {
	return path.Join(getCurrentPath(), ".candy")
}

func getCurrentPath() string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("%s\n", err)
		os.Exit(1)
	}

	return path
}

func fileWriter(filepath string, text string) {
	var fp *os.File
	var err error

	fp, err = os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	fp.WriteString(text + "\n")
}
