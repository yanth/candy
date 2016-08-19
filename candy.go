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
		take(make([]string, 0))

	case "takes":
		takes()

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

	confDir := path.Join(home, ".candy.d")
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
	fp, err := os.OpenFile(historyPath, os.O_CREATE|os.O_RDONLY, 0644)
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

func take(paths []string) {
	// 受け取った .candy パスの配列を対象に .candy を読みだす。
	// 何も引数に無ければカレントディレクトリcandyを取る
	if len(paths) <= 0 {
		paths = append(paths, path.Join(getCurrentPath(), ".candy"))
	}

	for _, p := range paths {
		fp, err := os.OpenFile(p, os.O_RDONLY, 0644)
		if err != nil {
			panic(err)
			os.Exit(1)
		}

		var line string
		scan := bufio.NewScanner(fp)

		fmt.Printf("=== %s\n", p)
		for scan.Scan() {
			if err = scan.Err(); err != nil {
				continue
			}
			line = scan.Text()
			fmt.Printf("  - %s\n", line)
		}

		fmt.Println()
	}
}

func takes() {
	// .candy\droped から取得する
	p := getHistoryPath()
	fp, err := os.OpenFile(p, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	defer fp.Close()

	lines := make([]string, 0)
	scan := bufio.NewScanner(fp)
	for scan.Scan() {
		lines = append(lines, scan.Text())

		if err = scan.Err(); err != nil {
			panic(err)
		}
	}

	take(lines)
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

	return path.Join(home, ".candy.d", "droped")
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
