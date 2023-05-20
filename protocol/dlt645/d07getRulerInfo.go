package dlt645

import (
	"errors"
	"fmt"
)

type D07RulerTrans func(D07TransFlg int, user *byte, frame *byte) int

/* 规约读写类型 */
const (
	D07READONLY  byte = iota /* 只读 */
	D07WRITEONLY             /* 只写 */
	D07READWRITE             /* 读写 */
)

//D07RulerInfo  规则对应的信息结构体
type D07RulerInfo struct {
	ID             uint32 //数据标识,非字符串形式
	BlockAddOffset int    //当前数据在块数据域内的偏移地址
	RulerAddOffset int    //当前变量在当前ID数据地址中的偏移地址
	Format         string //数据格式，字符串格式XXXX.XX....
	Len            byte   //数据长度
	Unit           string //单位
	Rdwr           byte   //读写  E_D07_RDWR_READ_ONLY--0,E_D07_RDWR_WRITE_ONLY--1,E_D07_RDWR_READ_WRITE--2
	Name           string //数据项目名称
}

func GetD07RulerInfo(rulerID uint32) ([]D07RulerInfo, error) {
	ucDi := make([]byte, 4)
	outRulerInfoSlice := make([]D07RulerInfo, 0)
	var name string
	var rulerInfo D07RulerInfo

	//分解数据标识
	ucDi[0] = byte(rulerID & 0xff)
	ucDi[1] = byte((rulerID >> 8) & 0xff)
	ucDi[2] = byte((rulerID >> 16) & 0xff)
	ucDi[3] = byte((rulerID >> 24) & 0xff)

	//根据数据标识，返回当前标识的详细信息
	rulerInfo.ID = rulerID
	switch ucDi[3] {

	/* 对应表A.1参变量数据标识编码表 */
	case 0: //{00}[*][*][*]
		{
			if ucDi[0] > 0x0c {
				return nil, errors.New("表A.1中数据标识的DI0不能大于0x0c")
			}

			// 封装表A.1结算日字符串
			if ucDi[0] == 0 {
				name = fmt.Sprintf("%s", "(当前)")
			} else {
				name = fmt.Sprintf("(上%d结算日)", ucDi[0])
			}

			// 对于表 A.1 相同的数据属性
			rulerInfo.Format = "XXXXXX.XX"
			rulerInfo.Len = 4
			rulerInfo.Rdwr = D07READONLY

			switch ucDi[2] {
			case 0x00: //[00][00]{*}[*]
				{
					rulerInfo.RulerAddOffset = 0
					rulerInfo.Unit = "kWh"

					if 0 == ucDi[1] { //[00][00]{00}[*]
						rulerInfo.BlockAddOffset = 0
						rulerInfo.Name = name + "组合有功总电能"
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)

					} else if ucDi[1] >= 0x01 && ucDi[1] <= 0x3F { //[00][00]{(1~3F)}[*]
						rulerInfo.BlockAddOffset = int(ucDi[1]) * 4
						rulerInfo.Name = name + fmt.Sprintf("组合有功费率%d电能", ucDi[1])
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)

					} else if 0xFF == ucDi[1] { //[00][01]{FF}[*]
						rulerInfo.BlockAddOffset = 0
						rulerInfo.Name = name + "组合有功总电能"
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)
						for i := 1; i <= 0x3F; i++ {
							rulerInfo.BlockAddOffset = i * 4
							rulerInfo.Name = name + fmt.Sprintf("组合有功费率%d电能", i)
							outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)
						}
					}
				}
			case 0x01: //[00][01]{*}[*]
				{
					rulerInfo.RulerAddOffset = 0
					rulerInfo.Unit = "kWh"

					if 0 == ucDi[1] { //[00][00]{00}[*]
						rulerInfo.BlockAddOffset = 0
						rulerInfo.Name = name + "正向有功总电能"
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)

					} else if ucDi[1] >= 0x01 && ucDi[1] <= 0x3F { //[00][00]{(1~3F)}[*]
						rulerInfo.BlockAddOffset = int(ucDi[1]) * 4
						rulerInfo.Name = name + fmt.Sprintf("正向有功费率%d电能", ucDi[1])
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)

					} else if 0xFF == ucDi[1] { //[00][01]{FF}[*]
						rulerInfo.BlockAddOffset = 0
						rulerInfo.Name = name + "正向有功总电能"
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)
						for i := 1; i <= 0x3F; i++ {
							rulerInfo.BlockAddOffset = i * 4
							rulerInfo.Name = name + fmt.Sprintf("正向有功费率%d电能", i)
							outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)
						}
					}
				}
			case 0x02: //[00]{02}[*][*]
				{
					rulerInfo.RulerAddOffset = 0
					rulerInfo.Unit = "kWh"

					if 0 == ucDi[1] { //[00][00]{00}[*]
						rulerInfo.BlockAddOffset = 0
						rulerInfo.Name = name + "反向有功总电能"
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)

					} else if ucDi[1] >= 0x01 && ucDi[1] <= 0x3F { //[00][00]{(1~3F)}[*]
						rulerInfo.BlockAddOffset = int(ucDi[1]) * 4
						rulerInfo.Name = name + fmt.Sprintf("反向有功费率%d电能", ucDi[1])
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)

					} else if 0xFF == ucDi[1] { //[00][01]{FF}[*]
						rulerInfo.BlockAddOffset = 0
						rulerInfo.Name = name + "反向有功总电能"
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)
						for i := 1; i <= 0x3F; i++ {
							rulerInfo.BlockAddOffset = i * 4
							rulerInfo.Name = name + fmt.Sprintf("反向有功费率%d电能", i)
							outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)
						}
					}
				}
			default:
			}
		}

	/* 对应表A.2参变量数据标识编码表 */
	case 1: //{01}[*][*][*]
		{
			if ucDi[0] > 0x0c {
				return nil, errors.New("表A.2中数据标识的DI0不能大于0x0c")
			}

			// 封装表A.2结算日字符串
			if ucDi[0] == 0 {
				name = fmt.Sprintf("%s", "(当前)")
			} else {
				name = fmt.Sprintf("(上%d结算日)", ucDi[0])
			}

			// 对于表 A.2 相同的数据属性
			rulerInfo.Len = 4
			rulerInfo.Rdwr = D07READONLY

			switch ucDi[2] {
			case 0x01: //[01][01]{*}[*]
				{
					if 0 == ucDi[1] { //[01][01]{00}[*]
						rulerInfo.BlockAddOffset = 0
						rulerInfo.RulerAddOffset = 0
						rulerInfo.Format = "XX.XXXX"
						rulerInfo.Unit = "kW"
						rulerInfo.Name = name + "正向有功总最大需量"
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)
						rulerInfo.RulerAddOffset = 4
						rulerInfo.Format = "YYMMDDhhmm"
						rulerInfo.Unit = "年月日时分"
						rulerInfo.Name = name + "正向有功总最大需量发生时间"
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)

					} else if ucDi[1] >= 0x01 && ucDi[1] <= 0x3F { //[01][00]{(1~3F)}[*]
						rulerInfo.BlockAddOffset = int(ucDi[1]) * 8
						rulerInfo.RulerAddOffset = 0
						rulerInfo.Format = "XX.XXXX"
						rulerInfo.Unit = "kW"
						rulerInfo.Name = name + fmt.Sprintf("正向有功费率%d最大需量", ucDi[1])
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)
						rulerInfo.RulerAddOffset = 4
						rulerInfo.Format = "YYMMDDhhmm"
						rulerInfo.Name = name + fmt.Sprintf("正向有功费率%d最大需量发生时间", ucDi[1])
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)

					} else if 0xFF == ucDi[1] { //[01][01]{FF}[*]
						rulerInfo.BlockAddOffset = 0
						rulerInfo.RulerAddOffset = 0
						rulerInfo.Format = "XX.XXXX"
						rulerInfo.Unit = "kW"
						rulerInfo.Name = name + "正向有功总最大需量"
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)
						rulerInfo.RulerAddOffset = 4
						rulerInfo.Format = "YYMMDDhhmm"
						rulerInfo.Name = name + "正向有功总最大需量发生时间"
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)

						for i := 1; i <= 0x3F; i++ {
							rulerInfo.BlockAddOffset = i * 8
							rulerInfo.RulerAddOffset = 0
							rulerInfo.Format = "XX.XXXX"
							rulerInfo.Unit = "kW"
							rulerInfo.Name = name + fmt.Sprintf("正向有功费率%d最大需量", i)
							outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)
							rulerInfo.RulerAddOffset = 4
							rulerInfo.Format = "YYMMDDhhmm"
							rulerInfo.Name = name + fmt.Sprintf("正向有功费率%d最大需量发生时间", i)
							outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)
						}
					}
				}
			case 0x02: //[01]{02}[*][*]
				{
					if 0 == ucDi[1] { //[01][02]{00}[*]
						rulerInfo.BlockAddOffset = 0
						rulerInfo.RulerAddOffset = 0
						rulerInfo.Format = "XX.XXXX"
						rulerInfo.Unit = "kW"
						rulerInfo.Name = name + "反向有功总最大需量"
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)
						rulerInfo.RulerAddOffset = 4
						rulerInfo.Name = name + "反向有功总最大需量发生时间"
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)

					} else if ucDi[1] >= 0x01 && ucDi[1] <= 0x3F { //[01][02]{(1~3F)}[*]
						rulerInfo.BlockAddOffset = int(ucDi[1]) * 8
						rulerInfo.RulerAddOffset = 0
						rulerInfo.Format = "XX.XXXX"
						rulerInfo.Unit = "kW"
						rulerInfo.Name = name + fmt.Sprintf("反向有功费率%d最大需量", ucDi[1])
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)
						rulerInfo.RulerAddOffset = 4
						rulerInfo.Name = name + fmt.Sprintf("反向有功费率%d最大需量发生时间", ucDi[1])
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)

					} else if 0xFF == ucDi[1] { //[01][02]{FF}[*]
						rulerInfo.BlockAddOffset = 0
						rulerInfo.RulerAddOffset = 0
						rulerInfo.Format = "XX.XXXX"
						rulerInfo.Unit = "kW"
						rulerInfo.Name = name + "反向有功总最大需量"
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)
						rulerInfo.RulerAddOffset = 4
						rulerInfo.RulerAddOffset = 4
						rulerInfo.Format = "YYMMDDhhmm"
						rulerInfo.Name = name + "反向有功总最大需量发生时间"
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)

						for i := 1; i <= 0x3F; i++ {
							rulerInfo.BlockAddOffset = i * 8
							rulerInfo.RulerAddOffset = 0
							rulerInfo.Name = name + fmt.Sprintf("反向有功费率%d最大需量", i)
							outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)
							rulerInfo.RulerAddOffset = 4
							rulerInfo.RulerAddOffset = 4
							rulerInfo.Format = "YYMMDDhhmm"
							rulerInfo.Name = name + fmt.Sprintf("反向有功费率%d最大需量发生时间", i)
							outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)
						}
					}
				}
			default:
			}

		}

	/* 对应表A.3参变量数据标识编码表 */
	case 2: //{02}[*][*][*]
		{
			switch ucDi[2] {
			case 0x01: //[02]{01}[*][*]
				{
					if 0 != ucDi[0] { //[02][01][*]{!(0)}
						return nil, errors.New("表A.3中数据标识[02][01][*]{!(0)}的DI0不能为非0")
					}

					rulerInfo.Format = "XXX.X"
					rulerInfo.Len = 2
					rulerInfo.Rdwr = D07READONLY
					rulerInfo.Unit = "V"
					rulerInfo.RulerAddOffset = 0

					if ucDi[1] == 0x01 || ucDi[1] == 0x0ff {
						rulerInfo.BlockAddOffset = 0
						rulerInfo.Name = name + "A相电压"
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)
					}

					if ucDi[1] == 0x02 || ucDi[1] == 0x0ff {
						rulerInfo.BlockAddOffset = 2
						rulerInfo.Name = name + "B相电压"
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)
					}

					if ucDi[1] == 0x03 || ucDi[1] == 0x0ff {
						rulerInfo.BlockAddOffset = 4
						rulerInfo.Name = name + "C相电压"
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)
					}
				}
			case 0x02: //[02]{02}[*][*]
				{
					if 0 != ucDi[0] { //[02][01][*]{!(0)}
						return nil, errors.New("表A.3中数据标识[02][02][*]{!(0)}的DI0不能为非0")
					}

					rulerInfo.Format = "XXX.XXX"
					rulerInfo.Len = 3
					rulerInfo.Rdwr = D07READONLY
					rulerInfo.Unit = "A"
					rulerInfo.RulerAddOffset = 0

					if ucDi[1] == 0x01 || ucDi[1] == 0x0ff {
						rulerInfo.BlockAddOffset = 0
						rulerInfo.Name = name + "A相电流"
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)
					}

					if ucDi[1] == 0x02 || ucDi[1] == 0x0ff {
						rulerInfo.BlockAddOffset = 3
						rulerInfo.Name = name + "B相电流"
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)
					}

					if ucDi[1] == 0x03 || ucDi[1] == 0x0ff {
						rulerInfo.BlockAddOffset = 6
						rulerInfo.Name = name + "C相电流"
						outRulerInfoSlice = append(outRulerInfoSlice, rulerInfo)
					}
				}
			default:
			}
		}

	/* 对应表A.4参变量数据标识编码表 */
	case 3: //{03}[*][*][*]
		{

		}

	/* 对应表A.5参变量数据标识编码表 **/
	case 4: //{04}[*][*][*]
		{

		}

	/* 表A.6 冻结数据标识编码表 */
	case 5: //{05}[*][*][*]
		{

		}

	/* 表 A.7负荷记录数据标识编码表 */
	case 6: //{06}[*][*][*]
		{

		}

	/* 用户扩展 */
	default: //{!(0~6)[*][*][*]}
		{
			return nil, errors.New("出错")
		}
	}

	return outRulerInfoSlice, nil
}

