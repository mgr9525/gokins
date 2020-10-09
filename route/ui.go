package route

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"errors"
	"gokins/comm"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var rder *zip.Reader

func getRdr() (*zip.Reader, error) {
	if rder != nil {
		return rder, nil
	}
	bts, err := base64.StdEncoding.DecodeString(comm.StaticPkg)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewReader(bts)
	r, err := zip.NewReader(buf, buf.Size())
	if err != nil {
		return nil, err
	}
	rder = r
	return rder, nil
}

func getFile(pth string) (*zip.File, error) {
	if pth == "" {
		return nil, errors.New("param err")
	}
	println("getFile:" + pth)
	r, err := getRdr()
	if err != nil {
		return nil, err
	}
	for _, f := range r.File {
		//println("f.name:", f.Name)
		if pth == f.Name {
			return f, nil
		}
	}
	return nil, errors.New("file not found")
}

func Statics(c *gin.Context) {
	pth := c.FullPath()
	if pth == "/static/*nm" {
		pth = "static" + c.Param("nm")
	}
	if pth == "/" {
		pth = "index.html"
	}
	r, err := getFile(pth)
	if err != nil {
		c.String(500, "rdr err:"+err.Error())
		return
	}
	rd, err := r.Open()
	if err != nil {
		c.String(500, "open err:"+err.Error())
		return
	}

	ext := filepath.Ext(r.Name)
	if ext == ".html" {
		c.Writer.Header().Set("Content-Type", "text/html")
	} else if ext == ".css" {
		c.Writer.Header().Set("Content-Type", "text/css")
	} else if ext == ".js" {
		c.Writer.Header().Set("Content-Type", "application/javascript")
	} else if ext == ".svg" {
		c.Writer.Header().Set("Content-Type", "text/xml")
	}
	c.Status(200)
	bts := make([]byte, 1024)
	for {
		select {
		case <-c.Done():
			goto end
		default:
			n, err := rd.Read(bts)
			if n > 0 {
				c.Writer.Write(bts[:n])
			}
			if err != nil {
				goto end
			}
		}
	}
end:
}
