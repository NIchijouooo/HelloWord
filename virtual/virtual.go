package virtual

import (
	"encoding/json"
	"errors"
	"fmt"
	"gateway/device"
	"gateway/device/eventBus"
	"gateway/setting"
	"gateway/utils"
	"math"
	"strconv"
	"strings"
	"time"
)

type VirtualDeviceTemplate struct {
	Nodes    map[string]*VirtualNodeTemplate `json:"nodes"`
	EventBus eventBus.Bus                    `json:"-"`
}

type VirtualNodeTemplate struct {
	Name       string                              `json:"name"`
	Label      string                              `json:"label"`
	Index      int                                 `json:"index"`
	Properties map[string]*VirtualPropertyTemplate `json:"properties"`
	CommStatus string                              `json:"commStatus"`
}

type VirtualPropertyTemplate struct {
	Name          string                         `json:"name"`
	Label         string                         `json:"label"`
	Type          int                            `json:"type"`
	Decimals      int                            `json:"decimals"`
	Unit          string                         `json:"unit"`
	Params        VirtualPropertyParamTemplate   `json:"params"`
	AlarmParams   VirtualAlarmParamTemplate      `json:"alarmParams,omitempty"`
	Value         []VirtualPropertyValueTemplate `json:"value,omitempty"`
	SaveValue     interface{}                    `json:"saveValue"`
	SaveValueFlag bool                           `json:"saveValueFlag"`
}

type VirtualPropertyParamTemplate struct {
	CollName     string `json:"collName"`
	DeviceName   string `json:"deviceName"`
	PropertyName string `json:"propertyName"`
	SaveEnable   bool   `json:"saveEnable"` //掉电保存
	Formula      string `json:"formula"`    //计算公示
}

type VirtualAlarmParamTemplate struct {
	Min             string `json:"min"`             //最小
	Max             string `json:"max"`             //最大
	MinMaxAlarm     bool   `json:"minMaxAlarm"`     //范围报警
	Step            string `json:"step"`            //步长
	StepAlarm       bool   `json:"stepAlarm"`       //阶跃报警
	DataLength      string `json:"dataLength"`      //字符串长度
	DataLengthAlarm bool   `json:"dataLengthAlarm"` //字符长度报警
}

type VirtualPropertyValueTemplate struct {
	Index     int         `json:"index"`
	Value     interface{} `json:"value"`   //变量值，不可以是字符串
	Explain   interface{} `json:"explain"` //变量值解释，必须是字符串
	TimeStamp time.Time   `json:"timeStamp"`
}

type VirtualEventTemplate struct {
	Topic         string   `json:"topic"` //事件主题，online，offline，update
	NodeName      string   `json:"nodeName"`
	PropertyNames []string `json:"propertyNames"` //属性名称
}

var VirtualDevice = &VirtualDeviceTemplate{
	Nodes: make(map[string]*VirtualNodeTemplate),
}

var writeTimer *time.Timer

func VirtualDeviceInit() {

	writeTimer = time.AfterFunc(time.Second, func() {
		VirtualDeviceWriteToJson()
	})
	writeTimer.Stop()

	err := VirtualDeviceReadFromJson()
	if err != nil {
		return
	}

	go VirtualDeviceSyncProperty()
}

func VirtualDeviceReadFromJson() error {

	VirtualDevice.EventBus = eventBus.NewBus()

	data, err := utils.FileRead("./selfpara/virtualDevice.json")
	if err != nil {
		setting.ZAPS.Errorf("虚拟设备配置json文件读取失败 %v", err)
		return err
	}

	err = json.Unmarshal(data, &VirtualDevice)
	if err != nil {
		setting.ZAPS.Errorf("虚拟设备配置json文件格式化失败 %v", err)
		return err
	}

	//清空Value
	for _, n := range VirtualDevice.Nodes {
		for _, p := range n.Properties {
			p.Value = p.Value[0:0]
		}
		n.CommStatus = "offLine"
	}

	setting.ZAPS.Info("虚拟设备配置json文件读取成功")

	return nil
}

func VirtualDeviceWriteParam() {
	writeTimer.Reset(time.Second)
}

