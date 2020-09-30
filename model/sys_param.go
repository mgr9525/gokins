package model

import (
	"time"
)

type SysParam struct {
	Id    int `xorm:"pk autoincr"`
	Key   string
	Cont  []byte
	Times time.Time
}
