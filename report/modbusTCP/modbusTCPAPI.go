package modbusTCP

import (
	"errors"
	"fmt"
	"gateway/setting"
	"gateway/utils"
)

func (r *ReportServiceMBTCPListTemplate) AddReportService(param ReportServiceMBTCPParamTemplate) error {

	index := -1
	for k, v := range r.ServiceList {
		if v.GWParam.ServiceName == param.ServiceName {
			index = k
			break
		}
	}

	if index != -1 {
		return errors.New("服务名称已经存在")
	}

	param.Param.CoilStatusRegisters = make(map[string]MBTCPRegisterTemplate)
	param.Param.InputStatusRegisters = make(map[string]MBTCPRegisterTemplate)
	param.Param.HoldingRegisters = make(map[string]MBTCPRegisterTemplate)
	param.Param.InputRegisters = make(map[string]MBTCPRegisterTemplate)

	mbTCP := &ReportServiceMBTCPTemplate{
		Index:   len(r.ServiceList),
		GWParam: param,
	}
	r.ServiceList = append(r.ServiceList, mbTCP)

	//创建新上报服务
	mbTCP.NewReportServiceMBTCP()

	ReportServiceMBTCPWriteParam()

	return nil
}

func (r *ReportServiceMBTCPListTemplate) ModifyReportService(param ReportServiceMBTCPParamTemplate) error {

	index := -1
	for k, v := range r.ServiceList {
		if v.GWParam.ServiceName == param.ServiceName {
			index = k
			break
		}
	}

	if index == -1 {
		return errors.New("服务名称不存在")
	}

	//旧上报服务退出
	r.ServiceList[index].CancelFunc()

	r.ServiceList[index].GWParam.IP = param.IP
	r.ServiceList[index].GWParam.Port = param.Port
	r.ServiceList[index].GWParam.Param.SlaveID = param.Param.SlaveID
	r.ServiceList[index].GWParam.Param.CoilStatusRegisterStart = param.Param.CoilStatusRegisterStart
	r.ServiceList[index].GWParam.Param.CoilStatusRegisterCnt = param.Param.CoilStatusRegisterCnt
	r.ServiceList[index].GWParam.Param.InputStatusRegisterStart = param.Param.InputStatusRegisterStart
	r.ServiceList[index].GWParam.Param.InputStatusRegisterCnt = param.Param.InputStatusRegisterCnt
	r.ServiceList[index].GWParam.Param.HoldingRegisterStart = param.Param.HoldingRegisterStart
	r.ServiceList[index].GWParam.Param.HoldingRegisterCnt = param.Param.HoldingRegisterCnt
	r.ServiceList[index].GWParam.Param.InputRegisterStart = param.Param.InputRegisterStart
	r.ServiceList[index].GWParam.Param.InputRegisterCnt = param.Param.InputRegisterCnt

	ReportServiceMBTCPWriteParam()

	//创建新上报服务
	r.ServiceList[index].NewReportServiceMBTCP()

	return nil
}

func (r *ReportServiceMBTCPListTemplate) DeleteReportService(serviceName string) {

	for k, v := range r.ServiceList {
		if v.GWParam.ServiceName == serviceName {
			//协程退出
			v.CancelFunc()

			r.ServiceList = append(r.ServiceList[:k], r.ServiceList[k+1:]...)
			ReportServiceMBTCPWriteParam()

			return
		}
	}
}

func (r *ReportServiceMBTCPParamTemplate) AddRegister(reg MBTCPRegisterTemplate) error {
	switch reg.RegType {
	case MBRegTypeCoilStatus:
		{

		}
	case MBRegTypeInputStatus:
		{

		}
	case MBRegTypeHoldingRegister:
		{
			_, ok := r.Param.HoldingRegisters[reg.RegName]
			if ok {
				return errors.New("寄存器名称已存在")
			}

			if len(r.Param.HoldingRegisters) == 0 {
				r.Param.HoldingRegisters = make(map[string]MBTCPRegisterTemplate)
			}

			reg.Index = len(r.Param.HoldingRegisters)
			r.Param.HoldingRegisters[reg.RegName] = reg

			ReportServiceMBTCPWriteParam()

			return nil
		}
	case MBRegTypeInputRegister:
		{

		}
	}

	return errors.New("寄存器类型不存在")
}