func VirtualDeviceWriteToJson() {
	utils.DirIsExist("./selfpara")

	sJson, _ := json.Marshal(VirtualDevice)
	err := utils.FileWrite("./selfpara/virtualDevice.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("虚拟设备配置json文件写入失败")
		return
	}
	setting.ZAPS.Debugf("虚拟设备配置json文件写入成功")
}

func (v *VirtualDeviceTemplate) VirtualDeviceAddNode(name string, label string) error {

	_, ok := v.Nodes[name]
	if ok {
		setting.ZAPS.Errorf("设备名称已经存在")
		return errors.New("设备名称已经存在")
	}

	node := &VirtualNodeTemplate{
		Name:       name,
		Label:      label,
		Index:      len(v.Nodes),
		Properties: make(map[string]*VirtualPropertyTemplate),
	}

	v.Nodes[name] = node

	VirtualDeviceWriteParam()
	return nil
}

func (v *VirtualDeviceTemplate) VirtualDeviceModifyNode(name string, label string) error {

	node, ok := v.Nodes[name]
	if !ok {
		setting.ZAPS.Errorf("设备名称不存在")
		return errors.New("设备名称不存在")
	}

	node.Label = label

	v.Nodes[name] = node

	VirtualDeviceWriteParam()
	return nil
}

func (v *VirtualDeviceTemplate) VirtualDeviceDeleteNodes(names []string) error {

	for _, n := range names {
		_, ok := v.Nodes[n]
		if !ok {
			continue
		}
		delete(v.Nodes, n)
		VirtualDeviceWriteParam()
	}
	return nil
}

func (v *VirtualNodeTemplate) VirtualDeviceAddProperty(property VirtualPropertyTemplate) error {

	_, ok := v.Properties[property.Name]
	if ok {
		setting.ZAPS.Errorf("属性名称已经存在")
		return errors.New("属性名称已经存在")
	}

	v.Properties[property.Name] = &property

	VirtualDeviceWriteParam()
	return nil
}

func (v *VirtualNodeTemplate) VirtualDeviceModifyProperty(property VirtualPropertyTemplate) error {

	_, ok := v.Properties[property.Name]
	if !ok {
		setting.ZAPS.Errorf("属性名称不存在")
		return errors.New("属性名称不存在")
	}

	v.Properties[property.Name] = &property

	VirtualDeviceWriteParam()
	return nil
}

func (v *VirtualNodeTemplate) VirtualDeviceDeleteProperties(names []string) error {

	for _, n := range names {
		_, ok := v.Properties[n]
		if !ok {
			continue
		}
		delete(v.Properties, n)
		VirtualDeviceWriteParam()
	}
	return nil
}

func VirtualDeviceSyncProperty() {

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			{
				//setting.ZAPS.Debugf("虚拟设备变量更新")
				for k, _ := range VirtualDevice.Nodes {
					VirtualDevice.Nodes[k].VirtualDeviceUpdateProperty()
					VirtualDevice.Nodes[k].VirtualDeviceUpdateSaveValue()
				}
				VirtualDevice.ProcessAlarmEvent()
			}
		}
	}
}

func (v *VirtualNodeTemplate) VirtualDeviceUpdateSaveValue() {
	for _, p := range v.Properties {
		if p.Params.SaveEnable == true && p.SaveValueFlag == false {
			coll, ok := device.CollectInterfaceMap.Coll[p.Params.CollName]
			if !ok {
				continue
			}
			node, ok := coll.DeviceNodeMap[p.Params.DeviceName]
			if !ok {
				continue
			}
			if node.CommStatus == "onLine" {
				p.SaveValueFlag = true
				property := make(map[string]interface{})
				property[p.Name] = p.SaveValue
				v.VirtualDeviceSetVariables(property)
			}
			//v.CommStatus = node.CommStatus
		}
	}
}

