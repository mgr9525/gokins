package sysServer

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	ruisUtil "github.com/mgr9525/go-ruisutil"
	"gokins/core"
	"gokins/service/sysService"
	"gokins/service/utilService"
	"time"
)

func LoginInfo(c *gin.Context) {
	rets := ruisUtil.NewMap()
	rets.Set("login", false)
	lguser := utilService.CurrUser(c)
	if lguser != nil {
		rets.Set("login", true)
		rets.Set("xid", lguser.Xid)
		rets.Set("name", lguser.Name)
		rets.Set("nick", lguser.Nick)
		rets.Set("avat", lguser.Avat)
	}

	c.JSON(200, rets)
}

func Login(c *gin.Context) {
	pars, err := core.BindMapJSON(c)
	if err != nil {
		c.String(500, "bind err:"+err.Error())
		return
	}
	name := pars.GetString("name")
	pass := pars.GetString("pass")
	if name == "" || pass == "" {
		c.String(500, "param err!")
		return
	}
	usr := sysService.FindUserName(name)
	if usr == nil {
		c.String(511, "未找到用户!")
		return
	}
	if usr.Pass != ruisUtil.Md5String(pass) {
		c.String(512, "密码错误!")
		return
	}

	tks, err := core.CreateToken(&jwt.MapClaims{
		"xid": usr.Xid,
	}, time.Hour*5)
	if err != nil {
		c.String(513, "获取token失败!")
		return
	}

	c.String(200, tks)
}
