package server

import (
	"gokins/service/sysService"

	"github.com/gin-gonic/gin"
	ruisUtil "github.com/mgr9525/go-ruisutil"
)

func CommInfo(c *gin.Context) {
	root := sysService.FindUser("admin")
	info := ruisUtil.NewMap()
	info.Set("need_install", root.Pass == "")
	c.JSON(200, info)
}