func (r *ReportServiceMBTCPParamTemplate) ModifyRegister(reg MBTCPRegisterTemplate) error {
	switch reg.RegType {
	case MBRegTypeCoilStatus:
		{

		}
	case MBRegTypeInputStatus:
		{

		}
	case MBRegTypeHoldingRegister:
		{
			hReg, ok := r.Param.HoldingRegisters[reg.RegName]
			if !ok {
				return errors.New("寄存器名称不存在")
			}
			index := hReg.Index
			hReg = reg
			hReg.Index = index
			r.Param.HoldingRegisters[reg.RegName] = hReg

			ReportServiceMBTCPWriteParam()

			return nil
		}
	case MBRegTypeInputRegister:
		{

		}
	}

	return errors.New("寄存器类型不存在")
}

func (r *ReportServiceMBTCPParamTemplate) DeleteRegisters(regs []string) error {

	for _, v := range regs {
		_, ok := r.Param.CoilStatusRegisters[v]
		if ok {
			delete(r.Param.CoilStatusRegisters, v)
			ReportServiceMBTCPWriteParam()
		}
		_, ok = r.Param.InputStatusRegisters[v]
		if ok {
			delete(r.Param.InputStatusRegisters, v)
			ReportServiceMBTCPWriteParam()
		}
		_, ok = r.Param.HoldingRegisters[v]
		if ok {
			delete(r.Param.HoldingRegisters, v)
			ReportServiceMBTCPWriteParam()
		}
		_, ok = r.Param.InputRegisters[v]
		if ok {
			delete(r.Param.InputRegisters, v)
			ReportServiceMBTCPWriteParam()
		}
	}

	return nil
}

func (r *ReportServiceMBTCPParamTemplate) GetRegisters() []MBTCPRegisterTemplate {

	regs := make([]MBTCPRegisterTemplate, 0)

	for _, v := range r.Param.CoilStatusRegisters {
		regs = append(regs, v)
	}

	for _, v := range r.Param.InputStatusRegisters {
		regs = append(regs, v)
	}

	for _, v := range r.Param.HoldingRegisters {
		regs = append(regs, v)
	}

	for _, v := range r.Param.InputRegisters {
		regs = append(regs, v)
	}

	return regs
}

func (r *ReportServiceMBTCPParamTemplate) ExportRegistersToXlsx() (bool, string) {

	//创建文件
	utils.DirIsExist("./tmp")
	fileName := "./tmp/" + r.ServiceName + ".xlsx"

	cellData := [][]string{
		{"寄存器名称", "寄存器标签", "属性类型", "采集接口名称", "设备名称", "属性名称", "寄存器地址", "寄存器数量", "寄存器类型", "寄存器解析规则"},
		{"regName", "label", "propertyType", "collName", "nodeName", "propertyName", "regAddr", "regCnt", "regType", "rule"},
	}

	regs := r.GetRegisters()
	for _, v := range regs {
		param := make([]string, 0)
		param = append(param, v.RegName)
		param = append(param, v.Label)
		if v.PropertyType == 0 {
			param = append(param, "uint32")
		} else if v.PropertyType == 1 {
			param = append(param, "int32")
		} else if v.PropertyType == 2 {
			param = append(param, "double")
		} else {
			param = append(param, "double")
		}
		param = append(param, v.CollName)
		param = append(param, v.NodeName)
		param = append(param, v.PropertyName)
		param = append(param, fmt.Sprintf("%d", v.RegAddr))
		param = append(param, fmt.Sprintf("%d", v.RegCnt))
		if v.RegType == MBRegTypeCoilStatus {
			param = append(param, "coilStatus")
		} else if v.RegType == MBRegTypeInputStatus {
			param = append(param, "inputStatus")
		} else if v.RegType == MBRegTypeHoldingRegister {
			param = append(param, "holdingRegister")
		} else if v.RegType == MBRegTypeInputRegister {
			param = append(param, "inputRegister")
		} else {
			param = append(param, "double")
		}
		param = append(param, v.Rule)
		cellData = append(cellData, param)
	}

	//写入数据
	err := setting.WriteExcel(fileName, cellData)
	if err != nil {
		return false, ""
	}

	return true, fileName
}
