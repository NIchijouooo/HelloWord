package service

import (
	"gateway/models"
	"gateway/models/ReturnModel"
	"gateway/utils"
	"strconv"
	"time"
)

func GetCharData(xAxisList []string, beginDt int64, endDt int64, interval int, intervalType int, historyList []*models.YcData, codeList []int, codeNameList []string) ReturnModel.CharData {
	xAxisList, dateHistoryMap := InitXAxisList(xAxisList, beginDt, endDt, interval, intervalType, historyList)
	var resYcData []ReturnModel.ResYcData

	for idx, code := range codeList { //按code分组
		var valList []string //存储结果值
		for _, xAxis := range xAxisList {
			ycHistoryList, exists := dateHistoryMap[xAxis] //根据时间获取值
			var val string
			if exists { //如果值不为空，再按code进行分组
				//将ycHistoruList按code进行分组收集
				codeMap := make(map[int][]models.YcData)
				// 按照 编码进行分组
				for _, person := range ycHistoryList {
					group := person.Code
					codeMap[group] = append(codeMap[group], person)
				}
				ycModels := codeMap[code] //将code对应的值取出来
				if len(ycModels) > 0 {
					//统计值
					//将返回值转成string prec保留几位小数，bitSize 32或64
					val = strconv.FormatFloat(utils.YcValueMax(ycModels), 'f', 3, 64)
					//	fmt.Println(val)
				} else {
					val = "-"
				}
			}
			valList = append(valList, val)
		}
		resYcData = append(resYcData, ReturnModel.ResYcData{Name: codeNameList[idx], Data: valList})
	}
	var returnMap ReturnModel.CharData
	returnMap.XAxisList = xAxisList
	returnMap.DataList = resYcData
	return returnMap
}

/**
 * 初始化x轴数据
 * @param xAxisList x轴集合
 * @param beginDt 开始时间
 * @param endDt 结束时间
 * @param interval 时间间隔
 * @param intervalType 间隔类型 1-秒;2-分;3-时;4-日;5-月;
 * @return x轴对应的历史数据集合,key=x轴数据,value=x轴对应的历史数据集合
 */
func InitXAxisList(xAxisList []string, beginDt int64, endDt int64, interval int, intervalType int, historyList []*models.YcData) ([]string, map[string][]models.YcData) {
	historyGroupMap := make(map[string][]models.YcData)
	//创建Calendar对象并设置时间
	calendar := time.Unix(0, beginDt*int64(time.Millisecond))
	//初始化Calendar对象的时间部分为0
	calendar = time.Date(calendar.Year(), calendar.Month(), calendar.Day(), 0, 0, 0, 0, calendar.Location())
	// 历史数据分组 map，key=x轴，value=x轴对应的历史数据
	//格式化日期后面遍历存储用
	dataFormat := utils.GetIntervalDateFormat(intervalType)
	//获取当前日期
	for i := beginDt; i <= endDt; {
		//先将毫秒转换成秒，再转成t.time对象
		t := time.Unix(0, i*int64(time.Millisecond))
		//转换日期作为前端展示
		format := t.Format(dataFormat)        //格式化日期
		xAxisList = append(xAxisList, format) //添加到x轴列表

		//计算long长度，增长
		var intervalLong int64
		//计算曾长长度
		intervalLong = utils.GetIntervalTime(calendar, intervalType, interval)
		var list []models.YcData
		intervalStart := t                                                   //开始时间等于当前遍历到的i时间
		intervalEnd := t.Add(time.Duration(intervalLong) * time.Millisecond) //当前时间加上长度，等于结束时间，用于后面遍历使用
		/*
		   获取当前时间间隔内的历史数据,有的数据不在x轴整点内,算到上个时间间隔里
		   如按两小时间隔查询历史数据,则x轴为0h,2h,4h...
		   历史集合里没有整点0点的数据,但是有1点的数据,将1点的数据算到0h里
		*/
		var newHistoryList []*models.YcData
		for _, item := range historyList { //遍历历史数据
			//历史数据要大于等于开始时间，并且小于结束时间
			if (item.Ts.Equal(intervalStart) || item.Ts.After(intervalStart)) && item.Ts.Before(intervalEnd) { //大于开始时间，小于结束时间
				//将符合条件的数据添加到list
				list = append(list, *item)
			} else {
				//如果不符合条件，那么就加进新的list后面重新赋值旧的historyList重新遍历 目的加快查询速度
				newHistoryList = append(newHistoryList, item)
			}
		}
		historyList = newHistoryList
		//添加到结果集
		historyGroupMap[format] = list
		//计算下一个时间段
		i += intervalLong
	}
	return xAxisList, historyGroupMap
}
