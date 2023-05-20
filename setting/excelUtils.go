package setting

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
)

func GetInt(field string) int {
	var r int
	var err error
	if r, err = strconv.Atoi(field); err != nil {
		return 0
	}
	return r
}

func GetString(field string) string {
	return field
}

func GetBool(field string) bool {
	if field == "true" || field == "TRUE" {
		return true
	} else {
		return false
	}
}

func ReadExcel(fileName string) (error, [][]string) {

	cells := make([][]string, 0)
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		ZAPS.Errorf("打开excel错误 %v", err)
		return err, nil
	}
	defer func() {
		if err := f.Close(); err != nil {
			ZAPS.Errorf("关闭excel错误 %v", err)
		}
	}()
	// 获取 Sheet1 上所有单元格
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		ZAPS.Errorf("读取excel错误 %v", err)
		return err, nil
	}
	for i := 2; i < len(rows); i++ {
		cells = append(cells, rows[i])
	}

	return nil, cells
}

func WriteExcel(fileName string, cells [][]string) error {

	f := excelize.NewFile()
	// 创建一个 sheet
	index := f.NewSheet("Sheet1")

	HStr := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W"}
	for i, v := range cells {
		for k, cell := range v {
			cellStr := HStr[k] + fmt.Sprintf("%d", i+1)
			_ = f.SetCellValue("Sheet1", cellStr, cell)
		}
	}

	// 设置文件打开后显示哪个 sheet， 0 表示 sheet1
	f.SetActiveSheet(index)
	// 保存到文件
	if err := f.SaveAs(fileName); err != nil {
		ZAPS.Errorf("写入xlsx文件错误 %v", err)
		return err
	}

	return nil
}