func (v *VirtualNodeTemplate) VirtualDeviceUpdateProperty() {

	event := VirtualEventTemplate{
		NodeName:      v.Name,
		PropertyNames: make([]string, 0),
	}

	for _, p := range v.Properties {
		coll, ok := device.CollectInterfaceMap.Coll[p.Params.CollName]
		if !ok {
			continue
		}
		node, ok := coll.DeviceNodeMap[p.Params.DeviceName]
		if !ok {
			continue
		}

		if (node.CommStatus != v.CommStatus) && (v.CommStatus != "") {
			if v.CommStatus == "offLine" {
				event.Topic = "onLine"
				_ = VirtualDevice.EventBus.Publish("onLine", event)
			} else if v.CommStatus == "onLine" {
				setting.ZAPS.Debugf("name %v,nodeStatus %v,virtualStatus %v", node.Name, node.CommStatus, v.CommStatus)
				event.Topic = "offLine"
				_ = VirtualDevice.EventBus.Publish("offLine", event)
			}
			v.CommStatus = node.CommStatus
		}

		for _, n := range node.Properties {
			if n.Name == p.Params.PropertyName {
				cnt := len(n.Value)
				//setting.ZAPS.Debugf("property %v", p.Value)
				if cnt > 0 {
					propertyValue := VirtualPropertyValueTemplate{
						Index:     len(n.Value),
						Explain:   n.Value[cnt-1].Explain,
						TimeStamp: n.Value[cnt-1].TimeStamp,
					}
					if p.Params.Formula != "" {
						fStr := p.Params.Formula
						if strings.Contains(fStr, "t") {
							switch n.Value[cnt-1].Value.(type) {
							case uint32:
								fStr = strings.ReplaceAll(fStr, "t", fmt.Sprintf("%d", n.Value[cnt-1].Value.(uint32)))
							case int32:
								fStr = strings.ReplaceAll(fStr, "t", fmt.Sprintf("%d", n.Value[cnt-1].Value.(int32)))
							case float64:
								fStr = strings.ReplaceAll(fStr, "t", fmt.Sprintf("%f", n.Value[cnt-1].Value.(float64)))
							case string:
							}
						}
						err, value := setting.FormulaRun(fStr)
						if err != nil {
							propertyValue.Value = n.Value[cnt-1].Value.(float64)
						} else {
							propertyValue.Value = value
						}
					} else {
						propertyValue.Value = n.Value[cnt-1].Value
					}

					if len(p.Value) == 2 {
						p.Value = append(p.Value[:0], p.Value[1:]...)
					}
					p.Value = append(p.Value, propertyValue)
				}
			}
		}
	}
}

func (v *VirtualNodeTemplate) VirtualDeviceGetVariables(propertiesName []string) {

	for _, p := range propertiesName {
		propery, ok := v.Properties[p]
		if !ok {
			continue
		}

		coll, ok := device.CollectInterfaceMap.Coll[propery.Params.CollName]
		if !ok {
			continue
		}

		_, ok = coll.DeviceNodeMap[propery.Params.DeviceName]
		if !ok {
			continue
		}
		//从采集服务中找到相应节点
		cmd := device.CommunicationCmdTemplate{}
		cmd.CollInterfaceName = propery.Params.CollName
		cmd.DeviceName = propery.Params.DeviceName
		cmd.FunName = "GetRealVariables"
		variableMap := make([]string, 0)
		variableMap = append(variableMap, propery.Params.PropertyName)
		paramStr, _ := json.Marshal(variableMap)
		cmd.FunPara = string(paramStr)

		coll.CommQueueManage.CommunicationManageAddEmergency(cmd)

		//采集服务数据更新到虚拟服务中
		v.VirtualDeviceUpdateProperty()
	}

}

func (v *VirtualNodeTemplate) VirtualDeviceSetVariables(propertiesName map[string]interface{}) map[string]interface{} {

	propertyMap := make(map[string]interface{})
	for k, p := range propertiesName {
		propery, ok := v.Properties[k]
		if !ok {
			continue
		}

		if propery.Params.SaveEnable == true {
			propery.SaveValue = p
			VirtualDeviceWriteParam()
		}

		coll, ok := device.CollectInterfaceMap.Coll[propery.Params.CollName]
		if !ok {
			continue
		}

		_, ok = coll.DeviceNodeMap[propery.Params.DeviceName]
		if !ok {
			continue
		}
		//从采集服务中找到相应节点
		cmd := device.CommunicationCmdTemplate{}
		cmd.CollInterfaceName = propery.Params.CollName
		cmd.DeviceName = propery.Params.DeviceName
		cmd.FunName = "SetVariables"
		propertyMap[propery.Params.PropertyName] = p
		paramStr, _ := json.Marshal(propertyMap)
		cmd.FunPara = string(paramStr)

		//setting.ZAPS.Debugf("虚拟服务写变量 %v", cmd)

		ackData := coll.CommQueueManage.CommunicationManageAddEmergency(cmd)
		if ackData.Status {
			propertyMap[k] = 0
		} else {
			propertyMap[k] = -1
		}
	}

	return propertyMap
}

