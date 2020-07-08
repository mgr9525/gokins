package models

import (
	"time"
)

type SysUser struct {
	Id      int       `xorm:"pk autoincr BIGINT(20)"`
	Xid     string    `xorm:"not null pk VARCHAR(64)"`
	Name    string    `xorm:"not null pk VARCHAR(50)"`
	Pass    string    `xorm:"VARCHAR(100)"`
	Nick    string    `xorm:"VARCHAR(100)"`
	Phone   string    `xorm:"comment('备用电话') index VARCHAR(50)"`
	Times   time.Time `xorm:"default 'CURRENT_TIMESTAMP' comment('创建时间') DATETIME"`
	Logintm time.Time `xorm:"comment('登陆时间') DATETIME"`
	Fwtm    time.Time `xorm:"comment('访问时间') DATETIME"`
	Avat    string    `xorm:"VARCHAR(500)"`
}
