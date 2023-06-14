package device

import (
	"encoding/json"
	"errors"
	"gateway/device/eventBus"
	"gateway/setting"
	"gateway/utils"
	"time"
)

type TSLDLT6452007ModelTemplate struct {
	Index int                                  `json:"index"`
	Name  string                               `json:"name"`  //名称，只可以是字母和数字的组合
	Label string                               `json:"label"` //名称标签
	Type  int                                  `json:"type"`  //模型类型
	Cmd   map[string]*TSLDLT6452007CmdTemplate `json:"cmd"`
	Event eventBus.Bus                         `json:"-"` //事件队列
}

type TSLDLT6452007PropertyTemplate struct {
	Name           string `json:"name"`
	Label          string `json:"label"`
	RulerId        string `json:"rulerId"` //数据标识
	Format         string `json:"format"`  //数据格式YYMMDDhhmm,XXXXXX.XX,XX.XXXX...
	Len            int    `json:"len"`     //数据长度
	Unit           string `json:"unit"`
	AccessMode     int    `json:"accessMode"`
	BlockAddOffset int    `json:"blockAddOffset"` //当前数据在块数据域内的偏移地址
	RulerAddOffset int    `json:"rulerAddOffset"` //当前变量在当前ID数据地址中的偏移地址
	Type           int    `json:"type"`           //float,uint32...
}

type TSLDLT6452007CmdTemplate struct {
	Name         string                           `json:"name"`
	Label        string                           `json:"label"`
	BlockRulerId string                           `json:"blockRulerId"` //块ID
	BlockRead    byte                             `json:"blockRead"`    //是否按块读写 0是单个标识读取  1是按块读取
	Properties   []*TSLDLT6452007PropertyTemplate `json:"properties"`
}

var TSLDLT6452007Map = make(map[string]*TSLDLT6452007ModelTemplate)
var dlt6452007WriteTimer *time.Timer

func TSLDLT6452007Init() {

	ReadTSLDLT6452007ParamFromJson()
	dlt6452007WriteTimer = time.AfterFunc(time.Second, func() {
		WriteTSLDLT6452007ParamToJson()
	})
	dlt6452007WriteTimer.Stop()
}

func ReadTSLDLT6452007ParamFromJson() bool {

	utils.DirIsExist("./selfpara")

	data, err := utils.FileRead("./selfpara/deviceTSLDLT6452007Param.json")
	if err != nil {
		setting.ZAPS.Debugf("打开采集模型[DLT645-2007]配置json文件失败 %v", err)
		return false
	}

	err = json.Unmarshal(data, &TSLDLT6452007Map)

	if err != nil {
		setting.ZAPS.Errorf("采集模型[DLT645-2007]配置json文件格式化失败 %v", err)
		for _, v := range TSLDLT6452007Map {
			v.Event = eventBus.NewBus()
		}
		return false
	}

	for _, v := range TSLDLT6452007Map {
		v.Event = eventBus.NewBus()
	}

	setting.ZAPS.Debugf("打开采集模型[DLT645-2007]配置json文件成功")
	return true
}

func WriteTSLDLT6452007ParamToJson() {

	utils.DirIsExist("./selfpara")

	sJson, _ := json.Marshal(TSLDLT6452007Map)
	err := utils.FileWrite("./selfpara/deviceTSLDLT6452007Param.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("采集模型[DLT645-2007]配置json文件写入失败")
		return
	}
	setting.ZAPS.Infof("采集模型[DLT645-2007]配置json文件写入成功")
}

func NewTSLDLT6452007(name string, label string, modelType int) *TSLDLT6452007ModelTemplate {
	return &TSLDLT6452007ModelTemplate{
		Index: len(TSLDLT6452007Map),
		Name:  name,
		Label: label,
		Type:  modelType,
		Cmd:   make(map[string]*TSLDLT6452007CmdTemplate),
		Event: eventBus.NewBus(),
	}
}

func (t *TSLDLT6452007ModelTemplate) GetTSLModelIndex() int {
	return t.Index
}

func (t *TSLDLT6452007ModelTemplate) GetTSLModelName() string {
	return t.Name
}

func (t *TSLDLT6452007ModelTemplate) GetTSLModelLabel() string {
	return t.Label
}

func (t *TSLDLT6452007ModelTemplate) GetTSLModelType() int {
	return t.Type
}

func (t *TSLDLT6452007ModelTemplate) GetTSLModelParam() interface{} {
	return ""
}

func (t *TSLDLT6452007ModelTemplate) TSLModelPropertiesGet() interface{} {

	properties := make([]*TSLDLT6452007PropertyTemplate, 0)
	for _, v := range t.Cmd {
		properties = append(properties, v.Properties...)
	}

	return properties
}

func (t *TSLDLT6452007ModelTemplate) GetTSLDLT6452007ModelProperties() []*TSLDLT6452007PropertyTemplate {

	properties := make([]*TSLDLT6452007PropertyTemplate, 0)
	for _, v := range t.Cmd {
		properties = append(properties, v.Properties...)
	}

	return properties
}

