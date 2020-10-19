package models

import (
	"gokins/comm"
	"gokins/model"
	"time"
)

type Trigger struct {
	Id     int `xorm:"pk autoincr"`
	Types  int //0 : 手动触发  1:git触发 2:定时器触发
	Name   string
	Desc   string
	Times  time.Time
	Config string //配置详情json
}

func (Trigger) TableName() string {
	return "t_trigger"
}

func (c *Trigger) Save() error {
	var err error
	if c.Id > 0 {
		_, err = comm.Db.Cols("types", "name", "desc", "config").Where("id=?", c.Id).Update(c)
	} else {
		c.Times = time.Now()
		_, err = comm.Db.Insert(c)
	}
	return err
}

func (c *Trigger) Del(id int) error {
	m := &model.TTrigger{Del: 1}
	_, err := comm.Db.Cols("del").Where("id=?", id).Update(m)
	return err
}
