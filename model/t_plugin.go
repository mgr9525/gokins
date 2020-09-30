package model

import (
	"time"
)

type TPlugin struct {
	Id    int `xorm:"pk autoincr"`
	Tid   int
	Type  int
	Title string
	Para  string
	Cont  string
	Times time.Time
	Sort  int
	Del   int
}
