package ReturnModel

type TreeDevice struct {
	Id              int          `json:"id" gorm:"primary_key"`
	Name            string       `json:"name"`
	Label           string       `json:"label"`
	DeviceType      string       `json:"deviceType"`
	ModelId         int          `json:"modelId"`
	CollInterfaceId int          `json:"collInterfaceId"`
	Addr            string       `json:"addr"`
	Data            string       `json:"data"`
	ConnectStatus   string       `json:"connectStatus"`
	ParentId        int          `json:"parentId"`
	Children        []TreeDevice `json:"children" gorm:"-"`
}

func (u *TreeDevice) TableName() string {
	return "em_device"
}
