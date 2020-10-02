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

func ModelList(c *gin.Context, req *ruisUtil.Map) {
	pg, _ := req.GetInt("page")
	q := req.GetString("q")
	ls := make([]*model.TModel, 0)
	ses := comm.Db.Where("del!='1'")
	if q != "" {
		ses.And("title like ?", "%"+q+"%")
	}
	page, err := core.XormFindPage(ses, &ls, pg, 20)
	if err != nil {
		c.String(500, "find err:"+err.Error())
		return
	}
	c.JSON(200, page)
}
func ModelEdit(c *gin.Context, req *models.Model) {
	if req.Title == "" {
		c.String(500, "param err")
		return
	}
	lguser := utilService.CurrMUser(c)
	req.Uid = lguser.Xid
	if err := req.Save(); err != nil {
		c.String(500, "save err:"+err.Error())
		return
	}
	c.String(200, fmt.Sprintf("%d", req.Id))
}
func ModelDel(c *gin.Context, req *ruisUtil.Map) {
	id, err := req.GetInt("id")
	if err != nil || id <= 0 {
		c.String(500, "param err")
		return
	}
	m := &models.Model{}
	if err := m.Del(int(id)); err != nil {
		c.String(500, "save err:"+err.Error())
		return
	}
	c.String(200, fmt.Sprintf("%d", m.Id))
}

func ModelRuns(c *gin.Context, req *ruisUtil.Map) {
	pg, _ := req.GetInt("page")
	tid, err := req.GetInt("tid")
	if err != nil || tid <= 0 {
		c.String(500, "param err")
		return
	}
	ls := make([]*models.ModelRun, 0)
	ses := comm.Db.Where("tid=?", tid).OrderBy("id DESC")
	page, err := core.XormFindPage(ses, &ls, pg, 20)
	if err != nil {
		c.String(500, "find err:"+err.Error())
		return
	}
	for _, v := range ls {
		v.ToUI()
	}
	c.JSON(200, page)
}
func ModelRun(c *gin.Context, req *ruisUtil.Map) {
	id, err := req.GetInt("id")
	if err != nil || id <= 0 {
		c.String(500, "param err")
		return
	}
	lgusr := utilService.CurrMUser(c)
	m := &models.ModelRun{}
	m.Tid = int(id)
	m.Uid = lgusr.Xid
	if err := m.Add(); err != nil {
		c.String(500, "add err:"+err.Error())
		return
	}
	mgr.ExecMgr.Refresh()
	c.String(200, fmt.Sprintf("%d", m.Id))
}
func ModelStop(c *gin.Context, req *ruisUtil.Map) {
	id, err := req.GetInt("id")
	if err != nil || id <= 0 {
		c.String(500, "param err")
		return
	}
	mgr.ExecMgr.StopTask(int(id))
	c.String(200, "ok")
}
