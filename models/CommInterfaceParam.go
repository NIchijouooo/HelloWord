package models

type commInterfaceBase struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type LocalSerial struct {
	commInterfaceBase
	Param LocalSerialParam `json:"param"`
}

type ModbusRTU struct {
	LocalSerial
}

type TCPClient struct {
	commInterfaceBase
	Param TCPParam `json:"param"`
}

type TCPServer struct {
	TCPClient
}

type ModbusTCP struct {
	TCPClient
}

type UDPClient struct {
	TCPClient
}

type UDPServer struct {
	TCPClient
}

type LocalSerialParam struct {
	Name            string `json:"name" binding:"required"`
	BaudRate        string `json:"baudRate" binding:"required"`
	DataBits        string `json:"dataBits" binding:"required"`
	StopBits        string `json:"stopBits" binding:"required"`
	Parity          string `json:"parity" binding:"required"`
	TimeOut         string `json:"timeOut"`
	Interval        string `json:"interval"`
	LedModuleEnable bool   `json:"ledModuleEnable"`
	LedGPIO         string `json:"ledGPIO"`
	LedOff          string `json:"ledOff"`
	LedOn           string `json:"ledOn"`
}

type TCPParam struct {
	IP       string `json:"ip" binding:"required"`
	Port     string `json:"port" binding:"required"`
	TimeOut  string `json:"timeOut"`
	Interval string `json:"interval"`
}
