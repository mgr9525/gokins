package model

import "time"

type TTrigger struct {
	Id     int `xorm:"pk autoincr"`
	Types  int //0 : 手动触发  1:git触发 2:定时器触发
	Name   string
	Desc   string
	Times  time.Time
	Del    int
	Config string //配置详情json
}
