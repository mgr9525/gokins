package server

import (
	"gokins/service/dbService"

	"github.com/gin-gonic/gin"
	ruisUtil "github.com/mgr9525/go-ruisutil"
)

func CommInfo(c *gin.Context) {
	root := dbService.FindUser("admin")
	info := ruisUtil.NewMap()
	info.Set("need_install", root == nil || root.Pass == "")
	c.JSON(200, info)
}
