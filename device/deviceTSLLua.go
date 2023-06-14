package device

import (
	"encoding/json"
	"gateway/device/eventBus"
	"gateway/setting"
	"gateway/utils"
	"github.com/pkg/errors"
	lua "github.com/yuin/gopher-lua"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 物模型 Thing Specification Language
type TSLLuaModelTemplate struct {
	Index      int                       `json:"index"`
	Name       string                    `json:"name"`  //名称，只可以是字母和数字的组合
	Label      string                    `json:"label"` //名称标签
	Type       int                       `json:"type"`  //模型类型
	Plugin     TSLLuaPluginTemplate      `json:"plugin"`
	Properties []*TSLLuaPropertyTemplate `json:"properties"` //属性
	Services   []*TSLLuaServiceTemplate  `json:"services"`   //服务
	Event      eventBus.Bus              `json:"-"`          //事件队列
}

type TSLLuaPluginTemplate struct {
	Name    string `json:"name"`
	Label   string `json:"label"`
	Version string `json:"version"`
	Author  string `json:"author"`
	Date    string `json:"date"`
	Message string `json:"message"`
}

type TSLLuaPropertyTemplate struct {
	Name       string                      `json:"name"`       //属性名称，只可以是字母和数字的组合
	Label      string                      `json:"label"`      //属性解释
	AccessMode int                         `json:"accessMode"` //读写属性
	Type       int                         `json:"type"`       //类型 uint32 int32 double string
	Unit       string                      `json:"unit"`       //单位
	Decimals   int                         `json:"decimals"`   //小数位数
	Params     TSLLuaPropertyParamTemplate `json:"params"`
}

type TSLLuaPropertyParamTemplate struct {
	Min             string `json:"min"`             //最小
	Max             string `json:"max"`             //最大
	MinMaxAlarm     bool   `json:"minMaxAlarm"`     //范围报警
	Step            string `json:"step"`            //步长
	StepAlarm       bool   `json:"stepAlarm"`       //阶跃报警
	DataLength      string `json:"dataLength"`      //字符串长度
	DataLengthAlarm bool   `json:"dataLengthAlarm"` //字符长度报警
}

type TSLLuaServiceTemplate struct {
	Name     string                 `json:"name"`     //服务名称
	Label    string                 `json:"explain"`  //服务名称说明
	CallType int                    `json:"callType"` //服务调用方式
	Params   map[string]interface{} `json:"params"`   //服务参数
}

var TSLLuaMap = make(map[string]*TSLLuaModelTemplate)
var luaWriteTimer *time.Timer

func TSLLuaInit() {

	ReadTSLLuaParamFromJson()
	luaWriteTimer = time.AfterFunc(time.Second, func() {
		WriteTSLLuaParamToJson()
	})
	luaWriteTimer.Stop()
}

func ReadTSLLuaParamFromJson() bool {

	utils.DirIsExist("./selfpara")

	data, err := utils.FileRead("./selfpara/deviceTSLParam.json")
	if err != nil {
		setting.ZAPS.Debugf("打开物模型配置json文件失败 %v", err)
		return false
	}

	err = json.Unmarshal(data, &TSLLuaMap)
	if err != nil {
		setting.ZAPS.Errorf("物模型配置json文件格式化失败 %v", err)
		for _, v := range TSLLuaMap {
			v.Event = eventBus.NewBus()
		}
		return false
	}

	for _, v := range TSLLuaMap {
		v.Event = eventBus.NewBus()
	}

	setting.ZAPS.Debugf("打开物模型配置json文件成功")
	return true
}

func WriteTSLLuaParamToJson() {

	utils.DirIsExist("./selfpara")

	sJson, _ := json.Marshal(TSLLuaMap)
	err := utils.FileWrite("./selfpara/deviceTSLParam.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("物模型配置json文件写入失败")
		return
	}
	setting.ZAPS.Infof("物模型配置json文件写入成功")
}

func NewTSLLua(name string, label string, modelType int, plugin TSLLuaPluginTemplate) *TSLLuaModelTemplate {
	return &TSLLuaModelTemplate{
		Index:      len(TSLLuaMap),
		Name:       name,
		Label:      label,
		Type:       modelType,
		Plugin:     plugin,
		Properties: make([]*TSLLuaPropertyTemplate, 0),
		Services:   make([]*TSLLuaServiceTemplate, 0),
		Event:      eventBus.NewBus(),
	}
}

func (t *TSLLuaModelTemplate) GetTSLModelIndex() int {
	return t.Index
}

func (t *TSLLuaModelTemplate) GetTSLModelName() string {
	return t.Name
}

func (t *TSLLuaModelTemplate) GetTSLModelLabel() string {
	return t.Label
}

func (t *TSLLuaModelTemplate) GetTSLModelType() int {
	return t.Type
}

func (t *TSLLuaModelTemplate) GetTSLModelParam() interface{} {
	return t.Plugin
}

func AddTSLModelLua(name string, label string, modelType int, plugin TSLLuaPluginTemplate) error {

	_, ok := TSLLuaMap[name]
	if !ok {
		TSLLuaMap[name] = NewTSLLua(name, label, modelType, plugin)
		WriteTSLLuaParamToJson()
		return nil
	}

	return errors.New("物模型模版已经存在")
}

func DeleteTSLModelLua(name string) error {

	for _, v := range CollectInterfaceMap.Coll {
		for _, d := range v.DeviceNodeMap {
			if d.TSL == name {
				return errors.New("物模型已被使用，不可以删除")
			}
		}
	}

	_, ok := TSLModelsName[name]
	if !ok {
		return errors.New("物模型不存在")
	}

	delete(TSLLuaMap, name)
	delete(TSLModelsName, name)
	WriteTSLLuaParamToJson()

	return nil
}

func ModifyTSLModelLua(name string, label string, plugin json.RawMessage) error {

	tsl, ok := TSLLuaMap[name]
	if !ok {
		return errors.New("物模型模版不存在")
	}

	param := TSLLuaPluginTemplate{}

	err := json.Unmarshal(plugin, &param)
	if err != nil {
		return errors.New("物模型模版参数不正确")
	}

	tsl.Label = label
	tsl.Type = TSLModelTypeLua
	tsl.Plugin = param
	setting.ZAPS.Infof("物模型[%s]Plugin[%s]插件发生变化", name, plugin)
	eventMsg := TSLEventTemplate{
		TSL:   name,
		Topic: "modify",
	}
	_ = tsl.Event.Publish("modify", eventMsg)
	WriteTSLLuaParamToJson()

	return nil
}

// 遍历plugin
func DeviceTSLTraversePlugin(path string, fileName []string) ([]string, error) {

	rd, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("readDir err,", err)
		return fileName, err
	}

	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := path + "/" + fi.Name()
			fileName, _ = DeviceTSLTraversePlugin(fullDir, fileName)
		} else {
			fullName := path + "/" + fi.Name()
			if strings.Contains(fi.Name(), ".json") {
				fileName = append(fileName, fullName)
			} else if strings.Contains(fi.Name(), ".lua") {
				fileName = append(fileName, fullName)
			}
		}
	}

	return fileName, nil
}

