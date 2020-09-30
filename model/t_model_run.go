package model

import (
	"time"
)

type TModelRun struct {
	Id     int `xorm:"pk autoincr"`
	Tid    int
	Uid    string
	Times  time.Time
	Timesd time.Time
	State  int
}
