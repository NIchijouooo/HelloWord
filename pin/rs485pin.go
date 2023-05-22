package pin

import (
	"fmt"
	"os"
)

const (
	Export = "/sys/class/gpio/export"
	GPIO   = "/sys/class/gpio/gpio"
)

const (
	Rs4851Name = "/dev/ttyS4"
	Rs4851Pin  = "95"

	Rs4852Name = "/dev/ttyS5"
	Rs4852Pin  = "94"
)

func Rs485PinInit(pin string) bool {
	GpioDir := GPIO + pin

	// 获取目录信息
	_, err := os.Stat(GpioDir)

	//  不存在则需要导出export
	if os.IsNotExist(err) {
		// 打开 export 文件
		file, err := os.OpenFile(Export, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("Failed to open export file:", err)
			return false
		}

		// 向 export 文件中写入 GPIO 口号
		_, err = file.WriteString(pin)
		if err != nil {
			fmt.Println("Failed to write GPIO number:", err)
			file.Close()
			return false
		}
		file.Close()
		fmt.Println("GPIO number exported successfully.")
	}

	filePath := GpioDir + "/direction"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open %s file:", filePath, err)
		return false
	}

	_, err = file.WriteString("out")
	if err != nil {
		fmt.Println("Failed to write direction -> out :", err)
		file.Close()
		return false
	}
	file.Close()

	return true
}

func writeValue(pin string, v string) bool {

	filePath := GPIO + pin + "/value"

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open %s file: %v ", filePath, err)
		return false
	}

	_, err = file.WriteString(v)
	if err != nil {
		fmt.Println("Failed to write direction -> out :", err)
		file.Close()
		return false
	}
	file.Close()

	return true
}

func Rs485xRX(pin string) {
	writeValue(pin, "0")
}

func Rs485xTX(pin string) {
	writeValue(pin, "1")
}
