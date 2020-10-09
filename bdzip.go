package main

import (
	"encoding/base64"
	"io/ioutil"
)

func main() {
	bts, _ := ioutil.ReadFile("uis/vue-admin/dist/1.zip")
	cont := base64.StdEncoding.EncodeToString(bts)
	ioutil.WriteFile("bin/zip.txt", []byte(cont), 0644)
}