// D07RulerPara 规约类型的额外数据
//type D07RulerPara struct {
//	D07ParaPayoff int // 结算日(不关心，当前，上(1-12)结算日)  E_D07_PAYOFF_NULL,E_D07_PAYOFF_CURRENT....E_D07_PAYOFF_PRE_12
//	D07ParaRate   int // 费率(不关心，1~63)  E_D07_RATE_NULL,E_D07_RATE_1...E_D07_RATE_64
//	D07ParaHarm   int // 谐波次(不关心，1~21) E_D07_HARM_NULL,E_D07_HARM_1...E_D07_HARM_21
//	D07ParaLast   int // 上(n)次记录(不关心，1-10)  E_D07_LAST_NULL,E_D07_LAST_1...E_D07_LAST_12
//}
//
////D07RulerInfo  规则对应的信息结构体
//type D07RulerInfo struct {
//	D07RulerType   int           // 规约类型  E_D07_RULER_TYPE_UNKNOWN,E_D07_RULER_TYPE_A1_MIN,E_D07_RULER_TYPE_COMB_HAVE_POWER_TOTAL...
//	D07RulerRdwr   int           // 规约数据的读写属性 E_D07_RDWR_READ_ONLY,E_D07_RDWR_WRITE_ONLY,E_D07_RDWR_READ_WRITE
//	D07RulerFormat int           // 规约数据的格式 E_D07_FMT_UNKOWN,E_D07_FMT_XXXX,E_D07_FMT_XX_2,E_D07_FMT_XXXXXX...
//	D07RulerTrans  D07RulerTrans // 数据域转换函数指针
//	D07RulerPara   D07RulerPara  // 规约类型的额外数据
//	D07Len         int           // 数据域字节长度
//	D07Name        string        // 该条规约数据项名称
//
//	ID             uint32 //数据标识,非字符串形式
//	BlockAddOffset int    //当前数据在块数据域内的偏移地址
//	RulerAddOffset int    //当前变量在当前ID数据地址中的偏移地址
//	Format         string //数据格式，字符串格式XXXX.XX....
//	Len            byte   //数据长度
//	Unit           string //单位
//	Rdwr           byte   //读写  E_D07_RDWR_READ_ONLY--0,E_D07_RDWR_WRITE_ONLY--1,E_D07_RDWR_READ_WRITE--2
//	Name           string //数据项目名称
//}

