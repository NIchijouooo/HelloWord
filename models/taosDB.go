package models

import (
	"database/sql"
	"fmt"
	_ "github.com/taosdata/driver-go/v3/taosRestful"
)

var taosDB *sql.DB

func InitTaosDB() {
	var taosUri = "root:taosdata@http(td1.iot.feisjy.com:6041)/realtimedata"
	taos, err := sql.Open("taosRestful", taosUri)
	if err != nil {
		fmt.Println("failed to connect TDengine, err:", err)
		return
	}
	taosDB = taos
}