func DeleteTSLDLT6452007(name string) error {

	for _, v := range CollectInterfaceMap.Coll {
		for _, d := range v.DeviceNodeMap {
			if d.TSL == name {
				return errors.New("采集模型[DLT645-2007]已被使用，不可以删除")
			}
		}
	}

	_, ok := TSLModelsName[name]
	if !ok {
		return errors.New("采集模型[DLT645-2007]不存在")
	}

	delete(TSLDLT6452007Map, name)
	delete(TSLModelsName, name)
	WriteTSLDLT6452007ParamToJson()

	return nil
}

func ModifyTSLDLT6452007(name string, label string) error {

	tsl, ok := TSLDLT6452007Map[name]
	if !ok {
		return errors.New("采集模型[DLT645-2007]模版不存在")
	}

	tsl.Label = label
	tsl.Type = TSLModelTypeDLT6452007
	setting.ZAPS.Infof("采集模型[DLT645-2007][%s]发生变化", name)
	eventMsg := TSLEventTemplate{
		TSL:   name,
		Topic: "modify",
	}
	_ = tsl.Event.Publish("modify", eventMsg)
	WriteTSLDLT6452007ParamToJson()

	return nil
}

func (d *TSLDLT6452007ModelTemplate) TSLModelCmdAdd(cmd TSLDLT6452007CmdTemplate) error {

	_, ok := d.Cmd[cmd.Name]
	if !ok {
		d.Cmd[cmd.Name] = &cmd
		eventMsg := TSLEventTemplate{
			TSL:   d.Name,
			Topic: "modify",
		}
		_ = d.Event.Publish("modify", eventMsg)
		dlt6452007WriteTimer.Reset(time.Second)
		return nil
	} else {
		return errors.New("命令名称已经存在")
	}
}

func (d *TSLDLT6452007ModelTemplate) TSLModelCmdDelete(cmdNames []string) error {

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
		dlt6452007WriteTimer.Reset(time.Second)
		eventMsg := TSLEventTemplate{
			TSL:   d.Name,
			Topic: "modify",
		}
		d.Event.Publish("modify", eventMsg)
		return nil
	}
	return errors.New("命令名称不存在")
}

func (d *TSLDLT6452007ModelTemplate) TSLModelCmdModify(cmd TSLDLT6452007CmdTemplate) error {

	_, ok := d.Cmd[cmd.Name]
	if !ok {
		return errors.New("命令名称不存在")
	}
	d.Cmd[cmd.Name].Label = cmd.Label
	d.Cmd[cmd.Name].BlockRulerId = cmd.BlockRulerId
	d.Cmd[cmd.Name].BlockRead = cmd.BlockRead
	d.Cmd[cmd.Name].Properties = cmd.Properties

	dlt6452007WriteTimer.Reset(time.Second)

	eventMsg := TSLEventTemplate{
		TSL:   d.Name,
		Topic: "modify",
	}
	d.Event.Publish("modify", eventMsg)

	return nil

}

func (d *TSLDLT6452007ModelTemplate) TSLModelPropertiesAdd(cmdName string, property TSLDLT6452007PropertyTemplate) error {

	cmd, ok := d.Cmd[cmdName]
	if !ok {
		return errors.New("命令名称不存在")
	}

	index := -1
	for k, v := range cmd.Properties {
		if v.Name == property.Name {
			index = k
			return errors.New("属性名称已经存在")
		}
	}
	if index == -1 {
		d.Cmd[cmdName].Properties = append(d.Cmd[cmdName].Properties, &property)
		eventMsg := TSLEventTemplate{
			TSL:   d.Name,
			Topic: "modify",
		}
		_ = d.Event.Publish("modify", eventMsg)
	}

	dlt6452007WriteTimer.Reset(time.Second)

	return nil
}

func (d *TSLDLT6452007ModelTemplate) TSLModelPropertiesDelete(cmdName string, propertiesName []string) error {

	cnt := 0
	cmd := d.Cmd[cmdName]
	for _, p := range propertiesName {
		for k := 0; k < len(cmd.Properties); k++ {
			if cmd.Properties[k].Name == p {
				cnt = 1
				cmd.Properties = append(cmd.Properties[:k], cmd.Properties[k+1:]...)
			}
		}
	}

	if cnt != 0 {
		dlt6452007WriteTimer.Reset(time.Second)
		eventMsg := TSLEventTemplate{
			TSL:   d.Name,
			Topic: "modify",
		}
		d.Event.Publish("modify", eventMsg)
		return nil
	}
	return errors.New("属性名称不存在")

}

func (d *TSLDLT6452007ModelTemplate) TSLModelPropertiesModify(cmdName string, property TSLDLT6452007PropertyTemplate) error {

	cmd, ok := d.Cmd[cmdName]
	if !ok {
		return errors.New("命令名称不存在")
	}

	for k, v := range cmd.Properties {
		if v.Name == property.Name {
			cmd.Properties[k] = &property
		}
	}

	index := -1
	for k, v := range cmd.Properties {
		if v.Name == property.Name {
			cmd.Properties[k] = &property
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
	dlt6452007WriteTimer.Reset(time.Second)

	return nil
}
