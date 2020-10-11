package main

import (
	"flag"
	"fmt"
	"gokins/comm"
	"gokins/core"
	"gokins/mgr"
	"gokins/route"
	"gokins/service/dbService"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var clearPass = ""

func init() {
	path, err := os.Executable()
	if err != nil {
		println("path err:" + err.Error())
		return
	}
	println("path:" + path)
	comm.Path = path
	dir, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		println("dir err:" + err.Error())
		return
	}
	println("dir:" + dir)
	flag.StringVar(&comm.Dir, "d", dir, "数据目录")
	flag.StringVar(&comm.Host, "bind", ":8030", "绑定地址")
	flag.StringVar(&clearPass, "clp", "", "清除用户密码")
	flag.Parse()
	comm.Gin = gin.Default()
}
func main() {
	err := comm.InitDb()
	if err != nil {
		println("InitDb err:" + err.Error())
		return
	}
	if clearUPass() {
		return
	}
	jwtKey := dbService.GetParam("jwt-key")
	jkey := jwtKey.GetString("key")
	if jkey == "" {
		jkey = core.RandomString(32)
		jwtKey.Set("key", jkey)
		dbService.SetParam("jwt-key", jwtKey)
	}
	core.JwtKey = jkey
	route.Init()
	mgr.ExecMgr.Start()
	err = comm.Gin.Run(comm.Host)
	if err != nil {
		println("gin run err:" + err.Error())
	}
	mgr.Cancel()
}

func clearUPass() bool {
	if clearPass != "" {
		usr := dbService.FindUserName(clearPass)
		if usr == nil {
			fmt.Printf("user(%s) not found\n", clearPass)
		} else {
			usr.Pass = ""
			_, err := comm.Db.Cols("pass").Where("id=?", usr.Id).Update(usr)
			if err != nil {
				fmt.Println("clear password err:" + err.Error())
			} else {
				fmt.Println("clear password ok")
			}
		}
		return true
	}
	return false
}
