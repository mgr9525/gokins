package route

import (
	"gokins/comm"
	"gokins/route/sysServer"
)

func Init() {
	gpComm := comm.Gin.Group("/comm")
	gpComm.Any("/info", sysServer.CommInfo)
	gpLogin := comm.Gin.Group("/login")
	gpLogin.Any("/info", sysServer.LoginInfo)
	gpLogin.Any("/lg", sysServer.Login)
}
