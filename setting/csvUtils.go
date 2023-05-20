package setting

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

type CsvTable struct {
	FileName string
	Records  []CsvRecord
}

type CsvRecord struct {
	Record map[string]string
}

func (c *CsvRecord) GetInt(field string) int {
	var r int
	var err error
	if r, err = strconv.Atoi(c.Record[field]); err != nil {
		ZAPS.Errorf("获取[%s]字段错误", field)
		return 0
	}
	return r
}

func (c *CsvRecord) GetString(field string) string {
	data, ok := c.Record[field]
	if ok {
		return data
	}

	ZAPS.Errorf("获取[%s]字段错误", field)
	return ""
}

func (c *CsvRecord) GetBool(field string) bool {
	data, ok := c.Record[field]
	if ok {
		if data == "true" || data == "TRUE" {
			return true
		} else {
			return false
		}
	}
	ZAPS.Errorf("获取[%s]字段错误", field)
	return false
}

func LoadCsvCfg(filename string, title int, row int, column int) *CsvTable {
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if reader == nil {
		ZAPS.Errorf("NewReader return nil, file:", file)
		return nil
	}
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return nil
	}
	if len(records) < row {
		return nil
	}
	//行数
	recordNum := len(records)
	ZAPS.Debugf("recordNum %d", recordNum)
	var allRecords []CsvRecord
	for i := row; i < recordNum; i++ {
		record := &CsvRecord{make(map[string]string)}
		//按照行内容进行解析，因为单元格内容可能会空
		for k := column; k < len(records[i]); k++ {
			record.Record[records[title][k]] = records[i][k]
		}
		allRecords = append(allRecords, *record)
	}
	var result = &CsvTable{
		filename,
		allRecords,
	}
	return result
}
