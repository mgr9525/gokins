package route

import (
	"gokins/comm"
	"gokins/route/sysServer"
)

func Init() {
	gpComm := comm.Gin.Group("/comm")
	gpComm.Any("/info", sysServer.CommInfo)
}
