package models

import (
	"database/sql"
	"fmt"

	_ "github.com/taosdata/driver-go/v3/taosRestful"
)

var TaosDB *sql.DB

func InitTaosDB() {
	var taosUri = "root:taosdata@http(192.168.31.162:6041)/realtimedata"
	taos, err := sql.Open("taosRestful", taosUri)
	if err != nil {
		fmt.Println("failed to connect TDengine, err:", err)
		return
	}
	TaosDB = taos
}
