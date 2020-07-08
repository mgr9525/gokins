package main

import (
	"github.com/gin-gonic/gin"
	ruisUtil "github.com/mgr9525/go-ruisutil"
	"gokins/comm"
	"gokins/utils"
	"html/template"
)

var testHTML = `
<html>
<head><title>{{AppName}}-测试</title></head>
<body>内容：{{.cont}}</body>
</html>
`

func main() {
	comm.Gin = gin.Default()
	comm.Gin.FuncMap = template.FuncMap{
		"AppName": func() string {
			return "mine app"
		},
	}
	comm.Gin.Any("/ping", func(c *gin.Context) {
		data := ruisUtil.NewMap()
		data.Set("cont", "你好啊world!")
		utils.RenderHTML(c, testHTML, data)
	})
	comm.Gin.Run(":8050")
}
