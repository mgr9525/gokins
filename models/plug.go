package models

import (
	"errors"
	"gokins/comm"
	"gokins/model"
	"time"
)

type Plugin struct {
	Id    int `xorm:"pk autoincr"`
	Tid   int
	Type  int
	Title string
	Para  string
	Cont  string
	Times time.Time
	Sort  int
	Exend int
}

func (Plugin) TableName() string {
	return "t_plugin"
}

func (c *Plugin) Save() error {
	var err error
	if c.Id > 0 {
		_, err = comm.Db.Cols("type", "title", "para", "cont", "sort", "exend").Where("id=?", c.Id).Update(c)
	} else {
		if c.Tid <= 0 {
			return errors.New("what?")
		}
		c.Times = time.Now()
		_, err = comm.Db.Insert(c)
	}
	return err
}
func (c *Plugin) Del(id int) error {
	m := &model.TPlugin{Del: 1}
	_, err := comm.Db.Cols("del").Where("id=?", id).Update(m)
	return err
}
