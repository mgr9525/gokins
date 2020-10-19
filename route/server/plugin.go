package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	ruisUtil "github.com/mgr9525/go-ruisutil"
	"gokins/comm"
	"gokins/model"
	"gokins/models"
	"gokins/service/dbService"
	"io/ioutil"
	"strconv"
	"time"
)

func PlugList(c *gin.Context, req *ruisUtil.Map) {
	tid, err := req.GetInt("tid")
	if err != nil || tid <= 0 {
		c.String(500, "param err")
		return
	}
	ls := make([]*model.TPlugin, 0)
	ses := comm.Db.Where("del!='1' and tid=?", tid).OrderBy("sort ASC,id ASC")
	err = ses.Find(&ls)
	if err != nil {
		c.String(500, "find err:"+err.Error())
		return
	}
	c.JSON(200, ls)
}
func PlugEdit(c *gin.Context, req *models.Plugin) {
	if req.Title == "" || req.Tid <= 0 {
		c.String(500, "param err")
		return
	}
	if err := req.Save(); err != nil {
		c.String(500, "save err:"+err.Error())
		return
	}
	c.String(200, fmt.Sprintf("%d", req.Id))
}
func PlugDel(c *gin.Context, req *ruisUtil.Map) {
	id, err := req.GetInt("id")
	if err != nil || id <= 0 {
		c.String(500, "param err")
		return
	}
	m := &models.Plugin{}
	if err := m.Del(int(id)); err != nil {
		c.String(500, "save err:"+err.Error())
		return
	}
	c.String(200, fmt.Sprintf("%d", m.Id))
}
func PlugRuns(c *gin.Context, req *ruisUtil.Map) {
	if req.GetBool("first") == false {
		time.Sleep(time.Second)
	}
	id, err := req.GetInt("id")
	if err != nil || id <= 0 {
		c.String(500, "param err")
		return
	}

	mr := dbService.GetModelRun(int(id))
	if mr == nil {
		c.String(404, "not found")
		return
	}

	ls := make([]*models.PluginUI, 0)
	ses := comm.Db.Where("del!='1' and tid=?", mr.Tid).OrderBy("sort ASC,id ASC")
	err = ses.Find(&ls)
	if err != nil {
		c.String(500, "find err:"+err.Error())
		return
	}
	for _, v := range ls {
		v.ToUI(mr.Id)
	}

	res := ruisUtil.NewMap()
	res.Set("list", ls)
	res.Set("tid", mr.Tid)
	res.Set("end", mr.State >= 2)
	plugLog(req.GetString("pid"), mr, res)
	c.JSON(200, res)
}
func plugLog(pids string, mr *model.TModelRun, ret *ruisUtil.Map) {
	pid, err := strconv.ParseInt(pids, 10, 64)
	if err != nil || pid <= 0 {
		return
	}
	rn := dbService.FindPluginRun(mr.Tid, mr.Id, int(pid))
	if rn == nil {
		return
	}
	res := ruisUtil.NewMap()
	res.Set("id", rn.Id)
	res.Set("up", true)
	res.Set("text", "")
	if rn.State >= 2 {
		res.Set("up", false)
	}
	logpth := fmt.Sprintf("%s/data/logs/%d/%d.log", comm.Dir, rn.Mid, rn.Id)
	outs, err := ioutil.ReadFile(logpth)
	if err == nil {
		res.Set("text", string(outs))
	}
	ret.Set("log", res)
}
