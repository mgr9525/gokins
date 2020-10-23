package model

import "time"

type TTrigger struct {
	Id     int `xorm:"pk autoincr"`
	Uid    string
	Types  string
	Title  string
	Desc   string
	Times  time.Time
	Config string //配置详情json
	Del    int
	Enable int
	Errs   string
	Mid    int
	Meid   int
}
