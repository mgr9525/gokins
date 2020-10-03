package route

import (
	"gokins/comm"
	"gokins/core"
	"gokins/route/server"
	"gokins/service/utilService"
)

func Init() {
	comm.Gin.Use(core.MidAccessAllow)
	gpComm := comm.Gin.Group("/comm")
	gpComm.Any("/info", server.CommInfo)
	gpLogin := comm.Gin.Group("/lg")
	gpLogin.Any("/info", server.LoginInfo)
	gpLogin.Any("/login", server.Login)
	gpLogin.Any("/install", server.Install)

	gpModel := comm.Gin.Group("/model")
	gpModel.Use(utilService.MidNeedLogin)
	gpModel.Any("/get", core.GinHandler(server.ModelGet))
	gpModel.Any("/list", core.GinHandler(server.ModelList))
	gpModel.Any("/edit", core.GinHandler(server.ModelEdit))
	gpModel.Any("/del", core.GinHandler(server.ModelDel))
	gpModel.Any("/runs", core.GinHandler(server.ModelRuns))
	gpModel.Any("/run", core.GinHandler(server.ModelRun))
	gpModel.Any("/stop", core.GinHandler(server.ModelStop))

	gpPlug := comm.Gin.Group("/plug")
	gpPlug.Use(utilService.MidNeedLogin)
	gpPlug.Any("/list", core.GinHandler(server.PlugList))
	gpPlug.Any("/edit", core.GinHandler(server.PlugEdit))
	gpPlug.Any("/del", core.GinHandler(server.PlugDel))
	gpPlug.Any("/runs", core.GinHandler(server.PlugRuns))
	gpPlug.Any("/log", core.GinHandler(server.PlugLog))
}
