package main

import (
	"github.com/gin-gonic/gin"
	ruisUtil "github.com/mgr9525/go-ruisutil"
	"gokins/comm"
	"gokins/route"
	"net/http"
	"os"
	"path/filepath"
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
	comm.Gin.Any("/test", func(c *gin.Context) {
		data := ruisUtil.NewMap()
		data.Set("cont", "你好啊world!")
		//utils.html(c, "test.html", data)
		c.HTML(200, "test.html", data)
	})
	route.Init()
	err = comm.Gin.Run(":8030")
	if err != nil {
		println("gin run err:" + err.Error())
	}
}
