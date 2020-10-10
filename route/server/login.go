package server

import (
	"gokins/comm"
	"gokins/core"
	"gokins/service/dbService"
	"gokins/service/utilService"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	ruisUtil "github.com/mgr9525/go-ruisutil"
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

func Login(c *gin.Context, req *ruisUtil.Map) {
	name := req.GetString("name")
	pass := req.GetString("pass")
	if name == "" || pass == "" {
		c.String(500, "param err!")
		return
	}
	usr := dbService.FindUserName(name)
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
func Install(c *gin.Context, req *ruisUtil.Map) {
	pass := req.GetString("newpass")
	if pass == "" {
		c.String(500, "param err!")
		return
	}
	usr := dbService.FindUser("admin")
	if usr == nil {
		c.String(511, "未找到用户!")
		return
	}
	if usr.Pass != "" {
		c.String(512, "what??!")
		return
	}
	usr.Pass = ruisUtil.Md5String(pass)
	_, err := comm.Db.Cols("pass").Where("id=?", usr.Id).Update(usr)
	if err != nil {
		c.String(511, "服务错误!")
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

func Uppass(c *gin.Context, req *ruisUtil.Map) {
	pass := req.GetString("pass")
	newpass := req.GetString("newpass")
	if pass == "" || newpass == "" {
		c.String(500, "param err!")
		return
	}
	usr := dbService.FindUser("admin")
	if usr == nil {
		c.String(511, "未找到用户!")
		return
	}
	if usr.Pass != ruisUtil.Md5String(pass) {
		c.String(512, "旧密码错误!")
		return
	}
	usr.Pass = ruisUtil.Md5String(newpass)
	_, err := comm.Db.Cols("pass").Where("id=?", usr.Id).Update(usr)
	if err != nil {
		c.String(511, "服务错误!")
		return
	}

	c.String(200, "ok")
}
