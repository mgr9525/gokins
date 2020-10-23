package server

import (
	"gokins/mgr"
	"gokins/service/dbService"
	"io/ioutil"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HookTrigger(c *gin.Context) {
	trid, err := strconv.Atoi(c.Param("trid"))
	if err != nil {
		c.String(500, "trigger id err")
		return
	}
	trg := dbService.GetTrigger(trid)
	if trg == nil || trg.Del == 1 {
		c.String(404, "trigger not found")
		return
	}
	if trg.Types != "hook" {
		c.String(404, "trigger types isn't exist")
		return
	}

	bodys, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(500, "read body err:"+err.Error())
		return
	}

	mgr.TriggerMgr.StartOne(trg, c.Request.URL.RawQuery, c.Request.Header, bodys)
	c.String(200, "ok")
}