func (v *VirtualDeviceTemplate) ProcessAlarmEvent() {

	properties := VirtualEventTemplate{
		Topic:         "update",
		PropertyNames: make([]string, 0),
	}

	for _, n := range v.Nodes {
		for _, p := range n.Properties {
			if p.AlarmParams.StepAlarm == true {
				valueCnt := len(p.Value)
				if valueCnt >= 2 { //阶跃报警必须是2个值
					switch p.Type {
					case device.PropertyTypeInt32:
						{
							pValueCur := p.Value[valueCnt-1].Value.(int32)
							pValuePre := p.Value[valueCnt-2].Value.(int32)
							step, _ := strconv.Atoi(p.AlarmParams.Step)
							if math.Abs(float64(pValueCur-pValuePre)) > float64(step) {
								setting.ZAPS.Infof("属性[%v]阶跃报警", p.Name)
								properties.PropertyNames = append(properties.PropertyNames, p.Name)
							}
						}
					case device.PropertyTypeUInt32:
						{
							pValueCur := p.Value[valueCnt-1].Value.(uint32)
							pValuePre := p.Value[valueCnt-2].Value.(uint32)
							step, _ := strconv.Atoi(p.AlarmParams.Step)
							if math.Abs(float64(pValueCur-pValuePre)) > float64(step) {
								setting.ZAPS.Infof("属性[%v]阶跃报警", p.Name)
								properties.PropertyNames = append(properties.PropertyNames, p.Name)
							}
						}
					case device.PropertyTypeDouble:
						{
							pValueCur := p.Value[valueCnt-1].Value.(float64)
							pValuePre := p.Value[valueCnt-2].Value.(float64)
							step, err := strconv.ParseFloat(p.AlarmParams.Step, 64)
							if err != nil {
								continue
							}
							if math.Abs(pValueCur-pValuePre) > float64(step) {
								setting.ZAPS.Infof("属性[%v]阶跃报警", p.Name)
								properties.PropertyNames = append(properties.PropertyNames, p.Name)
							}
						}
					}
				}
			}
			if p.AlarmParams.MinMaxAlarm == true {
				valueCnt := len(p.Value)
				if valueCnt > 0 {
					switch p.Value[valueCnt-1].Value.(type) {
					case uint32:
						{
							pValueCur := p.Value[valueCnt-1].Value.(uint32)
							min, _ := strconv.Atoi(p.AlarmParams.Min)
							max, _ := strconv.Atoi(p.AlarmParams.Max)
							if pValueCur < uint32(min) || pValueCur > uint32(max) {
								setting.ZAPS.Infof("属性[%v]范围报警", p.Name)
								properties.PropertyNames = append(properties.PropertyNames, p.Name)
							}
						}
					case int32:
						{
							pValueCur := p.Value[valueCnt-1].Value.(int32)
							min, _ := strconv.Atoi(p.AlarmParams.Min)
							max, _ := strconv.Atoi(p.AlarmParams.Max)
							if pValueCur < int32(min) || pValueCur > int32(max) {
								setting.ZAPS.Infof("属性[%v]范围报警", p.Name)
								properties.PropertyNames = append(properties.PropertyNames, p.Name)
							}
						}
					case float64:
						{
							pValueCur := p.Value[valueCnt-1].Value.(float64)
							min, err := strconv.ParseFloat(p.AlarmParams.Min, 64)
							if err != nil {
								continue
							}
							max, err := strconv.ParseFloat(p.AlarmParams.Max, 64)
							if err != nil {
								continue
							}
							setting.ZAPS.Debugf("value %v,min %v,max %v", pValueCur, min, max)
							if pValueCur < min || pValueCur > max {
								setting.ZAPS.Infof("属性[%v]范围报警", p.Name)
								properties.PropertyNames = append(properties.PropertyNames, p.Name)
							}
						}
					}
				}
			}
		}
	}

	if len(properties.PropertyNames) > 0 {
		_ = v.EventBus.Publish("update", properties)
	}

}
