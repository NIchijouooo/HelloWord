package models

import (
	"gateway/models"
	"gorm.io/gorm"
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

func (r *EmRepository) GetCollInterfaceByName(name string) (*models.EmCollInterface, error) {
	var emCollInterface models.EmCollInterface
	if err := r.db.Where("name = ?", name).First(&emCollInterface).Error; err != nil {
		return nil, err
	}
	return &emCollInterface, nil
}

// AddEmDevice EM设备
func (r *EmRepository) AddEmDevice(emDevice *models.EmDevice) error {
	return r.db.Create(emDevice).Error
}

func (r *EmRepository) UpdateEmDevice(emDevice *models.EmDevice) error {
	return r.db.Save(emDevice).Error
}

func (r *EmRepository) DelEmDevice(emDevice *models.EmDevice) error {
	return r.db.Delete(emDevice).Error
}

func (r *EmRepository) GetEmDeviceById(id int) (*models.EmDevice, error) {
	var emDevice models.EmDevice
	if err := r.db.First(&emDevice, id).Error; err != nil {
		return nil, err
	}
	return &emDevice, nil
}

// AddEmDeviceModel EM设备模型
func (r *EmRepository) AddEmDeviceModel(emDeviceModel *models.EmDeviceModel) error {
	return r.db.Create(emDeviceModel).Error
}

func (r *EmRepository) GetEmDeviceByName(name string) (*models.EmDevice, error) {
	var emDevice models.EmDevice
	if err := r.db.Where("name = ?", name).First(&emDevice).Error; err != nil {
		return nil, err
	}
	return &emDevice, nil
}

func (r *EmRepository) GetEmDeviceModelByName(name string) (*models.EmDeviceModel, error) {
	var emDeviceModel models.EmDeviceModel
	if err := r.db.Where("name = ?", name).First(&emDeviceModel).Error; err != nil {
		return nil, err
	}
	return &emDeviceModel, nil
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

//根据设备获取模型列表
func (r *EmRepository) GetEmDeviceModelCmdParamListByName(name string) ([]models.EmDeviceModelCmdParam, error) {
	var emDeviceModelCmdParam []models.EmDeviceModelCmdParam
	if err := r.db.Where("name = ?", name).Find(&emDeviceModelCmdParam).Error; err != nil {
		return nil, err
	}
	return emDeviceModelCmdParam, nil
}
