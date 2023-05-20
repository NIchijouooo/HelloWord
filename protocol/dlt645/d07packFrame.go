package dlt645

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
)

type D07PackFrame struct {
	RulerId  string /* 规约标签符 C0 C1  型如 0x40E3 */
	CtrlCode byte   /* 控制码 */
	DataLen  byte   /* 数据域字节数 包括规约和其它数据 */
	Address  string /* 地址 */
	Data     []byte /* 数据 */
}

//PackD07FrameByData 根据结构体组包，返回包的切片和错误码(正确返回 0 其它为错误类型)
func (d *D07PackFrame) PackD07FrameByData() ([]byte, error) {
	ucDi := make([]byte, 4)
	dataBuffer := make([]byte, 0)
	outBuffer := make([]byte, 0)

	if nil == d {
		fmt.Println("结构体源数据异常")
	}
	//fmt.Printf("结构体源数据[%v]\r\n", d)

	/* 1 准备地址域数据 */
	address, _ := D07Str2BCD(d.Address, 6)
	//fmt.Printf("Address[%s][%x]\n", d.Address, address)
	if len(address) != 6 || len(d.Address) != 12 {
		fmt.Println("地址域数据异常")
	}

	/* 2 准备数据标识数据 */
	rulerId, eer := strconv.ParseInt(d.RulerId, 16, 32) //将字符的数据标识转换成16进制的32位数据
	if eer != nil && rulerId >= 0 {
		fmt.Println("数据标识异常")
		return nil, errors.New("数据标识异常")
	}
	binary.BigEndian.PutUint32(ucDi, uint32(rulerId)) //把数据放置到切片中
	//fmt.Printf("ucDi[%s][%v]\r\n", d.RulerId, ucDi)
	for i, num := range ucDi {
		ucDi[i] = num + 0x33
	}

	/* 3 准备数据域数据 */
	//if len(d.Data)+4 != int(d.DataLen) {
	//	fmt.Println("数据域长度不对")
	//	return nil, errors.New("数据域长度不对")
	//}

	for _, num := range d.Data {
		dataBuffer = append(dataBuffer, num+0x33)
	}
	//fmt.Printf("dataBuffer[%v]\r\n", dataBuffer)

	/* 4 开始封帧 */
	outBuffer = append(outBuffer, 0x68)
	outBuffer = append(outBuffer, address...)
	outBuffer = append(outBuffer, 0x68)
	outBuffer = append(outBuffer, d.CtrlCode)
	outBuffer = append(outBuffer, d.DataLen+4)
	outBuffer = append(outBuffer, ucDi[3])
	outBuffer = append(outBuffer, ucDi[2])
	outBuffer = append(outBuffer, ucDi[1])
	outBuffer = append(outBuffer, ucDi[0])
	outBuffer = append(outBuffer, dataBuffer...)
	var sum byte = 0
	for _, num := range outBuffer {
		sum += num
	}
	outBuffer = append(outBuffer, sum)
	outBuffer = append(outBuffer, 0x16) // 9 结束符

	//fmt.Printf("封帧后数据[%v]\r\n", outBuffer)

	return outBuffer, nil
}
