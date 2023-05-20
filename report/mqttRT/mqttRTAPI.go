package mqttRT

import (
	"encoding/csv"
	"gateway/report/reportModel"
	"gateway/setting"
	"gateway/utils"
	"github.com/pkg/errors"
	"os"
)

func (s *ReportServiceParamListRTTemplate) AddReportService(param ReportServiceGWParamRTTemplate) error {

	index := -1
	for k, v := range s.ServiceList {
		if v.GWParam.ServiceName == param.ServiceName {
			index = k
			break
		}
	}

	if index != -1 {
		return errors.New("服务名称已经存在")
	}

	nodeList := make([]*ReportServiceNodeParamRTTemplate, 0)
	ReportServiceParam := NewReportServiceParamRT(param, nodeList)
	s.ServiceList = append(s.ServiceList, ReportServiceParam)

	ReportServiceRTWriteParamToJson()
	return nil
}

func (s *ReportServiceParamListRTTemplate) ModifyReportService(param ReportServiceGWParamRTTemplate) error {

	index := -1
	for k, v := range s.ServiceList {
		if v.GWParam.ServiceName == param.ServiceName {
			index = k
			break
		}
	}

	//不存在表示增加
	if index == -1 {
		return errors.New("服务名称不存在")
	}
	//存在相同的，表示修改;
	if s.ServiceList[index].GWParam.MQTTClientOptions != nil {
		s.ServiceList[index].GWParam.MQTTClientOptions.SetAutoReconnect(false)
	}
	if s.ServiceList[index].GWParam.MQTTClient != nil {
		if s.ServiceList[index].GWParam.MQTTClient.IsConnected() {
			s.ServiceList[index].GWParam.MQTTClient.Disconnect(0)
		}
	}

	//旧上报服务退出
	s.ServiceList[index].CancelFunc()
	s.ServiceList[index].GWParam = param
	ReportServiceRTWriteParamToJson()

	//启动新上报服务
	s.ServiceList[index] = NewReportServiceParamRT(param, s.ServiceList[index].NodeList)

	return nil
}

func (s *ReportServiceParamListRTTemplate) DeleteReportService(serviceName string) {

	for k, v := range s.ServiceList {
		if v.GWParam.ServiceName == serviceName {
			s.ServiceList = append(s.ServiceList[:k], s.ServiceList[k+1:]...)
			ReportServiceRTWriteParamToJson()
			//协程退出
			v.CancelFunc()
			return
		}
	}
}

func (r *ReportServiceParamRTTemplate) AddReportNode(param ReportServiceNodeParamRTTemplate) error {

	param.Index = len(r.NodeList)
	param.CommStatus = "offLine"
	param.ReportStatus = "offLine"
	param.ReportErrCnt = 0

	//节点存在则进行修改
	for _, v := range r.NodeList {
		//节点已经存在
		if v.Name == param.Name {
			return errors.New("设备名称已经存在")
		}
	}

	ReportServiceRTWriteParamToJson()

	//节点不存在则新建
	rModel, ok := reportModel.ReportModels[param.UploadModel]
	if ok {
		param.Properties = rModel.Properties
	}
	r.NodeList = append(r.NodeList, &param)

	return nil
}

func (r *ReportServiceParamRTTemplate) ModifyReportNode(param ReportServiceNodeParamRTTemplate) error {

	param.CommStatus = "offLine"
	param.ReportStatus = "offLine"
	param.ReportErrCnt = 0

	//节点存在则进行修改
	for k, v := range r.NodeList {
		//节点已经存在
		if v.Name == param.Name {
			rModel, ok := reportModel.ReportModels[param.UploadModel]
			if ok {
				param.Properties = rModel.Properties
			}
			r.NodeList[k] = &param
			ReportServiceRTWriteParamToJson()
			return nil
		}
	}

	//节点不存在则新建
	return errors.New("设备名称不存在")
}

func (r *ReportServiceParamRTTemplate) DeleteReportNode(name string) int {

	index := -1
	//节点存在则进行修改
	for k, v := range r.NodeList {
		//节点已经存在
		if v.Name == name {
			index = k
			r.NodeList = append(r.NodeList[:k], r.NodeList[k+1:]...)
			ReportServiceRTWriteParamToJson()
			return index
		}
	}
	return index
}

func (r *ReportServiceParamRTTemplate) ExportParamToCsv() (bool, string) {

	//创建文件
	utils.DirIsExist("./tmp")
	fileName := "./tmp/" + r.GWParam.ServiceName + ".csv"

	fs, err := os.Create(fileName)
	if err != nil {
		setting.ZAPS.Errorf("创建csv文件错误 %v", err)
		return false, ""
	}
	defer fs.Close()

	//创建一个新的写入文件流
	w := csv.NewWriter(fs)
	csvData := [][]string{
		{"上报服务名称", "采集接口名称", "设备名称", "设备标签", "设备通信地址", "上报模型", "上报服务协议", "传输编码"},
		{"ServiceName", "CollInterfaceName", "Name", "Label", "Addr", "UploadModel", "Protocol", "DeviceCode"},
	}

	for _, n := range r.NodeList {
		param := make([]string, 0)
		param = append(param, n.ServiceName)
		param = append(param, n.CollInterfaceName)
		param = append(param, n.Name)
		param = append(param, n.Label)
		param = append(param, n.Addr)
		param = append(param, n.UploadModel)
		param = append(param, n.Protocol)
		param = append(param, n.Param.DeviceCode)
		csvData = append(csvData, param)
	}

	//写入数据
	err = w.WriteAll(csvData)
	if err != nil {
		setting.ZAPS.Errorf("写csv文件错误 %v", err)
		return false, ""
	}
	w.Flush()

	return true, fileName
}

func (r *ReportServiceParamRTTemplate) ExportParamToXlsx() (bool, string) {

	//创建文件
	utils.DirIsExist("./tmp")
	fileName := "./tmp/" + r.GWParam.ServiceName + ".xlsx"

	//创建一个新的写入文件流
	cellData := [][]string{
		{"上报服务名称", "采集接口名称", "设备名称", "设备标签", "设备通信地址", "上报模型", "上报服务协议", "传输编码"},
		{"ServiceName", "CollInterfaceName", "Name", "Label", "Addr", "UploadModel", "Protocol", "DeviceCode"},
	}

	for _, n := range r.NodeList {
		param := make([]string, 0)
		param = append(param, n.ServiceName)
		param = append(param, n.CollInterfaceName)
		param = append(param, n.Name)
		param = append(param, n.Label)
		param = append(param, n.Addr)
		param = append(param, n.UploadModel)
		param = append(param, n.Protocol)
		param = append(param, n.Param.DeviceCode)
		cellData = append(cellData, param)
	}

	//写入数据
	err := setting.WriteExcel(fileName, cellData)
	if err != nil {
		return false, ""
	}

	return true, fileName
}
