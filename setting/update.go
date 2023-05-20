package setting

import (
	"log"
	"os"
	"path/filepath"
)

func Update(name string) bool {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	fileName := exeCurDir + "/config/" + name
	fileAbsoluteDir := exeCurDir + "/"
	ZAPS.Debugf("fileName %v\n", fileName)
	if err := UnZipFiles(fileName, fileAbsoluteDir); err != nil {
		log.Println(err)
		return false
	}
	err := os.Remove(fileName)
	if err != nil {
		ZAPS.Errorf("removeFile err,%s\n", fileName)
	}
	return true
}
