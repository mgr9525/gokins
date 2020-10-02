package models

import (
	"gokins/comm"
	"gokins/model"
	"time"
)

type Model struct {
	Id     int `xorm:"pk autoincr"`
	Uid    string
	Title  string
	Desc   string
	Times  time.Time
	Envs   string
	Wrkdir string
}

func (Model) TableName() string {
	return "t_model"
}

func (c *Model) Save() error {
	var err error
	if c.Id > 0 {
		_, err = comm.Db.Cols("title", "desc", "envs", "wrkdir").Where("id=?", c.Id).Update(c)
	} else {
		c.Times = time.Now()
		_, err = comm.Db.Insert(c)
	}
	return err
}
func (c *Model) Del(id int) error {
	m := &model.TModel{Del: 1}
	_, err := comm.Db.Cols("del").Where("id=?", id).Update(m)
	return err
}
