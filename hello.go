// https://go.dev/learn/

package main

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
)

func showDivider() {
	fmt.Println("==================================================================")
}

func main() {
	// https://tools.kalvinbg.cn/txt/ascii
	const asciiTitle string = "IF9fX18gICAgX18gICAgICAgICAgICAgICAgICAgIF9fX18gICAgX19fICAgICAgICAgICAgICAgICAgICAgICAKL1wgIF9gXCAvXCBcX18gICAgICAgICAgICAgICAgL1wgIF9gXCAvXF8gXCAgICAgICAgICAgICAgICAgICAgICAKXCBcLFxMXF9cIFwgLF9cICAgIF9fICAgICBfIF9fXCBcIFxMXCBcLy9cIFwgICAgIF9fXyAgICAgIF9fICAgICAKIFwvX1xfXyBcXCBcIFwvICAvJ19fYFwgIC9cYCdfX1wgXCAgXyA8J1wgXCBcICAgLyBfX2BcICAvJ18gYFwgICAKICAgL1wgXExcIFwgXCBcXy9cIFxMXC5cX1wgXCBcLyBcIFwgXExcIFxcX1wgXF8vXCBcTFwgXC9cIFxMXCBcICAKICAgXCBgXF9fX19cIFxfX1wgXF9fLy5cX1xcIFxfXCAgXCBcX19fXy8vXF9fX19cIFxfX19fL1wgXF9fX18gXCAKICAgIFwvX19fX18vXC9fXy9cL19fL1wvXy8gXC9fLyAgIFwvX19fLyBcL19fX18vXC9fX18vICBcL19fX0xcIFwKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIC9cX19fXy8KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIFxfL19fLyA"
	var titleDec, _ = base64.StdEncoding.DecodeString(asciiTitle)

	showDivider()
	fmt.Printf("%s\n", titleDec)
	showDivider()
	fmt.Println("Hive 数据导出工具")

	var (
		sql       string
		isExecute bool = false
	)

	for !isExecute {
		fmt.Println("请输入要执行的SQL语句，之后按回车确认。（输入 q 退出）")
		var sqlLength, err = fmt.Scanln(&sql)

		if err != nil {
			fmt.Printf("报错啦：%v\n", err)
			continue
		}

		if sqlLength == 0 {
			fmt.Println("请输入正确的SQL语句！")
			continue
		}

		if sql == "q" {
			os.Exit(0)
		}

		isExecute = true
	}

	dir, err := os.MkdirTemp("", "hive-out-")
	if err != nil {
		fmt.Printf("创建临时目录错误！%v\n", err)
		log.Fatal(err)
	}

	f, err := os.CreateTemp(dir, "hive.*.sql")
	if err != nil {
		fmt.Printf("创建临时文件错误！错误：%v\n", err)
		log.Fatal(err)
	}

	defer f.Close()
	defer os.Remove(f.Name())
	defer os.RemoveAll(dir)

	data := []byte(sql)
	if _, err := f.Write(data); err != nil {
		fmt.Println("写入文件失败！", err.Error())
		log.Fatal(err)
	}

	fmt.Println("创建临时文件", f.Name())

	// https://juejin.cn/post/7000925379145760782
	// https://github.com/google/uuid
	u1, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("创建UUID失败", err.Error())
		log.Fatal(err)
	}
	fmt.Println(u1.String())

	// https://www.cnblogs.com/wongbingming/p/13984538.html
	//fmt.Println("执行SQL！")
	//cmd := exec.Command("echo", f.Name())
	//out, err := cmd.CombinedOutput()
	//if err != nil {
	//	fmt.Printf("combined out:\n%s\n", string(out))
	//	log.Fatalf("cmd.Run() failed with %s\n", err)
	//}
	//fmt.Printf("combined out:\n%s\n", string(out))

}
