package model

import (
	"time"
)

type SysUser struct {
	Id      int `xorm:"pk autoincr"`
	Xid     string
	Name    string
	Pass    string
	Nick    string
	Phone   string
	Times   time.Time
	Logintm time.Time
	Fwtm    time.Time
	Avat    string
}
