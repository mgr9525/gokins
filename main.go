package main

import (
	"github.com/gin-gonic/gin"
	ruisUtil "github.com/mgr9525/go-ruisutil"
	"gokins/comm"
	"gokins/utils"
	"html/template"
)

func main() {
	comm.Gin = gin.Default()
	comm.Gin.FuncMap = template.FuncMap{
		"AppName": func() string {
			return "mine app"
		},
	}
	//comm.FileView=true
	comm.Gin.LoadHTMLGlob("view/*")
	comm.Gin.Any("/test", func(c *gin.Context) {
		data := ruisUtil.NewMap()
		data.Set("cont", "你好啊world!")
		utils.Render(c, "test.html", data)
	})

	err := utils.InitHtmls()
	if err != nil {
		println("InitHtmls err:" + err.Error())
	}
	comm.Gin.Run(":8050")
}
