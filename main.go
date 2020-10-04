package main

import (
	"gokins/comm"
	"gokins/core"
	"gokins/mgr"
	"gokins/route"
	"gokins/service/dbService"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func init() {
	path, err := os.Executable()
	if err != nil {
		println("path err:" + err.Error())
		return
	}
	println("path:" + path)
	dir, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		println("dir err:" + err.Error())
		return
	}
	println("dir:" + dir)
	comm.Path = path
	comm.Dir = dir
	comm.Gin = gin.Default()
	if len(os.Args) > 1 && os.Args[1] == "tests" {
		comm.Dir = "."
	}
	comm.Gin.StaticFS("/css", http.Dir(comm.Dir+"/ui/css"))
	comm.Gin.StaticFS("/js", http.Dir(comm.Dir+"/ui/js"))
	comm.Gin.StaticFS("/img", http.Dir(comm.Dir+"/ui/img"))
	comm.Gin.StaticFS("/fonts", http.Dir(comm.Dir+"/ui/fonts"))
	comm.Gin.StaticFile("/index.html", comm.Dir+"/ui/index.html")
	comm.Gin.StaticFile("/favicon.ico", comm.Dir+"/ui/favicon.ico")
	/*comm.Gin.FuncMap = template.FuncMap{
		"AppName": func() string {
			return "mine app"
		},
	}
	comm.Gin.LoadHTMLGlob("view/*")*/
	//comm.FileView=true
	/*comm.Gin.SetHTMLTemplate(utils.HtmlSource)
	err := utils.InitHtmls()
	if err != nil {
		println("InitHtmls err:" + err.Error())
		return
	}*/
}
func main() {
	err := comm.InitDb()
	if err != nil {
		println("InitDb err:" + err.Error())
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
	err = comm.Gin.Run(":8030")
	if err != nil {
		println("gin run err:" + err.Error())
	}
	mgr.Cancel()
}
