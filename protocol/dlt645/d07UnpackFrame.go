package dlt645

import (
	"errors"
)

type D07UserData struct {
	BlockAddOffset int //如果块读取时，当前数据在数据域内的偏移地址
	RulerAddOffset int //如果当前数据标识有多个变量，当前变量在当前数据地址中的偏移地址
	RulerInfo      D07RulerInfo
	User           interface{}
}

type D07UnpackFrame struct {
	LeadNum  byte   //前导字符 0xFE的个数
	Address  string //12位地址域数据
	CtrlCode byte   //控制码 字节型
	RulerID  uint32 //规约ID

	DataLen  int //数据域长
	FrameLen int //整个帧长
	Data     []byte
	UserData []D07UserData
}

func UnpackD07Frame(inBuffer []byte) (D07UnpackFrame, error) {
	var reUnpackFrame D07UnpackFrame

	inBufferLength := len(inBuffer)
	if inBufferLength < D07_FRAME_LEN_MIN {
		return D07UnpackFrame{}, errors.New("传入的帧长度小于最小帧长度，无法解析")
	}

	pos := 0 //当前字符位置

	/* 1、检查前导码FE的数量   */
	reUnpackFrame.LeadNum = 0
	for i := 0; i < 4; i++ {
		if 0xFE == inBuffer[i] {
			pos++
			reUnpackFrame.LeadNum++
		}
	}
	nCheckSumPosStart := pos //检验的开始位置

	/* 2、检查前导字符 0x68   */
	if 0x68 != inBuffer[pos] || 0x68 != inBuffer[pos+7] {
		return D07UnpackFrame{}, errors.New("字符0x68检查出错")
	}
	pos++

	/* 3、解析出地址 */
	addrBcd := make([]byte, 0)
	for i := 0; i < 6; i++ {
		addrBcd = append(addrBcd, inBuffer[pos])
		pos++
	}
	strBCD, err := D07BCD2Str(addrBcd, 6)
	if err != nil {
		return D07UnpackFrame{}, errors.New("解析数据包地址出错")
	}
	reUnpackFrame.Address = strBCD

	/* 4、解析控制码 */
	if inBufferLength < pos {
		return D07UnpackFrame{}, errors.New("数据包长度异常")
	}
	reUnpackFrame.CtrlCode = inBuffer[pos+1]
	pos += 2 //控制码前有一个0x68

	/* 5、解析数据域  */
	if inBufferLength < pos {
		return D07UnpackFrame{}, errors.New("数据包长度异常")
	}
	ucDataLen := inBuffer[pos]
	pos++
	nCheckSumPos := pos + int(ucDataLen) - nCheckSumPosStart
	nEndPos := pos + int(ucDataLen) + 2
	aucDataTmp := make([]byte, 0)
	aucRulerTmp := make([]byte, 0)

	for i := 0; i < int(ucDataLen); i++ {
		if inBufferLength < pos {
			return D07UnpackFrame{}, errors.New("数据包长度异常")
		}

		if i < 4 {
			aucRulerTmp = append(aucRulerTmp, inBuffer[pos])
		} else {
			aucDataTmp = append(aucDataTmp, inBuffer[pos])
		}

		pos++
	}
	reUnpackFrame.FrameLen = int(ucDataLen) + 12
	if ucDataLen >= 4 {
		reUnpackFrame.DataLen = int(ucDataLen) - 4
	}

	if len(aucRulerTmp) == 4 {
		reUnpackFrame.RulerID = uint32(aucRulerTmp[3]-0x33)<<24 | uint32(aucRulerTmp[2]-0x33)<<16 | uint32(aucRulerTmp[1]-0x33)<<8 | uint32(aucRulerTmp[0]-0x33)
	}
	reUnpackFrame.Data = append(reUnpackFrame.Data, aucDataTmp...)

	/* 5、查检checksum  */
	var ucCheckSum byte = 0
	if len(inBuffer) < nCheckSumPos {
		return D07UnpackFrame{}, errors.New("数据包长度异常")
	}

	for i := 0; i < nCheckSumPos; i++ {
		ucCheckSum += inBuffer[i+nCheckSumPosStart]
	}

	if ucCheckSum != inBuffer[nEndPos-2] {
		return D07UnpackFrame{}, errors.New("解析数据包检验出错")
	}

	/*6、查检结束符号  */
	if len(inBuffer) < nEndPos {
		return D07UnpackFrame{}, errors.New("数据包长度异常")
	}

	if 0x16 != inBuffer[nEndPos-1] {
		return D07UnpackFrame{}, errors.New("解析数据包结束符号0x16出错")
	}

	return reUnpackFrame, nil
}
