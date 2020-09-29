package route

import (
	"gokins/comm"
	"gokins/core"
	"gokins/route/sysServer"
)

func Init() {
	comm.Gin.Use(core.MidAccessAllow)
	gpComm := comm.Gin.Group("/comm")
	gpComm.Any("/info", sysServer.CommInfo)
	gpLogin := comm.Gin.Group("/login")
	gpLogin.Any("/info", sysServer.LoginInfo)
	gpLogin.Any("/lg", sysServer.Login)
}
