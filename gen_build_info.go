package main

import (
	"fmt"
	"gateway/utils"
	"os"
	"time"
)

func main() {

	// QJHui UPDATE 2023/6/7 修改为获取系统当前时间(本地时间)
	//buildTime := time.Now().UTC().Format(time.RFC3339)
	buildTime := time.Now().Format("2006-01-02 15:04:05")
	utils.DirIsExist("./buildInfo")
	f, err := os.Create("./buildInfo/build_info.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "package buildInfo\n\nconst BuildTime = %q\n", buildTime)
}
