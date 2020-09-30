package utilService

import (
	"gokins/core"
	"gokins/model"
	"gokins/service/dbService"

	"github.com/gin-gonic/gin"
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
	return dbService.FindUser(xid)
}
func CurrMUser(c *gin.Context) *model.SysUser {
	tu, ok := c.Get("lguser")
	if ok {
		lguser, ok := tu.(*model.SysUser)
		if ok {
			return lguser
		}
	}
	return nil
}
