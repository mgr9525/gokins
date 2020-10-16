package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	ruisUtil "github.com/mgr9525/go-ruisutil"
	"gokins/comm"
	"gokins/core"
	"gokins/model"
	"gokins/models"
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
	if req.Types < 0 || req.Types > 3 {
		c.String(500, "param err")
		return
	}
	if err := req.Save(); err != nil {
		c.String(500, "save err:"+err.Error())
		return
	}
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
