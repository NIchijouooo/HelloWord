package mqttFeisjy

import (
	"encoding/json"
	"fmt"
	"gateway/setting"
	"gateway/utils"
	"path/filepath"
	"time"
)

// 1 :"./selfpara" 2:"./selfpara" "./log"
const (
	upLoadNum = 2
)

var upLoadFileList = []string{
	"./selfpara",
	"./log",
}

type FileNameListFeisjyTemplate struct {
	Name       string `json:"name"`
	Crc        string `json:"crc"`
	DocType    string `json:"docType"` // 文件类型，目前暂定类型有：log /  data / config /
	UniqueCode string `json:"uniqueCode"`
	Code       int    `json:"code"`
	Url        string `json:"url"`
}

type FileListFeisjyTemplate struct {
	FileStep string                       `json:"fileStep"`
	FileList []FileNameListFeisjyTemplate `json:"fileList"`
	Url      string                       `json:"url"`
}

// 设备固件升级流程状态机
func (r *ReportServiceParamFeisjyTemplate) FeisjyFileListMachine(f FileListFeisjyTemplate) bool {

	switch f.FileStep {
	case "getFileList": // 获取文件列表
		return r.FeisjyFileListResult()
	case "readyToReceive": // 平台准备好接收文件
		{
			go func() {
				for i, v := range f.FileList {
					if i == 0 {
						setting.ZAPS.Infof("服务[FileList->%s] File[%s] 开始上传", f.FileStep, v.Name)
						err := setting.UpLoadFileFeisjy(v.Name, f.Url, v.UniqueCode, r.GWParam.Param.DeviceID)
						if err != nil {
							setting.ZAPS.Errorf("File[%s] 上传失败！:%v", v.Name, err)
						} else {
							setting.ZAPS.Infof("File[%s] 上传成功", v.Name)
						}
					}
				}
			}()
		}
	case "resultDownloadFile": // 设备下载文件
		{
			go func() {
				for _, v := range f.FileList { // 遍历文件列表 [{"code":xx,"uniqueCode":"xxx","name":"xxx","url":"xxx"},]
					if v.Code == 200 { // 等于200为文件正常可下载，其他为不可下载
						var ext string                          // 文件拓展名
						dir, fileName := filepath.Split(v.Name) //
						if len(dir) <= 0 {
							dir = "./other/"
						}
						if len(filepath.Ext(v.Name)) == 0 {
							ext = filepath.Ext(v.Url)
						}

						setting.ZAPS.Infof("服务[FileList->%s] File[%s] 开始下载", f.FileStep, fileName+ext)

						err := setting.DownloadFileFeisjy(dir, fileName+ext, v.Url, 3)
						if err != nil {
							setting.ZAPS.Errorf("File[%s] 下载失败！:%v", v.Name, err)
						} else {
							setting.ZAPS.Infof("File[%s] 下载成功", v.Name)
						}
					}
				}
			}()
		}
	}

	return false
}

func (r *ReportServiceParamFeisjyTemplate) FeisjyFileListResult() bool {

	status := false
	var fileLists []string
	var fileList []string
	var err error

	FileListFeisjy := FileListFeisjyTemplate{
		FileStep: "resultFileList",
	}
	FileNameListFeisjy := FileNameListFeisjyTemplate{
		DocType: "config",
	}

	for i := 0; i < upLoadNum; i++ {
		fmt.Println(upLoadFileList[i])
		fileList, err = utils.GetAllFileFormDir(upLoadFileList[i])
		if err == nil {
			fileLists = append(fileLists, fileList...)
		}
	}

	if err != nil {
		setting.ZAPS.Errorf("%v", err)
	}

	if len(fileLists) > 0 {

		for _, v := range fileLists {
			FileNameListFeisjy.Name = v
			FileNameListFeisjy.Crc = fmt.Sprintf("%d", utils.CalculateFileCRC16(v))
			FileListFeisjy.FileList = append(FileListFeisjy.FileList, FileNameListFeisjy)
		}

		sJson, _ := json.Marshal(FileListFeisjy)

		propertyPostTopic := fmt.Sprintf(FeisjyMQTTTopicRxFormat, r.GWParam.Param.AppKey, r.GWParam.Param.DeviceID, "fileList")

		if r.GWParam.MQTTClient != nil {
			if token := r.GWParam.MQTTClient.Publish(propertyPostTopic, 0, false, sJson); token.WaitTimeout(2000*time.Millisecond) && token.Error() != nil {
				status = false
				setting.ZAPS.Debugf("上报服务[%s]发布[%s]登录确认消息失败 %v", r.GWParam.ServiceName, r.GWParam.Param.DeviceID, token.Error())
			} else {
				status = true
				setting.ZAPS.Debugf("上报服务[%s]发布[%s]登录确认消息成功 内容%v", r.GWParam.ServiceName, r.GWParam.Param.DeviceID, string(sJson))
			}
		}
	}

	return status
}
