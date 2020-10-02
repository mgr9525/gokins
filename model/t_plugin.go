package model

import (
	"time"
)

type TPlugin struct {
	Id    int `xorm:"pk autoincr"`
	Tid   int //model id
	Type  int
	Title string
	Para  string
	Cont  string
	Times time.Time
	Sort  int
	Del   int
	Exend int
}

type TPluginRun struct {
	Id     int `xorm:"pk autoincr"`
	Pid    int //plugin id
	Mid    int //model id
	Tid    int //modelRun id
	Times  time.Time
	Timesd time.Time
	State  int //-1已停止，0等待，1运行，2运行失败，4运行成功
	Excode int
	Output string
}