func (d *TSLLuaModelTemplate) DeviceTSLExportPlugin(pluginName string) (bool, string) {

	//遍历文件
	pluginPath := "./plugin/" + pluginName
	fileNameMap := make([]string, 0)
	fileNameMap, _ = DeviceTSLTraversePlugin(pluginPath, fileNameMap)

	_ = utils.CompressFilesToZip(fileNameMap, "./tmp/"+pluginName+".zip")

	return true, "./tmp/" + pluginName + ".zip"
}

func (d *TSLLuaModelTemplate) DeviceTSLOpenPlugin() (error, *lua.LState) {

	lState := &lua.LState{}
	status := false

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	//遍历json和lua文件
	utils.DirIsExist("./plugin")
	pluginPath := exeCurDir + "/plugin"
	fileInfoMap, err := ioutil.ReadDir(pluginPath)
	if err != nil {
		setting.ZAPS.Errorf("打开plugin目录失败 %v", err)
		return err, nil
	}
	for _, v := range fileInfoMap {
		//文件夹并且文件名字和物模型plugin相同
		if (v.IsDir() == true) && (v.Name() == d.Plugin.Name) {
			fileDirName := pluginPath + "/" + v.Name()
			fileMap, _ := ioutil.ReadDir(fileDirName)
			index := -1
			for k, f := range fileMap {
				fileFullName := fileDirName + "/" + f.Name()
				if strings.Contains(f.Name(), ".lua") {
					//lua文件和设备模版名字一样
					if strings.EqualFold(f.Name(), d.Plugin.Name+".lua") == true {
						index = k
						lState, err = setting.LuaOpenFile(fileFullName)
						if err != nil {
							setting.ZAPS.Errorf("Lua主文件[%s]打开失败 %v", f.Name(), err)
							continue
						} else {
							setting.ZAPS.Debugf("Lua主文件[%s]打开成功", f.Name())
							status = true
						}
						lState.SetGlobal("GetCRCModbus", lState.NewFunction(setting.GetCRCModbus))
						lState.SetGlobal("CheckCRCModbus", lState.NewFunction(setting.CheckCRCModbus))
						lState.SetGlobal("GetCRCModbusLittleEndian", lState.NewFunction(setting.GetCRCModbusLittleEndian))
						lState.SetGlobal("GetCRCXmodem", lState.NewFunction(setting.GetCRCXmodem))
						lState.SetGlobal("GoLog", lState.NewFunction(setting.GoLog))
						break
					}
				}
			}
			if index == -1 {
				continue
			}

			for _, f := range fileMap {
				fileFullName := fileDirName + "/" + f.Name()
				if strings.Contains(f.Name(), ".lua") {
					//lua文件和设备模版名字不一样
					if strings.Contains(f.Name(), d.Plugin.Name) == false {
						err = lState.DoFile(fileFullName)
						if err != nil {
							setting.ZAPS.Errorf("Lua子文件[%s]打开失败 %v", f.Name(), err)
						} else {
							setting.ZAPS.Debugf("Lua子文件[%s]打开成功", f.Name())
						}
					}
				}
			}
		}
	}

	if status == false {
		lState = nil
	}

	return nil, lState
}

