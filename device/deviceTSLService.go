package device

import (
	"errors"
)

func AddTSLService(name string, Service TSLLuaServiceTemplate) error {

	tsl, ok := TSLLuaMap[name]
	if !ok {
		return errors.New("物模型模版名称不存在")
	}

	_, err := tsl.TSLLuaServicesAdd(Service)
	if err != nil {
		return err
	}

	eventMsg := TSLEventTemplate{
		TSL:   tsl.Name,
		Topic: "modify",
	}
	tsl.Event.Publish("modify", eventMsg)

	return nil
}

func ModifyTSLService(name string, Service TSLLuaServiceTemplate) error {

	tsl, ok := TSLLuaMap[name]
	if !ok {
		return errors.New("物模型模版名称不存在")
	}

	_, err := tsl.TSLLuaServicesModify(Service)
	if err != nil {
		return err
	}

	eventMsg := TSLEventTemplate{
		TSL:   tsl.Name,
		Topic: "modify",
	}
	tsl.Event.Publish("modify", eventMsg)

	return nil
}

func DeleteTSLServices(name string, ServiceNames []string) error {

	tsl, ok := TSLLuaMap[name]
	if !ok {
		return errors.New("物模型模版名称不存在")
	}

	_, err := tsl.TSLLuaServicesDelete(ServiceNames)
	if err != nil {
		return err
	}

	eventMsg := TSLEventTemplate{
		TSL:   tsl.Name,
		Topic: "modify",
	}
	tsl.Event.Publish("modify", eventMsg)

	return nil
}

func GetTSLService(name string) (error, []TSLLuaServiceTemplate) {

	services := make([]TSLLuaServiceTemplate, 0)

	tsl, ok := TSLLuaMap[name]
	if !ok {
		return errors.New("物模型模版名称不存在"), services
	}

	for _, v := range tsl.Services {
		services = append(services, *v)
	}

	return nil, services
}
