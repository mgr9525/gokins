package model

import (
	"time"
)

type TModel struct {
	Id    int `xorm:"pk autoincr"`
	Uid   string
	Title string
	Desc  string
	Times time.Time
	Del   int
}