func (d *TSLLuaModelTemplate) TSLModelPropertiesGet() interface{} {
	return d.Properties
}

func (d *TSLLuaModelTemplate) TSLModelPropertiesAdd(property json.RawMessage) error {

	luaProperty := TSLLuaPropertyTemplate{}

	err := json.Unmarshal(property, &luaProperty)
	if err != nil {
		return err
	}

	index := -1
	for k, v := range d.Properties {
		if v.Name == luaProperty.Name {
			d.Properties[k] = &luaProperty
			index = k
		}
	}
	if index == -1 {
		d.Properties = append(d.Properties, &luaProperty)
		eventMsg := TSLEventTemplate{
			TSL:   d.Name,
			Topic: "modify",
		}
		d.Event.Publish("modify", eventMsg)
	}

	luaWriteTimer.Reset(time.Second)

	return nil
}

func (d *TSLLuaModelTemplate) TSLModelPropertiesDelete(propertiesName []string) error {

	cnt := 0
	for _, p := range propertiesName {
		for k := 0; k < len(d.Properties); k++ {
			if d.Properties[k].Name == p {
				cnt = 1
				d.Properties = append(d.Properties[:k], d.Properties[k+1:]...)
			}
		}
	}
	if cnt != 0 {
		WriteTSLLuaParamToJson()
		eventMsg := TSLEventTemplate{
			TSL:   d.Name,
			Topic: "modify",
		}
		d.Event.Publish("modify", eventMsg)
		return nil
	}
	return errors.New("属性名称不存在")

}

func (d *TSLLuaModelTemplate) TSLModelPropertiesModify(property json.RawMessage) error {

	luaProperty := TSLLuaPropertyTemplate{}
	err := json.Unmarshal(property, &luaProperty)
	if err != nil {
		return err
	}

	properties := make([]*TSLLuaPropertyTemplate, 0)
	for _, v := range d.Properties {
		properties = append(properties, v)
	}
	index := -1
	for k, v := range properties {
		if v.Name == luaProperty.Name {
			properties[k] = &luaProperty
			index = k
		}
	}
	if index != -1 {
		d.Properties = properties
		WriteTSLLuaParamToJson()

		eventMsg := TSLEventTemplate{
			TSL:   d.Name,
			Topic: "modify",
		}
		d.Event.Publish("modify", eventMsg)

		return nil
	}

	return errors.New("属性名称不存在")
}

func (d *TSLLuaModelTemplate) TSLLuaServicesAdd(service TSLLuaServiceTemplate) (int, error) {

	index := -1
	for k, v := range d.Services {
		if v.Name == service.Name {
			d.Services[k] = &service
			index = k
		}
	}
	if index == -1 {
		d.Services = append(d.Services, &service)
	}
	WriteTSLLuaParamToJson()

	return 0, nil
}

func (d *TSLLuaModelTemplate) TSLLuaServicesDelete(servicesName []string) (int, error) {

	cnt := 0
	for _, p := range servicesName {
		for k := 0; k < len(d.Services); k++ {
			if d.Services[k].Name == p {
				setting.ZAPS.Debugf("Services index %v", k)
				cnt = 1
				d.Services = append(d.Services[:k], d.Services[k+1:]...)
				setting.ZAPS.Debugf("Services %v", d.Services)
			}
		}
	}
	if cnt != 0 {
		WriteTSLLuaParamToJson()
		return 0, nil
	}
	return 1, errors.New("servicesName is not exist")

}

func (d *TSLLuaModelTemplate) TSLLuaServicesModify(service TSLLuaServiceTemplate) (int, error) {

	for k, v := range d.Services {
		if v.Name == service.Name {
			d.Services[k] = &service
			WriteTSLLuaParamToJson()
			return 0, nil
		}
	}
	return 1, errors.New("service is not exist")
}
