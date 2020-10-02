package model

import (
	"time"
)

type TModel struct {
	Id     int `xorm:"pk autoincr"`
	Uid    string
	Title  string
	Desc   string
	Times  time.Time
	Del    int
	Envs   string
	Wrkdir string
}

type TModelRun struct {
	Id     int `xorm:"pk autoincr"`
	Tid    int //model id
	Uid    string
	Times  time.Time
	Timesd time.Time
	State  int //-1已停止，0等待，1运行，2运行失败，4运行成功
	Errs   string
}
