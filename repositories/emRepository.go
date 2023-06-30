package repositories

import (
	"gateway/models"
	"gorm.io/gorm"
	"strings"
)

type EmRepository struct {
	db *gorm.DB
}

func NewEmRepository() *EmRepository {
	return &EmRepository{
		db: models.DB,
	}
}

// GetCommInterfaceProtocolById 通信协议
func (r *EmRepository) GetCommInterfaceProtocolById(id int) (*models.CommInterfaceProtocol, error) {
	var commInterfaceProtocol models.CommInterfaceProtocol
	if err := r.db.First(&commInterfaceProtocol, id).Error; err != nil {
		return nil, err
	}
	return &commInterfaceProtocol, nil
}

func (r *EmRepository) GetCommInterfaceProtocolByName(name string) (*models.CommInterfaceProtocol, error) {
	var commInterfaceProtocol models.CommInterfaceProtocol
	if err := r.db.Where("name = ?", name).First(&commInterfaceProtocol).Error; err != nil {
		return nil, err
	}
	return &commInterfaceProtocol, nil
}

func (r *EmRepository) GetAllCommInterfaceProtocols() ([]*models.CommInterfaceProtocol, error) {
	var commInterfaceProtocol []*models.CommInterfaceProtocol
	if err := r.db.Find(&commInterfaceProtocol).Error; err != nil {
		return nil, err
	}
	return commInterfaceProtocol, nil
}

// GetAllCommInterfaces 通信接口
func (r *EmRepository) GetAllCommInterfaces() ([]*models.CommInterface, error) {
	var commInterface []*models.CommInterface
	if err := r.db.Find(&commInterface).Error; err != nil {
		return nil, err
	}
	return commInterface, nil
}

func (r *EmRepository) GetCommInterfaceById(id int) (*models.CommInterface, error) {
	var commInterface models.CommInterface
	if err := r.db.First(&commInterface, id).Error; err != nil {
		return nil, err
	}
	return &commInterface, nil
}

func (r *EmRepository) GetCommInterfaceByName(name string) (*models.CommInterface, error) {
	var commInterface models.CommInterface
	if err := r.db.Where("name = ?", name).First(&commInterface).Error; err != nil {
		return nil, err
	}
	return &commInterface, nil
}

func (r *EmRepository) AddCommInterface(commInterface *models.CommInterface) error {
	return r.db.Create(commInterface).Error
}

func (r *EmRepository) UpdateCommInterface(commInterface *models.CommInterface) error {
	var emCommInterface models.CommInterface
	name := commInterface.Name
	r.db.Where("name = ?", name).First(&emCommInterface)
	commInterface.Id = emCommInterface.Id
	return r.db.Save(commInterface).Error
}

func (r *EmRepository) DelCommInterface(id int) error {
	return r.db.Delete(&models.CommInterface{}, id).Error
}

// AddCollInterface 采集接口
func (r *EmRepository) AddCollInterface(emCollInterface *models.EmCollInterface) error {
	return r.db.Create(emCollInterface).Error
}

func (r *EmRepository) UpdateCollInterface(emCollInterface *models.EmCollInterface) error {
	return r.db.Save(emCollInterface).Error
}

func (r *EmRepository) DeleteCollInterface(id int) error {
	return r.db.Delete(&models.EmCollInterface{}, id).Error
}

func (r *EmRepository) GetCollInterfaceByName(name string) (*models.EmCollInterface, error) {
	var emCollInterface models.EmCollInterface
	if err := r.db.Where("name = ?", name).First(&emCollInterface).Error; err != nil {
		return nil, err
	}
	return &emCollInterface, nil
}

