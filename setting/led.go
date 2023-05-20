package setting

import (
	"gateway/utils"
	"github.com/goccy/go-json"
	"os"
	"time"
)

const (
	LEDON  int = 1
	LEDOFF int = 0
)

const (
	LEDNet    int = 0
	LEDServer int = 1
)

type LedGPIOTemplate struct {
	GPIO       string `json:"gpio"`
	SetValue   string `json:"setValue"`
	ResetValue string `json:"resetValue"`
}

type LedModuleTemplate struct {
	FlashTick  *time.Ticker    `json:"-"`
	SystemLed  LedGPIOTemplate `json:"systemLed,omitempty"`  //系统指示灯
	ServiceLed LedGPIOTemplate `json:"serviceLed,omitempty"` //服务指示灯
}

var LedFlashFlag int = 0
var LedModule = &LedModuleTemplate{}

func LedInit() {
	_ = ReadLedModuleParamFromJson()
}

func ReadLedModuleParamFromJson() bool {

	data, err := utils.FileRead("./selfpara/led.json")
	if err != nil {
		ZAPS.Debugf("打开LED配置json文件失败 %v", err)
		return false
	}
	err = json.Unmarshal(data, &LedModule)
	if err != nil {
		ZAPS.Errorf("LED配置json文件格式化失败 %v", err)
		return false
	}
	ZAPS.Infof("打开LED配置json文件成功")

	return true
}

func WriteLEDModuleParamToJson() {

	utils.DirIsExist("./selfpara")

	sJson, _ := json.Marshal(LedModule)
	err := utils.FileWrite("./selfpara/led.json", sJson)
	if err != nil {
		ZAPS.Errorf("写入LED配置json文件 %s %v", "失败", err)
		return
	}
	ZAPS.Infof("写入LEDjson文件 %s", "成功")
}

func LedOnOff(index int, cmd int) {
	ioStr := ""
	if index == LEDNet {
		ioStr = LedModule.SystemLed.GPIO
	} else if index == LEDServer {
		ioStr = LedModule.ServiceLed.GPIO
	}

	fd, err := os.OpenFile(ioStr, os.O_RDWR, 0777)
	if err != nil {
		//ZAPS.Errorf("LED接口[%s]打开失败 %v", ioStr, err)
		return
	}
	defer fd.Close()

	if cmd == LEDON {
		if index == LEDNet {
			_, err = fd.Write([]byte(LedModule.SystemLed.SetValue))
		} else if index == LEDServer {
			_, err = fd.Write([]byte(LedModule.ServiceLed.SetValue))
		}

	} else {
		if index == LEDNet {
			_, err = fd.Write([]byte(LedModule.SystemLed.ResetValue))
		} else if index == LEDServer {
			_, err = fd.Write([]byte(LedModule.ServiceLed.ResetValue))
		}
	}

	if err != nil {
		//ZAPS.Errorf("LED接口[%s]写入失败 %v", ioStr, err)
		return
	}
}

func LedFlash(index int) {
	ioName := ""
	if index == LEDNet {
		ioName = LedModule.SystemLed.GPIO
	} else if index == LEDServer {
		ioName = LedModule.ServiceLed.GPIO
	}

	go LedFlashRoutine(ioName)
}

func LedFlashRoutine(name string) {
	fd, err := os.OpenFile(name, os.O_RDWR, 0777)
	if err != nil {
		//ZAPS.Errorf("LED接口[%s]打开失败 %v", name, err)
		return
	}
	defer fd.Close()

	ledTimer := time.NewTicker(time.Millisecond * 500)
	for {
		select {
		case <-ledTimer.C:
			{
				if LedFlashFlag == 1 {
					_, err = fd.Write([]byte("0"))
					if err != nil {
						//ZAPS.Errorf("LED接口[%s]写入失败 %v", name, err)
					}
					LedFlashFlag = 0
				} else if LedFlashFlag == 0 {
					_, err = fd.Write([]byte("1"))
					if err != nil {
						//ZAPS.Errorf("LED接口[%s]写入失败 %v", name, err)
					}
					LedFlashFlag = 1
				}
			}
		}
	}
}
