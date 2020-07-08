package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"gokins/comm"
	"html/template"
)

func RenderHTML(c *gin.Context, htmls string, data interface{}) {
	t := template.New("test").Funcs(comm.Gin.FuncMap)
	/*tpe:=reflect.TypeOf(c).Elem()
	vae:=reflect.ValueOf(c).Elem()
	_,ok:=tpe.FieldByName("engine")
	if ok{
		val:=vae.FieldByName("engine").Pointer()
		eg:=(*gin.Engine)(unsafe.Pointer(val))
		//println("val:",eg)
		t.Funcs(eg.FuncMap)
	}*/
	_, err := t.Parse(htmls)
	if err != nil {
		c.String(500, "errs:"+err.Error())
		return
	}

	c.Render(200, render.HTML{
		Template: t,
		Name:     "",
		Data:     data,
	})
}