/*
*
获取全部采集接口
*/
func (r *EmRepository) GetAllCollInterface() ([]*models.EmCollInterface, error) {
	var list []*models.EmCollInterface
	if err := r.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// AddEmDevice EM设备
func (r *EmRepository) AddEmDevice(emDevice *models.EmDevice) error {
	return r.db.Create(emDevice).Error
}

func (r *EmRepository) UpdateEmDevice(emDevice *models.EmDevice) error {
	return r.db.Save(emDevice).Error
}

func (r *EmRepository) DeleteEmDevice(id int) error {
	return r.db.Delete(&models.EmDevice{}, id).Error
}

func (r *EmRepository) GetEmDeviceById(id int) (*models.EmDevice, error) {
	var emDevice models.EmDevice
	if err := r.db.First(&emDevice, id).Error; err != nil {
		return nil, err
	}
	return &emDevice, nil
}

func (r *EmRepository) GetEmDeviceByName(name string) (*models.EmDevice, error) {
	var emDevice models.EmDevice
	if err := r.db.Where("name = ?", name).First(&emDevice).Error; err != nil {
		return nil, err
	}
	return &emDevice, nil
}

func (r *EmRepository) GetEmDeviceByModelId(modelId int) ([]models.EmDevice, error) {
	var emDevice []models.EmDevice
	if err := r.db.Where("model_id = ?", modelId).Find(&emDevice).Error; err != nil {
		return nil, err
	}
	return emDevice, nil
}

// AddEmDeviceModel EM设备模型
func (r *EmRepository) AddEmDeviceModel(emDeviceModel *models.EmDeviceModel) error {
	return r.db.Create(emDeviceModel).Error
}

func (r *EmRepository) GetEmDeviceModelByName(name string) (*models.EmDeviceModel, error) {
	var emDeviceModel models.EmDeviceModel
	if err := r.db.Where("name = ?", name).First(&emDeviceModel).Error; err != nil {
		return nil, err
	}
	return &emDeviceModel, nil
}

/*
*
获取全部设备模型
*/
func (r *EmRepository) GetAllEmDeviceModel() ([]*models.EmDeviceModel, error) {
	var list []*models.EmDeviceModel
	if err := r.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *EmRepository) UpdateEmDeviceModel(emDeviceModel *models.EmDeviceModel) error {
	return r.db.Save(emDeviceModel).Error
}

func (r *EmRepository) DeleteEmDeviceModel(id int) error {
	return r.db.Delete(&models.EmDeviceModel{}, id).Error
}

// AddEmDeviceModelCmd EM设备模型命令
func (r *EmRepository) AddEmDeviceModelCmd(emDeviceModelCmd *models.EmDeviceModelCmd) error {
	return r.db.Create(emDeviceModelCmd).Error
}

func (r *EmRepository) GetEmDeviceModelCmdByName(name string) (*models.EmDeviceModelCmd, error) {
	var emDeviceModelCmd models.EmDeviceModelCmd
	if err := r.db.Where("name = ?", name).First(&emDeviceModelCmd).Error; err != nil {
		return nil, err
	}
	return &emDeviceModelCmd, nil
}

func (r *EmRepository) GetEmDeviceModelCmdByModelId(modelId int) ([]models.EmDeviceModelCmd, error) {
	var emDeviceModelCmd []models.EmDeviceModelCmd
	if err := r.db.Where("device_model_id = ?", modelId).Find(&emDeviceModelCmd).Error; err != nil {
		return nil, err
	}
	return emDeviceModelCmd, nil
}

func (r *EmRepository) UpdateEmDeviceModelCmd(emDeviceModelCmd *models.EmDeviceModelCmd) error {
	return r.db.Save(emDeviceModelCmd).Error
}

func (r *EmRepository) DeleteEmDeviceModelCmd(id int) error {
	return r.db.Delete(&models.EmDeviceModelCmd{}, id).Error
}

// AddEmDeviceModelCmdParam EM设备模型命令参数
func (r *EmRepository) AddEmDeviceModelCmdParam(emDeviceModelCmdParam *models.EmDeviceModelCmdParam) error {
	return r.db.Create(emDeviceModelCmdParam).Error
}

func (r *EmRepository) GetEmDeviceModelCmdParamByName(name string) (*models.EmDeviceModelCmdParam, error) {
	var emDeviceModelCmdParam models.EmDeviceModelCmdParam
	if err := r.db.Where("name = ?", name).First(&emDeviceModelCmdParam).Error; err != nil {
		return nil, err
	}
	return &emDeviceModelCmdParam, nil
}

func (r *EmRepository) UpdateEmDeviceModelCmdParam(emDeviceModelCmdParam *models.EmDeviceModelCmdParam) error {
	return r.db.Save(emDeviceModelCmdParam).Error
}

func (r *EmRepository) DeleteEmDeviceModelCmdParam(id int) error {
	return r.db.Delete(&models.EmDeviceModelCmdParam{}, id).Error
}

func (r *EmRepository) GetEmDeviceModelCmdParamByCmdId(cmdId int) ([]models.EmDeviceModelCmdParam, error) {
	var emDeviceModelCmdParam []models.EmDeviceModelCmdParam
	if err := r.db.Where("device_model_cmd_id = ?", cmdId).Find(&emDeviceModelCmdParam).Error; err != nil {
		return nil, err
	}
	return emDeviceModelCmdParam, nil
}

// GetEmDeviceModelCmdParamListByName 根据设备获取模型列表
func (r *EmRepository) GetEmDeviceModelCmdParamListByName(name string) ([]models.EmDeviceModelCmdParam, error) {
	var emDeviceModelCmdParam []models.EmDeviceModelCmdParam
	if err := r.db.Where("name = ?", name).Find(&emDeviceModelCmdParam).Error; err != nil {
		return nil, err
	}
	return emDeviceModelCmdParam, nil
}

// 根据设备id获取模型列表
func (r *EmRepository) GetEmDeviceModelCmdParamListByDeviceId(deviceId int) ([]models.EmDeviceModelCmdParam, error) {
	var emDeviceModelCmdParamList []models.EmDeviceModelCmdParam
	if err := r.db.Table("em_device_model_cmd_param as param").
		Select("param.id, param.device_model_cmd_id, param.name, param.label,param.data,param.iot_data_type").
		Joins("left join em_device_model_cmd as cmd on param.device_model_cmd_id = cmd.id").
		Joins("left join em_device_model as model on cmd.device_model_id = model.id").
		Joins("left join em_device as device on device.model_id = model.id").
		Where("device.id=?", deviceId).Scan(&emDeviceModelCmdParamList).Error; err != nil {
		return nil, err
	}
	return emDeviceModelCmdParamList, nil
}

// 根据设备id获取测点列表
func (r *EmRepository) GetEmDeviceModelCmdParamListByDeviceIdCodes(deviceId int, codes []string) ([]models.EmDeviceModelCmdParam, error) {
	codeString := "'" + strings.Join(codes, "','") + "'"
	var emDeviceModelCmdParamList []models.EmDeviceModelCmdParam
	if err := r.db.Table("em_device_model_cmd_param as param").
		Select("param.id, param.device_model_cmd_id, param.name, param.label,param.data,param.unit").
		Joins("left join em_device_model_cmd as cmd on param.device_model_cmd_id = cmd.id").
		Joins("left join em_device_model as model on cmd.device_model_id = model.id").
		Joins("left join em_device as device on device.model_id = model.id").
		Where("device.id=? and param.name in ("+codeString+")", deviceId).Scan(&emDeviceModelCmdParamList).Error; err != nil {
		return nil, err
	}
	return emDeviceModelCmdParamList, nil
}
func (r *EmRepository) GetCodesListByDeviceIdAndYxYc(deviceId int, iotDataType []string, like string) ([]models.EmDeviceModelCmdParam, error) {
	var emDeviceModelCmdParamList []models.EmDeviceModelCmdParam
	query := r.db.Table("em_device_model_cmd_param as param").
		Joins("left join em_device_model_cmd as cmd on param.device_model_cmd_id = cmd.id").
		Joins("left join em_device_model as model on cmd.device_model_id = model.id").
		Joins("left join em_device as device on device.model_id = model.id").
		Where("device.id=?", deviceId)
	if iotDataType != nil {
		query.Where("param.iot_data_type in (?)", iotDataType)
	}
	if like != "" && len(like) > 0 {
		query.Where(like)
	}
	if err := query.Select("param.id, param.device_model_cmd_id, param.name, param.label,param.data,param.iot_data_type,param.unit").Scan(&emDeviceModelCmdParamList).Error; err != nil {
		return nil, err
	}
	return emDeviceModelCmdParamList, nil
}

//设备包括设备台账
func (r *EmRepository) GetDeviceAndEAList(param models.DevicePageParam) ([]models.EmDeviceParamVO, int64, error) {
	var (
		emDeviceVOList []models.EmDeviceParamVO
		emDeviceList []models.EmDevice
		total        int64
	)
	//处理分页
	page := 1
	pageSize := 10
	if param.PageNum > 0 {
		page = param.PageNum
	}
	if param.PageSize > 0 {
		pageSize = param.PageSize
	}
	offset := (page - 1) * pageSize
	//拼接sql
	query := r.db.Model(&models.EmDevice{})
	query.Where("1 = 1")
	if param.KeyWord != "" {
		// 构建查询条件
		//condition1 := query.Where("name LIKE ?", "%"+param.KeyWork+"%")
		//condition2 := query.Where("label LIKE ?", "%"+param.KeyWork+"%")
		//// 使用Or方法进行条件查询
		//query.Where(condition1.Or(condition2))
		query = query.Where("name LIKE ? OR label LIKE ?", "%"+param.KeyWord+"%", "%"+param.KeyWord+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&emDeviceList).Error; err != nil {
		return nil, 0, err
	}

	dictRepo := NewDictDataRepository()
	dictDataList, _ := dictRepo.GetListByDiceType("device_type")
	m := make(map[string]string)

	// 遍历数组元素并将其添加到map中
	for _, item := range dictDataList {
		m[item.DictValue] = item.DictLabel
	}

	for _, emDevice := range emDeviceList {
		vo := models.EmDeviceParamVO{
			Id:				emDevice.Id,
			Name:			emDevice.Name,
			Label:			emDevice.Label,
			DeviceType:		emDevice.DeviceType,
			DeviceTypeKey:	m[emDevice.DeviceType],
		}

		ea, _ := r.GetEmDeviceEquipmentAccountByDevId(emDevice.Id)
		if ea.ID > 0{
			vo.Manufacturer = ea.Manufacturer
			vo.Polarity = ea.Polarity
			vo.FactoryModel = ea.FactoryModel
		}
		emDeviceVOList = append(emDeviceVOList, vo)
	}


	return emDeviceVOList, total, nil
}

//设备
func (r *EmRepository) GetDeviceList(deviceType string) ([]models.EmDevice, error) {
	var (
		emDeviceList []models.EmDevice
	)
	query := r.db.Model(&models.EmDevice{})
	if err := query.Where("device_type = ?", deviceType).Find(&emDeviceList).Error; err != nil {
		return nil, err
	}
	return emDeviceList, nil
}


// 获取单个设备台账
func (r *EmRepository) GetEmDeviceEquipmentAccountByDevId(deviceId int) (models.DeviceEquipmentAccountInfo, error) {
	var ea models.DeviceEquipmentAccountInfo
	query := r.db.Model(&models.DeviceEquipmentAccountInfo{})
	err := query.Where("device_id = ?", deviceId).First(&ea).Error
	return ea, err
}

// 修改字典类型
func (r *EmRepository) UpdateEmDeviceEquipmentAccountByDevId(ea *models.DeviceEquipmentAccountInfo) error {
	return r.db.Save(ea).Error
}
// 新增字典类型
func (r *EmRepository) CreateEmDeviceEquipmentAccountByDevId(ea *models.DeviceEquipmentAccountInfo) error {
	return r.db.Create(ea).Error
}
