package server

import (
	"fmt"
	"gokins/comm"
	"gokins/core"
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
		c.String(511, "参数错误")
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
