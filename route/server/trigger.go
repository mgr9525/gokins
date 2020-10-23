package server

import (
	"fmt"
	"gokins/comm"
	"gokins/core"
	"gokins/mgr"
	"gokins/model"
	"gokins/models"
	"gokins/service/utilService"

	"github.com/gin-gonic/gin"
	ruisUtil "github.com/mgr9525/go-ruisutil"
)

func TriggerList(c *gin.Context, req *ruisUtil.Map) {
	pg, _ := req.GetInt("page")
	q := req.GetString("q")
	ls := make([]*model.TTrigger, 0)
	ses := comm.Db.Where("del !='1'")
	if q != "" {
		ses.And("name like ?", "%"+q+"%")
	}
	page, err := core.XormFindPage(ses, &ls, pg, 20)
	if err != nil {
		c.String(500, "find err:"+err.Error())
		return
	}
	c.JSON(200, page)
}

func TriggerEdit(c *gin.Context, req *models.Trigger) {
	if req.Types == "" || req.Title == "" {
		c.String(500, "param err")
		return
	}
	if req.Types == "worked" && req.Mid == req.Meid {
		c.String(500, "两个流水线不能相等")
		return
	}
	lgusr := utilService.CurrMUser(c)
	req.Uid = lgusr.Xid
	if err := req.Save(); err != nil {
		c.String(500, "save err:"+err.Error())
		return
	}
	mgr.TriggerMgr.Refresh(req.Id)
	c.String(200, fmt.Sprintf("%d", req.Id))
}

func TriggerDel(c *gin.Context, req *ruisUtil.Map) {
	id, err := req.GetInt("id")
	if err != nil || id <= 0 {
		c.String(500, "param err")
		return
	}
	m := &models.Trigger{}
	if err := m.Del(int(id)); err != nil {
		c.String(500, "save err:"+err.Error())
		return
	}
	c.String(200, fmt.Sprintf("%d", m.Id))
}

func TriggerHooks(c *gin.Context) {
	c.JSON(200, mgr.HookjsMap)
}
