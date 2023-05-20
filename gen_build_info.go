package main

import (
	"fmt"
	"gateway/utils"
	"os"
	"time"
)

func main() {
	buildTime := time.Now().UTC().Format(time.RFC3339)
	utils.DirIsExist("./buildInfo")
	f, err := os.Create("./buildInfo/build_info.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "package buildInfo\n\nconst BuildTime = %q\n", buildTime)
}
