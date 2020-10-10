package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

func main() {
	bts, _ := ioutil.ReadFile("uis/vue-admin/dist/dist.zip")
	cont := base64.StdEncoding.EncodeToString(bts)
	ioutil.WriteFile("comm/vuefl.go",
		[]byte(fmt.Sprintf("package comm\n\nconst StaticPkg = \"%s\"", cont)),
		0644)
}
