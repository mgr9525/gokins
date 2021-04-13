package main

import (
	"flag"
	"gokins/comm"
	"gokins/core"
	"gokins/mgr"
	"gokins/route"
	"gokins/service"
	"gokins/service/dbService"
	"os"
	"path/filepath"

	ruisIo "github.com/mgr9525/go-ruisutil/ruisio"

	"github.com/gin-gonic/gin"
)

var (
	clearPass = ""
	mvData    = ""
	upgrade   = false
)

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
	flag.StringVar(&comm.Dir, "d", dir, "工作目录")
	flag.StringVar(&comm.Host, "bind", ":8030", "绑定地址")
	flag.BoolVar(&comm.NoUppass, "nps", false, "是否禁止修改密码")
	flag.IntVar(&comm.RunTaskLen, "rln", 5, "同时执行的流水线数量")
	flag.StringVar(&clearPass, "clp", "", "清除某个用户密码（请先关闭服务在执行）")
	flag.StringVar(&mvData, "mvdata", "", "转移某个库数据到本地（目前不转移运行日志记录）")
	flag.BoolVar(&upgrade, "up", false, "升级数据库（迁移流水线，触发器）")
	flag.Parse()
	comm.Gin = gin.Default()
}
func main() {
	if upgrade {
		os.RemoveAll(comm.Dir + "/data/logs")
		if !ruisIo.PathExists(comm.Dir + "/db.dat") {
			println("未找到原有数据库")
			return
		}
		os.Remove(comm.Dir + "/db.bak")
		os.Rename(comm.Dir+"/db.dat", comm.Dir+"/db.bak")
	}
	if !ruisIo.PathExists(comm.Dir + "/data") {
		err := os.MkdirAll(comm.Dir+"/data", 0755)
		if err != nil {
			println("Mkdir data err:" + err.Error())
			return
		}
	}
	err := comm.InitDb()
	if err != nil {
		println("InitDb err:" + err.Error())
		return
	}
	if clearPass != "" {
		service.ClearUPass(clearPass)
		return
	}
	if mvData != "" {
		service.MoveData(mvData)
		return
	}
	if upgrade {
		service.Upgrade()
		return
	}

	runWeb()
}
func runWeb() {
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
	mgr.TriggerMgr.Start()
	err := comm.Gin.Run(comm.Host)
	if err != nil {
		println("gin run err:" + err.Error())
	}
	mgr.Cancel()
}
