package comm

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

var (
	Dir  string
	Path string
	Gin  *gin.Engine
	Db   *xorm.Engine

	RunTaskLen int = 5
)
