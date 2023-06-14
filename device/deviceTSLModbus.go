package device

import (
	"encoding/json"
	"errors"
	"gateway/device/eventBus"
	"gateway/setting"
	"gateway/utils"
	"time"
)

type TSLModbusModelTemplate struct {
	Index int    `json:"index"`
	Name  string `json:"name"`  //名称，只可以是字母和数字的组合
	Label string `json:"label"` //名称标签
	Type  int    `json:"type"`  //模型类型
	//Properties []*TSLModbusPropertyTemplate     `json:"properties"` //属性
	Cmd   map[string]*TSLModbusCmdTemplate `json:"cmd"`
	Event eventBus.Bus                     `json:"-"` //事件队列
}

type TSLModbusPropertyTemplate struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	AccessMode  int    `json:"accessMode"`
	Type        int    `json:"type"`
	Decimals    int    `json:"decimals"`
	Unit        string `json:"unit"`
	RegAddr     int    `json:"regAddr"`
	RegCnt      int    `json:"regCnt"`
	RuleType    string `json:"ruleType"`
	Formula     string `json:"formula"`
	IotDataType string `json:"iotDataType"`
}

type TSLModbusCmdTemplate struct {
	Name         string                       `json:"name"`
	Label        string                       `json:"label"`
	FunCode      int                          `json:"funCode"`      //功能码
	StartRegAddr int                          `json:"startRegAddr"` //寄存器起始地址
	RegCnt       int                          `json:"regCnt"`       //寄存器数量
	Registers    []*TSLModbusPropertyTemplate `json:"registers"`
}

var TSLModbusMap = make(map[string]*TSLModbusModelTemplate)
var modbusWriteTimer *time.Timer

func TSLModbusInit() {

	ReadTSLModbusParamFromJson()
	modbusWriteTimer = time.AfterFunc(time.Second, func() {
		WriteTSLModbusParamToJson()
	})
	modbusWriteTimer.Stop()
}

func ReadTSLModbusParamFromJson() bool {

	utils.DirIsExist("./selfpara")

	data, err := utils.FileRead("./selfpara/deviceTSLModbusParam.json")
	if err != nil {
		setting.ZAPS.Debugf("打开采集模型[Modbus]配置json文件失败 %v", err)
		return false
	}

	err = json.Unmarshal(data, &TSLModbusMap)
	if err != nil {
		setting.ZAPS.Errorf("采集模型[Modbus]配置json文件格式化失败 %v", err)
		for _, v := range TSLModbusMap {
			v.Event = eventBus.NewBus()
		}
		return false
	}

	for _, v := range TSLModbusMap {
		v.Event = eventBus.NewBus()
	}

	setting.ZAPS.Debugf("打开采集模型[Modbus]配置json文件成功")
	return true
}

func WriteTSLModbusParamToJson() {

	utils.DirIsExist("./selfpara")

	sJson, _ := json.Marshal(TSLModbusMap)
	err := utils.FileWrite("./selfpara/deviceTSLModbusParam.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("采集模型[Modbus]配置json文件写入失败")
		return
	}
	setting.ZAPS.Infof("采集模型[Modbus]配置json文件写入成功")
}

func NewTSLModbus(name string, label string, modelType int) *TSLModbusModelTemplate {
	return &TSLModbusModelTemplate{
		Index: len(TSLModbusMap),
		Name:  name,
		Label: label,
		Type:  modelType,
		Cmd:   make(map[string]*TSLModbusCmdTemplate),
		//Properties: make([]*TSLModbusPropertyTemplate, 0),
		Event: eventBus.NewBus(),
	}
}

func (t *TSLModbusModelTemplate) GetTSLModelIndex() int {
	return t.Index
}

func (t *TSLModbusModelTemplate) GetTSLModelName() string {
	return t.Name
}

func (t *TSLModbusModelTemplate) GetTSLModelLabel() string {
	return t.Label
}

func (t *TSLModbusModelTemplate) GetTSLModelType() int {
	return t.Type
}

func (t *TSLModbusModelTemplate) GetTSLModelParam() interface{} {
	return ""
}

func (t *TSLModbusModelTemplate) TSLModelPropertiesGet() interface{} {

	properties := make([]*TSLModbusPropertyTemplate, 0)
	for _, v := range t.Cmd {
		properties = append(properties, v.Registers...)
	}

	return properties
}

func (t *TSLModbusModelTemplate) GetTSLModbusModelProperties() []*TSLModbusPropertyTemplate {

	properties := make([]*TSLModbusPropertyTemplate, 0)
	for _, v := range t.Cmd {
		properties = append(properties, v.Registers...)
	}

	return properties
}

func DeleteTSLModbus(name string) error {

	for _, v := range CollectInterfaceMap.Coll {
		for _, d := range v.DeviceNodeMap {
			if d.TSL == name {
				return errors.New("采集模型[Modbus]已被使用，不可以删除")
			}
		}
	}

	_, ok := TSLModelsName[name]
	if !ok {
		return errors.New("采集模型[Modbus]不存在")
	}

	delete(TSLModbusMap, name)
	delete(TSLModelsName, name)
	WriteTSLModbusParamToJson()

	return nil
}

