package utils

import (
	"gokins/comm"
	"html/template"
)

var HtmlSource = template.New("go")

func InitHtmls() error {
	HtmlSource.Funcs(comm.Gin.FuncMap)
	_, err := HtmlSource.New("head.html").Parse(`

    <meta http-equiv="X-UA-Compatible" content="IE=edge"/>
    <meta name="viewport" content="initial-scale=1,minimum-scale=1,maximum-scale=1,user-scalable=no,width=device-width"/>

`)
	if err != nil {
		return err
	}

	_, err = HtmlSource.New("test.html").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{AppName}}-测试</title>
    {{template "head.html" .}}
</head>
<body>
go 内容：{{.cont}}
</body>
</html>
`)
	if err != nil {
		return err
	}

	return nil
}
