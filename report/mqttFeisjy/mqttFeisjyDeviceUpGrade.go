package mqttFeisjy

import (
	"encoding/json"
	"errors"
	"fmt"
	"gateway/setting"
	"gateway/utils"
	"os"
	"time"
)

const (
	DeviceUpGradePath = "./"
	//DeviceUpGradeName = "armv7_gw_linux"
)

type MQTTFeisjyUpGradeTemplate struct {
	UpgradeStep      string `json:"upgradeStep"`
	WhetherToUpgrade int    `json:"whetherToUpgrade"`
	FilePath         string `json:"filePath"`
	Crc              string `json:"crc"`
}

type MQTTFeisjyUpGradeReadyInfoTemplate struct {
	UpgradeStep string `json:"upgradeStep"`
	ReadyInfo   int    `json:"readyInfo"`
}

type MQTTFeisjyUpGradeStatusTemplate struct {
	UpgradeStep string `json:"upgradeStep"`
	Status      string `json:"status"`
}

// 设备固件升级流程状态机
func (r *ReportServiceParamFeisjyTemplate) FeisjyDevUpGradeaMachine(upFrame *MQTTFeisjyUpGradeTemplate) bool {

	fmt.Println(upFrame.UpgradeStep)
	switch upFrame.UpgradeStep {
	case "upgradeConfirm": //收到平台发过来的登录确认帧
		{
			r.FeisjyDeviceUpgradeRedyInfo()
		}
	case "upgradeInfo": // 升级固件信息
		{
			// 为下载函数go一个携程，这样才不会阻塞当前携程的进行
			go func() bool {
				setting.ZAPS.Infof("服务[ deviceUpgrade->%s ] 固件包开始下载将存储在：[ %s ] 命名为: [ %s ]", upFrame.UpgradeStep, DeviceUpGradePath, setting.ExeName)

				err := setting.DownloadFileFeisjy(DeviceUpGradePath, setting.ExeName, upFrame.FilePath, 3)
				if err != nil {
					setting.ZAPS.Errorf("固件包下载失败! %v ", err)
				} else {
					setting.ZAPS.Errorf("固件包下载成功!")
				}

				if err == nil {

					crc16 := utils.CalculateFileCRC16(DeviceUpGradePath + setting.ExeName)
					println(crc16, upFrame.Crc)

					if crc16 != upFrame.Crc {
						err = errors.New("CrcFail")
					}
					return r.FeisjyDeviceUpgradeStatus(err)
				} else {
					return false
				}
			}()
		}
	}

	return false
}

// 回复平台固件升级信息
func (r *ReportServiceParamFeisjyTemplate) FeisjyDeviceUpgradeRedyInfo() bool {

	status := false

	upgradeReadyData := MQTTFeisjyUpGradeReadyInfoTemplate{
		UpgradeStep: "upgradeReady",
		ReadyInfo:   1,
	}

	sJson, _ := json.Marshal(upgradeReadyData)

	propertyPostTopic := fmt.Sprintf(FeisjyMQTTTopicRxFormat, r.GWParam.Param.AppKey, r.GWParam.Param.DeviceID, "deviceUpgrade")

	if r.GWParam.MQTTClient != nil {
		if token := r.GWParam.MQTTClient.Publish(propertyPostTopic, 0, false, sJson); token.WaitTimeout(2000*time.Millisecond) && token.Error() != nil {
			status = false
			setting.ZAPS.Debugf("上报服务[%s]发布[%s]登录确认消息失败 %v", r.GWParam.ServiceName, r.GWParam.Param.DeviceID, token.Error())
		} else {
			status = true
			setting.ZAPS.Debugf("上报服务[%s]发布[%s]登录确认消息成功 内容%v", r.GWParam.ServiceName, r.GWParam.Param.DeviceID, string(sJson))
		}
	}

	return status
}

// 回复平台设备升级状态
func (r *ReportServiceParamFeisjyTemplate) FeisjyDeviceUpgradeStatus(err error) bool {

	var status = false

	upgradeStatusData := MQTTFeisjyUpGradeStatusTemplate{
		UpgradeStep: "upgradeStatus",
	}

	if err != nil {
		upgradeStatusData.Status = "downloadFailed"
	} else if err == nil {
		upgradeStatusData.Status = "upgradeSuccess"
	}

	sJson, _ := json.Marshal(upgradeStatusData)
	propertyPostTopic := fmt.Sprintf(FeisjyMQTTTopicRxFormat, r.GWParam.Param.AppKey, r.GWParam.Param.DeviceID, "deviceUpgrade")

	if r.GWParam.MQTTClient != nil {
		if token := r.GWParam.MQTTClient.Publish(propertyPostTopic, 0, false, sJson); token.WaitTimeout(2000*time.Millisecond) && token.Error() != nil {
			status = false
			setting.ZAPS.Debugf("上报服务[%s]发布[%s]登录确认消息失败 %v", r.GWParam.ServiceName, r.GWParam.Param.DeviceID, token.Error())
		} else {
			status = true
			setting.ZAPS.Debugf("上报服务[%s]发布[%s]登录确认消息成功 内容%v", r.GWParam.ServiceName, r.GWParam.Param.DeviceID, string(sJson))
		}
	}

	//成功之后程序退出
	os.Exit(0)

	return status
}
