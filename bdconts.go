package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

func main() {
	bdsqls()
	bdzip()
}

func bdsqls() {
	bts, _ := ioutil.ReadFile("doc/sys.sql")
	ioutil.WriteFile("comm/dbfl.go",
		[]byte(fmt.Sprintf("package comm\n\nconst sqls = `\n%s\n`", string(bts))),
		0644)
}
func bdzip() {
	bts, _ := ioutil.ReadFile("uis/vue-admin/dist/dist.zip")
	cont := base64.StdEncoding.EncodeToString(bts)
	ioutil.WriteFile("comm/vuefl.go",
		[]byte(fmt.Sprintf("package comm\n\nconst StaticPkg = \"%s\"", cont)),
		0644)
}
