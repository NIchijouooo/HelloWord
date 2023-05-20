package setting

import (
	"encoding/json"
	"gateway/utils"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/beevik/ntp"
)

type NTPHostAddrTemplate struct {
	Enable     bool   `json:"enable"`
	TimeZone   string `json:"timeZone"`
	URLMaster  string `json:"urlMaster"`
	PortMaster int    `json:"portMaster"`
	URLSlave   string `json:"urlSlave"`
	PortSlave  int    `json:"portSlave"`
}

var NTPHostAddr = NTPHostAddrTemplate{
	Enable:     false,
	TimeZone:   "UTC+8",
	URLMaster:  "",
	PortMaster: 123,
	URLSlave:   "",
	PortSlave:  123,
}

func NTPInit() {
	ReadNTPHostAddrFromJson()
}

func NTPGetTime() error {
	if NTPHostAddr.Enable == true {
		//首先使用主服务器
		if NTPHostAddr.URLMaster != "" {
			options := ntp.QueryOptions{Timeout: 5 * time.Second, Port: NTPHostAddr.PortMaster}
			response, err := ntp.QueryWithOptions(NTPHostAddr.URLMaster, options)
			if err != nil {
				ZAPS.Errorf("ntp获取失败 %v", err)
				return err
			} else {
				ZAPS.Debugf("ntp获取成功 %v", response.Time.String())
				//ZAPS.Debugf("response %v,Now %v", response.Time.Unix(), time.Now().Unix())
				//注意：即使ntp返回是utc的时间，系统会自动将utc时间换算成cst的时间
				if math.Abs(float64(response.Time.Unix()-time.Now().Unix())) > 60 {
					//unix := response.Time.Unix() + 8*3600
					//SystemSetRTC(time.Unix(unix, 0).Format("2006-01-02 15:04:05"))
					SystemSetRTC(time.Unix(response.Time.Unix(), 0).Format("2006-01-02 15:04:05"))
				}
				return nil
			}
		} else if NTPHostAddr.URLSlave != "" {
			options := ntp.QueryOptions{Timeout: 5 * time.Second, Port: NTPHostAddr.PortSlave}
			response, err := ntp.QueryWithOptions(NTPHostAddr.URLSlave, options)
			if err != nil {
				return err
			} else {
				if math.Abs(float64(response.Time.Unix()-time.Now().Unix())) > 60 {
					unix := response.Time.Unix() + 8*3600
					SystemSetRTC(time.Unix(unix, 0).Format("2006-01-02 15:04:05"))
				}
				return nil
			}
		}
	}
	return nil
}

func GetNTP() {
	_ = NTPGetTime()
}

func ReadNTPHostAddrFromJson() bool {

	fileDir := "./selfpara/ntpHostAddr.json"

	utils.DirIsExist("./selfpara")

	fp, err := os.OpenFile(fileDir, os.O_RDONLY, 0777)
	if err != nil {
		ZAPS.Debugf("打开NTP服务配置json文件失败 %v", err)
		return false
	}
	defer fp.Close()

	data := make([]byte, 20480)
	dataCnt, err := fp.Read(data)

	err = json.Unmarshal(data[:dataCnt], &NTPHostAddr)
	if err != nil {
		ZAPS.Info("格式化NTP服务配置json文件失败 %v", err)
		return false
	}
	return true
}

func WriteNTPHostAddrToJson() {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	fileDir := exeCurDir + "/selfpara/ntpHostAddr.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println("open ntpHostAddr.json err", err)
		return
	}
	defer fp.Close()

	sJson, _ := json.Marshal(NTPHostAddr)

	_, err = fp.Write(sJson)
	if err != nil {
		ZAPS.Errorf("写入NTP服务配置参数错误 %v", err)
		return
	}
	ZAPS.Infof("写入NTP服务配置参数成功")
}
