package setting

import (
	"gateway/utils"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

var (
	AppMode        string
	HttpPort       string
	LogLevel       string
	LogToFile      bool
	LogFile        string
	LogFileMaxSize int
	LogFileBackup  int

	ReportNet string
	ExeName   string
)

func LoadServer(file *ini.File) {
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	HttpPort = file.Section("server").Key("HttpPort").MustString(":8080")

	LogLevel = file.Section("logger").Key("LogLevel").MustString("debug")
	LogToFile = file.Section("logger").Key("LogToFile").MustBool(false)
	LogFile = file.Section("logger").Key("LogFile").MustString("/log/gateway.log")
	LogFileMaxSize = file.Section("logger").Key("LogFileMaxSize").MustInt(5)
	LogFileBackup = file.Section("logger").Key("LogFileBackup").MustInt(3)

	ReportNet = file.Section("ops").Key("ReportNet").MustString("usb0")
	ExeName = file.Section("ops").Key("ExeName").MustString("armv7_gw_linux")
}

/**************获取配置信息************************/
func GetConf() {
	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	path := exeCurDir + "/config/config.ini"
	iniFile, err := ini.Load(path)
	if err != nil {
		log.Printf("读取config.ini失败 %v", err)

		cfg := ini.Empty()

		AppMode = "debug"
		HttpPort = ":8080"
		cfg.Section("server").Key("AppMode").SetValue("debug")
		cfg.Section("server").Key("HttpPort").SetValue(":8080")

		LogLevel = "debug"
		LogToFile = false
		LogFile = "./log/gateway.log"
		cfg.Section("logger").Key("LogLevel").SetValue(LogLevel)
		cfg.Section("logger").Key("LogToFile").MustBool(LogToFile)
		cfg.Section("logger").Key("LogFile").SetValue(LogFile)
		cfg.Section("logger").Key("LogFileMaxSize").MustInt(LogFileMaxSize)
		cfg.Section("logger").Key("LogFileBackup").MustInt(LogFileBackup)

		ReportNet = "usb0"
		ExeName = "armv7_gw_linux"
		cfg.Section("ops").Key("ReportNet").SetValue("usb0")
		cfg.Section("ops").Key("ExeName").SetValue("armv7_gw_linux")

		utils.DirIsExist("./config")
		err = cfg.SaveTo(path)
		if err != nil {
			log.Printf("写入config.ini失败 %v", err)
		}
		return
	}

	LoadServer(iniFile)
}
