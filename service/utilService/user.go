package utilService

import (
	"github.com/gin-gonic/gin"
	"gokins/core"
	"gokins/model"
	"gokins/service/sysService"
)

func CurrUser(c *gin.Context) *model.SysUser {
	tks := core.GetToken(c)
	if tks == nil {
		return nil
	}
	xid, ok := tks["xid"].(string)
	if !ok || xid == "" {
		return nil
	}
	return sysService.FindUser(xid)
}
