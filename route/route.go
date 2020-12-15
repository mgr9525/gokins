package route

import (
	"gokins/comm"
	"gokins/core"
	"gokins/route/server"
	"gokins/service/utilService"
)

func Init() {
	comm.Gin.GET("/", Statics)
	comm.Gin.GET("/static/*nm", Statics)

	comm.Gin.Use(core.MidAccessAllow)

	// comm模块路由
	{
		gpComm := comm.Gin.Group("/comm")
		gpComm.Any("/info", server.CommInfo)
		gpLogin := comm.Gin.Group("/lg")
		gpLogin.Any("/info", server.LoginInfo)
		gpLogin.Any("/login", core.GinHandler(server.Login))
		gpLogin.Any("/install", core.GinHandler(server.Install))
		gpLogin.Any("/uppass", core.GinHandler(server.Uppass))
	}

	// 流水线路由
	{
		gpModel := comm.Gin.Group("/model")
		gpModel.Use(utilService.MidNeedLogin)
		gpModel.Any("/get", core.GinHandler(server.ModelGet))
		gpModel.Any("/list", core.GinHandler(server.ModelList))
		gpModel.Any("/edit", core.GinHandler(server.ModelEdit))
		gpModel.Any("/del", core.GinHandler(server.ModelDel))
		gpModel.Any("/runs", core.GinHandler(server.ModelRuns))
		gpModel.Any("/run", core.GinHandler(server.ModelRun))
		gpModel.Any("/stop", core.GinHandler(server.ModelStop))
		gpModel.Any("/copy", core.GinHandler(server.ModelCopy))
	}

	// 插件路由
	{
		gpPlug := comm.Gin.Group("/plug")
		gpPlug.Use(utilService.MidNeedLogin)
		gpPlug.Any("/list", core.GinHandler(server.PlugList))
		gpPlug.Any("/edit", core.GinHandler(server.PlugEdit))
		gpPlug.Any("/del", core.GinHandler(server.PlugDel))
		gpPlug.Any("/runs", core.GinHandler(server.PlugRuns))
		gpPlug.Any("/log", core.GinHandler(server.PlugLog))
	}

	// 触发器路由
	{
		gpTimer := comm.Gin.Group("/trigger")
		gpTimer.Use(utilService.MidNeedLogin)
		gpTimer.Any("/list", core.GinHandler(server.TriggerList))
		gpTimer.Any("/edit", core.GinHandler(server.TriggerEdit))
		gpTimer.Any("/del", core.GinHandler(server.TriggerDel))
		gpTimer.Any("/hooks", server.TriggerHooks)
	}

	// hook
	{
		gpHook := comm.Gin.Group("/hook")
		gpHook.Any("/trigger/:trid", server.HookTrigger)
	}
}