//GetD07RulerInfo  通过规约ID获得对应规约的详细信息   返回D07RulerInfo信息结构体和错误类型(正确返回 0 其它为错误类型)
//func GetD07RulerInfo(rulerID int32) ([]D07RulerInfo, int) {
//	ucDi := make([]byte, 4)
//	var strPayOff string
//	var name string
//	outRulerInfoSlice := make([]D07RulerInfo, 0)
//
//	outRulerInfo := D07RulerInfo{
//		D07RulerType:   E_D07_RULER_TYPE_UNKNOWN,                                                           // 规约类型  E_D07_RULER_TYPE_UNKNOWN
//		D07RulerRdwr:   E_D07_RDWR_READ_ONLY,                                                               // 规约数据的读写属性
//		D07RulerFormat: E_D07_FMT_UNKOWN,                                                                   // 规约数据的格式
//		D07RulerTrans:  nil,                                                                                // 数据域转换函数指针
//		D07RulerPara:   D07RulerPara{E_D07_PAYOFF_NULL, E_D07_RATE_NULL, E_D07_HARM_NULL, E_D07_LAST_NULL}, // 规约类型的额外数据
//		D07Len:         0,
//		D07Name:        "",
//		BlockAddOffset: 0,
//		RulerAddOffset: 0,
//	}
//
//	ucDi[0] = byte(rulerID & 0xff)
//	ucDi[1] = byte((rulerID >> 8) & 0xff)
//	ucDi[2] = byte((rulerID >> 16) & 0xff)
//	ucDi[3] = byte((rulerID >> 24) & 0xff)
//
//	switch ucDi[3] {
//	case 0: //{00}[*][*][*]  /* 对应表A.1参变量数据标识编码表 */
//		{
//			if ucDi[0] > 0x0c {
//				return nil, E_D07_ERRO_UNKOWN_ID
//			}
//
//			// 封装结算日字符串
//			if ucDi[0] == 0 {
//				strPayOff = fmt.Sprintf("%s", "(当前)")
//			} else {
//				strPayOff = fmt.Sprintf("(上%d结算日)", ucDi[0])
//			}
//
//			/* 对于表 A.1 相同的数据属性 */
//			outRulerInfo.D07RulerPara.D07ParaPayoff = int(ucDi[0] + 1)
//			outRulerInfo.D07RulerPara.D07ParaRate = E_D07_RATE_NULL
//			outRulerInfo.D07RulerPara.D07ParaHarm = E_D07_HARM_NULL
//			outRulerInfo.D07RulerPara.D07ParaLast = E_D07_LAST_NULL
//			outRulerInfo.D07RulerRdwr = E_D07_RDWR_READ_ONLY
//			outRulerInfo.D07RulerFormat = E_D07_FMT_XXXXXX_XX
//			outRulerInfo.D07Len = 4
//			outRulerInfo.D07RulerTrans = TransD07DataXXXXXX_XX_TEST
//
//			switch ucDi[2] {
//			case 0x00: //[00][00]{*}[*]
//				{
//					if 0 == ucDi[1] { //[00][00]{00}[*]
//						name = fmt.Sprintf("%s", "正向有功总电能")
//						outRulerInfo.D07Name = fmt.Sprintf("%s%s", strPayOff, name)
//						outRulerInfoSlice = append(outRulerInfoSlice, outRulerInfo)
//
//					}
//				}
//			case 0x01: //[00][01]{*}[*]
//				{
//					if 0 == ucDi[1] { //[00][01]{00}[*]
//						outRulerInfo.D07RulerType = E_D07_RULER_TYPE_FORTH_HAVE_POWER_TOTAL
//						name = fmt.Sprintf("%s", "正向有功总电能")
//						outRulerInfo.D07Name = fmt.Sprintf("%s%s", strPayOff, name)
//						outRulerInfoSlice = append(outRulerInfoSlice, outRulerInfo)
//
//					} else if ucDi[1] >= 0x01 && ucDi[1] <= 0x3F { //[00][01]{(1~3F)}[*]
//						outRulerInfo.D07RulerType = E_D07_RULER_TYPE_FORTH_HAVE_POWER_RATE
//						outRulerInfo.D07RulerPara.D07ParaRate = int(ucDi[1])
//						name = fmt.Sprintf("正向有功费率%d电能", ucDi[1])
//						outRulerInfo.D07Name = fmt.Sprintf("%s%s", strPayOff, name)
//						outRulerInfoSlice = append(outRulerInfoSlice, outRulerInfo)
//
//					} else if 0xFF == ucDi[1] { //[00][01]{FF}[*]
//						outRulerInfo.D07RulerType = E_D07_RULER_TYPE_FORTH_HAVE_POWER_BLOCK
//						name = fmt.Sprintf("%s", "正向有功电能数据块")
//						outRulerInfo.D07Name = fmt.Sprintf("%s%s", strPayOff, name)
//						outRulerInfoSlice = append(outRulerInfoSlice, outRulerInfo)
//					} else { //[00][01]{(!(0-3F,FF))}[*]
//						return nil, E_D07_ERRO_UNKOWN_ID
//					}
//				}
//			case 0x02: //[00]{02}[*][*]
//				{
//
//				}
//
//			}
//		}
//	case 1: //{01}[*][*][*] /* 对应表A.2参变量数据标识编码表 */
//		{
//
//		}
//	case 2: //{02}[*][*][*] /* 对应表A.3参变量数据标识编码表 */
//		{
//
//		}
//	case 3: //{03}[*][*][*] /* 对应表A.4参变量数据标识编码表 */
//		{
//
//		}
//	case 4: //{04}[*][*][*]  /* 对应表A.5参变量数据标识编码表 */
//		{
//
//		}
//	case 5: //{05}[*][*][*]  /* 表A.6 冻结数据标识编码表 */
//		{
//
//		}
//	case 6: //{06}[*][*][*]  /* 表 A.7负荷记录数据标识编码表 */
//		{
//
//		}
//	default: //{!(0~6)[*][*][*]}  /* 用户扩展 */
//		{
//			return nil, E_D07_ERRO_UNKOWN_ID
//		}
//	}
//
//	// 合成最后的名字
//
//	return outRulerInfoSlice, E_D07_OK
//}
