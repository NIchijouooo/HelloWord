package models

type LimitConfigVo struct {
	DeviceLabel  string `json:"deviceLabel"`
	PropertyCode string `json:"propertyCode"`
	EnableFlag   string `json:"enableFlag"`
	ModelJson    string `json:"modelJson"`
	NotifyMin    string `json:"notifyMin"`
	NotifyMax    string `json:"notifyMax"`
	SecondaryMin string `json:"secondaryMin"`
	SecondaryMax string `json:"secondaryMax"`
	SeriousMin   string `json:"seriousMin"`
	SeriousMax   string `json:"seriousMax"`
	UrgentMin    string `json:"urgentMin"`
	UrgentMax    string `json:"urgentMax"`
}
