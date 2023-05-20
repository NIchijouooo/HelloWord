package device

import (
	"fmt"
	"gateway/setting"
	"math"
	"strconv"
	"time"

	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

//设备模板
type DeviceNodeTemplate struct {
	Index          int                     `json:"index"`          //设备偏移量
	Name           string                  `json:"name"`           //设备名称
	Label          string                  `json:"label"`          //设备标签
	Addr           string                  `json:"addr"`           //设备地址
	TSL            string                  `json:"tsl"`            //设备物模型
	LastCommRTC    string                  `json:"lastCommRTC"`    //最后一次通信时间戳
	CommTotalCnt   int                     `json:"commTotalCnt"`   //通信总次数
	CommSuccessCnt int                     `json:"commSuccessCnt"` //通信成功次数
	CurCommFailCnt int                     `json:"-"`              //当前通信失败次数
	CommStatus     string                  `json:"commStatus"`     //通信状态
	Properties     []TSLPropertiesTemplate `json:"properties"`     //属性列表
	Services       []TSLLuaServiceTemplate `json:"services"`       //服务
}

type ReceiveVariableTemplate struct {
	Index   int
	Name    string
	Label   string
	Type    string
	Value   interface{}
	Explain string
}

type ReceiveDataTemplate struct {
	Status     bool
	Properties []ReceiveVariableTemplate
}

type NodeVariableTemplate struct {
	Name  string                   `json:"name"`
	Value TSLPropertyValueTemplate `json:"value"`
}

const (
	VariableMaxCnt = 2
)

func (d *DeviceNodeTemplate) NewVariables() []TSLPropertiesTemplate {

	properties := make([]TSLPropertiesTemplate, 0)

	property := TSLPropertiesTemplate{}
	for _, v := range TSLLuaMap {
		if v.Name == d.TSL {
			for k, p := range v.Properties {
				property.Index = k
				property.Name = p.Name
				property.Label = p.Label
				property.Type = p.Type
				property.AccessMode = p.AccessMode
				properties = append(properties, property)
			}
		}
	}

	for _, v := range TSLModelS7Map {
		if v.Name == d.TSL {
			for k, p := range v.Properties {
				property.Index = k
				property.Name = p.Name
				property.Label = p.Label
				property.Type = p.Type
				property.AccessMode = p.AccessMode
				properties = append(properties, property)
			}
		}
	}

	for _, v := range TSLModbusMap {
		if v.Name == d.TSL {
			modbusProperties := v.GetTSLModbusModelProperties()
			for k, p := range modbusProperties {
				property.Index = k
				property.Name = p.Name
				property.Label = p.Label
				property.Type = p.Type
				property.AccessMode = p.AccessMode
				properties = append(properties, property)
			}
		}
	}

	for _, v := range TSLDLT6452007Map {
		if v.Name == d.TSL {
			d07Properties := v.GetTSLDLT6452007ModelProperties()
			for k, p := range d07Properties {
				property.Index = k
				property.Name = p.Name
				property.Label = p.Label
				property.Type = p.Type
				property.AccessMode = p.AccessMode
				properties = append(properties, property)
			}
		}
	}

	return properties
}

func (d *DeviceNodeTemplate) NewServices() []TSLLuaServiceTemplate {

	services := make([]TSLLuaServiceTemplate, 0)

	for _, v := range TSLLuaMap {
		if v.Name == d.TSL {
			for _, s := range v.Services {
				services = append(services, *s)
			}
		}
	}

	return services
}

func (d *DeviceNodeTemplate) GenerateGetRealVariables(lState *lua.LState, sAddr string, step int, variables string) ([]byte, bool, bool) {

	type LuaVariableMapTemplate struct {
		Status   string `json:"Status"`
		Variable []*byte
	}
	if lState != nil {
		//调用GenerateGetRealVariables
		err := lState.CallByParam(lua.P{
			Fn:      lState.GetGlobal("GenerateGetRealVariables"),
			NRet:    1,
			Protect: true,
		}, lua.LString(sAddr), lua.LNumber(step), lua.LString(variables))
		if err != nil {
			setting.ZAPS.Warnf("设备[%s]GenerateGetRealVariables 调用函数错误 %v", d.Name, err)
			return nil, false, false
		}

		//获取返回结果
		ret := lState.Get(-1)
		_, ok := ret.(*lua.LTable)
		if ok == false {
			setting.ZAPS.Warnf("设备[%s]GenerateGetRealVariables 返回结果错误 %v", d.Name, err)
			return nil, false, false
		}
		lState.Pop(1)

		LuaVariableMap := LuaVariableMapTemplate{}
		if err := gluamapper.Map(ret.(*lua.LTable), &LuaVariableMap); err != nil {
			setting.ZAPS.Warnf("设备[%s]GenerateGetRealVariables 返回结果格式化错误 %v", d.Name, err)
			return nil, false, false
		}

		result := false
		continuous := false //后续是否有报文
		nBytes := make([]byte, 0)
		if len(LuaVariableMap.Variable) > 0 {
			result = true
			for _, v := range LuaVariableMap.Variable {
				nBytes = append(nBytes, *v)
			}
			if LuaVariableMap.Status == "0" {
				continuous = false
			} else {
				continuous = true
			}
		} else {
			//setting.ZAPS.Debug("GenerateGetRealVariables 返回结果长度为0")
			result = true
		}
		return nBytes, result, continuous
	}
	return nil, false, false
}

func (d *DeviceNodeTemplate) DeviceCustomCmd(lState *lua.LState, sAddr string,
	cmdName string, cmdParam string,
	step int, variables string) ([]byte, bool, bool) {

	type LuaVariableMapTemplate struct {
		Status   string  `json:"Status"`
		Variable []*byte `json:"Variable"`
	}

	//log.Printf("cmdParam %+v\n", cmdParam)
	//setting.ZAPS.Debugf("cmdName %v", cmdName)
	if lState != nil {
		var err error
		var ret lua.LValue

		//调用DeviceCustomCmd
		err = lState.CallByParam(lua.P{
			Fn:      lState.GetGlobal("DeviceCustomCmd"),
			NRet:    1,
			Protect: true,
		}, lua.LString(sAddr),
			lua.LString(cmdName),
			lua.LString(cmdParam),
			lua.LNumber(step),
			lua.LString(variables))
		if err != nil {
			setting.ZAPS.Warnf("设备[%s]DeviceCustomCmd err %v", d.Name, err)
			return nil, false, false
		}

		//获取返回结果
		ret = lState.Get(-1)
		_, ok := ret.(*lua.LTable)
		if ok == false {
			setting.ZAPS.Warnf("设备[%s]DeviceCustomCmd 返回结果错误 %v", d.Name, err)
			return nil, false, false
		}
		lState.Pop(1)

		LuaVariableMap := LuaVariableMapTemplate{}
		if err := gluamapper.Map(ret.(*lua.LTable), &LuaVariableMap); err != nil {
			setting.ZAPS.Warnf("设备[%s]DeviceCustomCmd gluamapper.Map err %v", d.Name, err)
			return nil, false, false
		}

		result := false
		continuous := false //后续是否有报文
		if LuaVariableMap.Status == "0" {
			continuous = false
		} else {
			continuous = true
		}
		nBytes := make([]byte, 0)
		if len(LuaVariableMap.Variable) > 0 {
			result = true
			for _, v := range LuaVariableMap.Variable {
				nBytes = append(nBytes, *v)
			}
		} else {
			result = true
		}
		return nBytes, result, continuous
	}
	return nil, false, false
}

func (d *DeviceNodeTemplate) AnalysisRx(lState *lua.LState, sAddr string, variables []TSLPropertiesTemplate, rxBuf []byte, rxBufCnt int) chan ReceiveDataTemplate {

	result := make(chan ReceiveDataTemplate, 1)

	type LuaVariableMapTemplate struct {
		Status   string `json:"Status"`
		Variable []ReceiveVariableTemplate
	}
	if lState != nil {
		tbl := lua.LTable{}
		for _, v := range rxBuf {
			tbl.Append(lua.LNumber(v))
		}
		lState.SetGlobal("rxBuf", luar.New(lState, &tbl))

		//AnalysisRx
		err := lState.CallByParam(lua.P{
			Fn:      lState.GetGlobal("AnalysisRx"),
			NRet:    1,
			Protect: true,
		}, lua.LString(sAddr), lua.LNumber(rxBufCnt))
		if err != nil {
			setting.ZAPS.Warnf("设备[%s]AnalysisRx 执行接收解析错误 %v", d.Name, err)
			return result
		}

		//获取返回结果
		ret := lState.Get(-1)
		_, ok := ret.(*lua.LTable)
		if ok == false {
			setting.ZAPS.Warnf("设备[%s]AnalysisRx 返回结果错误 %v", d.Name, err)
			return result
		}
		lState.Pop(1)

		LuaVariableMap := LuaVariableMapTemplate{}

		if err := gluamapper.Map(ret.(*lua.LTable), &LuaVariableMap); err != nil {
			setting.ZAPS.Warnf("设备[%s]AnalysisRx gluamapper.Map err %v", d.Name, err)
			return result
		}

		timeNow := time.Now()
		rx := ReceiveDataTemplate{
			Status: true,
		}
		value := TSLPropertyValueTemplate{}
		if LuaVariableMap.Status == "0" {
			if len(LuaVariableMap.Variable) > 0 {
				for _, lv := range LuaVariableMap.Variable {
					for k, p := range variables {
						if lv.Name == p.Name {
							//value.Index = lv.Index
							//当值为nil时会触发溢出
							if lv.Value == nil {
								continue
							}
							value.Index = k
							switch p.Type {
							case PropertyTypeInt32:
								value.Value = (int32)(lv.Value.(float64))
							case PropertyTypeUInt32:
								value.Value = (uint32)(lv.Value.(float64))
							case PropertyTypeDouble:
								{
									dec := 0
									tsl, ok := TSLLuaMap[d.TSL]
									if !ok {
										continue
									}
									for _, v := range tsl.Properties {
										if v.Name == p.Name {
											dec = v.Decimals
										}
									}
									decimals := dec
									if decimals > 0 {
										value.Value = lv.Value.(float64) / math.Pow10(decimals)
									} else {
										value.Value, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", lv.Value.(float64)), 64)
									}
								}
							case PropertyTypeString:
								value.Value = lv.Value.(string)
							}

							value.Explain = lv.Explain
							value.TimeStamp = timeNow

							if len(variables[k].Value) < VariableMaxCnt {
								variables[k].Value = append(variables[k].Value, value)
							} else {
								variables[k].Value = variables[k].Value[1:]
								variables[k].Value = append(variables[k].Value, value)
							}
							rx.Properties = append(rx.Properties, lv)
						}
					}
				}
			}
			//setting.ZAPS.Debugf("rx %+v", rx)
			result <- rx
		}
	}
	return result
}

func (d *DeviceNodeTemplate) GetTSLModelUnit(name string) string {

	unit := ""

	model, ok := TSLLuaMap[d.TSL]
	if ok {
		for _, v := range model.Properties {
			if v.Name == name {
				unit = v.Unit
			}
		}
	}

	modelS7, ok := TSLModelS7Map[d.TSL]
	if ok {
		for _, v := range modelS7.Properties {
			if v.Name == name {
				unit = v.Unit
			}
		}
	}

	return unit
}
