package httpServer

import (
	"gateway/controllers"
	"gateway/httpServer/contorl"
	"gateway/httpServer/middleware"
	"gateway/setting"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func RouterWeb(port string) {

	err := middleware.PrivilegeInit()
	if err != nil {
		setting.ZAPS.Errorf("权限配置文件加载失败 %v", err)
	}

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	//router := gin.New()

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	router.Static("/assets", exeCurDir+"/webroot/assets")
	//router.Static("/em-web", exeCurDir+"/webroot/em-web")
	router.StaticFile("/", exeCurDir+"/webroot/index.html")
	router.StaticFile("/favicon.ico", exeCurDir+"/webroot/favicon.ico")
	router.StaticFile("/config.json", exeCurDir+"/webroot/config.json")

	loginRouter := router.Group("/api/v2/account")
	{
		loginRouter.POST("/login", contorl.ApiLogin)
	}

	nodeRouter := router.Group("/api/v2/device/node")
	{
		//查看节点变量
		nodeRouter.GET("/variable/cache", contorl.ApiV2GetNodeVariableFromCache)

		//查看节点历史变量
		//nodeRouter.GET("/variable/histroy", contorl.ApiV2GetNodeHistoryVariable)

		//查看节点变量实时值
		nodeRouter.GET("/variable/realtime", contorl.ApiV2GetNodeRealtimeVariable)

		//调用节点服务
		nodeRouter.POST("/service", contorl.ApiV2InvokeNodeService)
	}

	wsRouter := router.Group("/api")
	{
		wsRouter.GET("/ws", contorl.InitWebsocket)
	}
	router.Use(middleware.JWTAuth())
	{
		//		router.Use(middleware.Privilege())
		//		{

		emController := controllers.NewEMController()
		emController.RegisterRoutes(&router.RouterGroup)

		/**
		20230605
		*/
		/**
		字典
		*/
		realtimeDataController := controllers.NewRealtimeDataController()
		realtimeDataController.RegisterRoutes(&router.RouterGroup)
		dictTypeController := controllers.NewDictTypeController()
		dictTypeController.RegisterRoutes(&router.RouterGroup)
		dictDataController := controllers.NewDictDataController()
		dictDataController.RegisterRoutes(&router.RouterGroup)

		projectInfoController := controllers.NewProjectInfoController()
		projectInfoController.RegisterRoutes(&router.RouterGroup)

		auxiliaryController := controllers.NewAuxiliaryController()
		auxiliaryController.RegisterRoutes(&router.RouterGroup)
		bmsController := controllers.NewBmsController()
		bmsController.RegisterRoutes(&router.RouterGroup)
		ycController := controllers.NewYcController()
		ycController.RegisterRoutes(&router.RouterGroup)

		configurationController := controllers.NewConfigurationCenterController()
		configurationController.RegisterRoutes(&router.RouterGroup)

		accountRouter := router.Group("/api/v2/account")
		{
			accountRouter.GET("/permissions", contorl.ApiGetPermissions)

			accountRouter.PUT("/role/password", contorl.ApiModifyLoginPassword)
		}
		systemRouter := router.Group("/api/v2/system")
		{
			systemRouter.POST("/reboot/system", contorl.ApiSystemReboot)

			systemRouter.POST("/reboot/service", contorl.ApiSystemRebootService)

			systemRouter.GET("/params", contorl.ApiGetSystemStatus)

			systemRouter.GET("/params/cpu", contorl.ApiSystemCPUUseList)

			systemRouter.GET("/params/memory", contorl.ApiSystemMemoryUseList)

			systemRouter.GET("/params/disk", contorl.ApiSystemDiskUseList)

			systemRouter.GET("/params/device/online", contorl.ApiSystemDeviceOnlineList)

			systemRouter.GET("/params/device/packetLoss", contorl.ApiSystemDevicePacketLossList)

			systemRouter.POST("/backup", contorl.ApiBackupFiles)

			//systemRouter.GET("/backup/remote/rt", contorl.ApiBackupFilesToRemote)

			systemRouter.POST("/recover", contorl.ApiRecoverFiles)

			//systemRouter.GET("/recover/remote/rt/files/info", contorl.ApiRecoverFileInfoFromRemote)

			//systemRouter.GET("/recover/remote/rt/file", contorl.ApiRecoverFileFromRemote)

			systemRouter.POST("/update", contorl.ApiSystemUpdate)

			//systemRouter.GET("/login/password", contorl.ApiSystemGetLoginPassword)

			systemRouter.POST("/systemRTC", contorl.ApiSystemSetSystemRTC)

			systemRouter.POST("/ping", contorl.ApiSystemSendPing)

			systemRouter.POST("/commTool", contorl.ApiSendDirectDataToCollInterface)

			systemRouter.GET("/auth", contorl.ApiSystemExportAuth)

			systemRouter.POST("/auth", contorl.ApiSystemImportAuth)

		}

		productRouter := router.Group("/api/v2/product")
		{
			productRouter.POST("/sn", contorl.ApiSystemSetSN)

			productRouter.GET("/sn", contorl.ApiSystemGetSN)
		}

		logRouter := router.Group("/api/v1/log")
		{
			logRouter.GET("/param", contorl.ApiGetLogParam)

			logRouter.POST("/param", contorl.ApiSetLogParam)

			logRouter.GET("/filesInfo", contorl.ApiGetLogFilesInfo)

			logRouter.DELETE("/files", contorl.ApiDeleteLogFile)

			logRouter.GET("/file", contorl.ApiGetLogFile)
		}

		ntpRouter := router.Group("/api/v2/ntp")
		{
			ntpRouter.PUT("/param", contorl.ApiSystemSetNTPParam)

			ntpRouter.POST("/cmd", contorl.ApiSystemSetNTPCmd)

			ntpRouter.GET("/param", contorl.ApiSystemGetNTPParam)
		}

		networkRouter := router.Group("/api/v2/network")
		{
			networkRouter.POST("/param", contorl.ApiAddNetwork)

			networkRouter.PUT("/param", contorl.ApiModifyNetwork)

			networkRouter.DELETE("/param", contorl.ApiDeleteNetwork)

			networkRouter.GET("/params", contorl.ApiGetNetworkParams)

			//networkRouter.GET("/names", contorl.ApiGetNetworkLinkState)

			networkRouter.GET("/names", contorl.ApiGetNetworkNames)
		}

		mobileRouter := router.Group("/api/v2/network")
		{
			mobileRouter.GET("/mobiles", contorl.ApiGetMobileModuleParam)

			mobileRouter.POST("/mobile", contorl.ApiAddMobileModuleParam)

			mobileRouter.PUT("/mobile", contorl.ApiModifyMobileModuleParam)

			mobileRouter.DELETE("/mobile", contorl.ApiDeleteMobileModuleParam)
		}

		commRouter := router.Group("/api/v2/interface")
		{
			//获取通信接口
			commRouter.GET("/communication", contorl.ApiGetCommInterface)

			//增加通信接口
			commRouter.POST("/communication", contorl.ApiAddCommInterface)

			//修改通信接口
			commRouter.PUT("/communication", contorl.ApiModifyCommInterface)

			//删除通信接口
			commRouter.DELETE("/communication", contorl.ApiDeleteCommInterface)

			//获取支持通信接口
			commRouter.GET("/communication/protocol", contorl.ApiGetCommInterfaceProtocol)
		}

		collectRouter := router.Group("/api/v2/interface")
		{
			//增加采集接口
			collectRouter.POST("/collect", contorl.ApiAddCollectInterface)

			//修改采集接口
			collectRouter.PUT("/collect", contorl.ApiModifyCollectInterface)

			//删除采集接口
			collectRouter.DELETE("/collect", contorl.ApiDeleteCollectInterface)

			//获取单个接口信息
			collectRouter.GET("/collect", contorl.ApiGetCollectInterface)

			//获取所有接口信息
			collectRouter.GET("/collects", contorl.ApiGetCollectInterfaces)
		}

		deviceRouter := router.Group("/api/v2/interface/collect")
		{
			//增加节点
			deviceRouter.POST("/device", contorl.ApiAddNode)

			//从CSV文件中批量增加节点
			deviceRouter.POST("/devices/csv", contorl.ApiAddNodesFromCSV)

			//从Xlsx文件中批量增加节点
			deviceRouter.POST("/devices/xlsx/import", contorl.ApiAddNodesFromXlsx)

			//批量导出节点到CSV文件中
			deviceRouter.GET("/devices/csv", contorl.ApiExportNodesToCSV)

			//批量导出节点到Xlsx文件中
			deviceRouter.POST("/devices/xlsx/export", contorl.ApiExportNodesToXlsx)

			//修改单个节点
			deviceRouter.PUT("/device", contorl.ApiModifyNode)

			//修改多个节点
			deviceRouter.PUT("/devices", contorl.ApiModifyNodes)

			//删除节点
			deviceRouter.DELETE("/devices", contorl.ApiDeleteNode)

			//查看节点
			deviceRouter.GET("/device", contorl.ApiGetNode)

			//查看所有节点
			deviceRouter.POST("/devices", contorl.ApiGetNodes)

			//查看节点变量
			deviceRouter.GET("/variable/cache", contorl.ApiGetNodeVariableFromCache)

			//查看节点历史变量
			deviceRouter.GET("/nodeHistoryVariable", contorl.ApiGetNodeHistoryVariableFromCache)

			//查看节点变量实时值
			deviceRouter.GET("/variable/real", contorl.ApiGetNodeRealVariable)

			//调用节点服务
			deviceRouter.POST("/node/service", contorl.ApiInvokeService)
		}

		tslRouter := router.Group("/api/v2")
		{
			//增加设备物模型
			tslRouter.POST("/tsl", contorl.ApiAddDeviceTSL)

			//删除设备物模型
			tslRouter.DELETE("/tsl", contorl.ApiDeleteDeviceTSL)

			//修改设备物模型
			tslRouter.PUT("/tsl", contorl.ApiModifyDeviceTSL)

			//查看设备物模型
			tslRouter.GET("/tsls", contorl.ApiGetDeviceTSLs)

			//查看设备物模型
			//tslRouter.GET("/tsl", contorl.ApiGetDeviceTSL)

			//导入设备物模型插件
			tslRouter.POST("/tsl/plugin/file", contorl.ApiImportTSLLuaPlugin)

			//导出设备物模型插件
			tslRouter.GET("/tsl/plugin/file", contorl.ApiExportTSLLuaPlugin)

			//获取物模型插件内容
			tslRouter.GET("/tsl/plugin/param", contorl.ApiGetTSLLuaPluginParam)

			//删除设备物模型插件
			//tslRouter.DELETE("/tsl/plugin/file", contorl.ApiDeleteTSLLuaPlugin)

			//查看设备物模型内容
			tslRouter.GET("/tsl/contents", contorl.ApiGetDeviceTSLContents)

			//批量导入设备物模型内容
			tslRouter.POST("/tsl/contents/csv", contorl.ApiImportDeviceTSLContents)

			//批量导入设备物模型内容
			tslRouter.POST("/tsl/contents/xlsx", contorl.ApiImportDeviceTSLContentsFromXlsx)

			//从plugin中xlsx文件中同步设备物模型内容
			tslRouter.GET("/tsl/contents/plugin/xlsx", contorl.ApiGetDeviceTSLContentsFromPluginXlsx)

			//批量导出设备物模型内容
			tslRouter.GET("/tsl/contents/csv", contorl.ApiExportDeviceTSLContentsToCSV)

			//批量导出设备物模型内容到xlsx
			tslRouter.GET("/tsl/contents/xlsx", contorl.ApiExportDeviceTSLContentsToXlsx)

			//导出设备物模型内容模板
			tslRouter.GET("/tsl/contents/template", contorl.ApiExportDeviceTSLContentsTemplate)

			//增加设备物模型属性
			tslRouter.POST("/tsl/content/property", contorl.ApiAddTSLProperty)

			//修改设备物模型属性
			tslRouter.PUT("/tsl/content/property", contorl.ApiModifyTSLProperty)

			//删除设备物模型属性
			tslRouter.DELETE("/tsl/content/properties", contorl.ApiDeleteTSLProperties)

			//查看设备物模型属性
			tslRouter.GET("/tsl/content/properties", contorl.ApiGetTSLProperties)

			//增加设备物模型服务
			tslRouter.POST("/tsl/content/service", contorl.ApiAddTSLService)

			//修改设备物模型服务
			tslRouter.PUT("/tsl/content/service", contorl.ApiModifyTSLService)

			//删除设备物模型服务
			tslRouter.DELETE("/tsl/content/services", contorl.ApiDeleteTSLServices)

			//批量导出modbus模型命令列表
			tslRouter.GET("/tsl/modbus/block/xlsx", contorl.ApiExportTSLModbusCmdToXlsx)

			//批量导入modbus模型命令列表
			tslRouter.POST("/tsl/modbus/block/xlsx", contorl.ApiAddTSLModbusCmdFromXlsx)

			//增加modbus采集模型命令
			tslRouter.POST("/tsl/modbus/cmd", contorl.ApiAddTSLModbusCmd)

			//修改modbus采集模型命令
			tslRouter.PUT("/tsl/modbus/cmd", contorl.ApiModifyTSLModbusCmd)

			//删除modbus采集模型命令
			tslRouter.DELETE("/tsl/modbus/cmd", contorl.ApiDeleteTSLModbusCmd)

			//查看modbus采集模型命令
			tslRouter.GET("/tsl/modbus/cmd", contorl.ApiGetTSLModbusCmd)

			//增加modbus采集模型命令参数
			tslRouter.POST("/tsl/modbus/cmd/param", contorl.ApiAddTSLModbusCmdProperty)

			//批量导入modbus采集模型命令参数
			tslRouter.POST("/tsl/modbus/cmd/param/xlsx", contorl.ApiAddTSLModbusCmdPropertyFromXlsx)

			//修改modbus采集模型命令参数
			tslRouter.PUT("/tsl/modbus/cmd/param", contorl.ApiModifyTSLModbusCmdProperty)

			//删除modbus采集模型命令参数
			tslRouter.DELETE("/tsl/modbus/cmd/params", contorl.ApiDeleteTSLModbusCmdProperties)

			//查看modbus采集模型单个命令的参数
			tslRouter.GET("/tsl/modbus/cmd/params", contorl.ApiGetTSLModbusCmdProperties)

			//批量导出modbus采集模型单个命令的参数
			tslRouter.GET("/tsl/modbus/cmd/param/xlsx", contorl.ApiExportTSLModbusCmdPropertiesToXlsx)

			//查看modbus采集模型所有属性
			tslRouter.GET("/tsl/modbus/properties", contorl.ApiGetTSLModbusProperties)

			//批量导出dlt64507模型命令列表
			tslRouter.GET("/tsl/dlt64507/block/xlsx", contorl.ApiExportTSLD07CmdToXlsx)

			//批量导入dlt64507模型命令列表
			tslRouter.POST("/tsl/dlt64507/block/xlsx", contorl.ApiAddTSLD07CmdFromXlsx)

			//增加dlt64507采集模型命令
			tslRouter.POST("/tsl/dlt64507/cmd", contorl.ApiAddTSLD07Cmd)

			//修改dlt64507采集模型命令
			tslRouter.PUT("/tsl/dlt64507/cmd", contorl.ApiModifyTSLD07Cmd)

			//删除dlt64507采集模型命令
			tslRouter.DELETE("/tsl/dlt64507/cmd", contorl.ApiDeleteTSLD07Cmd)

			//查看dlt64507采集模型命令
			tslRouter.GET("/tsl/dlt64507/cmd", contorl.ApiGetTSLD07Cmd)

			//增加dlt64507采集模型命令参数
			tslRouter.POST("/tsl/dlt64507/cmd/param", contorl.ApiAddTSLD07CmdProperty)

			//批量导入dlt64507采集模型命令参数
			tslRouter.POST("/tsl/dlt64507/cmd/param/xlsx", contorl.ApiAddTSLD07CmdPropertyFromXlsx)

			//修改dlt64507采集模型命令参数
			tslRouter.PUT("/tsl/dlt64507/cmd/param", contorl.ApiModifyTSLD07CmdProperty)

			//删除modbus采集模型命令参数
			tslRouter.DELETE("/tsl/dlt64507/cmd/params", contorl.ApiDeleteTSLD07CmdProperties)

			//查看dlt64507采集模型单个命令的参数
			tslRouter.GET("/tsl/dlt64507/cmd/params", contorl.ApiGetTSLD07CmdProperties)

			//批量导出dlt64507集模型单个命令的参数
			tslRouter.GET("/tsl/dlt64507/cmd/param/xlsx", contorl.ApiExportTSLD07CmdPropertiesToXlsx)

			//查看dlt64507采集模型所有属性
			tslRouter.GET("/tsl/dlt64507/properties", contorl.ApiGetTSLD07Properties)
		}

		reportRouter := router.Group("/api/v2/service/report")
		{
			reportRouter.GET("/protocol", contorl.ApiGetReportProtocol)

			reportRouter.POST("/gateway", contorl.ApiAddReportGWParam)

			reportRouter.PUT("/gateway", contorl.ApiModifyReportGWParam)

			reportRouter.GET("/gateways", contorl.ApiGetReportGWParam)

			reportRouter.DELETE("/gateway", contorl.ApiDeleteReportGWParam)

			reportRouter.POST("/model", contorl.ApiAddReportModel)

			reportRouter.PUT("/model", contorl.ApiModifyReportModel)

			reportRouter.GET("/models", contorl.ApiGetReportModel)

			reportRouter.DELETE("/model", contorl.ApiDeleteReportModel)

			reportRouter.POST("/model/property", contorl.ApiAddReportModelProperty)

			reportRouter.POST("/model/properties", contorl.ApiAddReportModelPropertyes) //ltg add 2023-05-26

			reportRouter.PUT("/model/property", contorl.ApiModifyReportModelProperty)

			reportRouter.GET("/model/properties", contorl.ApiGetReportModelProperties)

			reportRouter.DELETE("/model/properties", contorl.ApiDeleteReportModelProperties)

			reportRouter.POST("/device", contorl.ApiAddReportNodeWParam)

			reportRouter.PUT("/device", contorl.ApiModifyReportNodeWParam)

			reportRouter.POST("/devices/csv", contorl.ApiBatchAddReportNodeParam)

			reportRouter.POST("/devices/xlsx", contorl.ApiBatchAddReportNodeParamFromXlsx)

			reportRouter.GET("/devices", contorl.ApiGetReportNodeWParam)

			reportRouter.GET("/devices/csv", contorl.ApiBatchExportReportNodeWParam)

			reportRouter.GET("/devices/xlsx", contorl.ApiBatchExportReportNodeWParamToXlsx)

			reportRouter.DELETE("/devices", contorl.ApiDeleteReportNodeWParam)

			reportRouter.POST("/device/cmd/report", contorl.ApiSetReportNodeReport)
		}

		reportMBTCPRouter := router.Group("/api/v2/service/report/mbTCP")
		{
			reportMBTCPRouter.POST("/register", contorl.ApiAddReportMBTCPRegister)

			reportMBTCPRouter.PUT("/register", contorl.ApiModifyReportMBTCPRegister)

			reportMBTCPRouter.DELETE("/registers", contorl.ApiDeleteReportMBTCPRegisters)

			reportMBTCPRouter.GET("/registers", contorl.ApiGetReportMBTCPRegisters)

			reportMBTCPRouter.GET("/register/xlsx", contorl.ApiBatchExportReportMBTCPRegistersToXlsx)

			reportMBTCPRouter.POST("/register/xlsx", contorl.ApiBatchImportReportMBTCPRegistersToXlsx)
		}

		virtualRouter := router.Group("/api/v2/virtual")
		{
			virtualRouter.POST("/device", contorl.ApiAddVirtualDevice)

			virtualRouter.PUT("/device", contorl.ApiModifyVirtualDevice)

			virtualRouter.DELETE("/devices", contorl.ApiDeleteVirtualDevice)

			virtualRouter.GET("/devices", contorl.ApiGetVirtualDevice)

			virtualRouter.POST("/property", contorl.ApiAddVirtualDeviceProperty)

			virtualRouter.PUT("/property", contorl.ApiModifyVirtualDeviceProperty)

			virtualRouter.DELETE("/properties", contorl.ApiDeleteVirtualDeviceProperties)

			virtualRouter.GET("/properties", contorl.ApiGetVirtualDeviceProperties)
		}

		eventRouter := router.Group("/api/v1/event")
		{
			eventRouter.GET("/report", contorl.ApiGetReportEvent)

			eventRouter.DELETE("/report", contorl.ApiDeleteReportEvent)
		}
		//		}
	}

	setting.ZAPS.Infof("gin 初始化端口[%s]", port)
	err = router.Run(port)
	if err != nil {
		setting.ZAPS.Errorf("gin启动错误 %v", err)
		os.Exit(0)
	}
}
