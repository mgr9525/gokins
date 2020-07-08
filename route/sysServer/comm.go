package sysServer

import (
	"github.com/gin-gonic/gin"
	ruisUtil "github.com/mgr9525/go-ruisutil"
	"gokins/service/sysService"
)

func CommInfo(c *gin.Context) {
	root := sysService.FindUser("admin")
	info := ruisUtil.NewMap()
	info.Set("need_set_root_pass", root.Pass == "")
	c.JSON(200, info)
}