func ModifyTSLModbus(name string, label string) error {

	tsl, ok := TSLModbusMap[name]
	if !ok {
		return errors.New("采集模型[Modbus]模版不存在")
	}

	tsl.Label = label
	tsl.Type = TSLModelTypeModbus
	setting.ZAPS.Infof("采集模型[Modbus][%s]发生变化", name)
	eventMsg := TSLEventTemplate{
		TSL:   name,
		Topic: "modify",
	}
	_ = tsl.Event.Publish("modify", eventMsg)
	WriteTSLModbusParamToJson()

	return nil
}

func (d *TSLModbusModelTemplate) TSLModelCmdAdd(cmd TSLModbusCmdTemplate) error {

	_, ok := d.Cmd[cmd.Name]
	if !ok {
		d.Cmd[cmd.Name] = &cmd
		eventMsg := TSLEventTemplate{
			TSL:   d.Name,
			Topic: "modify",
		}
		_ = d.Event.Publish("modify", eventMsg)
		modbusWriteTimer.Reset(time.Second)
		return nil
	} else {
		return errors.New("命令名称已经存在")
	}
}

func (d *TSLModbusModelTemplate) TSLModelCmdDelete(cmdNames []string) error {

	cnt := 0
	for _, p := range cmdNames {
		_, ok := d.Cmd[p]
		if !ok {
			continue
		} else {
			delete(d.Cmd, p)
			cnt++
		}
	}
	if cnt != 0 {
		modbusWriteTimer.Reset(time.Second)
		eventMsg := TSLEventTemplate{
			TSL:   d.Name,
			Topic: "modify",
		}
		d.Event.Publish("modify", eventMsg)
		return nil
	}
	return errors.New("命令名称不存在")
}

func (d *TSLModbusModelTemplate) TSLModelCmdModify(cmd TSLModbusCmdTemplate) error {

	_, ok := d.Cmd[cmd.Name]
	if !ok {
		return errors.New("命令名称不存在")
	}
	d.Cmd[cmd.Name].Label = cmd.Label
	d.Cmd[cmd.Name].FunCode = cmd.FunCode
	d.Cmd[cmd.Name].StartRegAddr = cmd.StartRegAddr
	d.Cmd[cmd.Name].RegCnt = cmd.RegCnt

	modbusWriteTimer.Reset(time.Second)

	eventMsg := TSLEventTemplate{
		TSL:   d.Name,
		Topic: "modify",
	}
	d.Event.Publish("modify", eventMsg)

	return nil

}

func (d *TSLModbusModelTemplate) TSLModelPropertiesAdd(cmdName string, property TSLModbusPropertyTemplate) error {

	cmd, ok := d.Cmd[cmdName]
	if !ok {
		return errors.New("命令名称不存在")
	}

	index := -1
	for k, v := range cmd.Registers {
		if v.Name == property.Name {
			index = k
			return errors.New("属性名称已经存在")
		}
	}
	if index == -1 {
		d.Cmd[cmdName].Registers = append(d.Cmd[cmdName].Registers, &property)
		//d.Properties = append(d.Properties, &property)
		eventMsg := TSLEventTemplate{
			TSL:   d.Name,
			Topic: "modify",
		}
		_ = d.Event.Publish("modify", eventMsg)
	}

	modbusWriteTimer.Reset(time.Second)

	return nil
}

func (d *TSLModbusModelTemplate) TSLModelPropertiesDelete(cmdName string, propertiesName []string) error {

	cnt := 0
	cmd := d.Cmd[cmdName]
	for _, p := range propertiesName {
		for k := 0; k < len(cmd.Registers); k++ {
			if cmd.Registers[k].Name == p {
				cnt = 1
				//setting.ZAPS.Debugf("properties %v", d.Properties)
				cmd.Registers = append(cmd.Registers[:k], cmd.Registers[k+1:]...)
			}
		}
	}

	if cnt != 0 {
		modbusWriteTimer.Reset(time.Second)
		eventMsg := TSLEventTemplate{
			TSL:   d.Name,
			Topic: "modify",
		}
		d.Event.Publish("modify", eventMsg)
		return nil
	}
	return errors.New("属性名称不存在")

}

func (d *TSLModbusModelTemplate) TSLModelPropertiesModify(cmdName string, property TSLModbusPropertyTemplate) error {

	cmd, ok := d.Cmd[cmdName]
	if !ok {
		return errors.New("命令名称不存在")
	}

	for k, v := range cmd.Registers {
		if v.Name == property.Name {
			cmd.Registers[k] = &property
		}
	}

	index := -1
	for k, v := range cmd.Registers {
		if v.Name == property.Name {
			cmd.Registers[k] = &property
			index = k
		}
	}
	if index == -1 {
		return errors.New("参数名称不存在")
	}

	eventMsg := TSLEventTemplate{
		TSL:   d.Name,
		Topic: "modify",
	}
	_ = d.Event.Publish("modify", eventMsg)
	modbusWriteTimer.Reset(time.Second)

	return nil
}
