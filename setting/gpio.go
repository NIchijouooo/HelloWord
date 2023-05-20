package setting

import "os"

func GPIOWriteCmd(gpio string, value string) {

	if gpio == "" {
		return
	}
	fd, err := os.OpenFile(gpio, os.O_RDWR, 0777)
	if err != nil {
		//ZAPS.Errorf("GPIO接口[%s]打开失败 %v", gpio, err)
		return
	}
	defer fd.Close()

	if value == "" {
		return
	}
	_, err = fd.Write([]byte(value))

	if err != nil {
		//ZAPS.Errorf("GPIO接口[%s]写入失败 %v", gpio, err)
		return
	}
}
