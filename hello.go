// https://go.dev/learn/

package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/exec"
	"time"
)

func showDivider() {
	fmt.Println("==================================================================")
}

func main() {
	// https://tools.kalvinbg.cn/txt/ascii
	const asciiTitle string = "IF9fX18gICAgX18gICAgICAgICAgICAgICAgICAgIF9fX18gICAgX19fICAgICAgICAgICAgICAgICAgICAgICAKL1wgIF9gXCAvXCBcX18gICAgICAgICAgICAgICAgL1wgIF9gXCAvXF8gXCAgICAgICAgICAgICAgICAgICAgICAKXCBcLFxMXF9cIFwgLF9cICAgIF9fICAgICBfIF9fXCBcIFxMXCBcLy9cIFwgICAgIF9fXyAgICAgIF9fICAgICAKIFwvX1xfXyBcXCBcIFwvICAvJ19fYFwgIC9cYCdfX1wgXCAgXyA8J1wgXCBcICAgLyBfX2BcICAvJ18gYFwgICAKICAgL1wgXExcIFwgXCBcXy9cIFxMXC5cX1wgXCBcLyBcIFwgXExcIFxcX1wgXF8vXCBcTFwgXC9cIFxMXCBcICAKICAgXCBgXF9fX19cIFxfX1wgXF9fLy5cX1xcIFxfXCAgXCBcX19fXy8vXF9fX19cIFxfX19fL1wgXF9fX18gXCAKICAgIFwvX19fX18vXC9fXy9cL19fL1wvXy8gXC9fLyAgIFwvX19fLyBcL19fX18vXC9fX18vICBcL19fX0xcIFwKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIC9cX19fXy8KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIFxfL19fLyA"
	var titleDec, _ = base64.StdEncoding.DecodeString(asciiTitle)

	fmt.Printf("%s\n", titleDec)
	showDivider()
	fmt.Println("Hive 数据导出工具")
	fmt.Println("https://github.com/Deali-Axy/hive-output-tool")
	showDivider()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err.Error())
	}

	var (
		sql       string
		isExecute bool = false
	)

	// https://books.studygolang.com/The-Golang-Standard-Library-by-Example/chapter01/01.4.html
	reader := bufio.NewReader(os.Stdin)
	for !isExecute {
		fmt.Println("请输入要执行的SQL语句，之后按回车确认。（输入 q 退出）")
		line, err := reader.ReadBytes('\n')

		if err != nil {
			fmt.Printf("报错啦：%v\n", err)
			continue
		}

		if len(line) == 0 {
			fmt.Println("请输入正确的SQL语句！")
			continue
		}

		line = bytes.TrimRight(line, "\r\n")
		sql = string(line)

		if sql == "q" {
			os.Exit(0)
		}

		isExecute = true
	}

	fmt.Println()

	dir, err := os.MkdirTemp("", "hive-out-")
	if err != nil {
		fmt.Printf("创建临时目录错误！%v\n", err)
		log.Fatal(err)
	}

	tempFile, err := os.CreateTemp(dir, "hive.*.sql")
	if err != nil {
		fmt.Printf("创建临时文件错误！错误：%v\n", err)
		log.Fatal(err)
	}

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())
	defer os.RemoveAll(dir)

	data := []byte(sql)
	if _, err := tempFile.Write(data); err != nil {
		fmt.Println("写入文件失败！", err.Error())
		log.Fatal(err)
	}

	log.Printf("创建临时文件 %s\n\n", tempFile.Name())

	// https://juejin.cn/post/7000925379145760782
	// https://github.com/google/uuid
	u1, err := uuid.NewUUID()
	if err != nil {
		log.Printf("创建UUID失败: %v\n", err)
		log.Fatal(err)
	}
	log.Printf("Session ID: %v\n\n", u1)

	// https://www.cnblogs.com/mayanan/p/15342214.html
	t := time.Now()
	timeStr := fmt.Sprintf("%04d-%02d-%02d_%02d-%02d-%02d-%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
	outputPath := fmt.Sprintf("%s.csv", timeStr)

	// https://www.cnblogs.com/wongbingming/p/13984538.html
	log.Printf("正在执行SQL，请稍等...\n\n")
	cmd := exec.Command("beeline",
		"-u", os.Getenv("CONN_STR"),
		"-n", os.Getenv("USERNAME"),
		"-p", os.Getenv("PASSWORD"),
		"--incremental=true",
		"--verbose=true",
		"--fastConnect=true",
		"--silent=true",
		"--outputformat="+os.Getenv("OUTPUT_FORMAT"),
		"-f", tempFile.Name(),
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("combined out: %s\n", string(out))
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("创建文件失败: %v\n", err)
	}
	if _, err := outputFile.Write(out); err != nil {
		log.Printf("写入文件失败: %v\n", err)
		log.Fatal(err)
	}

	log.Printf("导出数据到：%s\n\n", outputPath)
	log.Println("搞定！")
}
